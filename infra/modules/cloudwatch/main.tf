# CloudWatch ロググループ（Lambda用）
resource "aws_cloudwatch_log_group" "lambda_logs" {
  for_each = toset(var.lambda_function_names)

  name              = "/aws/lambda/${each.value}"
  retention_in_days = 14
}

# CloudWatch アラーム（Lambda エラー率）
resource "aws_cloudwatch_metric_alarm" "lambda_errors" {
  for_each = toset(var.lambda_function_names)

  alarm_name          = "${each.value}-errors"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "Errors"
  namespace           = "AWS/Lambda"
  period              = 300
  statistic           = "Sum"
  threshold           = 1
  alarm_description   = "Lambda function error rate"
  alarm_actions       = var.sns_topic_arn != null ? [var.sns_topic_arn] : []

  dimensions = {
    FunctionName = each.value
  }
}

# CloudWatch アラーム（Lambda 実行時間）
resource "aws_cloudwatch_metric_alarm" "lambda_duration" {
  for_each = toset(var.lambda_function_names)

  alarm_name          = "${each.value}-duration"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "Duration"
  namespace           = "AWS/Lambda"
  period              = 300
  statistic           = "Average"
  threshold           = 25000  # 25秒
  alarm_description   = "Lambda function execution duration"
  alarm_actions       = var.sns_topic_arn != null ? [var.sns_topic_arn] : []

  dimensions = {
    FunctionName = each.value
  }
}

# CloudWatch ダッシュボード
resource "aws_cloudwatch_dashboard" "main" {
  dashboard_name = "${var.project}-${var.environment}-dashboard"

  dashboard_body = jsonencode({
    widgets = [
      {
        type   = "metric"
        x      = 0
        y      = 0
        width  = 12
        height = 6

        properties = {
          metrics = [
            for func_name in var.lambda_function_names : [
              "AWS/Lambda",
              "Invocations",
              "FunctionName",
              func_name
            ]
          ]
          period = 300
          stat   = "Sum"
          region = data.aws_region.current.name
          title  = "Lambda Invocations"
        }
      },
      {
        type   = "metric"
        x      = 12
        y      = 0
        width  = 12
        height = 6

        properties = {
          metrics = [
            for func_name in var.lambda_function_names : [
              "AWS/Lambda",
              "Errors",
              "FunctionName",
              func_name
            ]
          ]
          period = 300
          stat   = "Sum"
          region = data.aws_region.current.name
          title  = "Lambda Errors"
        }
      }
    ]
  })
}

# データソース
data "aws_region" "current" {}
