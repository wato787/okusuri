# DynamoDB 設計書

## 📋 概要

このドキュメントは、Okusuri アプリケーションのデータベースを PostgreSQL から DynamoDB に移行する際の詳細設計書です。

**重要**: 認証システムとユーザー基本情報は AWS Cognito に移行するため、認証関連のテーブル（session, account, verification, user）は DynamoDB には含まれません。DynamoDB では**単一テーブル**内でアプリケーション固有データ（服用履歴、通知設定）のみを管理します。

## 🎯 設計方針

### **基本方針**

- **単一テーブル設計**: 1 つのテーブル内で複数エンティティを管理し、パフォーマンスとスケーラビリティを最適化
- **guregu/dynamo 使用**: Go 言語での使いやすさとパフォーマンスの両立
- **アクセスパターン最適化**: 現在のクエリパターンに基づいたキー設計
- **Cognito 連携**: ユーザー認証・基本情報は Cognito に委譲、アプリケーションデータのみ管理

### **技術スタック**

- **データベース**: Amazon DynamoDB
- **認証・ユーザー管理**: AWS Cognito User Pool
- **Go ライブラリ**: `github.com/guregu/dynamo/v2`
- **インフラ**: Terraform 管理

## 🗄️ テーブル設計

### **メインテーブル: `okusuri-table`**

#### **テーブル構造**

```go
type OkusuriTable struct {
    // プライマリキー
    PK        string    `dynamo:"PK,hash"`                    // Partition Key
    SK        string    `dynamo:"SK,range"`                   // Sort Key

    // インデックス用フィールド（DateIndexのみ）
    Date      string    `dynamo:"Date,index:DateIndex,hash"`   // GSI1: 日付検索

    // データフィールド
    Data      map[string]interface{} `dynamo:"Data"`           // エンティティ固有のデータ

    // メタデータ
    CreatedAt string    `dynamo:"CreatedAt"`                   // 作成日時 (ISO8601)
    UpdatedAt string    `dynamo:"UpdatedAt"`                   // 更新日時 (ISO8601)
    TTL       *int64    `dynamo:"TTL,omitempty"`               // TTL（必要に応じて）
}
```

#### **キー設計（単一テーブル内の異なるエンティティ）**

##### **1. 服用履歴**

```
PK: "USER#{cognitoUserId}"
SK: "MEDICATION#{date}#{id}"
Type: "medication_log"
Data: {
    "hasBleeding": false,
    "createdAt": "2025-08-30T10:00:00Z",
    "updatedAt": "2025-08-30T10:00:00Z"
}
```

##### **2. 通知設定**

```
PK: "USER#{cognitoUserId}"
SK: "NOTIFICATION#{platform}"
Type: "notification_setting"
Data: {
    "platform": "web",
    "isEnabled": true,
    "subscription": "webpush_subscription_json",
    "createdAt": "2025-08-30T10:00:00Z",
    "updatedAt": "2025-08-30T10:00:00Z"
}
```

## 🔍 GSI（Global Secondary Index）設計

### **GSI1: DateIndex（唯一必要な GSI）**

```
PK: Date
SK: UserID#MEDICATION#{id}
用途: 日付ベースの服用履歴検索
例: 連続服用日数の計算、特定日付の履歴取得
```

**注意**: 他の GSI は不要です。理由は以下の通り：

- **TypeIndex**: 全服用履歴・全通知設定の取得は実用的でない
- **UserIndex**: ユーザー別データは PK（USER#{cognitoUserId}）で十分取得可能
- **DateIndex**: 連続服用日数計算で必要（日付順ソート）

## 📊 アクセスパターンとクエリ例

### **1. ユーザー情報取得（Cognito から取得）**

#### **Cognito User ID からユーザー情報取得**

```go
// ユーザー基本情報は Cognito から取得
func (r *UserRepository) GetUserByCognitoID(cognitoUserID string) (*model.User, error) {
    // Cognito の AdminGetUser API を呼び出し
    cognitoUser, err := r.cognitoClient.AdminGetUser(cognitoUserID)
    if err != nil {
        return nil, err
    }

    return unmarshalCognitoUser(cognitoUser), nil
}
```

#### **メールアドレスでユーザー検索（Cognito から取得）**

```go
// メールアドレスでの検索は Cognito の ListUsers API を使用
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
    // Cognito の ListUsers API を呼び出し
    cognitoUser, err := r.cognitoClient.ListUsers(email)
    if err != nil {
        return nil, err
    }

    return unmarshalCognitoUser(cognitoUser), nil
}
```

### **2. 服用履歴関連**

#### **ユーザーの服用履歴一覧取得**

```go
func (r *MedicationRepository) GetLogsByUserID(cognitoUserID string) ([]model.MedicationLog, error) {
    var results []OkusuriTable
    err := r.table.Get("PK", "USER#"+cognitoUserID).
        Filter("begins_with(SK, 'MEDICATION#')").
        All(&results)
    if err != nil {
        return nil, err
    }

    return unmarshalMedicationLogs(results), nil
}
```

#### **連続服用日数の計算**

```go
func (r *MedicationRepository) GetConsecutiveDays(cognitoUserID string) (int, error) {
    var results []OkusuriTable
    err := r.table.Get("PK", "USER#"+cognitoUserID).
        Filter("begins_with(SK, 'MEDICATION#')").
        All(&results)
    if err != nil {
        return 0, err
    }

    // 日付順でソートして連続日数を計算
    return calculateConsecutiveDays(results), nil
}
```

### **3. 通知設定関連**

#### **ユーザーの通知設定取得**

```go
func (r *NotificationRepository) GetSetting(cognitoUserID string, platform string) (*model.NotificationSetting, error) {
    var result OkusuriTable
    err := r.table.Get("PK", "USER#"+cognitoUserID).Range("SK", "NOTIFICATION#"+platform).One(&result)
    if err != nil {
        return nil, err
    }

    return unmarshalNotificationSetting(result), nil
}
```

## 🔐 Cognito 連携の詳細

### **認証フローの変更**

#### **現在のフロー（PostgreSQL + カスタム認証）**

```
1. Google OAuth → カスタムJWT生成 → セッションテーブル保存
2. リクエスト時: JWT検証 → セッションテーブル参照 → ユーザー情報取得
```

#### **新しいフロー（Cognito + DynamoDB）**

```
1. Google OAuth → Cognito Identity Provider → Cognito User Pool
2. リクエスト時: Cognito JWT検証 → Cognitoからユーザー基本情報取得 → DynamoDBからアプリケーションデータ取得
```

### **Cognito 設定項目**

#### **User Pool 設定**

- **認証プロバイダー**: Google OAuth
- **ユーザー属性**: email, name, picture, email_verified
- **アプリクライアント**: フロントエンド用
- **トークン有効期限**: アクセストークン（1 時間）、リフレッシュトークン（30 日）

#### **Identity Pool 設定**

- **認証プロバイダー**: Cognito User Pool
- **IAM ロール**: 認証済み・未認証ユーザー用
- **DynamoDB アクセス権限**: 認証済みユーザーのみ

## 🔧 実装詳細

### **依存関係**

```go
import (
    "github.com/guregu/dynamo/v2"
    "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)
```

### **テーブル初期化**

```go
type Repository struct {
    table         dynamo.Table
    cognitoClient *cognitoidentityprovider.Client
}

func NewRepository(db *dynamo.DB, cognitoClient *cognitoidentityprovider.Client) *Repository {
    return &Repository{
        table:         db.Table("okusuri-table"),
        cognitoClient: cognitoClient,
    }
}
```

### **データマーシャリング**

```go
func unmarshalCognitoUser(cognitoUser *cognitoidentityprovider.AdminGetUserOutput) *model.User {
    return &model.User{
        ID:            *cognitoUser.Username,
        Name:          getAttributeValue(cognitoUser.UserAttributes, "name"),
        Email:         getAttributeValue(cognitoUser.UserAttributes, "email"),
        EmailVerified: getAttributeValue(cognitoUser.UserAttributes, "email_verified") == "true",
        Image:         getAttributeValue(cognitoUser.UserAttributes, "picture"),
        CreatedAt:     cognitoUser.UserCreateDate,
        UpdatedAt:     cognitoUser.UserLastModifiedDate,
    }
}

func getAttributeValue(attributes []cognitoidentityprovider.AttributeType, name string) string {
    for _, attr := range attributes {
        if *attr.Name == name {
            return *attr.Value
        }
    }
    return ""
}
```

## 📈 パフォーマンス考慮事項

### **1. ホットパーティション対策**

- Cognito User ID をパーティションキーとして使用
- 各ユーザーのデータが分散配置される

### **2. クエリ最適化**

- 必要最小限の GSI（DateIndex のみ）でコスト削減
- 複雑なクエリは複数のシンプルなクエリに分割
- ユーザー別データは PK で直接取得（GSI 不要）

### **3. データサイズ管理**

- 大きなデータは S3 に保存し、DynamoDB には参照のみ
- TTL を使用して古いデータを自動削除

## 🚀 移行戦略

### **フェーズ 1: インフラ準備**

1. Cognito User Pool 作成
2. DynamoDB テーブル作成
3. GSI 設定
4. IAM 権限設定

### **フェーズ 2: コード移行**

1. Cognito クライアント統合
2. リポジトリ層の書き換え
3. データマーシャリング関数の実装
4. エラーハンドリングの調整

### **フェーズ 3: データ移行**

1. 既存ユーザーデータを Cognito に移行
2. アプリケーションデータを DynamoDB に移行
3. データ整合性チェック

### **フェーズ 4: 切り替え**

1. 段階的なトラフィック移行
2. 動作確認とロールバック準備
3. 完全切り替え

## ⚠️ 注意点・制約事項

### **1. DynamoDB の制限**

- アイテムサイズ: 400KB
- パーティションあたりのスループット制限
- GSI の更新遅延

### **2. Cognito の制限**

- ユーザープールあたりの最大ユーザー数
- カスタム属性の制限
- トークンサイズの制限

### **3. 移行時のリスク**

- データ整合性の確保
- ダウンタイムの最小化
- ロールバック手順の準備

### **4. コスト管理**

- 読み書きユニットの適切な設定
- GSI のコスト影響
- Cognito の料金体系

## 📝 次のステップ

1. **Cognito 設定の詳細設計**
2. **Terraform でのテーブル定義作成**
3. **リポジトリ層の実装**
4. **データ移行スクリプトの作成**

---

_このドキュメントは移行計画の第 2 段階「インフラ移行」の一部として作成されました。_
_作成日: 2025 年 8 月 30 日_
_更新日: 2025 年 8 月 30 日_
