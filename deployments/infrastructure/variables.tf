variable "resource_project_id" {
  type = string
}

variable "resource_region" {
   type = string
   default = "asia-southeast2"
}

variable "subnet_1_name" {
  type = string
  default = "car-valuation-agent-subnet-1"
}

variable "subnet_1_cidr_range" {
  type = string
}

variable "subnet_2_name" {
  type = string
  default = "car-valuation-agent-subnet-2"
}

variable "subnet_2_cidr_range" {
  type = string
}
