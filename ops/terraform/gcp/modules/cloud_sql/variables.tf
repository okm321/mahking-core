variable "project_id" {
  description = "Google CloudプロジェクトID"
  type        = string
}

variable "region" {
  description = "Cloud SQLインスタンスのリージョン"
  type        = string
  default     = "asia-northeast1"
}

variable "instance_name" {
  description = "Cloud SQLインスタンスの名前"
  type        = string
}

variable "database_name" {
  description = "Cloud SQLデータベースの名前"
  type        = string
}

variable "db_user" {
  description = "Cloud SQLデータベースのユーザー名"
  type        = string
}

variable "db_password" {
  description = "Secret Managerに保存されているCloud SQLデータベースのパスワードのシークレット名"
  type        = string
}

variable "postgres_version" {
  description = "Cloud SQLインスタンスのPostgreSQLバージョン"
  type        = string
  default     = "POSTGRES_16"
}

variable "tier" {
  description = "Cloud SQLインスタンスのマシンタイプ"
  type        = string
  default     = "db-f1-micro"
}

variable "disk_size" {
  description = "Cloud SQLインスタンスのディスクサイズ（GB単位）"
  type        = number
  default     = 10
}

variable "authorized_networks" {
  description = "Cloud SQLインスタンスにアクセスを許可するネットワークのリスト"
  type = list(object({
    name  = string
    value = string
  }))
  default = []
}

variable "use_backup" {
  description = "Cloud SQLインスタンスのバックアップを有効にするかどうか"
  type        = bool
  default     = true
}
