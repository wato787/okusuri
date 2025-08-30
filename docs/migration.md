🚀 移行計画

1. 準備フェーズ

現行環境の棚卸し

現在のフロントエンド・バックエンド・通知処理のコード、インフラ設定を確認

DB スキーマのエクスポート（DynamoDB や RDS など）

Terraform ステート管理準備

S3 バケット + DynamoDB で remote state を準備

CI/CD 確認

GitHub Actions / Vercel / AWS デプロイパイプラインを設定

2. インフラ移行

Cognito 導入

ユーザー認証基盤を現行から Cognito に移行

ユーザーデータを既存 DB から Cognito ユーザープールへ移行スクリプト作成

DynamoDB セットアップ

新テーブル作成（現行のスキーマに合わせて）

データ移行バッチ実行

API Gateway + Lambda バックエンド

Go アプリを Lambda 化

Terraform でデプロイ

3. アプリケーション移行

Backend API

既存 API を新 API Gateway 経由に切り替え

ステージング環境で疎通確認

Notification Lambda

EventBridge でスケジュール設定

通知送信処理の動作確認（メール/SNS/Push）

Frontend (Vercel)

Next.js を Vercel にデプロイ

Backend API のエンドポイントを新環境に向ける

Cognito 認証との統合確認

4. データ移行

ユーザー移行

既存ユーザーデータを Cognito にインポート

アプリデータ移行

DynamoDB に旧 DB のデータを投入

整合性確認

5. 切替・リリース

ステージング → 本番環境

ドメイン切替（Route53 / Vercel）

API エンドポイント切替

モニタリング導入

CloudWatch Logs + Metrics 監視

エラー通知（Slack/Webhook）

6. 移行後フォロー

ユーザー影響調査（ログイン・通知）

不具合修正の即応体制

ドキュメント更新（docs/ に残す）

📌 移行時のリスクと対策
リスク 対策
ユーザーのログイン不可 移行前に全ユーザーのテストアカウント作成・動作確認
データ欠損 移行スクリプトでダブルチェック、差分確認
通知が送信されない ステージングで EventBridge + Lambda を実稼働テスト
切替後の不具合 フィーチャーフラグ or 段階的リリースで影響を最小化
