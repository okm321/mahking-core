output "ip_address" {
  description = "Load BalancerのグローバルIPアドレス"
  value       = google_compute_global_address.main.address
}

output "http_url" {
  description = "HTTP URL"
  value       = "http://${google_compute_global_address.main.address}"
}

output "https_url" {
  description = "HTTPS URL（ドメイン設定時）"
  value       = var.domain != null ? "https://${var.domain}" : null
}

output "backend_service_id" {
  description = "Backend ServiceのID（Cloud Armor設定用）"
  value       = google_compute_backend_service.main.id
}
