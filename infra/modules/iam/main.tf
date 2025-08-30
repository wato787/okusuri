# Lambda 実行ロール
resource "aws_iam_role" "lambda_role" {
  name = "${var.project}-${var.environment}-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = var.common_tags
}

# Lambda 基本実行ポリシー
resource "aws_iam_role_policy_attachment" "lambda_basic" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# DynamoDB アクセスポリシー - wildcardベースの権限
resource "aws_iam_policy" "dynamodb_access" {
  name = "${var.project}-${var.environment}-dynamodb-access"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Query",
          "dynamodb:Scan"
        ]
        Resource = [
          "arn:aws:dynamodb:*:*:table/${var.project}-${var.environment}-*",
          "arn:aws:dynamodb:*:*:table/${var.project}-${var.environment}-*/index/*"
        ]
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_dynamodb" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.dynamodb_access.arn
}

# Cognito アクセスポリシー - wildcardベースの権限
resource "aws_iam_policy" "cognito_access" {
  name = "${var.project}-${var.environment}-cognito-access"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "cognito-idp:AdminGetUser",
          "cognito-idp:ListUsers"
        ]
        Resource = "arn:aws:cognito-idp:*:*:userpool/*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_cognito" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.cognito_access.arn
}

# EventBridge ロール
resource "aws_iam_role" "eventbridge_role" {
  name = "${var.project}-${var.environment}-eventbridge-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "scheduler.amazonaws.com"
        }
      }
    ]
  })

  tags = var.common_tags
}

# EventBridge Lambda 実行ポリシー - wildcardベースの権限
resource "aws_iam_policy" "eventbridge_lambda" {
  name = "${var.project}-${var.environment}-eventbridge-lambda"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "lambda:InvokeFunction"
        ]
        Resource = "arn:aws:lambda:*:*:function:${var.project}-${var.environment}-*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "eventbridge_lambda" {
  role       = aws_iam_role.eventbridge_role.name
  policy_arn = aws_iam_policy.eventbridge_lambda.arn
}
