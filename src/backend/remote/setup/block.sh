#!/bin/bash
set -euo pipefail  # Exit on errors and prevent unintentional script failures

STATS_FILE="ufw_stats_output.txt"

# Check if the stats file exists
if [ ! -f "$STATS_FILE" ]; then
    echo "Error: Statistics file '$STATS_FILE' not found!"
    exit 1
fi

echo "Blocking Top Offending IPs..."

# Extract the offending IPs and block them
grep -A 10 "Top 10 Offending IPs:" "$STATS_FILE" | \
    tail -n 10 | awk '{print $2}' | \
    while read -r ip; do
        if [[ -n "$ip" ]]; then
            echo "Blocking IP: $ip"
            sudo ufw deny from "$ip"
        fi
    done

echo "All offending IPs have been blocked."

# Reload UFW to apply the changes
sudo ufw reload
echo "UFW rules reloaded successfully."
