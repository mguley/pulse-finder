# The unique ID of the EC2 instance.
output "instance_id" {
  value       = aws_instance.ec2_instance.id
  description = "The unique ID of the created EC2 instance."
}

# The public IP address of the EC2 instance.
output "instance_public_ip" {
  value       = aws_instance.ec2_instance.public_ip
  description = "The public IP address of the created EC2 instance."
}

# The private IP address of the EC2 instance.
output "instance_private_ip" {
  value       = aws_instance.ec2_instance.private_ip
  description = "The private IP address of the created EC2 instance."
}

# The ID of the subnet where the EC2 instance is deployed.
output "instance_subnet_id" {
  value       = aws_instance.ec2_instance.subnet_id
  description = "The ID of the subnet where the EC2 instance is placed."
}
