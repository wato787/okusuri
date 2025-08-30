#!/bin/bash
set -e

echo "🚀 Deploying notification Lambda..."

cd infra

# Lambda関数更新
echo "🔄 Updating Lambda functions..."
terraform apply -auto-approve -target=module.lambda
echo "✅ Lambda functions updated successfully"
