variable "project" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "common_tags" {
  description = "Common tags for all resources"
  type        = map(string)
}

variable "user_pool_name" {
  description = "Cognito User Pool name"
  type        = string
  default     = "okusuri-user-pool"
}

variable "client_name" {
  description = "Cognito App Client name"
  type        = string
  default     = "okusuri-app"
}

variable "domain_prefix" {
  description = "Cognito domain prefix"
  type        = string
  default     = null
}

variable "google_client_id" {
  description = "Google OAuth Client ID"
  type        = string
  sensitive   = true
}

variable "google_client_secret" {
  description = "Google OAuth Client Secret"
  type        = string
  sensitive   = true
}

variable "callback_urls" {
  description = "Callback URLs for OAuth"
  type        = list(string)
  default = [
    "http://localhost:3000/auth/callback",
    "https://okusuri.vercel.app/auth/callback"
  ]
}

variable "logout_urls" {
  description = "Logout URLs"
  type        = list(string)
  default = [
    "http://localhost:3000/",
    "https://okusuri.vercel.app/"
  ]
}
