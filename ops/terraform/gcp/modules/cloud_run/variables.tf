variable "project_id" {
  description = "GCPプロジェクトID"
  type        = string
}

variable "region" {
  description = "リージョン"
  type        = string
  default     = "asia-northeast1"
}

variable "service_name" {
  description = "Cloud Runサービス名"
  type        = string
}

variable "image" {
  description = "コンテナイメージ"
  type        = string
}

variable "port" {
  description = "コンテナのポート"
  type        = number
  default     = 8080
}

variable "cpu" {
  description = "CPU（例: 1, 2, 4）"
  type        = string
  default     = "1"
}

variable "memory" {
  description = "メモリ（例: 512Mi, 1Gi）"
  type        = string
  default     = "512Mi"
}

variable "min_instances" {
  description = "最小インスタンス数"
  type        = number
  default     = 0
}

variable "max_instances" {
  description = "最大インスタンス数"
  type        = number
  default     = 10
}

variable "cpu_idle" {
  description = "リクエストがないときにCPUを割り当てるか（falseでコールドスタートが発生しやすい）"
  type        = bool
  default     = true
}

variable "startup_cpu_boost" {
  description = "起動時にCPUをブーストするか"
  type        = bool
  default     = true
}

# ==============================================================================
# VPC設定（Direct VPC Egress用）
# ==============================================================================

variable "vpc_network" {
  description = "VPCネットワーク名（Direct VPC Egress用）"
  type        = string
  default     = null
}

variable "vpc_subnetwork" {
  description = "サブネット名（Direct VPC Egress用）"
  type        = string
  default     = null
}

variable "vpc_egress" {
  description = "VPCへのトラフィック制御（ALL_TRAFFIC: 全トラフィック、PRIVATE_RANGES_ONLY: プライベートIPのみ）"
  type        = string
  default     = "PRIVATE_RANGES_ONLY"
}

# ==============================================================================
# アクセス制御
# ==============================================================================

variable "ingress" {
  description = "外部からのアクセス制御（INGRESS_TRAFFIC_ALL, INGRESS_TRAFFIC_INTERNAL_ONLY, INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER）"
  type        = string
  default     = "INGRESS_TRAFFIC_ALL"
}

variable "allow_unauthenticated" {
  description = "未認証アクセスを許可するか"
  type        = bool
  default     = true
}

variable "service_account_email" {
  description = "Cloud Runサービスアカウントのメールアドレス"
  type        = string
  default     = null
}

# ==============================================================================
# 環境変数
# ==============================================================================

variable "env_vars" {
  description = "環境変数のマップ"
  type        = map(string)
  default     = {}
}

variable "secret_env_vars" {
  description = "Secret Managerから取得する環境変数"
  type = map(object({
    secret_id = string
    version   = string
  }))
  default = {}
}

# ==============================================================================
# カスタムドメイン
# ==============================================================================

variable "domain" {
  description = "カスタムドメイン名（nullの場合はドメインマッピングなし）"
  type        = string
  default     = null
}
