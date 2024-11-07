# Digital Ocean API token
variable "digitalocean_token" {
  description = "Digital Ocean API token for authentication"
  type        = string
  sensitive   = true
}

# SSH key fingerprint for accessing the droplet
variable "ssh_fingerprint" {
  description = "Fingerprint of SSH key to access the droplet"
  type        = string
}
