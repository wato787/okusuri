# 基本設定
variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "ap-northeast-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "dev"
}

variable "project" {
  description = "Project name"
  type        = string
  default     = "okusuri"
}

# Google OAuth 設定
variable "google_client_id" {
  description = "Google OAuth Client ID"
  type        = string
  sensitive   = true
}

variable "google_client_secret" {
  description = "Google OAuth Client Secret"
  type        = string
  sensitive   = true
}

# Cognito 設定
variable "cognito_user_pool_name" {
  description = "Cognito User Pool name"
  type        = string
  default     = "okusuri-user-pool"
}

variable "cognito_client_name" {
  description = "Cognito App Client name"
  type        = string
  default     = "okusuri-app"
}

# DynamoDB 設定
variable "dynamodb_table_name" {
  description = "DynamoDB table name"
  type        = string
  default     = "okusuri-table"
}

variable "dynamodb_billing_mode" {
  description = "DynamoDB billing mode"
  type        = string
  default     = "PAY_PER_REQUEST"
}

# Lambda 設定
variable "lambda_runtime" {
  description = "Lambda runtime"
  type        = string
  default     = "provided.al2"
}

variable "lambda_timeout" {
  description = "Lambda timeout in seconds"
  type        = number
  default     = 30
}

variable "lambda_memory_size" {
  description = "Lambda memory size in MB"
  type        = number
  default     = 512
}

# API Gateway 設定
variable "api_gateway_name" {
  description = "API Gateway name"
  type        = string
  default     = "okusuri-api"
}

variable "api_gateway_stage_name" {
  description = "API Gateway stage name"
  type        = string
  default     = "v1"
}

# EventBridge 設定
variable "notification_schedule" {
  description = "Notification schedule (cron expression)"
  type        = string
  default     = "cron(0 9 * * ? *)"  # 毎日9:00
}

# タグ設定
variable "common_tags" {
  description = "Common tags for all resources"
  type        = map(string)
  default = {
    Project     = "okusuri"
    Environment = "dev"
    ManagedBy   = "terraform"
  }
}

# Lambda イメージ・パッケージ設定
variable "api_image_uri" {
  description = "ECR image URI for API Lambda"
  type        = string
  default     = ""  # デプロイ時に上書きされる
}

variable "notification_zip_path" {
  description = "Path to notification Lambda deployment package"
  type        = string
  default     = "dist/lambda/notification.zip"
}
