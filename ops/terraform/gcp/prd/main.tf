# ==============================================================================
# Artifact Registry
# ==============================================================================
module "artifact_registry_docker" {
  source        = "../modules/artifact_registry"
  repository_id = local.ar_mahking.repository_id
  env           = local.env
}

# ==============================================================================
# Project Services (API有効化)
# ==============================================================================
module "project_services" {
  source     = "../modules/project_services"
  project_id = local.gcp_project
  services   = local.services
}

# ==============================================================================
# Secret Manager
# ==============================================================================
module "secrets" {
  source   = "../modules/secret_manager"
  for_each = local.secrets

  project_id      = local.gcp_project
  secret_id       = each.value.secret_id
  auto_generate   = lookup(each.value, "auto_generate", true)
  secret_data     = lookup(each.value, "secret_data", "")
  password_length = lookup(each.value, "password_length", 24)

  depends_on = [module.project_services]
}

# ==============================================================================
# VPC Network
# ==============================================================================
module "vpc" {
  source = "../modules/vpc"

  project_id   = local.gcp_project
  region       = local.region
  network_name = local.vpc.network_name
  subnet_name  = local.vpc.subnet_name
  subnet_cidr  = local.vpc.subnet_cidr

  depends_on = [module.project_services]
}

# ==============================================================================
# Cloud SQL (Private IP)
# ==============================================================================
module "cloud_sql" {
  source = "../modules/cloud_sql"

  project_id             = local.gcp_project
  region                 = local.region
  instance_name          = local.cloud_sql.instance_name
  database_name          = local.cloud_sql.database_name
  db_user                = local.cloud_sql.db_user
  db_password            = module.secrets["db_password"].secret_data
  tier                   = local.cloud_sql.tier
  disk_size              = local.cloud_sql.disk_size
  postgres_version       = local.cloud_sql.postgres_version
  use_backup             = local.cloud_sql.use_backup
  deletion_protection    = local.cloud_sql.deletion_protection
  network_id             = module.vpc.network_id
  private_vpc_connection = module.vpc.private_vpc_connection
  enable_public_ip       = false

  depends_on = [module.project_services, module.vpc]
}

# ==============================================================================
# Cloud Run (Direct VPC Egress)
# ==============================================================================
module "cloud_run" {
  source = "../modules/cloud_run"

  project_id    = local.gcp_project
  region        = local.region
  service_name  = local.cloud_run.service_name
  image         = local.cloud_run.image
  port          = local.cloud_run.port
  cpu           = local.cloud_run.cpu
  memory        = local.cloud_run.memory
  min_instances = local.cloud_run.min_instances
  max_instances = local.cloud_run.max_instances

  # VPC設定（Direct VPC Egress）
  vpc_network    = module.vpc.network_name
  vpc_subnetwork = module.vpc.subnet_name
  vpc_egress     = "PRIVATE_RANGES_ONLY"

  # アクセス制御
  ingress               = local.cloud_run.ingress
  allow_unauthenticated = local.cloud_run.allow_unauthenticated

  # 環境変数（DBへの接続情報）
  env_vars = {
    PG_HOST   = module.cloud_sql.private_ip_address
    PG_PORT   = "5432"
    PG_DBNAME = module.cloud_sql.database_name
    PG_USER   = module.cloud_sql.db_user
    PG_SCHEMA = "public"
  }

  # Secret Managerからパスワードを取得
  secret_env_vars = {
    PG_PASS = {
      secret_id = module.secrets["db_password"].secret_id
      version   = "latest"
    }
  }

  depends_on = [module.project_services, module.vpc, module.cloud_sql]
}

# ==============================================================================
# Load Balancer
# ==============================================================================
module "load_balancer" {
  source = "../modules/load_balancer"

  project_id             = local.gcp_project
  region                 = local.region
  name                   = local.load_balancer.name
  cloud_run_service_name = module.cloud_run.service_name
  domain                 = local.load_balancer.domain
  enable_cdn             = local.load_balancer.enable_cdn
  enable_logging         = local.load_balancer.enable_logging

  depends_on = [module.cloud_run]
}

# ==============================================================================
# Workload Identity Federation (GitHub Actions)
# ==============================================================================
module "github_actions_wif" {
  source = "../modules/github_actions_wif"

  project_id = local.gcp_project
  env        = local.env

  depends_on = [module.project_services]
}
