# Cognito User Pool
resource "aws_cognito_user_pool" "main" {
  name = var.user_pool_name

  # パスワードポリシー（Google OAuthのみ使用のため最小設定）
  password_policy {
    minimum_length    = 8
    require_lowercase = false
    require_numbers   = false
    require_symbols   = false
    require_uppercase = false
  }

  # ユーザー属性
  username_attributes = ["email"]
  auto_verified_attributes = ["email"]

  # 多要素認証（Google OAuthのみ使用のため無効）
  mfa_configuration = "OFF"

  # アカウント復旧設定
  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_email"
      priority = 1
    }
  }

  # ユーザープール属性
  schema {
    attribute_data_type = "String"
    name                = "email"
    required            = true
    mutable             = true
  }

  # カスタム属性
  schema {
    attribute_data_type = "String"
    name                = "name"
    required            = false
    mutable             = true
  }

  # 検証メッセージ
  verification_message_template {
    default_email_option = "CONFIRM_WITH_CODE"
    email_subject        = "メールアドレスの確認"
    email_message        = "確認コード: {####}"
  }

  # メール設定
  email_configuration {
    email_sending_account = "COGNITO_DEFAULT"
  }

  # タグ
  tags = var.common_tags
}

# Google Identity Provider
resource "aws_cognito_identity_provider" "google" {
  user_pool_id = aws_cognito_user_pool.main.id
  provider_name = "Google"
  provider_type = "Google"

  provider_details = {
    client_id        = var.google_client_id
    client_secret    = var.google_client_secret
    authorize_scopes = "email profile openid"
  }

  attribute_mapping = {
    email    = "email"
    name     = "name"
    username = "sub"
  }
}

# Cognito App Client
resource "aws_cognito_user_pool_client" "main" {
  name = var.client_name
  user_pool_id = aws_cognito_user_pool.main.id

  # 認証フロー（Google OAuthのみ）
  generate_secret = false
  explicit_auth_flows = [
    "ALLOW_REFRESH_TOKEN_AUTH"
  ]

  # コールバックURL
  callback_urls = var.callback_urls
  logout_urls   = var.logout_urls

  # サポートするIDプロバイダー（Google OAuthのみ）
  supported_identity_providers = ["Google"]

  # トークン有効期限
  access_token_validity  = 1    # 1時間
  id_token_validity      = 1    # 1時間
  refresh_token_validity = 30   # 30日

  # セキュリティ設定
  prevent_user_existence_errors = "ENABLED"
}

# Cognito Domain
resource "aws_cognito_user_pool_domain" "main" {
  domain       = var.domain_prefix != null ? var.domain_prefix : "${var.project}-${var.environment}"
  user_pool_id = aws_cognito_user_pool.main.id
}

# User Pool Group（管理者用）
resource "aws_cognito_user_group" "admin" {
  name         = "admin"
  user_pool_id = aws_cognito_user_pool.main.id
  description  = "管理者グループ"
  precedence   = 1
}
