output "table_name" {
  description = "DynamoDB table name"
  value       = aws_dynamodb_table.main.name
}

output "table_arn" {
  description = "DynamoDB table ARN"
  value       = aws_dynamodb_table.main.arn
}

output "table_id" {
  description = "DynamoDB table ID"
  value       = aws_dynamodb_table.main.id
}

# 個人用のため無効化
# output "backup_vault_arn" {
#   description = "Backup vault ARN"
#   value       = aws_backup_vault.dynamodb.arn
# }

# output "backup_plan_arn" {
#   description = "Backup plan ARN"
#   value       = aws_backup_plan.dynamodb.arn
# }
