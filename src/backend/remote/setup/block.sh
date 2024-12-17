#!/bin/bash
set -euo pipefail  # Enhanced error handling: -e (exit on error), -u (treat unset variables as errors), -o pipefail (catch errors in pipelines)

# ======================================================================== #
# VARIABLES
# ======================================================================== #

STATS_FILE="ufw_stats_output.txt"

# ======================================================================== #
# FUNCTIONS
# ======================================================================== #

# Check if the statistics file exists
check_stats_file() {
    if [ ! -f "$STATS_FILE" ]; then
        echo "Error: Statistics file '$STATS_FILE' not found! Please generate it first."
        exit 1
    fi
}

# Extract the top offending IPs from the UFW section
extract_offending_ips() {
    echo "Extracting Top 10 Offending IPs..."
    grep -A 10 "Top 10 Offending IPs:" "$STATS_FILE" | \
        awk '{if ($2 ~ /^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$/) print $2}' | sort -u
}

# Extract the top IPs from Failed SSH Login Attempts
extract_failed_ssh_ips() {
    echo "Extracting Top 10 Failed SSH Login Attempt IPs..."
    grep -A 10 "Top 10 Failed SSH Login Attempts" "$STATS_FILE" | \
        awk '{if ($2 ~ /^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$/) print $2}' | sort -u
}

# Block a list of IPs using UFW
block_ips() {
    local description="$1"
    shift
    local ips=("$@")

    if [[ ${#ips[@]} -eq 0 ]]; then
        echo "No IPs found to block for $description."
        return
    fi

    echo "Blocking IPs from $description:"
    for ip in "${ips[@]}"; do
        if [[ -n "$ip" && "$ip" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            echo "Blocking IP: $ip"
            sudo ufw deny from "$ip" comment "Blocked due to $description"
        else
            echo "Skipping invalid IP: $ip"
        fi
    done
}

# Reload UFW to apply the changes
reload_ufw() {
    echo "Reloading UFW to apply changes..."
    sudo ufw reload
    echo "UFW rules reloaded successfully."
}

# Main function to orchestrate IP blocking
block_offending_and_ssh_ips() {
    check_stats_file

    # Extract and block IPs
    mapfile -t offending_ips < <(extract_offending_ips)
    block_ips "Top Offending IPs" "${offending_ips[@]}"

    mapfile -t failed_ssh_ips < <(extract_failed_ssh_ips)
    block_ips "Failed SSH Login Attempts" "${failed_ssh_ips[@]}"

    # Reload UFW
    reload_ufw

    echo "All offending and failed SSH login IPs have been successfully blocked."
}

# ======================================================================== #
# MAIN SCRIPT
# ======================================================================== #

main() {
    echo "==========================="
    echo " UFW IP Blocking Script"
    echo "==========================="
    block_offending_and_ssh_ips
    echo "Script execution completed."
}

main "$@"
