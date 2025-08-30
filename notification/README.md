# 通知用 Lambda

## 📋 概要

EventBridge Scheduler で定期実行される通知送信用 Lambda 関数です。

## 🏗️ アーキテクチャ

- **実行環境**: AWS Lambda
- **デプロイ方法**: ZIP パッケージ
- **データベース**: DynamoDB（単一テーブル設計）
- **認証**: AWS Cognito 連携
- **通知**: WebPush（VAPID）

## 🔧 必要な環境変数

```bash
# Cognito設定
COGNITO_USER_POOL_ID=us-east-1_xxxxxxxxx

# VAPID鍵（WebPush通知用）
VAPID_PUBLIC_KEY=BPxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
VAPID_PRIVATE_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# AWS設定（Lambda実行環境では自動設定）
AWS_REGION=us-east-1
```

## 🚀 ビルドとデプロイ

### ローカルビルド

```bash
# ルートディレクトリから
task build:lambda:notification
```

### デプロイ

```bash
# Terraformでデプロイ
cd infra
terraform apply
```

## 📊 データフロー

1. **EventBridge Scheduler** → Lambda 関数実行
2. **Cognito** → ユーザー一覧取得
3. **DynamoDB** → 通知設定・服用履歴取得
4. **WebPush** → ブラウザ通知送信

## 🗄️ DynamoDB テーブル設計

### テーブル名: `okusuri-table`

#### 通知設定

```
PK: "USER#{cognitoUserId}"
SK: "NOTIFICATION#{platform}"
Data: {
    "platform": "web",
    "isEnabled": true,
    "subscription": "webpush_subscription_json"
}
```

#### 服用履歴

```
PK: "USER#{cognitoUserId}"
SK: "MEDICATION#{date}#{id}"
Data: {
    "hasBleeding": false,
    "createdAt": "2025-08-30T10:00:00Z"
}
```

## 🔍 アクセスパターン

- **ユーザー別データ**: PK（USER#{cognitoUserId}）で直接取得
- **日付検索**: DateIndex GSI を使用
- **通知設定**: SK（NOTIFICATION#{platform}）で取得

## ⚠️ 注意事項

- 通知送信の重複防止（5 分間隔）
- DynamoDB の単一テーブル設計に準拠
- Cognito ユーザー情報との連携
- WebPush 通知の配信保証なし
