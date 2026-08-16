output "gke_cluster_name" {
  type = string
  value = google_container_cluster.gke_cluster.name
}

output "gke_cluster_description" {
  type = string
  value = google_container_cluster.gke_cluster.description
}