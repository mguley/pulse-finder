### Deployment and Hosting Guide for AWS EC2 Instance

This guide outlines the application deployment process on an AWS EC2 instance using Ubuntu.

### Infrastructure Setup

This deployment uses:
- **AWS EC2 instance** running Ubuntu
- **API Gateway** as the reverse proxy for the Go application that exposes API endpoints.

---
#### Step 1: Initialize Infrastructure with Terraform
- We assume that AWS CLI is configured and ready to use.
- Use Terraform scripts to provision the AWS EC2 instance.
- Commands remain unchanged:
```
    terraform init
    terraform plan
    terraform apply
```

---
#### Step 2: Note the EC2 Instance Public IP
- After Terraform applies the changes, note the `Instance Public IP address`.

---
#### Step 3: Update `.envrc` with EC2 Public IP
- Update the `.envrc` file with the EC2 Public IP address for deployment commands

---
#### Step 4: Run the Setup Script
- Copy the `02.sh` script to the EC2 instance
```bash
rsync -rP --delete -e "ssh -i ~/.ssh/pulse-finder-key-pair" ./setup ubuntu@<INSTANCE-IP>:/home/ubuntu
```

- Run the script
```bash
ssh -i ~/.ssh/pulse-finder-key-pair -t ubuntu@<INSTANCE-IP> "bash ./setup/02.sh"
```

---
#### Step 5: Initial Login as `pulsefinder` user
After the script completes, log in as the `pulsefinder` user for additional setup.
```bash
ssh -i ~/.ssh/pulse-finder-key-pair pulsefinder@<INSTANCE-IP>
```
On your first login, you'll be prompted to set a password for the `pulsefinder` user.