variable "project_id" {
  description = "Google CloudプロジェクトID"
  type        = string
}

variable "secret_id" {
  description = "Secret ManagerのシークレットID"
  type        = string
}

variable "auto_generate" {
  description = "パスワードを自動生成するかどうか"
  type        = bool
  default     = true
}

variable "secret_data" {
  description = "シークレットに保存するデータ（auto_generateがfalseの場合に使用）"
  type        = string
  default     = ""
  sensitive   = true

  validation {
    condition     = var.auto_generate || var.secret_data != ""
    error_message = "auto_generate=false の場合は secret_data が必須です"
  }
}

variable "password_length" {
  description = "生成するパスワードの長さ"
  type        = number
  default     = 24
}
