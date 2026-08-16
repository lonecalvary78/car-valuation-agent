output "network_name" {
    type = string
    value = google_compute_network.gcp_network.name
}

output "subnets" {
  value = [
    google_compute_subnetwork.gcp_subnetwork_1.name, 
    google_compute_subnetwork.gcp_subnetwork_2.name]
}