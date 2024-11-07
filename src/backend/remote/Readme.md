### Deployment and Hosting Guide

This guide outlines the application deployment process on a Digital Ocean droplet.

### Infrastructure Setup

This deployment uses:
- **Digital Ocean Droplet** running Ubuntu
- **Nginx** for reverse proxy configuration

---
#### Step 1: Initialize Infrastructure with Terraform

Initialize configuration, preview change, and apply the changes.

```bash
terraform init
```
```bash
terraform plan
```
```bash
terraform apply
```

---
#### Step 2: Note the Droplet IP

After applying the changes, note the `Droplet IP address`. You will need this IP address for the following steps.

---
#### Step 3: Update `.envrc` with Droplet IP

Copy the `Droplet IP address` and update the `.envrc` file to reflect the new server IP.

This will allow future deployment commands to work correctly.

---
#### Step 4. Run the Setup Script

Copy the `01.sh` setup script to the new droplet. This script will install necessary software and configure the system.
```bash
rsync -rP --delete ./remote/setup root@<DROPLET_IP>:/root
```

Log into the droplet as root and run the script
```bash
ssh -t root@<DROPLET_IP> "bash /root/setup/01.sh"
```

This will:
- set up the firewall
- install required software
- configure Nginx

---
#### Step 5. Initial Login as `pulsefinder` User

After the script completes, log in as the `pulsefinder` user for additional setup.
```bash
ssh pulsefinder@<DROPLET_IP>
```
On your first login, you'll be prompted to set a password for the `pulsefinder` user.

---
#### Step 6. Verify Setup

After logging in, verify the setup by running the following commands:
- Check Nginx status
```bash
sudo systemctl status nginx
```  

---
#### Step 7. Update `api-nginx.conf` with Droplet IP

Before deploying the application, update `api-nginx.conf` with the droplet IP. This ensures Nginx serves the API on
the correct API address.

---
#### Step 8. Deploy the application
Use `make` to deploy the application to the production server.
```bash
make production/deploy/api
```  

This command will:
- Deploy the new binary, Nginx configuration files and service files
- Restart services as necessary