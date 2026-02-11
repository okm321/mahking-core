resource "random_password" "main" {
  count   = var.auto_generate ? 1 : 0 // 自動生成の時だけ作る
  length  = var.password_length
  special = false
}

resource "google_secret_manager_secret" "main" {
  project   = var.project_id
  secret_id = var.secret_id

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "main" {
  secret      = google_secret_manager_secret.main.id
  secret_data = var.auto_generate ? random_password.main[0].result : var.secret_data

  depends_on = [random_password.main] # 順序を明示
}

# ==============================================================================
# IAM: シークレットへのアクセス権限を付与
# ==============================================================================
resource "google_secret_manager_secret_iam_member" "accessor" {
  for_each = toset(var.accessor_service_accounts)

  project   = var.project_id
  secret_id = google_secret_manager_secret.main.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${each.value}"
}
