output "schedule_arn" {
  description = "EventBridge schedule ARN"
  value       = aws_scheduler_schedule.notification.arn
}

output "schedule_name" {
  description = "EventBridge schedule name"
  value       = aws_scheduler_schedule.notification.name
}

output "rule_arn" {
  description = "EventBridge rule ARN"
  value       = aws_cloudwatch_event_rule.notification.arn
}

output "rule_name" {
  description = "EventBridge rule name"
  value       = aws_cloudwatch_event_rule.notification.name
}

output "target_id" {
  description = "EventBridge target ID"
  value       = aws_cloudwatch_event_target.notification.target_id
}
