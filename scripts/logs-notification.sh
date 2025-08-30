#!/bin/bash
set -e

cd infra

# 通知Lambdaログを表示
LAMBDA_NAME=$(terraform output -raw notification_lambda_function_name)
echo "📊 Viewing logs for: $LAMBDA_NAME"
aws logs tail /aws/lambda/$LAMBDA_NAME --follow
