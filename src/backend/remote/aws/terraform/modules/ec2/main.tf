terraform {
  required_version = ">= 1.0.0" # Ensure Terraform version is 1.0.0 or higher for compatibility.

  required_providers {
    aws = {
      source  = "hashicorp/aws" # AWS provider sourced from the HashiCorp Terraform registry.
      version = "~> 4.0"        # Use version 4.x of the AWS provider for compatibility.
    }
  }
}

# AWS EC2 instance resource to create and manage a virtual server.
resource "aws_instance" "ec2_instance" {
  ami             = var.ami_id              # Amazon Machine Image (AMI) ID for the EC2 instance.
  instance_type   = var.instance_type       # EC2 instance type (e.g., t2.micro for free tier).
  key_name        = var.key_name            # Name of the SSH key pair attached to the instance.
  subnet_id       = var.subnet_id           # Subnet ID where the EC2 instance will be deployed.
  security_groups = [var.security_group_id] # Security group for controlling network access.

  # Configuration for the root block storage attached to the EC2 instance.
  root_block_device {
    volume_size           = var.volume_size          # Size of the root volume in GB.
    volume_type           = var.volume_type          # Type of the volume (e.g., gp3, gp2).
    delete_on_termination = true                     # Deletes the volume when the instance is terminated.
  }

  # Tags for identifying and organizing the EC2 instance.
  tags = {
    Name        = var.instance_name    # Name tag for the instance.
    Environment = var.environment      # Environment tag (e.g., dev, prod, test).
  }
}
