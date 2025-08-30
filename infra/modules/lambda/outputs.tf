output "api_function_name" {
  description = "API Lambda function name"
  value       = aws_lambda_function.api.function_name
}

output "api_function_arn" {
  description = "API Lambda function ARN"
  value       = aws_lambda_function.api.arn
}

output "api_function_url" {
  description = "API Lambda function URL"
  value       = aws_lambda_function_url.api.function_url
}

output "notification_function_name" {
  description = "Notification Lambda function name"
  value       = aws_lambda_function.notification.function_name
}

output "notification_function_arn" {
  description = "Notification Lambda function ARN"
  value       = aws_lambda_function.notification.arn
}

output "notification_function_url" {
  description = "Notification Lambda function URL"
  value       = aws_lambda_function_url.notification.function_url
}

output "api_log_group_name" {
  description = "API Lambda CloudWatch log group name"
  value       = aws_cloudwatch_log_group.api.name
}

output "notification_log_group_name" {
  description = "Notification Lambda CloudWatch log group name"
  value       = aws_cloudwatch_log_group.notification.name
}

output "api_log_group_arn" {
  description = "API Lambda CloudWatch log group ARN"
  value       = aws_cloudwatch_log_group.api.arn
}

output "notification_log_group_arn" {
  description = "Notification Lambda CloudWatch log group ARN"
  value       = aws_cloudwatch_log_group.notification.arn
}
