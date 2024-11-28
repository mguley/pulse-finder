# The AWS region where resources will be deployed.
variable "region" {
  description = "AWS region to deploy resources."
  default     = "eu-central-1" # Default AWS region is Frankfurt.
}

# Name prefix for the security group.
variable "security_group_name" {
  type        = string
  default     = "group name #1"
  description = "Name of the security group."
}

# Ingress rule variables for SSH access.
variable "ssh_port" {
  description = "Port number for SSH access."
  default     = 22
}

variable "ssh_protocol" {
  description = "Protocol for SSH access (e.g., tcp)."
  default     = "tcp"
}

variable "ssh_cidr_blocks" {
  description = "List of CIDR blocks allowed for SSH access."
  default     = ["0.0.0.0/0"] # Default allows access from all IPs (use with caution).
}

# Ingress rule variables for HTTP access.
variable "http_port" {
  description = "Port number for HTTP access."
  default     = 80
}

variable "http_protocol" {
  description = "Protocol for HTTP access (e.g., tcp)."
  default     = "tcp"
}

variable "http_cidr_blocks" {
  description = "List of CIDR blocks allowed for HTTP access."
  default     = ["0.0.0.0/0"] # Default allows access from all IPs (use with caution).
}

# Egress rule variables to allow outbound traffic.
variable "egress_from_port" {
  description = "Starting port for egress traffic."
  default     = 0
}

variable "egress_to_port" {
  description = "Ending port for egress traffic."
  default     = 0
}

variable "egress_protocol" {
  description = "Protocol for egress traffic (e.g., '-1' for all)."
  default     = "-1"
}

variable "egress_cidr" {
  description = "CIDR block for egress traffic."
  default     = ["0.0.0.0/0"] # Default allows access to all destinations (use with caution).
}

# Tag to help identify the security group.
variable "security_group_tag" {
  description = "Tag for the security group."
  default     = "security_group_tag"
}
