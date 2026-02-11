variable "project_id" {
  description = "GCPプロジェクトID"
  type        = string
}

variable "env" {
  description = "環境名（dev / prd）"
  type        = string
}

variable "github_repository" {
  description = "GitHubリポジトリ（owner/repo）"
  type        = string
  default     = "okm321/mahking-core"
}
