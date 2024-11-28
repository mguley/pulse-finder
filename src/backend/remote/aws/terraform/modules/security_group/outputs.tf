# The name of the created security group.
output "security_group_name" {
  value       = aws_security_group.sg.name
  description = "The name of the created security group."
}

# The unique ID assigned to the security group by AWS.
output "security_group_id" {
  value       = aws_security_group.sg.id
  description = "The unique ID of the created security group."
}
