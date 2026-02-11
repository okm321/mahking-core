output "instance_name" {
  description = "Cloud SQLインスタンス名"
  value       = google_sql_database_instance.main.name
}

output "instance_connection_name" {
  description = "Cloud SQLインスタンスの接続名（project:region:instance）"
  value       = google_sql_database_instance.main.connection_name
}

output "private_ip_address" {
  description = "Cloud SQLインスタンスのプライベートIPアドレス"
  value       = google_sql_database_instance.main.private_ip_address
}

output "public_ip_address" {
  description = "Cloud SQLインスタンスのパブリックIPアドレス"
  value       = google_sql_database_instance.main.public_ip_address
}

output "database_name" {
  description = "データベース名"
  value       = google_sql_database.name.name
}

output "db_user" {
  description = "データベースユーザー名"
  value       = google_sql_user.name.name
}
