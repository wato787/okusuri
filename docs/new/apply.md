# Okusuri インフラデプロイ手順

## 環境設定

### terraform.tfvars
```hcl
google_client_id     = "your-google-client-id"
google_client_secret = "your-google-client-secret"
vapid_public_key     = "your-vapid-public-key"
vapid_private_key    = "your-vapid-private-key"
```

**注意**: cognito_user_pool_idは基盤リソース作成後に動的に取得して設定

## デプロイ順序

```bash
cd infra

# 1. 初期化
terraform init

# 2. 基盤リソース
terraform apply -auto-approve \
  -target=module.iam \
  -target=module.dynamodb \
  -target=module.ecr \
  -target=module.cognito \
  -target=module.cloudwatch

# Cognito User Pool IDを取得してtfvarsに追記
echo "cognito_user_pool_id = \"$(terraform output -raw cognito_user_pool_id)\"" >> terraform.tfvars

# 3. API Lambda準備
cd ..
task build:lambda:api
cd infra
ECR_REPO_URL=$(terraform output -raw ecr_repository_url)
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin ${ECR_REPO_URL%/*}
docker tag okusuri-backend:latest $ECR_REPO_URL:latest
docker push $ECR_REPO_URL:latest

# 4. Notification Lambda準備
cd ..
task build:lambda:notification

# 5. Lambda関数作成
cd infra
terraform apply -auto-approve \
  -target=module.lambda_api \
  -target=module.lambda_notification

# 6. API Gateway・EventBridge
terraform apply -auto-approve \
  -target=module.apigateway \
  -target=module.eventbridge

# 7. 最終確認
terraform apply -auto-approve
```

## フロントエンド設定

```bash
terraform output frontend_config
```

出力をフロントエンドの環境変数に設定