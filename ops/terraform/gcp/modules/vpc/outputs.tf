output "network_id" {
  description = "VPCネットワークのID"
  value       = google_compute_network.main.id
}

output "network_name" {
  description = "VPCネットワーク名"
  value       = google_compute_network.main.name
}

output "network_self_link" {
  description = "VPCネットワークのself_link"
  value       = google_compute_network.main.self_link
}

output "subnet_id" {
  description = "サブネットのID"
  value       = google_compute_subnetwork.main.id
}

output "subnet_name" {
  description = "サブネット名"
  value       = google_compute_subnetwork.main.name
}

output "subnet_self_link" {
  description = "サブネットのself_link"
  value       = google_compute_subnetwork.main.self_link
}

output "private_vpc_connection" {
  description = "Private Service Connectionの依存関係用"
  value       = google_service_networking_connection.private_vpc_connection.id
}
