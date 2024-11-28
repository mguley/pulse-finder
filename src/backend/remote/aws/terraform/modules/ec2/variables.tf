# The AWS region where the EC2 instance will be deployed.
variable "region" {
  type        = string
  default     = "eu-central-1" # Default region is Frankfurt.
  description = "AWS region to deploy resources."
}

# The AMI ID for the EC2 instance's operating system (Ubuntu in this case).
variable "ami_id" {
  description = "AMI ID for the Ubuntu operating system."
  type        = string
  default     = "ami-0084a47cc718c111a" # AMI ID for Ubuntu Server 24.04 LTS (HVM), SSD Volume Type.
}

# The EC2 instance type (e.g., t2.micro for free tier eligibility).
variable "instance_type" {
  description = "EC2 instance type (e.g., t2.micro for free-tier eligibility)."
  type        = string
  default     = "t2.micro"
}

# The name of the SSH key pair to be attached to the EC2 instance for secure access.
variable "key_name" {
  description = "Name of the key pair to attach to the EC2 instance for SSH access."
  type        = string
}

# The ID of the security group to associate with the EC2 instance.
variable "security_group_id" {
  description = "ID of the security group to attach to the EC2 instance."
  type        = string
}

# The name tag for identifying the EC2 instance.
variable "instance_name" {
  description = "Name tag for the EC2 instance."
  type        = string
  default     = "ubuntu-instance"
}

# The environment tag for the EC2 instance (e.g., dev, prod, test).
variable "environment" {
  description = "Environment tag for the EC2 instance (e.g., dev, prod, test)."
  type        = string
  default     = "dev"
}

# The size of the root volume in GB for the EC2 instance.
variable "volume_size" {
  description = "The size of the root volume in GB."
  type        = number
  default     = 8
}

# The type of volume to attach to the EC2 instance (e.g., gp3, gp2, io1).
variable "volume_type" {
  description = "The type of volume to attach to the instance."
  type        = string
  default     = "gp3"
}

# The subnet ID where the EC2 instance will be deployed.
variable "subnet_id" {
  description = "The ID of the subnet where the EC2 instance will be placed."
  type        = string
}
