# HTTP API Gateway
resource "aws_apigatewayv2_api" "main" {
  name          = var.api_name
  protocol_type = "HTTP"
  description   = "Okusuri API Gateway"

  cors_configuration {
    allow_credentials = true
    allow_headers     = ["*"]
    allow_methods     = ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    allow_origins     = ["*"]
    expose_headers    = ["*"]
    max_age           = 300
  }

  tags = var.common_tags
}

# API Gateway Stage
resource "aws_apigatewayv2_stage" "main" {
  api_id      = aws_apigatewayv2_api.main.id
  name        = var.stage_name
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gateway.arn
    format = jsonencode({
      requestId      = "$context.requestId"
      ip             = "$context.identity.sourceIp"
      requestTime    = "$context.requestTime"
      httpMethod     = "$context.httpMethod"
      routeKey       = "$context.routeKey"
      status         = "$context.status"
      protocol       = "$context.protocol"
      responseLength = "$context.responseLength"
      integrationLatency = "$context.integrationLatency"
      responseLatency    = "$context.responseLatency"
    })
  }

  tags = var.common_tags
}

# API Gateway の CloudWatch Logs
resource "aws_cloudwatch_log_group" "api_gateway" {
  name              = "/aws/apigateway/${aws_apigatewayv2_api.main.name}/${var.stage_name}"
  retention_in_days = 14
  tags              = var.common_tags
}

# Lambda 統合
resource "aws_apigatewayv2_integration" "lambda" {
  api_id           = aws_apigatewayv2_api.main.id
  integration_type = "AWS_PROXY"

  integration_uri    = var.lambda_function_arn
  integration_method = "POST"
  payload_format_version = "2.0"
}

# 認証なしのルート（ヘルスチェック用）
resource "aws_apigatewayv2_route" "health" {
  api_id    = aws_apigatewayv2_api.main.id
  route_key = "GET /health"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

# 認証ありのルート（メディケーション関連）
resource "aws_apigatewayv2_route" "medications" {
  api_id    = aws_apigatewayv2_api.main.id
  route_key = "GET /medications"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
  authorization_type = "JWT"
  authorizer_id = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "medication_logs" {
  api_id    = aws_apigatewayv2_api.main.id
  route_key = "GET /medication-logs"
  target    = "integrations/${aws_apigatewayv2_api.main.id}"
  authorization_type = "JWT"
  authorizer_id = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "medication_logs_post" {
  api_id    = aws_apigatewayv2_api.main.id
  route_key = "POST /medication-logs"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
  authorization_type = "JWT"
  authorizer_id = aws_apigatewayv2_authorizer.cognito.id
}

# 通知設定関連のルート
resource "aws_apigatewayv2_route" "notifications" {
  api_id    = aws_apigatewayv2_api.main.id
  route_key = "GET /notifications"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
  authorization_type = "JWT"
  authorizer_id = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "notifications_post" {
  api_id    = aws_apigatewayv2_api.main.id
  route_key = "POST /notifications"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
  authorization_type = "JWT"
  authorizer_id = aws_apigatewayv2_authorizer.cognito.id
}

# Cognito JWT Authorizer
resource "aws_apigatewayv2_authorizer" "cognito" {
  api_id           = aws_apigatewayv2_api.main.id
  authorizer_type  = "JWT"
  identity_sources = ["$request.header.Authorization"]
  name             = "cognito-authorizer"

  jwt_configuration {
    audience = [var.cognito_client_id]
    issuer   = "https://cognito-idp.${data.aws_region.current.name}.amazonaws.com/${var.cognito_user_pool_id}"
  }
}

# Lambda 実行権限
resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.main.execution_arn}/*/*"
}

# データソース
data "aws_region" "current" {}
