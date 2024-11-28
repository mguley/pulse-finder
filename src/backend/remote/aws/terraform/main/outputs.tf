# Outputs for Key Pair

# The name of the key pair created by the module.
output "key_pair_name" {
  value       = module.key_pair.key_name
  description = "The name of the created key pair."
}

# The fingerprint of the key pair created by the module (useful for validation).
output "key_pair_fingerprint" {
  value       = module.key_pair.fingerprint
  description = "The fingerprint of the created key pair."
}