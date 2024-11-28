# The AWS region where the resources will be deployed.
variable "region" {
  type        = string
  default     = "eu-central-1" # Default AWS region is Frankfurt; can be overridden.
  description = "AWS region to deploy resources."
}