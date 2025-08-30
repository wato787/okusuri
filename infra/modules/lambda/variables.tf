variable "project" {
  description = "Project name"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "ap-northeast-1"
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

variable "iam_role_arn" {
  description = "IAM role ARN for Lambda execution"
  type        = string
}

variable "cognito_user_pool_id" {
  description = "Cognito User Pool ID"
  type        = string
  default     = ""
}

variable "dynamodb_table_name" {
  description = "DynamoDB table name"
  type        = string
}

variable "api_image_uri" {
  description = "ECR image URI for API Lambda function"
  type        = string
  default     = ""
}

variable "notification_zip_path" {
  description = "Path to notification Lambda deployment package"
  type        = string
  default     = ""
}

variable "timeout" {
  description = "Lambda timeout in seconds"
  type        = number
  default     = 30
}

variable "memory_size" {
  description = "Lambda memory size in MB"
  type        = number
  default     = 512
}

variable "vapid_public_key" {
  description = "VAPID public key for push notifications"
  type        = string
  sensitive   = true
  default     = ""
}

variable "vapid_private_key" {
  description = "VAPID private key for push notifications"
  type        = string
  sensitive   = true
  default     = ""
}
