resource "google_container_cluster" "gke_cluster" {
    name = "car-valuation-ai-cluster"
    enable_autopilot = true
    location = var.target_region
    
    release_channel {
      channel = "rapid"
    }
    network = module.network.network_selflink
    subnetwork = var.assigned_subnets
}

