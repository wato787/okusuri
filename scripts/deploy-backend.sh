#!/bin/bash
set -e

echo "🚀 Backend Lambda をデプロイしています..."

# Lambda関数名を取得
FUNCTION_NAME=$(cd infra && terraform output -raw lambda_function_name 2>/dev/null || echo "okusuri-backend-api")
echo "Lambda関数名: $FUNCTION_NAME"

# ECRリポジトリURIを取得
ECR_URI=$(cd infra && terraform output -raw ecr_repository_url 2>/dev/null)
if [ -z "$ECR_URI" ]; then
  echo "❌ ECRリポジトリが見つかりません。Terraformを先に実行してください。"
  exit 1
fi

# ECR認証
echo "🔐 ECRに認証中..."
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin $ECR_URI

# イメージタグ付けとプッシュ
echo "📦 イメージをプッシュ中..."
docker tag okusuri-backend:latest $ECR_URI:latest
docker push $ECR_URI:latest

# Lambda関数を更新
echo "🔄 Lambda関数を更新中..."
aws lambda update-function-code \
  --function-name $FUNCTION_NAME \
  --image-uri $ECR_URI:latest

echo "✅ Backendデプロイ完了"