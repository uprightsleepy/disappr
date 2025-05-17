output "disappr_url" {
  description = "Public URL for the disappr Cloud Run service"
  value       = google_cloud_run_service.disappr.status[0].url
}
