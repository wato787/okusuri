# EventBridge Scheduler
resource "aws_scheduler_schedule" "notification" {
  name                = "${var.project}-${var.environment}-notification-schedule"
  group_name          = "default"
  flexible_time_window {
    mode = "OFF"
  }

  schedule_expression = var.schedule_expression

  target {
    arn      = var.lambda_function_arn
    role_arn = var.iam_role_arn

    input = jsonencode({
      source = "eventbridge.scheduler"
      time   = "$$.State.EnteredTime"
    })
  }
}

# EventBridge ルール（オプション：より柔軟なスケジュール制御用）
resource "aws_cloudwatch_event_rule" "notification" {
  name                = "${var.project}-${var.environment}-notification-rule"
  description         = "Trigger notification Lambda function"
  schedule_expression = var.schedule_expression
  
  tags = var.common_tags
}

# EventBridge ターゲット
resource "aws_cloudwatch_event_target" "notification" {
  rule      = aws_cloudwatch_event_rule.notification.name
  target_id = "NotificationTarget"
  arn       = var.lambda_function_arn
  role_arn  = var.iam_role_arn

  input = jsonencode({
    source = "eventbridge.rule"
    time   = "$$.time"
  })
}

# Lambda 実行権限（EventBridge用）
resource "aws_lambda_permission" "eventbridge" {
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.notification.arn
}
