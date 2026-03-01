# ==============================================================================
# Outputs
# ==============================================================================

output "cloud_sql_private_ip" {
  description = "Cloud SQLのプライベートIPアドレス"
  value       = module.cloud_sql.private_ip_address
}

output "cloud_sql_connection_name" {
  description = "Cloud SQLの接続名"
  value       = module.cloud_sql.instance_connection_name
}

output "cloud_run_url" {
  description = "Cloud RunのデフォルトURL"
  value       = module.cloud_run.service_uri
}

output "vpc_network_name" {
  description = "VPCネットワーク名"
  value       = module.vpc.network_name
}

output "vpc_subnet_name" {
  description = "サブネット名"
  value       = module.vpc.subnet_name
}
