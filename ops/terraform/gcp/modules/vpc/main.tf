# ==============================================================================
# VPC Network
# ==============================================================================
resource "google_compute_network" "main" {
  project                 = var.project_id
  name                    = var.network_name
  auto_create_subnetworks = false # カスタムサブネットを使用
  routing_mode            = "REGIONAL"
}

# ==============================================================================
# Subnet
# ==============================================================================
resource "google_compute_subnetwork" "main" {
  project       = var.project_id
  name          = var.subnet_name
  ip_cidr_range = var.subnet_cidr
  region        = var.region
  network       = google_compute_network.main.id

  # Cloud Run Direct VPC Egress用に必要
  private_ip_google_access = true
}

# ==============================================================================
# Private Service Connection (Cloud SQL Private IP用)
# ==============================================================================

# Google サービス用のプライベートIPレンジを予約
resource "google_compute_global_address" "private_ip_range" {
  project       = var.project_id
  name          = "${var.network_name}-private-ip-range"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.main.id
}

# VPC PeeringでGoogleサービスに接続
resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = google_compute_network.main.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_range.name]
}
