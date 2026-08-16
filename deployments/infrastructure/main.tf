terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "6.8.0"
    
    }
  }
}

provider "google" {
  project = var.resource_project_id
  region = var.resource_region
}

module "network" {
  source = "./modules/network"
  network_name = "car-valuation-agent-net"
  resource_region = var.resource_region
  subnet_1_name = var.subnet_1_name
  subnet_1_cidr_range = var.subnet_1_cidr_range
  subnet_2_name = var.subnet_2_name
  subnet_2_cidr_range = var.subnet_2_cidr_range
}

module "k8s" {
   source = "./modules/k8s"
   assigned_subnets = module.network.subnets
}