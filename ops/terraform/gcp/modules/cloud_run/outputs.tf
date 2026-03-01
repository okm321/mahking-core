output "service_name" {
  description = "Cloud Runサービス名"
  value       = google_cloud_run_v2_service.main.name
}

output "service_id" {
  description = "Cloud RunサービスID"
  value       = google_cloud_run_v2_service.main.id
}

output "service_uri" {
  description = "Cloud RunサービスのデフォルトURL"
  value       = google_cloud_run_v2_service.main.uri
}

output "service_location" {
  description = "Cloud Runサービスのロケーション"
  value       = google_cloud_run_v2_service.main.location
}

output "domain_mapping_records" {
  description = "カスタムドメインのDNSレコード情報"
  value       = var.domain != null ? google_cloud_run_domain_mapping.main[0].status[0].resource_records : []
}
