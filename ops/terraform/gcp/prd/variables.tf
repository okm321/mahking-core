locals {
  # ============================================================================
  # Common
  # ============================================================================
  env         = "prd"
  gcp_project = "mahking-prd" # 本番用のプロジェクトID
  region      = "asia-northeast1"

  # ============================================================================
  # Artifact Registry
  # ============================================================================
  ar_mahking = {
    repository_id = "mahking-${local.env}"
  }

  # ============================================================================
  # Secret Manager
  # ============================================================================
  secrets = {
    db_password = {
      secret_id       = "mahking-${local.env}-db-password"
      auto_generate   = true
      password_length = 32 # 本番はより長いパスワード
    }
  }

  # ============================================================================
  # Project Services (有効化するAPI)
  # ============================================================================
  services = [
    "sqladmin.googleapis.com",          # Cloud SQL
    "compute.googleapis.com",           # Compute Engine / VPC
    "secretmanager.googleapis.com",     # Secret Manager
    "run.googleapis.com",               # Cloud Run
    "servicenetworking.googleapis.com", # Private Service Connection
    "iam.googleapis.com",               # Workload Identity Federation用
  ]

  # ============================================================================
  # VPC
  # ============================================================================
  vpc = {
    network_name = "mahking-${local.env}-vpc"
    subnet_name  = "mahking-${local.env}-subnet"
    subnet_cidr  = "10.1.0.0/24" # devとは異なるCIDR
  }

  # ============================================================================
  # Cloud SQL
  # ============================================================================
  cloud_sql = {
    instance_name       = "mahking-${local.env}-db"
    database_name       = "mahking_${local.env}"
    db_user             = "app_user"
    tier                = "db-custom-1-3840" # 本番は1vCPU, 3.75GB RAM
    disk_size           = 20
    postgres_version    = "POSTGRES_18"
    use_backup          = true # 本番はバックアップ有効
    deletion_protection = true # 本番は削除保護有効
  }

  # ============================================================================
  # Cloud Run
  # ============================================================================
  cloud_run = {
    service_name = "mahking-${local.env}-api"
    # 初回はプレースホルダーイメージを使用（後でCI/CDで上書き）
    image                 = "us-docker.pkg.dev/cloudrun/container/hello"
    port                  = 8080
    cpu                   = "2"   # 本番はより多くのCPU
    memory                = "1Gi" # 本番はより多くのメモリ
    min_instances         = 1     # 本番は常時1台稼働（コールドスタート回避）
    max_instances         = 100   # 本番はより多くスケール可能
    ingress               = "INGRESS_TRAFFIC_ALL"
    allow_unauthenticated = true
    domain                = null # 独自ドメインを設定する場合はここに指定（例: "api.example.com"）
  }
}
