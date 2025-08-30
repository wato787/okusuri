terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket         = "okusuri-terraform-state"
    key            = "infra/terraform.tfstate"
    region         = "ap-northeast-1"
    encrypt        = true
  }
}

# AWS プロバイダー設定
provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "okusuri"
      ManagedBy   = "terraform"
    }
  }
}

# データソース
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

# モジュール呼び出し（依存関係順）
module "cognito" {
  source = "./modules/cognito"
  
  environment = var.environment
  project     = var.project
  common_tags = var.common_tags
  
  google_client_id     = var.google_client_id
  google_client_secret = var.google_client_secret
}

module "dynamodb" {
  source = "./modules/dynamodb"
  
  environment = var.environment
  project     = var.project
  common_tags = var.common_tags
}

module "iam" {
  source = "./modules/iam"
  
  environment = var.environment
  project     = var.project
  common_tags = var.common_tags
}

module "ecr" {
  source = "./modules/ecr"
  
  environment = var.environment
  project     = var.project
  common_tags = var.common_tags
}

# Backend Lambda用（Cognito依存なし）
module "lambda_api" {
  source = "./modules/lambda"
  
  environment = var.environment
  project     = var.project
  common_tags = var.common_tags
  aws_region  = var.aws_region
  
  api_image_uri        = "${module.ecr.repository_url}:latest"
  dynamodb_table_name  = module.dynamodb.table_name
  iam_role_arn         = module.iam.lambda_role_arn
}

# Notification Lambda用（Cognito依存あり）
module "lambda_notification" {
  source = "./modules/lambda"
  
  environment = var.environment
  project     = var.project
  common_tags = var.common_tags
  aws_region  = var.aws_region
  
  notification_zip_path = var.notification_zip_path
  cognito_user_pool_id = module.cognito.user_pool_id
  dynamodb_table_name  = module.dynamodb.table_name
  iam_role_arn         = module.iam.lambda_role_arn
  
  vapid_public_key     = var.vapid_public_key
  vapid_private_key    = var.vapid_private_key
}

module "apigateway" {
  source = "./modules/apigateway"
  
  environment = var.environment
  project     = var.project
  common_tags = var.common_tags
  
  cognito_user_pool_id = module.cognito.user_pool_id
  cognito_client_id    = module.cognito.client_id
  lambda_function_arn  = module.lambda_api.api_function_arn
  lambda_function_name = module.lambda_api.api_function_name
}

module "eventbridge" {
  source = "./modules/eventbridge"
  
  environment = var.environment
  project     = var.project
  common_tags = var.common_tags
  
  lambda_function_arn  = module.lambda_notification.notification_function_arn
  lambda_function_name = module.lambda_notification.notification_function_name
  iam_role_arn         = module.iam.eventbridge_role_arn
}

module "cloudwatch" {
  source = "./modules/cloudwatch"
  
  environment = var.environment
  project     = var.project
  common_tags = var.common_tags
  
  lambda_function_names = [
    module.lambda_api.api_function_name,
    module.lambda_notification.notification_function_name
  ]
}
