# DynamoDB テーブル（単一テーブル設計）
resource "aws_dynamodb_table" "main" {
  name           = var.table_name
  billing_mode   = var.billing_mode
  hash_key       = "PK"
  range_key      = "SK"

  attribute {
    name = "PK"
    type = "S"
  }

  attribute {
    name = "SK"
    type = "S"
  }

  attribute {
    name = "Date"
    type = "S"
  }

  # GSI1: DateIndex（日付検索用）
  global_secondary_index {
    name            = "DateIndex"
    hash_key        = "Date"
    range_key       = "SK"
    projection_type = "ALL"
  }

  # ポイントインタイムリカバリー
  point_in_time_recovery {
    enabled = true
  }

  # 暗号化設定
  server_side_encryption {
    enabled = true
  }

  tags = var.common_tags
}

# DynamoDB バックアップ設定
resource "aws_backup_vault" "dynamodb" {
  name = "${var.project}-${var.environment}-dynamodb-backup"
  tags = var.common_tags
}

resource "aws_backup_plan" "dynamodb" {
  name = "${var.project}-${var.environment}-dynamodb-backup-plan"

  rule {
    rule_name         = "daily_backup"
    target_vault_name = aws_backup_vault.dynamodb.name
    schedule          = "cron(0 2 * * ? *)"  # 毎日2:00

    lifecycle {
      delete_after = 35  # 35日間保持
    }
  }

  tags = var.common_tags
}

# DynamoDB バックアップ選択
resource "aws_backup_selection" "dynamodb" {
  name         = "${var.project}-${var.environment}-dynamodb-backup-selection"
  plan_id      = aws_backup_plan.dynamodb.id
  iam_role_arn = aws_iam_role.backup_role.arn

  resources = [
    aws_dynamodb_table.main.arn
  ]
}

# バックアップ用IAMロール
resource "aws_iam_role" "backup_role" {
  name = "${var.project}-${var.environment}-backup-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "backup.amazonaws.com"
        }
      }
    ]
  })

  tags = var.common_tags
}

resource "aws_iam_role_policy_attachment" "backup_policy" {
  role       = aws_iam_role.backup_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSBackupServiceRolePolicyForBackup"
}
