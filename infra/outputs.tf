# Cognito 出力
output "cognito_user_pool_id" {
  description = "Cognito User Pool ID"
  value       = module.cognito.user_pool_id
}

output "cognito_user_pool_arn" {
  description = "Cognito User Pool ARN"
  value       = module.cognito.user_pool_arn
}

output "cognito_client_id" {
  description = "Cognito App Client ID"
  value       = module.cognito.client_id
}

output "cognito_domain" {
  description = "Cognito Domain"
  value       = module.cognito.domain
}

# DynamoDB 出力
output "dynamodb_table_name" {
  description = "DynamoDB table name"
  value       = module.dynamodb.table_name
}

output "dynamodb_table_arn" {
  description = "DynamoDB table ARN"
  value       = module.dynamodb.table_arn
}

# ECR 出力
output "ecr_repository_url" {
  description = "ECR repository URL"
  value       = module.ecr.repository_url
}

output "ecr_repository_name" {
  description = "ECR repository name"
  value       = module.ecr.repository_name
}

# Lambda 出力
output "api_lambda_function_name" {
  description = "API Lambda function name"
  value       = module.lambda.api_function_name
}

output "api_lambda_function_arn" {
  description = "API Lambda function ARN"
  value       = module.lambda.api_function_arn
}

output "notification_lambda_function_name" {
  description = "Notification Lambda function name"
  value       = module.lambda.notification_function_name
}

output "notification_lambda_function_arn" {
  description = "Notification Lambda function ARN"
  value       = module.lambda.notification_function_arn
}

# API Gateway 出力
output "api_gateway_id" {
  description = "API Gateway ID"
  value       = module.apigateway.api_id
}

output "api_gateway_url" {
  description = "API Gateway URL"
  value       = module.apigateway.api_url
}

output "api_gateway_execution_arn" {
  description = "API Gateway execution ARN"
  value       = module.apigateway.execution_arn
}

# EventBridge 出力
output "eventbridge_schedule_arn" {
  description = "EventBridge schedule ARN"
  value       = module.eventbridge.schedule_arn
}

# IAM 出力
output "lambda_role_arn" {
  description = "Lambda execution role ARN"
  value       = module.iam.lambda_role_arn
}

output "eventbridge_role_arn" {
  description = "EventBridge role ARN"
  value       = module.iam.eventbridge_role_arn
}

# CloudWatch 出力
output "cloudwatch_log_groups" {
  description = "CloudWatch log group names"
  value       = module.cloudwatch.log_group_names
}

# フロントエンド設定用
output "frontend_config" {
  description = "Frontend configuration"
  value = {
    api_url        = module.apigateway.api_url
    cognito_domain = module.cognito.domain
    client_id      = module.cognito.client_id
    region         = var.aws_region
  }
}
