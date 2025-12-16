terraform {
  required_version = ">= 1.10.5"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 7.1.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 7.1.0"
    }
  }
}

provider "google" {
  project = local.gcp_project
}

provider "google-beta" {
  project = local.gcp_project
}
