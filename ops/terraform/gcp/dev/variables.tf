# variable "db_password" {
#   type      = string
#   sensitive = true
# }

locals {
  /****************************************
    Common
  ****************************************/
  env         = "dev"
  gcp_project = "mahking-dev"

  /****************************************
    Artifact Registry
  ****************************************/
  ar_mahking = {
    repository_id = "mahking-${local.env}"
  }

  /****************************************
    Secret Manager
  ****************************************/
  secrets = {
    db_password = {
      secret_id     = "mahking-${local.env}-db-password"
      auto_generate = true
    }
  }

  /****************************************
    Project Services
  ****************************************/
  services = [
    "sqladmin.googleapis.com",
    "compute.googleapis.com",
    "secretmanager.googleapis.com"
  ]

  /****************************************
    Cloud SQL
  ****************************************/
  cloud_sql = {
    instance_name    = "mahking-${local.env}-db"
    database_name    = "mahking_${local.env}"
    db_user          = "app_user"
    tier             = "db-f1-micro"
    disk_size        = 10
    postgres_version = "POSTGRES_18"
    use_backup       = false
  }

  authorized_networks = []
}
