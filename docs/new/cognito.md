# AWS Cognito 設計書

## 📋 概要

このドキュメントは、Okusuri アプリケーションの認証システムを PostgreSQL のカスタム認証から AWS Cognito に移行する際の詳細設計書です。

**重要**: 認証・ユーザー管理を Cognito に委譲することで、セキュリティの向上と運用負荷の軽減を実現します。

## 🎯 設計方針

### **基本方針**

- **Cognito User Pool**: ユーザー認証・管理の中心
- **Google OAuth 統合**: **唯一の認証方法**（ユーザー名・パスワード認証は無効）
- **JWT トークン**: フロントエンド・バックエンド間の認証
- **セキュアな設計**: 最小権限の原則に基づく IAM 設定

### **技術スタック**

- **認証サービス**: AWS Cognito User Pool
- **OAuth プロバイダー**: Google OAuth 2.0
- **トークン形式**: JWT (JSON Web Token)
- **統合**: AWS SDK for Go v2

## 🏗️ アーキテクチャ設計

### **全体構成**

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   API Gateway   │    │   Backend       │    │   AWS Cognito   │
│   (Next.js)     │◄──►│   + Cognito     │◄──►│   (Go + Gin)    │◄──►│   User Pool     │
│                 │    │   Authorizer     │    │                 │    │                 │
│ - Google OAuth  │    │ - 認証・認可     │    │ - ビジネスロジック│    │ - ユーザー管理  │
│ - JWT保存       │    │ - レート制限     │    │ - DynamoDB操作  │    │ - 認証フロー    │
│ - 認証状態管理  │    │ - CORS制御       │    │ - ユーザー情報  │    │ - トークン管理  │
└─────────────────┘    └─────────────────┘    └─────────────────┘    └─────────────────┘
```

### **認証フロー**

#### **1. ログインフロー**

```
1. ユーザーがGoogle OAuthでログイン
2. CognitoがGoogleから認証情報を取得
3. Cognitoがユーザープールにユーザーを作成/更新
4. CognitoがJWTトークン（ID, Access, Refresh）を発行
5. フロントエンドがトークンを保存
6. 以降のAPIリクエストでJWTを使用
```

#### **2. API 認証フロー（API Gateway 経由）**

```
1. フロントエンドがJWTをAuthorizationヘッダーに設定
2. API GatewayがCognito User Pool AuthorizerでJWTを検証
3. 認証成功時のみ、バックエンドにリクエストを転送
4. バックエンドは認証済みリクエストのみ受信（JWT検証不要）
5. JWTからユーザー情報（sub, email等）を抽出
6. アプリケーションデータをDynamoDBから取得
```

## 🔧 Cognito User Pool 設定

### **基本設定**

```yaml
UserPoolName: 'okusuri-user-pool'
AutoVerifiedAttributes: ['email']
UsernameAttributes: ['email']
MfaConfiguration: 'OFF'
AccountRecoverySetting:
  RecoveryMechanisms:
    - Name: 'verified_email'
      Priority: 1
```

### **ユーザー属性設定**

#### **標準属性**

```yaml
StandardAttributes:
  email:
    Required: true
    Mutable: true
  name:
    Required: false
    Mutable: true
  picture:
    Required: false
    Mutable: true
  email_verified:
    Required: false
    Mutable: false
```

### **トークン設定**

```yaml
TokenValidityUnits:
  AccessToken: 'hours'
  IdToken: 'hours'
  RefreshToken: 'days'

TokenValidity:
  AccessToken: 1 # 1時間
  IdToken: 1 # 1時間
  RefreshToken: 30 # 30日
```

## 🌐 アプリクライアント設定

### **フロントエンド用クライアント**

```yaml
ClientName: 'okusuri-app'
GenerateSecret: false
# Google OAuthのみ許可（ユーザー名・パスワード認証は無効）
ExplicitAuthFlows:
  - 'ALLOW_REFRESH_TOKEN_AUTH' # リフレッシュトークンのみ

SupportedIdentityProviders:
  - 'Google' # Google OAuthのみ

# Cognitoのユーザー名・パスワード認証は無効
# - 'ALLOW_USER_PASSWORD_AUTH'
# - 'ALLOW_USER_SRP_AUTH'
# - 'ALLOW_ADMIN_USER_PASSWORD_AUTH'

CallbackURLs:
  - 'http://localhost:3000/auth/callback'
  - 'https://yourdomain.com/auth/callback'

LogoutURLs:
  - 'http://localhost:3000/'
  - 'https://yourdomain.com/'

AllowedOAuthFlows:
  - 'code' # Authorization Code Flowのみ

AllowedOAuthScopes:
  - 'email'
  - 'openid'
  - 'profile'

AllowedOAuthFlowsUserPoolClient: true
```

## 🔗 Google OAuth 統合

### **Google OAuth 設定**

```yaml
IdentityProvider:
  ProviderName: 'Google'
  ProviderType: 'Google'
  ProviderDetails:
    client_id: '${GOOGLE_CLIENT_ID}'
    client_secret: '${GOOGLE_CLIENT_SECRET}'
    authorize_scopes: 'email profile openid'
    attributes_request_method: 'GET'
    oidc_issuer: 'https://accounts.google.com'
    authorize_url: 'https://accounts.google.com/o/oauth2/v2/auth'
    token_url: 'https://oauth2.googleapis.com/token'
    attributes_url: 'https://www.googleapis.com/oauth2/v3/userinfo'
    jwks_uri: 'https://www.googleapis.com/oauth2/v3/certs'
```

### **属性マッピング**

```yaml
AttributeMapping:
  email: 'email'
  email_verified: 'email_verified'
  name: 'name'
  picture: 'picture'
  given_name: 'given_name'
  family_name: 'family_name'
```

### **認証フロー（Google OAuth のみ）**

#### **1. ログインフロー**

```
1. ユーザーがGoogle OAuthでログイン
2. CognitoがGoogleから認証情報を取得
3. Cognitoがユーザープールにユーザーを作成/更新
4. CognitoがJWTトークン（ID, Access, Refresh）を発行
5. フロントエンドがトークンを保存
6. 以降のAPIリクエストでJWTを使用
```

#### **2. API 認証フロー（API Gateway 経由）**

```
1. フロントエンドがJWTをAuthorizationヘッダーに設定
2. API GatewayがCognito User Pool AuthorizerでJWTを検証
3. 認証成功時のみ、バックエンドにリクエストを転送
4. バックエンドは認証済みリクエストのみ受信（JWT検証不要）
5. JWTからユーザー情報（sub, email等）を抽出
6. アプリケーションデータをDynamoDBから取得
```

## 🌐 API Gateway 設定

### **基本設定**

```yaml
API:
  Name: 'okusuri-api'
  Description: 'Okusuri application API'
  ProtocolType: 'HTTP'
  CorsConfiguration:
    AllowOrigins:
      - 'http://localhost:3000'
      - 'https://yourdomain.com'
    AllowMethods:
      - 'GET'
      - 'POST'
      - 'PUT'
      - 'DELETE'
      - 'OPTIONS'
    AllowHeaders:
      - 'Content-Type'
      - 'Authorization'
      - 'X-Requested-With'
    AllowCredentials: true
    MaxAge: 86400
```

### **Cognito User Pool Authorizer**

```yaml
Authorizer:
  Name: 'okusuri-cognito-authorizer'
  Type: 'COGNITO_USER_POOLS'
  IdentitySource: 'method.request.header.Authorization'
  UserPoolArn: 'arn:aws:cognito-idp:ap-northeast-1:123456789012:userpool/ap-northeast-1_xxxxx'
  UserPoolClientId: 'client-id'
  AuthorizerResultTtlInSeconds: 300
```

### **API ルート設定**

```yaml
Routes:
  # 認証不要のエンドポイント
  - Path: '/health'
    Method: 'GET'
    Authorization: 'NONE'
    Integration: 'health-check'

  # 認証必須のエンドポイント
  - Path: '/api/medication'
    Method: 'GET'
    Authorization: 'COGNITO_USER_POOLS'
    Integration: 'medication-service'

  - Path: '/api/medication'
    Method: 'POST'
    Authorization: 'COGNITO_USER_POOLS'
    Integration: 'medication-service'

  - Path: '/api/notification'
    Method: 'GET'
    Authorization: 'COGNITO_USER_POOLS'
    Integration: 'notification-service'
```

### **統合設定**

```yaml
Integrations:
  - Name: 'medication-service'
    Type: 'HTTP_PROXY'
    IntegrationUri: 'https://your-backend-service.amazonaws.com/api/medication'
    # または AWS Lambda の場合
    # IntegrationUri: 'arn:aws:lambda:ap-northeast-1:123456789012:function:okusuri-backend'
    IntegrationMethod: 'ANY'
    ConnectionType: 'INTERNET'

  - Name: 'notification-service'
    Type: 'HTTP_PROXY'
    IntegrationUri: 'https://your-backend-service.amazonaws.com/api/notification'
    # または AWS Lambda の場合
    # IntegrationUri: 'arn:aws:lambda:ap-northeast-1:123456789012:function:okusuri-backend'
    IntegrationMethod: 'ANY'
    ConnectionType: 'INTERNET'
```

## 🔐 JWT トークン設計

### **トークン構造**

#### **ID Token (JWT)**

```json
{
  "sub": "cognito-user-sub",
  "aud": "client-id",
  "email": "user@example.com",
  "email_verified": true,
  "name": "ユーザー名",
  "picture": "https://...",
  "iat": 1693382400,
  "exp": 1693386000,
  "iss": "https://cognito-idp.ap-northeast-1.amazonaws.com/ap-northeast-1_xxxxx"
}
```

#### **Access Token (JWT)**

```json
{
  "sub": "cognito-user-sub",
  "aud": "client-id",
  "scope": "openid email profile",
  "iat": 1693382400,
  "exp": 1693386000,
  "iss": "https://cognito-idp.ap-northeast-1.amazonaws.com/ap-northeast-1_xxxxx"
}
```

### **重要な変更点: API Gateway 経由統一の認証方式**

#### **認証方式**

**全環境共通**: API Gateway + Cognito User Pool Authorizer

- **開発環境**: API Gateway → Backend(localhost:8080)
- **本番環境**: API Gateway → Backend(AWS Lambda/ECS)

**バックエンドでの JWT 検証は一切不要**

#### **バックエンドでの処理**

```go
// API Gatewayが事前に認証済みのため、JWT検証は不要
func (h *Handler) GetMedicationLogs(c *gin.Context) {
    // API Gatewayが設定する認証情報を取得
    userID := c.GetHeader("X-User-ID")        // Cognito User Sub
    email := c.GetHeader("X-User-Email")      // ユーザーメールアドレス

    // ビジネスロジックのみ実行
    logs, err := h.medicationService.GetLogsByUserID(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, logs)
}
```

#### **JWT 関連の完全削除**

**バックエンドに実装しない項目**:

- JWT 検証ロジック
- JWT 解析処理
- 認証ミドルウェア
- JWT ライブラリ依存関係

**API Gateway が提供する情報のみ使用**:

- ユーザー ID（Cognito User Sub）
- メールアドレス
- その他の認証済み情報

## 🛡️ セキュリティ設定

### **IAM ロール・ポリシー**

#### **認証済みユーザー用ロール**

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "dynamodb:GetItem",
        "dynamodb:PutItem",
        "dynamodb:UpdateItem",
        "dynamodb:DeleteItem",
        "dynamodb:Query",
        "dynamodb:Scan"
      ],
      "Resource": [
        "arn:aws:dynamodb:ap-northeast-1:*:table/okusuri-table",
        "arn:aws:dynamodb:ap-northeast-1:*:table/okusuri-table/index/*"
      ],
      "Condition": {
        "StringEquals": {
          "dynamodb:LeadingKeys": ["${cognito-identity.amazonaws.com:sub}"]
        }
      }
    }
  ]
}
```

#### **未認証ユーザー用ロール**

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Deny",
      "Action": "*",
      "Resource": "*"
    }
  ]
}
```

### **CORS 設定**

```yaml
CORS:
  AllowOrigins:
    - 'http://localhost:3000'
    - 'https://yourdomain.com'
  AllowMethods:
    - 'GET'
    - 'POST'
    - 'PUT'
    - 'DELETE'
    - 'OPTIONS'
  AllowHeaders:
    - 'Content-Type'
    - 'Authorization'
    - 'X-Requested-With'
  AllowCredentials: true
  MaxAge: 86400
```

## 🔄 移行戦略（個人用簡素版）

### **基本方針**

**個人用アプリケーションのため、複雑なマイグレーションは不要**

- 既存データは必要に応じて手動で再入力
- ユーザーは初回 Google OAuth ログイン時に自動で Cognito に作成
- シンプルな再作成ベースの移行

### **フェーズ 1: インフラ環境構築**

1. **Cognito User Pool 作成**

   - 基本設定
   - Google OAuth 統合
   - アプリクライアント設定

2. **DynamoDB テーブル作成**

   - 単一テーブル設計
   - GSI 設定
   - IAM 権限設定

3. **API Gateway 設定**
   - Cognito User Pool Authorizer
   - ルート設定
   - 統合設定

### **フェーズ 2: フロントエンド統合**

1. **認証フロー変更**

   - Google OAuth → Cognito フロー
   - JWT トークン管理
   - 認証状態管理

2. **API クライアント更新**
   - 本番 API Gateway エンドポイント使用
   - Authorization ヘッダー設定

### **フェーズ 3: 動作確認・データ再入力**

1. **基本動作確認**

   - ログイン・ログアウト
   - API 呼び出し
   - 認証フロー

2. **必要データの再入力**
   - 服用履歴（必要に応じて）
   - 通知設定
   - その他の個人設定

### **移行時の注意点**

#### **1. データの扱い**

- **既存データ**: バックアップは取るが、移行は行わない
- **新データ**: 必要に応じて手動で入力
- **ユーザー情報**: Google OAuth から自動取得

#### **2. ダウンタイム**

- **最小限**: フロントエンドの切り替えのみ
- **段階的**: 動作確認後に本番切り替え

#### **3. ロールバック**

- **簡単**: フロントエンドの設定を戻すだけ
- **リスク**: 最小限

## 📊 監視・ログ設定

### **CloudWatch メトリクス**

```yaml
Metrics:
  - AuthenticationSuccesses
  - AuthenticationFailures
  - SignUpSuccesses
  - SignUpFailures
  - TokenRefreshSuccesses
  - TokenRefreshFailures
  - UserPoolQuota
  - UserPoolQuotaUsage
```

### **CloudWatch ログ**

```yaml
Logs:
  - UserPoolLogs
  - UserPoolEvents
  - AuthenticationLogs
  - TokenLogs
```

### **アラート設定**

```yaml
Alerts:
  - AuthenticationFailureRate:
      Threshold: 5.0
      Period: 300
      EvaluationPeriods: 2
      ComparisonOperator: 'GreaterThanThreshold'

  - UserPoolQuotaUsage:
      Threshold: 80.0
      Period: 300
      EvaluationPeriods: 1
      ComparisonOperator: 'GreaterThanThreshold'
```

## 💰 コスト管理

### **料金体系**

- **User Pool**: 月額 $0.0055/MAU（月間アクティブユーザー）
- **認証**: $0.0055/認証
- **MFA**: $0.06/認証
- **詳細分析**: $0.15/MAU

### **コスト最適化**

1. **MAU 削減**

   - 非アクティブユーザーの定期的なクリーンアップ
   - セッションタイムアウトの最適化

2. **認証回数削減**
   - トークン有効期限の延長
   - リフレッシュトークンの効率的な活用

## 📝 次のステップ

1. **Terraform での Cognito + API Gateway リソース定義**
2. **フロントエンド認証フローの実装**
3. **本番環境での動作確認**
4. **必要データの手動再入力**

---

_このドキュメントは移行計画の第 2 段階「インフラ移行」の一部として作成されました。_
_作成日: 2025 年 8 月 30 日_
_更新日: 2025 年 8 月 30 日_
