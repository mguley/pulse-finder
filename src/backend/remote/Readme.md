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
- Check Migrate CLI version
```bash
migrate -version
```  
- Test PostgreSQL connection
```bash
psql $DB_DSN
``` 
---
#### Step 7. Configure DNS Records and Wait for Propagation

Before proceeding to generate SSL certificates, ensure that your domain is correctly pointed to the Digital Ocean droplet
by updating the DNS records through your domain registrar (e.g., GoDaddy):

1. `Update the A Record`: Add or update the `A` record for your domain to point to the public IP address of the Digital
    Ocean droplet
2. `Update the www Subdomain`: Similarly, update the `A` record for the `www`subdomain to point to the same public IP address.
3. `Wait for DNS Propagation`: DNS changes may take some time to propagate (typically a few minutes to several hours).
   Use tools like `WhatsMyDNS` (https://www.whatsmydns.net/) to confirm the propagation.

Once the domain resolves to the droplet's public IP address, you can proceed to the next step.

---
#### Step 8. Generate SSL Certificates and Deploy Nginx (Run Only Once)

After DNS propagation is complete, generate SSL certificates and deploy Nginx by running the following command:
```bash
make production/deploy/nginx-and-ssl
```

This command will:
- Deploy a temporary Nginx configuration to allow Certbot to generate SSL certificates.
- Generate new SSL certificates for the domain using Let's Encrypt.
- Apply the final Nginx configuration using the generated certificates.

Once completed, your server will be configured to serve traffic over HTTPS.

---
#### Step 9. Deploy the application
Use `make` to deploy the application to the production server.
```bash
make production/deploy/api
```  

This command will:
- Deploy the new binary, Nginx configuration files and service files
- Restart services as necessary

---
#### Summary

- The `SSL certificate generation` step is required `only once` after setting up the infrastructure and updating DNS records.
- Subsequent deployments can proceed directly with ```make production/deploy/api ``` without repeating the SSL setup.