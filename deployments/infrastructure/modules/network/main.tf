resource "google_compute_network" "gcp_network" {
  name = var.network_name
  auto_create_subnetworks = true
}

resource "google_compute_subnetwork" "gcp_subnetwork_1" {
  name = var.subnet_1_name
  network = google_compute_network.gcp_network.name
  ip_cidr_range = var.subnet_1_cidr_range
  region = var.resource_region
}

resource "google_compute_subnetwork" "gcp_subnetwork_2" {
  name = var.subnet_2_name
  network = google_compute_network.gcp_network.self_link
  region = var.resource_region
  ip_cidr_range = var.subnet_2_cidr_range
}
