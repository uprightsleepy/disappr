variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
  default     = "us-central1"
}

variable "container_image" {
  description = "Docker image deployed to Cloud Run"
  type        = string
}

variable "firebase_project_id" {
  description = "Firebase project ID for JWT validation"
  type        = string
}
