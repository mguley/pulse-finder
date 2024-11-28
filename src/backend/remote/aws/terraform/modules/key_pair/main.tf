terraform {
  required_version = ">= 1.0.0" # Ensures Terraform version compatibility.

  required_providers {
    aws = {
      source  = "hashicorp/aws" # Specifies the AWS provider source.
      version = "~> 4.0"        # Locks the provider version to avoid breaking changes.
    }
  }
}

# AWS Key Pair resource for managing SSH keys.
resource "aws_key_pair" "key_pair" {
  key_name   = var.key_name   # Specifies the name of the key pair.
  public_key = var.public_key # The public key material (in PEM format) to associate with the key pair.

  # Tags for the key pair to help with resource identification and organization.
  tags = {
    Name        = var.key_name    # Key pair name tag for easy identification.
    Environment = var.environment # Environment tag to distinguish between dev, prod, etc.
  }
}