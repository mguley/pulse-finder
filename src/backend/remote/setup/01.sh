#!/bin/bash
set -euo pipefail  # Enhanced error handling: -e (exit on error), -u (treat unset variables as errors), -o pipefail (catch errors in pipelines)

# ======================================================================== #
# VARIABLES
# ======================================================================== #

TIMEZONE="Europe/Berlin"             # Set the timezone
USERNAME="pulsefinder"               # Name of the new user to create
DB_PASSWORD=""                       # Password for Postgres pulse_finder user (prompted)
JWT_SECRET="meewuZ8Hei2theefaiK9ahphie5oiDai1eiX6ehaeveeThae1oocha2sooMeeguog"

SSL_DIR="/etc/nginx/ssl"             # Directory to store SSL certificates
SSL_KEY="${SSL_DIR}/selfsigned.key"  # Path to SSL private key
SSL_CERT="${SSL_DIR}/selfsigned.crt" # Path to SSL certificate

# Prompt user for the DB password
read -rsp "Enter password for pulse_finder DB user: " DB_PASSWORD
echo ""

# Export locale to avoid any locale-related errors.
export LC_ALL=en_US.UTF-8

# ======================================================================== #
# FUNCTIONS
# ======================================================================== #

# Generate a self-signed SSL certificate
generate_self_signed_cert() {
    echo "Generating a self-signed SSL certificate..."

    # Create the SSL directory if it doesn't exist
    mkdir -p "${SSL_DIR}"

    # Generate the SSL certificate and private key
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout "${SSL_KEY}" -out "${SSL_CERT}" \
        -subj "/CN=pulse-finder-api"

    echo "Self-signed SSL certificate created at:"
    echo "  Certificate: ${SSL_CERT}"
    echo "  Key: ${SSL_KEY}"

    # Adjust permissions
    chown -R ${USERNAME}:${USERNAME} "${SSL_DIR}"
    chmod 600 "${SSL_DIR}"/*
}

# Remove default Nginx configuration
remove_default_nginx_config() {
    echo "Removing default Nginx configuration..."
    rm -f /etc/nginx/sites-enabled/default /etc/nginx/sites-available/default || true
    echo "Default Nginx configuration removed."
}

# Validate Nginx configuration
validate_nginx_config() {
    echo "Validating Nginx configuration..."
    if ! nginx -t; then
        echo "Nginx configuration test failed! Exiting."
        exit 1
    fi
    echo "Nginx configuration is valid."
}

# Update and upgrade system packages
update_system() {
    echo "Updating system packages..."
    apt-get update -q
}

# Enable necessary repositories
enable_repositories() {
    echo "Enabling universe repository..."
    add-apt-repository --yes universe
}

# Set timezone and install locales
setup_time_and_locale() {
    echo "Setting timezone to ${TIMEZONE}..."
    timedatectl set-timezone "${TIMEZONE}"
    echo "Installing all locales..."
    apt-get --yes install locales-all
}

# Create a new user with sudo privileges and SSH access
create_user() {
    if id "${USERNAME}" &>/dev/null; then
        echo "User ${USERNAME} already exists. Skipping creation."
    else
        echo "Creating user ${USERNAME}..."
        useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"
        passwd --delete "${USERNAME}"
        chage --lastday 0 "${USERNAME}"

        echo "Copying SSH keys to new user..."
        if [ -d "/root/.ssh" ]; then
            rsync --archive --chown="${USERNAME}:${USERNAME}" /root/.ssh /home/"${USERNAME}"
        else
            echo "No SSH keys found in /root/.ssh. Skipping SSH key copy."
        fi
    fi
}

# Configure firewall to allow SSH, HTTPS
configure_firewall() {
    echo "Configuring firewall to allow SSH, HTTPS..."
    ufw allow 22           # SSH
    ufw allow 443/tcp      # HTTPS
    ufw --force enable
}

# Set environment variables
set_environment_variables() {
    echo "Adding environment variables to /etc/environment..."
    # Safeguard existing environment variables by adding only if they are not present
    grep -qxF "PORT=4000" /etc/environment || echo "PORT=4000" >> /etc/environment
    grep -qxF "ENV=production" /etc/environment || echo "ENV=production" >> /etc/environment
    grep -qxF "DB_DSN='postgres://pulse_finder:${DB_PASSWORD}@localhost/pulse_finder'" /etc/environment || echo "DB_DSN='postgres://pulse_finder:${DB_PASSWORD}@localhost/pulse_finder'" >> /etc/environment
    grep -qxF "JWT_SECRET=${JWT_SECRET}" /etc/environment || echo "JWT_SECRET=${JWT_SECRET}" >> /etc/environment
}

# Install migrate CLI for database migrations
install_migrate_cli() {
  echo "Installing the migrate CLI tool..."
  curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz && \
  mv migrate /usr/local/bin/migrate
}

# Install Postgres and set up pulse_finder database
install_and_setup_postgresql() {
    echo "Installing Postgres..."
    apt-get --yes install postgresql

    echo "Setting up pulse_finder database and user..."
    sudo -i -u postgres psql -c "CREATE DATABASE pulse_finder"
    sudo -i -u postgres psql -d pulse_finder -c "CREATE EXTENSION IF NOT EXISTS citext"
    sudo -i -u postgres psql -c "CREATE ROLE pulse_finder WITH LOGIN PASSWORD '${DB_PASSWORD}'"
    sudo -i -u postgres psql -c "ALTER DATABASE pulse_finder OWNER TO pulse_finder"
}

# Install and configure Nginx as the reverse proxy server
install_nginx() {
    echo "Installing Nginx web server..."
    apt-get update -q
    apt-get --yes install nginx

    # Enable Nginx to start on boot and start the service
    systemctl enable nginx
    systemctl start nginx
}

# Restart Nginx service
restart_nginx() {
    echo "Restarting Nginx to apply changes..."
    systemctl reload nginx
    echo "Nginx restarted successfully."
}

# Install and configure NATS server
install_nats() {
    echo "Installing NATS server..."
    # Download the latest NATS Server (nats-server) binary
    curl -L https://github.com/nats-io/nats-server/releases/download/v2.10.22/nats-server-v2.10.22-linux-amd64.tar.gz | tar xz
    # Move the binary to /usr/local/bin
    mv nats-server-v2.10.22-linux-amd64/nats-server /usr/local/bin/
    chmod +x /usr/local/bin/nats-server
}

# Upgrade system packages
upgrade_system() {
    echo "Upgrading all system packages..."
    apt-get --yes -o Dpkg::Options::="--force-confnew" upgrade
}

# ======================================================================== #
# MAIN SCRIPT
# ======================================================================== #

main() {
  enable_repositories
  update_system
  setup_time_and_locale
  create_user
  configure_firewall
  install_migrate_cli
  install_and_setup_postgresql
  set_environment_variables
  install_nginx
  remove_default_nginx_config
  install_nats
  upgrade_system
  generate_self_signed_cert
  validate_nginx_config
  restart_nginx

  echo "Script complete! Rebooting..."
  reboot
}

main "$@"
