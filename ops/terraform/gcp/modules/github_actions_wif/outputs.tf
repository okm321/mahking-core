output "workload_identity_provider" {
  description = "Workload Identity Provider のリソース名（GitHub Actions で使用）"
  value       = google_iam_workload_identity_pool_provider.github.name
}

output "service_account_email" {
  description = "GitHub Actions 用サービスアカウントのメール"
  value       = google_service_account.github_actions.email
}
