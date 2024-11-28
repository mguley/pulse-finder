terraform {
  required_version = ">= 1.0.0" # Ensures Terraform version compatibility.

  required_providers {
    aws = {
      source  = "hashicorp/aws" # Specifies the AWS provider source.
      version = "~> 4.0"        # Locks the provider version to avoid breaking changes.
    }
  }
}

# AWS Security Group resource to manage ingress and egress rules for network traffic.
resource "aws_security_group" "sg" {
  name_prefix = var.security_group_name # Prefix for the security group name.

  # Ingress rule to allow SSH access.
  ingress {
    from_port   = var.ssh_port          # Starting port for SSH access.
    to_port     = var.ssh_port          # Ending port for SSH access.
    protocol    = var.ssh_protocol      # Protocol to allow (e.g., tcp).
    cidr_blocks = var.ssh_cidr_blocks   # List of allowed CIDR blocks for SSH.
  }

  # Ingress rule to allow HTTP access.
  ingress {
    from_port   = var.http_port         # Starting port for HTTP access.
    to_port     = var.http_port         # Ending port for HTTP access.
    protocol    = var.http_protocol     # Protocol to allow (e.g., tcp).
    cidr_blocks = var.http_cidr_blocks  # List of allowed CIDR blocks for HTTP.
  }

  # Egress rule to allow outbound traffic.
  egress {
    from_port   = var.egress_from_port  # Starting port for outbound traffic.
    to_port     = var.egress_to_port    # Ending port for outbound traffic.
    protocol    = var.egress_protocol   # Protocol to allow (e.g., "-1" for all).
    cidr_blocks = var.egress_cidr       # List of CIDR blocks for outbound traffic.
  }

  # Tags for identifying the security group.
  tags = {
    Name = var.security_group_tag # Name tag for the security group.
  }
}
