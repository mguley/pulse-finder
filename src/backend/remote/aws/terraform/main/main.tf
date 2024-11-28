terraform {
  required_version = ">= 1.0.0" # Ensures Terraform version compatibility.

  required_providers {
    aws = {
      source  = "hashicorp/aws" # Specifies the AWS provider source.
      version = "~> 4.0"        # Locks the provider version to avoid breaking changes.
    }
  }
}

# Configures the AWS provider to use the specified region.
provider "aws" {
  region = var.region # The AWS region for deploying resources.
}

# Data block to fetch the default VPC in the AWS region.
data "aws_vpc" "default" {
  default = true # Ensures the module uses the default VPC.
}

# Data block to fetch all subnets in the default VPC.
data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id] # Filter subnets belonging to the default VPC.
  }
}

# Select the first subnet from the list.
locals {
  selected_subnet_id = data.aws_subnets.default.ids[0] # Select the first subnet ID.
}

# Key Pair Module: Manages the creation of an SSH key pair.
module "key_pair" {
  source       = "../modules/key_pair"                     # Path to the key pair module.
  key_name     = "pulse-finder-key-pair"                   # Name assigned to the key pair.
  public_key   = file("~/.ssh/pulse-finder-key-pair.pub")  # Public key file path for the SSH key.
  environment  = "prod"                                    # Specifies the environment (e.g., prod, dev).
  region       = var.region                                # The AWS region passed to the module.
}

# Security Group Module: Creates a security group with rules for EC2 instance access.
module "security_group_ec2" {
  source              = "../modules/security_group"       # Path to the security group module.
  security_group_name = "ec2-security-group"              # Name prefix for the security group.
  security_group_tag  = "security_group_ec2"              # Tag for identifying the security group.
  http_cidr_blocks    = ["0.0.0.0/0"]                     # CIDR blocks allowed for HTTP access. Restrict access per IP in production.
  ssh_cidr_blocks     = ["0.0.0.0/0"]                     # CIDR blocks allowed for SSH access. Restrict access per IP in production.
  region              = var.region                        # The AWS region passed to the module.
}

# EC2 Module: Launches an EC2 instance with specified configuration.
module "ec2_instance" {
  source              = "../modules/ec2"                            # Path to the EC2 module.
  key_name            = module.key_pair.key_name                    # Key pair used for SSH access.
  security_group_id   = module.security_group_ec2.security_group_id # Security group to control network access.
  subnet_id           = local.selected_subnet_id                    # Subnet ID selected from available subnets.
  instance_name       = "pulse-finder-api"                          # Name tag for the EC2 instance.
  environment         = "prod"                                      # Environment tag (e.g., prod, dev).
  region              = var.region                                  # The AWS region passed to the module.
}