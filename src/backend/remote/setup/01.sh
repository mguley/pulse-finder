#!/bin/bash
set -eu

# ======================================================================== #
# VARIABLES
# ======================================================================== #

TIMEZONE="Europe/Berlin"             # Set the timezone
USERNAME="pulsefinder"               # Name of the new user to create

# Export locale to avoid any locale-related errors.
export LC_ALL=en_US.UTF-8

# ======================================================================== #
# FUNCTIONS
# ======================================================================== #

# Update and upgrade system packages
update_system() {
    echo "Updating system packages..."
    apt update
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
    apt --yes install locales-all
}

# Create a new user with sudo privileges and SSH access
create_user() {
    echo "Creating user ${USERNAME}..."
    useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"
    passwd --delete "${USERNAME}"
    chage --lastday 0 "${USERNAME}"

    echo "Copying SSH keys to new user..."
    rsync --archive --chown="${USERNAME}:${USERNAME}" /root/.ssh /home/"${USERNAME}"
}

# Configure firewall to allow SSH, HTTP, and HTTPS
configure_firewall() {
    echo "Configuring firewall to allow SSH, HTTP, and HTTPS..."
    ufw allow 22           # SSH
    ufw allow 80/tcp       # HTTP
    ufw allow 443/tcp      # HTTPS
    ufw --force enable
}

# Set environment variables
set_environment_variables() {
    echo "Adding environment variables to /etc/environment..."

    # Write all variables at once into /etc/environment
    {
        # Application settings
        echo "PORT=4000"
        echo "ENV=production"
    } >> /etc/environment
}

# Install and configure Nginx as the reverse proxy server
install_nginx() {
    echo "Installing Nginx web server..."
    apt update
    apt --yes install nginx

    # Enable Nginx to start on boot and start the service
    systemctl enable nginx
    systemctl start nginx
}

# Upgrade system packages
upgrade_system() {
    echo "Upgrading all system packages..."
    apt --yes -o Dpkg::Options::="--force-confnew" upgrade
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
  set_environment_variables
  install_nginx
  upgrade_system

  echo "Script complete! Rebooting..."
  reboot
}

main "$@"
