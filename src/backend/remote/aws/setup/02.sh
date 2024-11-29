#!/bin/bash
set -euo pipefail

# ======================================================================== #
# VARIABLES
# ======================================================================== #

TIMEZONE="Europe/Berlin"             # Set the timezone
USERNAME="pulsefinder"               # Name of the new user to create
DB_PASSWORD=""                       # Password for Postgres pulse_finder user (prompted)
JWT_SECRET="meewuZ8Hei2theefaiK9ahphie5oiDai1eiX6ehaeveeThae1oocha2sooMeeguog"

# Prompt user for the DB password
read -rsp "Enter password for pulse_finder DB user: " DB_PASSWORD
echo ""

# Export locale to avoid any locale-related errors
export LC_ALL=en_US.UTF-8

# ======================================================================== #
# FUNCTIONS
# ======================================================================== #

# Update and upgrade system packages
update_system() {
    echo "Updating system packages..."
    sudo apt-get update -q
}

# Set timezone and install locales
setup_time_and_locale() {
    echo "Setting timezone to ${TIMEZONE}..."
    sudo timedatectl set-timezone "${TIMEZONE}"
    echo "Installing all locales..."
    sudo apt-get --yes install locales-all
}

# Create a new user with sudo privileges and SSH access
create_user() {
    if id "${USERNAME}" &>/dev/null; then
        echo "User ${USERNAME} already exists. Skipping creation."
    else
        echo "Creating user ${USERNAME}..."
        sudo useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"
        sudo passwd --delete "${USERNAME}"
        sudo chage --lastday 0 "${USERNAME}"

        echo "Copying SSH keys to new user..."
        if [ -d "/home/ubuntu/.ssh" ]; then
            sudo mkdir -p /home/"${USERNAME}"/.ssh
            sudo rsync --archive --chown="${USERNAME}:${USERNAME}" /home/ubuntu/.ssh/ /home/"${USERNAME}"/.ssh/
            sudo chmod 700 /home/"${USERNAME}"/.ssh
            sudo chmod 600 /home/"${USERNAME}"/.ssh/authorized_keys
        else
            echo "No SSH keys found in /home/ubuntu/.ssh. Skipping SSH key copy."
        fi
    fi
}

# Configure firewall to allow SSH
configure_firewall() {
    echo "Configuring firewall to allow SSH..."
    sudo ufw allow 22       # SSH
    sudo ufw --force enable
}

# Set environment variables
set_environment_variables() {
    echo "Adding environment variables to /etc/environment..."
    grep -qxF "PORT=4000" /etc/environment || echo "PORT=4000" | sudo tee -a /etc/environment >/dev/null
    grep -qxF "ENV=production" /etc/environment || echo "ENV=production" | sudo tee -a /etc/environment >/dev/null
    grep -qxF "DB_DSN='postgres://pulse_finder:${DB_PASSWORD}@localhost/pulse_finder'" /etc/environment || echo "DB_DSN='postgres://pulse_finder:${DB_PASSWORD}@localhost/pulse_finder'" | sudo tee -a /etc/environment >/dev/null
    grep -qxF "JWT_SECRET=${JWT_SECRET}" /etc/environment || echo "JWT_SECRET=${JWT_SECRET}" | sudo tee -a /etc/environment >/dev/null

    echo "Reloading environment variables..."
    source /etc/environment

    echo "Environment variables added successfully!"
}

# Install migrate CLI for database migrations
install_migrate_cli() {
    echo "Installing the migrate CLI tool..."
    sudo curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz && \
    sudo mv migrate /usr/local/bin/migrate
}

# Install Postgres and set up pulse_finder database
install_and_setup_postgresql() {
    echo "Installing Postgres..."
    sudo apt-get --yes install postgresql

    echo "Setting up pulse_finder database and user..."
    sudo -i -u postgres psql -c "CREATE DATABASE pulse_finder"
    sudo -i -u postgres psql -d pulse_finder -c "CREATE EXTENSION IF NOT EXISTS citext"
    sudo -i -u postgres psql -c "CREATE ROLE pulse_finder WITH LOGIN PASSWORD '${DB_PASSWORD}'"
    sudo -i -u postgres psql -c "ALTER DATABASE pulse_finder OWNER TO pulse_finder"
}

# Install and configure NATS server
install_nats() {
    echo "Installing NATS server..."
    sudo curl -L https://github.com/nats-io/nats-server/releases/download/v2.10.22/nats-server-v2.10.22-linux-amd64.tar.gz | tar xz
    sudo mv nats-server-v2.10.22-linux-amd64/nats-server /usr/local/bin/
    sudo chmod +x /usr/local/bin/nats-server
}

# Upgrade system packages
upgrade_system() {
    echo "Upgrading all system packages..."
    sudo apt-get --yes -o Dpkg::Options::="--force-confnew" upgrade
}

# ======================================================================== #
# MAIN SCRIPT
# ======================================================================== #

main() {
    update_system
    setup_time_and_locale
    create_user
    configure_firewall
    install_migrate_cli
    install_and_setup_postgresql
    set_environment_variables
    install_nats
    upgrade_system

    echo "Script complete! Rebooting..."
    sudo reboot
}

main "$@"