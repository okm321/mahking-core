# ==============================================================================
# Cloud Run Service
# ==============================================================================
resource "google_cloud_run_v2_service" "main" {
  project  = var.project_id
  name     = var.service_name
  location = var.region
  ingress  = var.ingress # 外部からのアクセス制御

  template {
    # ==================================================================
    # スケーリング設定
    # ==================================================================
    scaling {
      min_instance_count = var.min_instances
      max_instance_count = var.max_instances
    }

    # ==================================================================
    # Direct VPC Egress設定（VPCコネクタ不要でVPCに接続）
    # ==================================================================
    dynamic "vpc_access" {
      for_each = var.vpc_network != null ? [1] : []
      content {
        network_interfaces {
          network    = var.vpc_network
          subnetwork = var.vpc_subnetwork
        }
        egress = var.vpc_egress # ALL_TRAFFIC or PRIVATE_RANGES_ONLY
      }
    }

    # サービスアカウント
    service_account = var.service_account_email

    # ==================================================================
    # コンテナ設定
    # ==================================================================
    containers {
      image = var.image

      ports {
        container_port = var.port
      }

      resources {
        limits = {
          cpu    = var.cpu
          memory = var.memory
        }
        cpu_idle          = var.cpu_idle          # CPUをアイドル時に割り当てるか
        startup_cpu_boost = var.startup_cpu_boost # 起動時にCPUをブースト
      }

      # ==================================================================
      # 環境変数
      # ==================================================================
      dynamic "env" {
        for_each = var.env_vars
        content {
          name  = env.key
          value = env.value
        }
      }

      # ==================================================================
      # Secret Manager から環境変数を設定
      # ==================================================================
      dynamic "env" {
        for_each = var.secret_env_vars
        content {
          name = env.key
          value_source {
            secret_key_ref {
              secret  = env.value.secret_id
              version = env.value.version
            }
          }
        }
      }
    }
  }

  # トラフィック設定（最新リビジョンに100%）
  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }

  lifecycle {
    ignore_changes = [
      # CI/CDでイメージが更新されるため、Terraformでは無視
      template[0].containers[0].image,
      # GCP APIがサービスレベルのscalingを返すが、Terraformで管理しないため無視
      scaling,
    ]
  }
}

# ==============================================================================
# IAM: 未認証アクセスを許可（パブリックAPI用）
# ==============================================================================
resource "google_cloud_run_v2_service_iam_member" "public_access" {
  count = var.allow_unauthenticated ? 1 : 0

  project  = var.project_id
  location = var.region
  name     = google_cloud_run_v2_service.main.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
