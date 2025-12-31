output "secret_id" {
  value = google_secret_manager_secret.main.secret_id
}

output "secret_data" {
  value     = var.auto_generate ? random_password.main[0].result : var.secret_data
  sensitive = true
}
