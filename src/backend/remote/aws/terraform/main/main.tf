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

# Key Pair Module: Manages the creation of an SSH key pair.
module "key_pair" {
  source       = "../modules/key_pair"                     # Path to the key pair module.
  key_name     = "pulse-finder-key-pair"                   # Name assigned to the key pair.
  public_key   = file("~/.ssh/pulse-finder-key-pair.pub")  # Public key file path for the SSH key.
  environment  = "prod"                                    # Specifies the environment (e.g., prod, dev).
  region       = var.region                                # The AWS region passed to the module.
}