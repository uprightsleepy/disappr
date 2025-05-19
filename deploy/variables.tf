variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "Region for resources"
  type        = string
  default     = "us-central1"
}

variable "container_image" {
  description = "Container image URL for Cloud Run"
  type        = string
}

variable "firebase_project_id" {
  description = "Firebase project ID for JWT verification"
  type        = string
}

variable "project_number" {
  description = "GCP project number"
  type        = string
}
