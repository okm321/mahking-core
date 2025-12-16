module "artifact_registry_docker" {
  source        = "../modules/artifact_registry"
  repository_id = local.ar_mahking.repository_id
  env           = local.env
}
