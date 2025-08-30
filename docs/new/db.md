# DynamoDB 設計書

## 📋 概要

このドキュメントは、Okusuri アプリケーションのデータベースを PostgreSQL から DynamoDB に移行する際の詳細設計書です。

## 🎯 設計方針

### **基本方針**

- **単一テーブル設計**: パフォーマンスとスケーラビリティの最適化
- **guregu/dynamo 使用**: Go 言語での使いやすさとパフォーマンスの両立
- **アクセスパターン最適化**: 現在のクエリパターンに基づいたキー設計

### **技術スタック**

- **データベース**: Amazon DynamoDB
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

    // インデックス用フィールド
    Type      string    `dynamo:"Type,index:TypeIndex,hash"`  // GSI1: タイプ別検索
    UserID    string    `dynamo:"UserID,index:UserIndex,hash"` // GSI2: ユーザー別検索
    Email     string    `dynamo:"Email,index:EmailIndex,hash"` // GSI3: メール検索
    Token     string    `dynamo:"Token,index:TokenIndex,hash"` // GSI4: トークン検索
    Date      string    `dynamo:"Date,index:DateIndex,hash"`   // GSI5: 日付検索

    // データフィールド
    Data      map[string]interface{} `dynamo:"Data"`           // エンティティ固有のデータ

    // メタデータ
    CreatedAt string    `dynamo:"CreatedAt"`                   // 作成日時 (ISO8601)
    UpdatedAt string    `dynamo:"UpdatedAt"`                   // 更新日時 (ISO8601)
    TTL       *int64    `dynamo:"TTL,omitempty"`               // TTL（必要に応じて）
}
```

#### **キー設計**

##### **1. ユーザー情報**

```
PK: "USER#{userId}"
SK: "PROFILE"
Type: "user"
Data: {
    "name": "ユーザー名",
    "email": "user@example.com",
    "emailVerified": true,
    "image": "https://...",
    "createdAt": "2025-08-30T10:00:00Z",
    "updatedAt": "2025-08-30T10:00:00Z"
}
```

##### **2. セッション情報**

```
PK: "USER#{userId}"
SK: "SESSION#{sessionId}"
Type: "session"
Data: {
    "expiresAt": "2025-09-06T10:00:00Z",
    "token": "jwt_token_here",
    "ipAddress": "192.168.1.1",
    "userAgent": "Mozilla/5.0...",
    "createdAt": "2025-08-30T10:00:00Z",
    "updatedAt": "2025-08-30T10:00:00Z"
}
```

##### **3. OAuth 認証情報**

```
PK: "USER#{userId}"
SK: "ACCOUNT#{providerId}"
Type: "account"
Data: {
    "accountId": "google_account_id",
    "providerId": "google",
    "accessToken": "access_token_here",
    "refreshToken": "refresh_token_here",
    "idToken": "id_token_here",
    "accessTokenExpiresAt": "2025-08-30T11:00:00Z",
    "scope": "openid profile email",
    "createdAt": "2025-08-30T10:00:00Z",
    "updatedAt": "2025-08-30T10:00:00Z"
}
```

##### **4. 服用履歴**

```
PK: "USER#{userId}"
SK: "MEDICATION#{date}#{id}"
Type: "medication_log"
Data: {
    "hasBleeding": false,
    "createdAt": "2025-08-30T10:00:00Z",
    "updatedAt": "2025-08-30T10:00:00Z"
}
```

##### **5. 通知設定**

```
PK: "USER#{userId}"
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

##### **6. メール認証**

```
PK: "VERIFICATION#{identifier}"
SK: "VERIFICATION#{value}"
Type: "verification"
Data: {
    "expiresAt": "2025-08-30T11:00:00Z",
    "createdAt": "2025-08-30T10:00:00Z",
    "updatedAt": "2025-08-30T10:00:00Z"
}
```

## 🔍 GSI（Global Secondary Index）設計

### **GSI1: TypeIndex**

```
PK: Type
SK: PK
用途: タイプ別の一覧取得
例: 全ユーザー取得、全セッション取得
```

### **GSI2: UserIndex**

```
PK: UserID
SK: SK
用途: 特定ユーザーの全データ取得
例: ユーザーの服用履歴一覧、通知設定一覧
```

### **GSI3: EmailIndex**

```
PK: Email
SK: UserID
用途: メールアドレスでのユーザー検索
例: ログイン時のユーザー特定
```

### **GSI4: TokenIndex**

```
PK: Token
SK: UserID
用途: セッショントークンでのユーザー検索
例: 認証ミドルウェアでのユーザー特定
```

### **GSI5: DateIndex**

```
PK: Date
SK: UserID#MEDICATION#{id}
用途: 日付ベースの服用履歴検索
例: 連続服用日数の計算、特定日付の履歴取得
```

## 📊 アクセスパターンとクエリ例

### **1. ユーザー認証フロー**

#### **トークンからユーザー情報取得**

```go
func (r *UserRepository) GetUserByToken(token string) (*model.User, error) {
    var result OkusuriTable
    err := r.table.Get("Token", token).Index("TokenIndex").One(&result)
    if err != nil {
        return nil, err
    }

    // PKからuserIdを抽出: "USER#{userId}" → userId
    userID := strings.TrimPrefix(result.PK, "USER#")

    // ユーザー情報を取得
    var userData OkusuriTable
    err = r.table.Get("PK", "USER#"+userID).Range("SK", "PROFILE").One(&userData)
    if err != nil {
        return nil, err
    }

    return unmarshalUser(userData), nil
}
```

#### **メールアドレスでユーザー検索**

```go
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
    var result OkusuriTable
    err := r.table.Get("Email", email).Index("EmailIndex").One(&result)
    if err != nil {
        return nil, err
    }

    // UserIDからユーザー情報を取得
    var userData OkusuriTable
    err = r.table.Get("PK", "USER#"+result.UserID).Range("SK", "PROFILE").One(&userData)
    if err != nil {
        return nil, err
    }

    return unmarshalUser(userData), nil
}
```

### **2. 服用履歴関連**

#### **ユーザーの服用履歴一覧取得**

```go
func (r *MedicationRepository) GetLogsByUserID(userID string) ([]model.MedicationLog, error) {
    var results []OkusuriTable
    err := r.table.Get("UserID", userID).Index("UserIndex").
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
func (r *MedicationRepository) GetConsecutiveDays(userID string) (int, error) {
    var results []OkusuriTable
    err := r.table.Get("UserID", userID).Index("UserIndex").
        Filter("begins_with(SK, 'MEDICATION#')").
        Order(true). // 降順
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
func (r *NotificationRepository) GetSetting(userID string, platform string) (*model.NotificationSetting, error) {
    var result OkusuriTable
    err := r.table.Get("PK", "USER#"+userID).Range("SK", "NOTIFICATION#"+platform).One(&result)
    if err != nil {
        return nil, err
    }

    return unmarshalNotificationSetting(result), nil
}
```

## 🔧 実装詳細

### **依存関係**

```go
import (
    "github.com/guregu/dynamo/v2"
    "github.com/guregu/dynamo/v2/dynamoattribute"
)
```

### **テーブル初期化**

```go
type Repository struct {
    table dynamo.Table
}

func NewRepository(db *dynamo.DB) *Repository {
    return &Repository{
        table: db.Table("okusuri-table"),
    }
}
```

### **データマーシャリング**

```go
func unmarshalUser(data OkusuriTable) *model.User {
    userData := data.Data
    return &model.User{
        ID:            strings.TrimPrefix(data.PK, "USER#"),
        Name:          userData["name"].(string),
        Email:         userData["email"].(string),
        EmailVerified: userData["emailVerified"].(bool),
        Image:         userData["image"].(*string),
        CreatedAt:     parseTime(userData["createdAt"].(string)),
        UpdatedAt:     parseTime(userData["updatedAt"].(string)),
    }
}
```

## 📈 パフォーマンス考慮事項

### **1. ホットパーティション対策**

- ユーザー ID をパーティションキーとして使用
- 各ユーザーのデータが分散配置される

### **2. クエリ最適化**

- 必要な GSI のみを作成
- 複雑なクエリは複数のシンプルなクエリに分割

### **3. データサイズ管理**

- 大きなデータは S3 に保存し、DynamoDB には参照のみ
- TTL を使用して古いデータを自動削除

## 🚀 移行戦略

### **フェーズ 1: インフラ準備**

1. DynamoDB テーブル作成
2. GSI 設定
3. IAM 権限設定

### **フェーズ 2: コード移行**

1. リポジトリ層の書き換え
2. データマーシャリング関数の実装
3. エラーハンドリングの調整

### **フェーズ 3: データ移行**

1. 既存データのエクスポート
2. DynamoDB 形式への変換
3. データ投入と整合性チェック

### **フェーズ 4: 切り替え**

1. 段階的なトラフィック移行
2. 動作確認とロールバック準備
3. 完全切り替え

## ⚠️ 注意点・制約事項

### **1. DynamoDB の制限**

- アイテムサイズ: 400KB
- パーティションあたりのスループット制限
- GSI の更新遅延

### **2. 移行時のリスク**

- データ整合性の確保
- ダウンタイムの最小化
- ロールバック手順の準備

### **3. コスト管理**

- 読み書きユニットの適切な設定
- GSI のコスト影響
- データ転送料金

## 📝 次のステップ

1. **Terraform でのテーブル定義作成**
2. **リポジトリ層の実装**
3. **データ移行スクリプトの作成**
4. **テスト環境での動作確認**

---

_このドキュメントは移行計画の第 2 段階「インフラ移行」の一部として作成されました。_
_作成日: 2025 年 8 月 30 日_
_更新日: 2025 年 8 月 30 日_
