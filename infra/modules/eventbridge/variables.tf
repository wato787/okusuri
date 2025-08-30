variable "project" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "production"
}

variable "common_tags" {
  description = "Common tags for all resources"
  type        = map(string)
}

variable "schedule_expression" {
  description = "EventBridge schedule expression (cron or rate)"
  type        = string
  default     = "cron(0 9 * * ? *)"  # 毎日9:00
}

variable "lambda_function_arn" {
  description = "Lambda function ARN to trigger"
  type        = string
}

variable "lambda_function_name" {
  description = "Lambda function name for permission"
  type        = string
}

variable "iam_role_arn" {
  description = "IAM role ARN for EventBridge execution"
  type        = string
}
