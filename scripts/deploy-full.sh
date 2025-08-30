#!/bin/bash
set -e

echo "🚀 Starting full deployment process..."

echo "📦 Step 1: Building Lambda functions completed"
echo "🐳 Step 2: Pushing API Lambda to ECR..."

# ECRログイン
echo "🔐 Getting ECR login token..."
cd infra
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

# インフラ作成・更新
echo "🔧 Step 3: Creating infrastructure with Terraform..."
terraform apply -auto-approve

echo "✅ Deployment completed successfully!"
echo "📊 Check CloudWatch dashboard for monitoring"
echo "🔔 Notification Lambda will run every night at 22:00"
