variable "project_id" {
  description = "GCPプロジェクトID"
  type        = string
}

variable "region" {
  description = "リージョン"
  type        = string
  default     = "asia-northeast1"
}

variable "network_name" {
  description = "VPCネットワーク名"
  type        = string
}

variable "subnet_name" {
  description = "サブネット名"
  type        = string
}

variable "subnet_cidr" {
  description = "サブネットのCIDR範囲"
  type        = string
  default     = "10.0.0.0/24"
}
