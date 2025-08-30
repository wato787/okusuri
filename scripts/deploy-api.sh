#!/bin/bash
set -e

echo "🚀 Deploying API Lambda..."

cd infra

# ECRログイン
echo "🔐 Getting ECR login token..."
ECR_REPO_URL=$(terraform output -raw ecr_repository_url)
ECR_REGISTRY=$(echo $ECR_REPO_URL | cut -d'/' -f1)
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin $ECR_REGISTRY

# Dockerイメージのビルドとプッシュ
echo "🐳 Building and pushing API Lambda image..."
cd ../backend
docker build -t okusuri-api .
cd ../infra
docker tag okusuri-api:latest $ECR_REPO_URL:latest
docker push $ECR_REPO_URL:latest
echo "✅ API Lambda image pushed to ECR successfully"

# Lambda関数更新
echo "🔄 Updating Lambda functions..."
terraform apply -auto-approve -target=module.lambda
echo "✅ Lambda functions updated successfully"
