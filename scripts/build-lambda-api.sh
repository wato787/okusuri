#!/bin/bash
set -e

echo "🐳 API Lambda用Dockerイメージをビルド中..."

# プロジェクトルートに移動
cd "$(dirname "$0")/.."

# ECRリポジトリURLを取得
ECR_URL=$(cd infra && terraform output -raw ecr_repository_url 2>/dev/null || echo "")

if [ -z "$ECR_URL" ]; then
  echo "❌ ECRリポジトリが見つかりません。Terraformを先に実行してください。"
  exit 1
fi

echo "📦 ECRリポジトリ: $ECR_URL"

# Dockerイメージをビルド
cd backend
docker build -t okusuri-backend:latest .

# ECRにタグ付け
docker tag okusuri-backend:latest $ECR_URL:latest

# ECRにプッシュ
echo "🚀 ECRにプッシュ中..."
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin $ECR_URL
docker push $ECR_URL:latest

echo "✅ Dockerイメージビルド・プッシュ完了: $ECR_URL:latest"