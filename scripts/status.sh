#!/bin/bash
set -e

cd infra

echo "📊 Infrastructure Status:"
terraform show -json | jq -r '.values.outputs | to_entries[] | "\(.key): \(.value.value)"'

echo ""
echo "💰 Cost estimation (monthly):"
echo "DynamoDB: \$2-5, Lambda: \$1-3, ECR: \$1-2, CloudWatch: \$2-5, EventBridge: \$0.01"
echo "Total: \$6-15/month"
