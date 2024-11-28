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

# Outputs for Security Group

# The ID of the security group created for EC2.
output "security_group_ec2_id" {
  value       = module.security_group_ec2.security_group_id
  description = "The ID of the security group (EC2)."
}

# The name of the security group created for EC2.
output "security_group_ec2_name" {
  value       = module.security_group_ec2.security_group_name
  description = "The name of the security group (EC2)."
}