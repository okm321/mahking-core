resource "google_sql_database_instance" "main" {
  project          = var.project_id
  name             = var.instance_name
  database_version = var.postgres_version
  region           = var.region

  settings {
    tier            = var.tier
    edition         = "ENTERPRISE" # ← これを追加
    disk_size       = var.disk_size
    disk_type       = "PD_SSD"
    disk_autoresize = true

    ip_configuration {
      ipv4_enabled = true

      dynamic "authorized_networks" {
        for_each = var.authorized_networks
        content {
          name  = aauthorized_networks.value.name
          value = aauthorized_networks.value.value
        }
      }
    }

    backup_configuration {
      enabled    = var.use_backup
      start_time = "03:00"
    }

    insights_config {
      query_insights_enabled  = true
      query_string_length     = 1024 # クエリ文字列の最大長
      record_application_tags = true # アプリタグ記録
      record_client_address   = true # クライアントIP記録
      query_plans_per_minute  = 20   # 1分あたりのクエリプラン数
    }
  }

  deletion_protection = false
}

resource "google_sql_database" "name" {
  project  = var.project_id
  name     = var.database_name
  instance = google_sql_database_instance.main.name
}

resource "google_sql_user" "name" {
  project  = var.project_id
  name     = var.db_user
  instance = google_sql_database_instance.main.name
  password = var.db_password
}

# resource "google_sql_user" "additional" {
#   for_each = { for u in var.additional_users : u.name => u }
#
#   project  = var.project_id
#   name     = each.value.name
#   instance = google_sql_database_instance.main.name
#   password = each.value.password
# }
