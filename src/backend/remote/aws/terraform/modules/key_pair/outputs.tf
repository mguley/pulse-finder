# The name of the created key pair (useful for reference in other modules or outputs).
output "key_name" {
  value       = aws_key_pair.key_pair.key_name
  description = "The name of the created key pair."
}

# The unique fingerprint of the created key pair (useful for verification purposes).
output "fingerprint" {
  value       = aws_key_pair.key_pair.fingerprint
  description = "The fingerprint of the created key pair."
}

# The unique ID assigned to the key pair by AWS.
output "key_pair_id" {
  value       = aws_key_pair.key_pair.id
  description = "The unique ID of the created key pair."
}
