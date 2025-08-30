Terraform リソース一覧（想定）
カテゴリ サービス / リソース 用途 備考
API AWS Lambda Go（Gin/LWA 対応）で実装した API 本体 API Gateway から呼び出し
API Gateway HTTP API フロント（Vercel）からのエントリーポイント Cognito で認証連携
CloudWatch Logs Lambda / API のログ管理 Lambda 実行ログ確認用
認証 Amazon Cognito User Pool Google ログイン連携 自分だけのユーザー管理想定
Cognito Identity Provider Google OAuth 連携設定 Google API Console と紐付け
DB DynamoDB 永続データ保存（通知内容・ユーザーサブスクリプション情報など） on-demand 課金で OK
通知 Amazon SNS（オプション） Web Push に直接は不要。将来モバイル Push 用途なら 今回は Lambda が直接 VAPID で Push 送信
スケジュール EventBridge Scheduler 毎日同じ時間に Lambda を起動 「通知送信用 Lambda」を呼び出す
その他 IAM Role & Policy Lambda 実行権限、DynamoDB 読書き、Cognito 連携など 最小権限で付与
S3（オプション） 静的アセットやバックアップ用 Vercel 利用がメインなら不要
デプロイの流れ（Terraform 管理）

IAM Role → Lambda が実行できるロールを作成

DynamoDB → テーブル作成（ユーザー情報 / 通知内容など）

Cognito → User Pool 作成 + Google 連携設定

Lambda (API) → Go Gin + Lambda Web Adapter でデプロイ

API Gateway → Lambda をバックエンドに設定 + Cognito 認証を適用

Lambda (通知処理) → DB を参照して通知送信処理を実装

EventBridge Scheduler → 通知処理 Lambda を毎日実行

CloudWatch Logs → 各 Lambda のログ確認用

これを Terraform でモジュール分割して管理するときは、

networking/（必要なら VPC）

iam/

dynamodb/

cognito/

lambda/

apigateway/

eventbridge/

のようにディレクトリ分けすると見通しが良くなります 👍
