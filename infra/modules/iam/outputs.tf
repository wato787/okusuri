output "lambda_role_arn" {
  description = "Lambda execution role ARN"
  value       = aws_iam_role.lambda_role.arn
}

output "lambda_role_name" {
  description = "Lambda execution role name"
  value       = aws_iam_role.lambda_role.name
}

output "eventbridge_role_arn" {
  description = "EventBridge role ARN"
  value       = aws_iam_role.eventbridge_role.arn
}

output "eventbridge_role_name" {
  description = "EventBridge role name"
  value       = aws_iam_role.eventbridge_role.name
}

output "dynamodb_policy_arn" {
  description = "DynamoDB access policy ARN"
  value       = aws_iam_policy.dynamodb_access.arn
}

output "cognito_policy_arn" {
  description = "Cognito access policy ARN"
  value       = aws_iam_policy.cognito_access.arn
}
