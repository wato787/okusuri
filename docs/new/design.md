[Vercel Frontend]
↓ (ユーザー認証リクエスト)
[Amazon Cognito (Google 認証)]
↓ (ID トークン)
[API Gateway] → [API Lambda (Gin+LWA)] → DynamoDB

[EventBridge Scheduler (毎日 09:00)]
↓
[Notify Lambda] → DynamoDB → Web Push (via web-push library)

📘 アプリ全体設計まとめ

1. フロントエンド (Vercel)

デプロイ先: Vercel

構成要素:

Next.js が v1。React vite が v2。React vite に移行予定

Google ログインボタン → Cognito Hosted UI にリダイレクト

ログイン後、Cognito から受け取った ID トークン を保持

API Gateway 経由でバックエンドにアクセス

Web Push API で通知を受信

2. 認証基盤 (Cognito + Google)

Amazon Cognito が Google OAuth2 と連携

フロー:

ユーザーが Google でログイン

Cognito が認証 → JWT (ID トークン) を返す

フロントエンドは API にリクエストするときにトークンを付与

メリット:

API Gateway がトークン検証を代行

バックエンド実装をシンプル化

3. バックエンド API

構成: API Gateway + Lambda (Go Gin) + DynamoDB

詳細:

API Gateway (HTTP API)

Cognito トークンを検証

Lambda にリクエストをルーティング

Lambda (Gin + Lambda Web Adapter)

Gin で実装した Go API を Lambda で動作させる

DynamoDB へデータの CRUD を実行

DynamoDB

ユーザー情報

通知サブスクリプション情報

通知に必要なデータ

4. 通知基盤

EventBridge Scheduler

毎日 09:00 に通知処理 Lambda を起動

通知 Lambda

DynamoDB から通知対象ユーザーと内容を取得

Web Push を実行 (Service Worker 経由でフロントへ配信)

5. フローまとめ

ログイン

Vercel (フロント) → Cognito (Google 認証) → ID トークン取得

通常操作

フロントから API Gateway 経由で API Lambda にアクセス

DynamoDB に保存・取得

通知

EventBridge Scheduler → 通知 Lambda 起動

DynamoDB を参照し通知内容を決定

Web Push でフロントへ通知

6. 使用サービス一覧

Vercel: フロントエンドホスティング

Amazon Cognito: Google 認証・ユーザー管理

API Gateway: API リクエストのエントリーポイント

Lambda (Gin + LWA): バックエンド API

DynamoDB: データストア

EventBridge Scheduler: 定時ジョブ実行

通知 Lambda: DynamoDB を見て通知を配信

Web Push API: ユーザー通知
