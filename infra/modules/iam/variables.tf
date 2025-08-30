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

variable "dynamodb_table_arn" {
  description = "DynamoDB table ARN"
  type        = string
  default     = ""
}

variable "cognito_user_pool_arn" {
  description = "Cognito User Pool ARN"
  type        = string
  default     = ""
}

variable "lambda_function_arn" {
  description = "Lambda function ARN for EventBridge"
  type        = string
  default     = ""
}
