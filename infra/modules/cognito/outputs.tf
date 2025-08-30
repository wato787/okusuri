output "user_pool_id" {
  description = "Cognito User Pool ID"
  value       = aws_cognito_user_pool.main.id
}

output "user_pool_arn" {
  description = "Cognito User Pool ARN"
  value       = aws_cognito_user_pool.main.arn
}

output "user_pool_name" {
  description = "Cognito User Pool name"
  value       = aws_cognito_user_pool.main.name
}

output "client_id" {
  description = "Cognito App Client ID"
  value       = aws_cognito_user_pool_client.main.id
}

output "client_name" {
  description = "Cognito App Client name"
  value       = aws_cognito_user_pool_client.main.name
}

output "domain" {
  description = "Cognito domain"
  value       = aws_cognito_user_pool_domain.main.domain
}

output "domain_arn" {
  description = "Cognito domain ARN"
  value       = aws_cognito_user_pool_domain.main.id
}

output "google_provider_name" {
  description = "Google Identity Provider name"
  value       = aws_cognito_identity_provider.google.provider_name
}

output "admin_group_name" {
  description = "Admin group name"
  value       = aws_cognito_user_group.admin.name
}
