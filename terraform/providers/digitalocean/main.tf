# Define your droplet
resource "digitalocean_droplet" "pulse_finder" {
  image     = "ubuntu-24-04-x64"           # Specify Ubuntu 24.04
  name      = "pulse-finder-production"    # Droplet name
  region    = "fra1"                       # Frankfurt data center
  size      = "s-1vcpu-512mb-10gb"         # Small size; adjust as needed
  ssh_keys  = [var.ssh_fingerprint]        # SSH key fingerprint to access the droplet
  tags      = ["pulse", "finder", "production"]
}
