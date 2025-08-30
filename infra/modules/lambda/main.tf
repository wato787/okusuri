# API用 Lambda 関数（コンテナイメージ）
resource "aws_lambda_function" "api" {
  function_name    = "${var.project}-${var.environment}-api"
  role            = var.iam_role_arn
  package_type    = "Image"
  image_uri       = var.api_image_uri
  timeout         = var.timeout
  memory_size     = var.memory_size

  environment {
    variables = {
      ENVIRONMENT        = var.environment
      DYNAMODB_TABLE    = var.dynamodb_table_name
      COGNITO_USER_POOL = var.cognito_user_pool_id
      LOG_LEVEL         = "INFO"
    }
  }

  tags = var.common_tags
}

# 通知用 Lambda 関数（ZIPパッケージ）
resource "aws_lambda_function" "notification" {
  filename         = var.notification_zip_path
  function_name    = "${var.project}-${var.environment}-notification"
  role            = var.iam_role_arn
  handler         = "notification"
  runtime         = "provided.al2"
  timeout         = 30
  memory_size     = 512
  source_code_hash = fileexists(var.notification_zip_path) ? filebase64sha256(var.notification_zip_path) : null

  environment {
    variables = {
      ENVIRONMENT        = var.environment
      DYNAMODB_TABLE    = var.dynamodb_table_name
      COGNITO_USER_POOL = var.cognito_user_pool_id
      LOG_LEVEL         = "INFO"
    }
  }

  tags = var.common_tags
}

# API Lambda の CloudWatch Logs
resource "aws_cloudwatch_log_group" "api" {
  name              = "/aws/lambda/${aws_lambda_function.api.function_name}"
  retention_in_days = 14
  tags              = var.common_tags
}

# 通知 Lambda の CloudWatch Logs
resource "aws_cloudwatch_log_group" "notification" {
  name              = "/aws/lambda/${aws_lambda_function.notification.function_name}"
  retention_in_days = 14
  tags              = var.common_tags
}

# Lambda 関数 URL（API用）
resource "aws_lambda_function_url" "api" {
  function_name      = aws_lambda_function.api.function_name
  authorization_type = "NONE"  # API Gateway経由でアクセスするため

  cors {
    allow_credentials = true
    allow_origins     = ["*"]
    allow_methods     = ["*"]
    allow_headers     = ["*"]
    expose_headers    = ["keep-alive", "date"]
  }
}

# Lambda 関数 URL（通知用）
resource "aws_lambda_function_url" "notification" {
  function_name      = aws_lambda_function.notification.function_name
  authorization_type = "NONE"  # EventBridge経由でアクセスするため
}

# Lambda 関数の更新設定
resource "aws_lambda_function_event_invoke_config" "api" {
  function_name                = aws_lambda_function.api.function_name
  maximum_event_age_in_seconds = 60  # 最小値は60秒
  maximum_retry_attempts       = 0
}

resource "aws_lambda_function_event_invoke_config" "notification" {
  function_name                = aws_lambda_function.notification.function_name
  maximum_event_age_in_seconds = 60  # 最小値は60秒
  maximum_retry_attempts       = 0
}
