output "log_group_names" {
  description = "CloudWatch log group names"
  value       = [for log_group in aws_cloudwatch_log_group.lambda_logs : log_group.name]
}

output "log_group_arns" {
  description = "CloudWatch log group ARNs"
  value       = [for log_group in aws_cloudwatch_log_group.lambda_logs : log_group.arn]
}

output "error_alarm_names" {
  description = "Lambda error alarm names"
  value       = [for alarm in aws_cloudwatch_metric_alarm.lambda_errors : alarm.alarm_name]
}

output "duration_alarm_names" {
  description = "Lambda duration alarm names"
  value       = [for alarm in aws_cloudwatch_metric_alarm.lambda_duration : alarm.alarm_name]
}

output "dashboard_name" {
  description = "CloudWatch dashboard name"
  value       = aws_cloudwatch_dashboard.main.dashboard_name
}

output "dashboard_arn" {
  description = "CloudWatch dashboard ARN"
  value       = aws_cloudwatch_dashboard.main.dashboard_arn
}
