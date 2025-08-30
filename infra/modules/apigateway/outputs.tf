output "api_id" {
  description = "API Gateway ID"
  value       = aws_apigatewayv2_api.main.id
}

output "api_name" {
  description = "API Gateway name"
  value       = aws_apigatewayv2_api.main.name
}

output "api_url" {
  description = "API Gateway URL"
  value       = "${aws_apigatewayv2_stage.main.invoke_url}"
}

output "stage_name" {
  description = "API Gateway stage name"
  value       = aws_apigatewayv2_stage.main.name
}

output "execution_arn" {
  description = "API Gateway execution ARN"
  value       = aws_apigatewayv2_api.main.execution_arn
}

output "authorizer_id" {
  description = "Cognito JWT Authorizer ID"
  value       = aws_apigatewayv2_authorizer.cognito.id
}

output "log_group_name" {
  description = "API Gateway CloudWatch log group name"
  value       = aws_cloudwatch_log_group.api_gateway.name
}

output "log_group_arn" {
  description = "API Gateway CloudWatch log group ARN"
  value       = aws_cloudwatch_log_group.api_gateway.arn
}
