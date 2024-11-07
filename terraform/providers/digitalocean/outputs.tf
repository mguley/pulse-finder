# Display the IP address of the created droplet
output "pulse_finder_ip" {
  description = "The public IP address of the pulse_finder droplet"
  value       = digitalocean_droplet.pulse_finder.ipv4_address
}

output "pulse_finder_price_monthly" {
  description = "Droplet monthly price"
  value       = digitalocean_droplet.pulse_finder.price_monthly
}

output "pulse_finder_price_hourly" {
  description = "Droplet hourly price"
  value       = digitalocean_droplet.pulse_finder.price_hourly
}
