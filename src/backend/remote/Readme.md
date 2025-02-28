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
#### Step 6. Configure SSH access

Allow SSH access for your subnet
```bash
sudo ufw allow from <ALLOWED-SUBNET>/16 to any port 22 comment "Allow SSH from our subnet"
```

Deny SSH for all other IPs
```bash
sudo ufw deny 22 comment "Deny SSH for all other IPs"
```

Status
```bash
sudo ufw status numbered
```
Reload
```bash
sudo ufw reload
```

---
#### Step 7. Verify Setup

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
#### Step 8. Configure DNS Records and Wait for Propagation

Before proceeding to generate SSL certificates, ensure that your domain is correctly pointed to the Digital Ocean droplet
by updating the DNS records through your domain registrar (e.g., GoDaddy):

1. `Update the A Record`: Add or update the `A` record for your domain to point to the public IP address of the Digital
    Ocean droplet
2. `Update the www Subdomain`: Similarly, update the `A` record for the `www`subdomain to point to the same public IP address.
3. `Wait for DNS Propagation`: DNS changes may take some time to propagate (typically a few minutes to several hours).
   Use tools like `WhatsMyDNS` (https://www.whatsmydns.net/) to confirm the propagation.

Once the domain resolves to the droplet's public IP address, you can proceed to the next step.

---
#### Step 9. Generate SSL Certificates and Deploy Nginx (Run Only Once)

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
#### Step 10. Deploy the application (API)
Use `make` to deploy the application to the production server.
```bash
make production/deploy/api
```  

This command will:
- Deploy the new binary, Nginx configuration files and service files
- Restart services as necessary

---
#### Step 11. Generate TLS certificates (gRPC)
```bash
make production/request-tls-certificates
```

---
#### Step 12. Deploy the application (gRPC)
```bash
make production/deploy/grpc
```

---
#### Step 13. Block access to all ports (incoming) except 443 (Nginx)

```bash
    # IPv4: Block first and second halves of the address space
    ufw deny from 0.0.0.0/1 to any comment "Block first half of IPv4 space"
    ufw deny from 128.0.0.0/1 to any comment "Block second half of IPv4 space"

    # IPv6: Block first and second halves of the address space
    ufw deny from ::/1 to any comment "Block first half of IPv6 space"
    ufw deny from 8000::/1 to any comment "Block second half of IPv6 space"
```

---
#### Summary

- The `SSL certificate generation` step is required `only once` after setting up the infrastructure and updating DNS records.
- Subsequent deployments can proceed directly with ```make production/deploy/api ``` without repeating the SSL setup.

---
### IP Blocking Feature

The application now supports banning specific IP addresses at the Nginx level.

##### How it works:

- A geo map named `$ban_ip` is defined in `api-nginx.conf` that checks if an IP is listed in `/etc/nginx/blocked_ips.conf`.
- If `$ban_ip` is set to `1`, Nginx returns a `403 Forbidden` for that client's requests.
- If an IP is not listed in the file (the default case), `$ban_ip` remains `0`, and the request proceeds normally.

##### Where is `blocked_ips.conf` located?

- On the production server, it is stored at `/etc/nginx/blocked_ips.conf`
- This file contains one or more lines mapping IP addresses to `1`, for example:

```
1.2.3.4 1;
1.2.3.5 1;
```

This would block both `1.2.3.4` and `1.2.3.5`.

##### Setup

- Open `/etc/nginx/blocked_ips.conf` in a text editor.
```bash
sudo nano /etc/nginx/blocked_ips.conf
```

- Add the IP address you wish to block.
```
# This file contains dynamically banned IPs
1.2.3.4 1;
```

- Save and exit.
- Reload or restart Nginx to apply the changes
```bash
sudo systemctl reload nginx
```


---
### Blackhole Routing for Flooding Mitigation

When dealing with SYN floods or other denial-of-service (DoS) attacks, one effective approach is to use blackhole routing.
It prevents your server from even attempting to respond to unwanted traffic, reducing resource consumption.

##### Problem: SYN Floods and Half-Open TCP Connections
- A SYN flood attack leaves many connections in a `SYN_RECV` state, consuming server resources.
- The goal is to block or drop these connections before they reach critical services like Nginx.

##### Mitigation Strategies

---
1. Verify or Strengthen UFW Rules

`Ensure Correct UFW Configuration:`

- Use `ufw` to block the attacking subnet:
```bash
sudo ufw deny from 1.2.0.0/16 to any
```
- To block specific ports, like HTTPS (port 443):
```bash
sudo ufw deny from 1.2.0.0/16 to any port 443
```

`Verify Rules:`

- Check Rule Presence:
```bash
sudo ufw status numbered
```
- Ensure Correct Rule Order: UFW processes rules sequentially. Ensure no broader `ALLOW` rules precede your `DENY` rule.
- Reload UFW
```bash
sudo ufw reload
```

---
2. Use Blackhole Routes

Adding a blackhole route silently discards all traffic from a specific subnet or IP.
```bash
sudo ip route add blackhole 1.2.0.0/16
```
This bypasses connection tracking, reducing overhead by immediately dropping traffic.

`View Blackhole Routes:`
```bash
ip route | grep blackhole
```

`Remove a Blackhole Route:`
```bash
sudo ip route del blackhole 1.2.0.0/16
```
Use blackhole routes with care, as they indiscriminately drop all traffic from the specified subnet.

---
3. Drop Traffic at the System Level (iptables)

If UFW is insufficient, use direct `iptables` rules for more control.
```bash
sudo iptables -I INPUT -s 1.2.0.0/16 -j DROP
```

This rule:
- Inserts a DROP rule at the top of the `INPUT` chain.
- Stops processing for traffic from the specified subnet.

`View Rules:`
```bash
sudo iptables -L -n
```

---
4. Rate Limiting with UFW or iptables

Instead of outright blocking, limit connection attempts:
```bash
sudo ufw limit https/tcp
```

iptables Example:
```bash
sudo iptables -A INPUT -p tcp --syn --dport 443 -m limit --limit 10/second --limit-burst 20 -j ACCEPT
sudo iptables -A INPUT -p tcp --syn --dport 443 -j DROP
```

This limits new connections to 10 per second (with a burst of up to 20), dropping excessive requests.