provider "google" {
  project = var.project_id
  region  = var.region
}

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 6.35.0"
    }
  }
}

# Enable required APIs
resource "google_project_service" "required" {
  for_each = toset([
    "run.googleapis.com",
    "firestore.googleapis.com",
    "cloudbuild.googleapis.com",
    "logging.googleapis.com",
    "artifactregistry.googleapis.com"
  ])
  service = each.key
}

# Firestore native DB setup
resource "google_firestore_database" "default" {
  project     = var.project_id
  name        = "(default)"
  location_id = var.region
  type        = "FIRESTORE_NATIVE"
  depends_on  = [google_project_service.required]
}

# Cloud Run service
resource "google_cloud_run_service" "disappr" {
  name     = "disappr"
  location = var.region

  template {
    spec {
      containers {
        image = var.container_image
        env {
          name  = "GCP_PROJECT"
          value = var.project_id
        }
        env {
          name  = "FIREBASE_PROJECT_ID"
          value = var.firebase_project_id
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

# Allow public access to Cloud Run
resource "google_cloud_run_service_iam_member" "invoker" {
  service  = google_cloud_run_service.disappr.name
  location = var.region
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# Cloud Build trigger using GitHub (Gen 1)
resource "google_cloudbuild_trigger" "disappr_trigger" {
  name     = "disappr-deploy"
  location = "us-central1"

  github {
    owner = "uprightsleepy"
    name  = "disappr"

    push {
      branch = "^main$"
    }
  }

  filename = "cloudbuild.yaml"
}
