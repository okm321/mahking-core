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
  region      = "asia-northeast1"

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
    "secretmanager.googleapis.com",
    "servicenetworking.googleapis.com",  # Private Service Connection用
    "run.googleapis.com",                # Cloud Run用
    "certificatemanager.googleapis.com", # Certificate Manager用
    "iam.googleapis.com",                # Workload Identity Federation用
  ]

  /****************************************
    VPC
  ****************************************/
  vpc = {
    network_name = "mahking-${local.env}-vpc"
    subnet_name  = "mahking-${local.env}-subnet"
    subnet_cidr  = "10.0.0.0/24"
  }

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

  /****************************************
    Cloud Run
  ****************************************/
  cloud_run = {
    service_name          = "mahking-${local.env}-api"
    image                 = "us-docker.pkg.dev/cloudrun/container/hello" # 初回はサンプルイメージ
    port                  = 8080
    cpu                   = "1"
    memory                = "512Mi"
    min_instances         = 0
    max_instances         = 10
    ingress               = "INGRESS_TRAFFIC_ALL"
    allow_unauthenticated = true
  }

  /****************************************
    Load Balancer
  ****************************************/
  load_balancer = {
    name       = "mahking-${local.env}-lb"
    domain     = "mahking-api.okmkm.dev" # カスタムドメイン（HTTPS有効）
    enable_cdn = false
  }
}
