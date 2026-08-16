variable "resource_region" {
  type = string
  default = "asia-southeast-2"
}

variable "network_name" {
  type = string
  default = "car-valuation-ai-net"
}

variable "subnet_1_name" {
  type = string
  default = "car-valuation-ai-subnet-1"
}

variable "subnet_1_cidr_range" {
  type = string
}


variable "subnet_2_name" {
  type = string
  default = "car-valuation-ai-subnet-2"
}

variable "subnet_2_cidr_range" {
  type = string
}


