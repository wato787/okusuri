variable "project" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "common_tags" {
  description = "Common tags for all resources"
  type        = map(string)
}

variable "api_name" {
  description = "API Gateway name"
  type        = string
  default     = "okusuri-api"
}

variable "stage_name" {
  description = "API Gateway stage name"
  type        = string
  default     = "v1"
}

variable "lambda_function_arn" {
  description = "Lambda function ARN for integration"
  type        = string
}

variable "lambda_function_name" {
  description = "Lambda function name for permission"
  type        = string
}

variable "cognito_user_pool_id" {
  description = "Cognito User Pool ID"
  type        = string
}

variable "cognito_client_id" {
  description = "Cognito App Client ID"
  type        = string
}
