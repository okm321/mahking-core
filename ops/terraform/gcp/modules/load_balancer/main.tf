# ==============================================================================
# Global Static IP Address
# ==============================================================================
resource "google_compute_global_address" "main" {
  project = var.project_id
  name    = "${var.name}-ip"
}

# ==============================================================================
# Serverless NEG (Network Endpoint Group)
# Cloud RunサービスをLoad Balancerのバックエンドとして登録
# ==============================================================================
resource "google_compute_region_network_endpoint_group" "main" {
  project               = var.project_id
  name                  = "${var.name}-neg"
  network_endpoint_type = "SERVERLESS"
  region                = var.region

  cloud_run {
    service = var.cloud_run_service_name
  }
}

# ==============================================================================
# Backend Service
# バックエンドの設定（ヘルスチェック、タイムアウトなど）
# ==============================================================================
resource "google_compute_backend_service" "main" {
  project = var.project_id
  name    = "${var.name}-backend"

  protocol    = "HTTP"
  port_name   = "http"
  timeout_sec = var.timeout_sec

  # Serverless NEGをバックエンドとして追加
  backend {
    group = google_compute_region_network_endpoint_group.main.id
  }

  # Cloud CDNを有効化（オプション）
  enable_cdn = var.enable_cdn

  # ロギング設定
  log_config {
    enable      = var.enable_logging
    sample_rate = var.log_sample_rate
  }
}

# ==============================================================================
# URL Map
# URLパスとバックエンドのマッピング
# ==============================================================================
resource "google_compute_url_map" "main" {
  project = var.project_id
  name    = "${var.name}-urlmap"

  default_service = google_compute_backend_service.main.id
}

# ==============================================================================
# Certificate Manager - 証明書（新しい方式）
# ==============================================================================
resource "google_certificate_manager_certificate" "main" {
  count   = var.domain != null ? 1 : 0
  project = var.project_id
  name    = "${var.name}-cert"

  managed {
    domains = [var.domain]
  }
}

# ==============================================================================
# Certificate Manager - 証明書マップ
# 複数の証明書をまとめて管理するためのマップ
# ==============================================================================
resource "google_certificate_manager_certificate_map" "main" {
  count   = var.domain != null ? 1 : 0
  project = var.project_id
  name    = "${var.name}-cert-map"
}

# ==============================================================================
# Certificate Manager - 証明書マップエントリ
# どのドメインにどの証明書を使うかを定義
# ==============================================================================
resource "google_certificate_manager_certificate_map_entry" "main" {
  count        = var.domain != null ? 1 : 0
  project      = var.project_id
  name         = "${var.name}-cert-map-entry"
  map          = google_certificate_manager_certificate_map.main[0].name
  certificates = [google_certificate_manager_certificate.main[0].id]
  hostname     = var.domain
}

# ==============================================================================
# HTTPS Target Proxy（Certificate Manager使用）
# ==============================================================================
resource "google_compute_target_https_proxy" "main" {
  count   = var.domain != null ? 1 : 0
  project = var.project_id
  name    = "${var.name}-https-proxy"

  url_map         = google_compute_url_map.main.id
  certificate_map = "//certificatemanager.googleapis.com/${google_certificate_manager_certificate_map.main[0].id}"
}

# ==============================================================================
# HTTP Target Proxy (HTTPSへリダイレクト用、またはHTTPのみの場合)
# ==============================================================================
resource "google_compute_target_http_proxy" "main" {
  project = var.project_id
  name    = "${var.name}-http-proxy"

  url_map = var.domain != null ? google_compute_url_map.redirect[0].id : google_compute_url_map.main.id
}

# ==============================================================================
# URL Map for HTTP to HTTPS Redirect
# ==============================================================================
resource "google_compute_url_map" "redirect" {
  count   = var.domain != null ? 1 : 0
  project = var.project_id
  name    = "${var.name}-redirect"

  default_url_redirect {
    https_redirect         = true
    strip_query            = false
    redirect_response_code = "MOVED_PERMANENTLY_DEFAULT"
  }
}

# ==============================================================================
# Forwarding Rules (Frontend)
# ==============================================================================

# HTTPS Forwarding Rule (port 443)
resource "google_compute_global_forwarding_rule" "https" {
  count      = var.domain != null ? 1 : 0
  project    = var.project_id
  name       = "${var.name}-https"
  target     = google_compute_target_https_proxy.main[0].id
  port_range = "443"
  ip_address = google_compute_global_address.main.address
}

# HTTP Forwarding Rule (port 80)
resource "google_compute_global_forwarding_rule" "http" {
  project    = var.project_id
  name       = "${var.name}-http"
  target     = google_compute_target_http_proxy.main.id
  port_range = "80"
  ip_address = google_compute_global_address.main.address
}
