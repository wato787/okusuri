#!/bin/bash
set -e

echo "🚀 Deploying notification Lambda..."

cd infra

# Lambda関数を作成・更新
echo "🔧 Creating/updating notification Lambda function..."
terraform apply -auto-approve -target=module.lambda
echo "✅ Notification Lambda function updated successfully"

echo "🎉 Notification Lambda deployment completed!"
echo "📊 Check logs with: task logs:notification"
echo "🔔 Lambda will run every night at 22:00 via EventBridge"
