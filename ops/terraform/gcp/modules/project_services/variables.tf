variable "project_id" {
  description = "Google CloudプロジェクトID"
  type        = string
}

variable "services" {
  description = "有効化するGoogle Cloudサービスのリスト"
  type        = list(string)
  default     = []
}
