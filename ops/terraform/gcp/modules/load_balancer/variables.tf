variable "project_id" {
  description = "GCPプロジェクトID"
  type        = string
}

variable "region" {
  description = "Cloud Runサービスのリージョン"
  type        = string
  default     = "asia-northeast1"
}

variable "name" {
  description = "Load Balancer関連リソースの名前プレフィックス"
  type        = string
}

variable "cloud_run_service_name" {
  description = "バックエンドとなるCloud Runサービス名"
  type        = string
}

variable "domain" {
  description = "SSL証明書用のドメイン名（nullの場合はHTTPのみ）"
  type        = string
  default     = null
}

variable "timeout_sec" {
  description = "バックエンドへのリクエストタイムアウト（秒）"
  type        = number
  default     = 30
}

variable "enable_cdn" {
  description = "Cloud CDNを有効化するか"
  type        = bool
  default     = false
}

variable "enable_logging" {
  description = "アクセスログを有効化するか"
  type        = bool
  default     = true
}

variable "log_sample_rate" {
  description = "ログのサンプリングレート（0.0〜1.0）"
  type        = number
  default     = 1.0
}
