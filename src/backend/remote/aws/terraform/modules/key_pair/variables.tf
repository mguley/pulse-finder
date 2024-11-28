# The AWS region where resources will be deployed.
variable "region" {
  type        = string
  description = "AWS region to deploy resources."
  default     = "eu-central-1" # Default is Frankfurt region; can be overridden.
}

# The name of the key pair to be created.
variable "key_name" {
  type        = string
  description = "Name of the key pair."
}

# The public key to be associated with the key pair. Must be in RSA PEM format.
variable "public_key" {
  type        = string
  description = "The public key material to import. It must be an RSA key in PEM format."
}

# An environment tag to help distinguish key pairs for different environments.
variable "environment" {
  type        = string
  description = "Environment tag for the key pair (e.g., dev, prod, test)."
  default     = "dev" # Default is "dev" environment.
}
