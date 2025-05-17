provider "google" {
  project = var.project_id
  region  = var.region
}

resource "google_project_service" "required" {
  for_each = toset([
    "run.googleapis.com",
    "firestore.googleapis.com",
    "cloudbuild.googleapis.com",
    "logging.googleapis.com"
  ])
  service = each.key
}

resource "google_firestore_database" "default" {
  project     = var.project_id
  name        = "(default)"
  location_id = var.region
  type        = "FIRESTORE_NATIVE"
  depends_on  = [google_project_service.required]
}

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

resource "google_cloud_run_service_iam_member" "invoker" {
  service  = google_cloud_run_service.disappr.name
  location = var.region
  role     = "roles/run.invoker"
  member   = "allUsers"
}
