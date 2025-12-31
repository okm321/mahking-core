module "artifact_registry_docker" {
  source        = "../modules/artifact_registry"
  repository_id = local.ar_mahking.repository_id
  env           = local.env
}

module "project_services" {
  source     = "../modules/project_services"
  project_id = local.gcp_project
  services   = local.services
}

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

module "cloud_sql" {
  source = "../modules/cloud_sql"

  project_id          = local.gcp_project
  instance_name       = local.cloud_sql.instance_name
  database_name       = local.cloud_sql.database_name
  db_user             = local.cloud_sql.db_user
  db_password         = module.secrets["db_password"].secret_data
  tier                = local.cloud_sql.tier
  disk_size           = local.cloud_sql.disk_size
  postgres_version    = local.cloud_sql.postgres_version
  authorized_networks = local.authorized_networks
  use_backup          = local.cloud_sql.use_backup

  depends_on = [module.project_services]
}
