#!/bin/bash
set -uo pipefail  # Avoid stopping script execution on empty output

LOG_FILE="/var/log/ufw.log"
OUTPUT_FILE="ufw_stats_output.txt"

if [ ! -f "$LOG_FILE" ]; then
    echo "UFW log file not found at $LOG_FILE. Please ensure UFW logging is enabled."
    exit 1
fi

# Redirect all output to the file
exec > "$OUTPUT_FILE" 2>&1

echo "==========================="
echo " UFW Activity Statistics"
echo "==========================="

# Total Blocked Connections
echo ""
echo "Blocked Connections: $(grep -c 'BLOCK' "$LOG_FILE" || echo "0")"

# Total Allowed Connections
echo ""
echo "Allowed Connections: $(grep -c 'ALLOW' "$LOG_FILE" || echo "0")"

# Top 10 Blocked IPs
echo ""
echo "Top 10 Blocked IPs:"
grep 'BLOCK' "$LOG_FILE" | awk '{for(i=1;i<=NF;i++) if($i ~ /^SRC=/) print $i}' | cut -d'=' -f2 | sort | uniq -c | sort -nr | head || echo "No blocked IPs found."

# Top 10 Accessed Ports
echo ""
echo "Top 10 Accessed Ports:"
grep 'BLOCK' "$LOG_FILE" | awk '{for(i=1;i<=NF;i++) if($i ~ /^DPT=/) print $i}' | cut -d'=' -f2 | sort | uniq -c | sort -nr | head || echo "No accessed ports found."

# Top Offending IPs
echo ""
echo "Top 10 Offending IPs:"
grep 'BLOCK' "$LOG_FILE" | awk '{for(i=1;i<=NF;i++) if($i ~ /^SRC=/) print $i}' | cut -d'=' -f2 | sort | uniq -c | sort -nr | head || echo "No offending IPs found."

# Top Targeted Ports
echo ""
echo "Top 10 Targeted Ports:"
grep 'BLOCK' "$LOG_FILE" | grep -oP 'DPT=\K[0-9]+' | sort | uniq -c | sort -nr | head || echo "No targeted ports found."

# Suspicious Small Packets (Likely Port Scans or Probes)
echo ""
echo "Suspicious Small Packets (Likely Probes):"
grep 'BLOCK' "$LOG_FILE" | awk '{for(i=1;i<=NF;i++) if($i ~ /^LEN=/) print $i}' | cut -d'=' -f2 | sort -n | uniq -c | head || echo "No suspicious small packets found."

# SYN Flood Indicators (Multiple SYN Packets from Same IP)
echo ""
echo "Possible SYN Flood IPs (Top 5):"
grep 'BLOCK' "$LOG_FILE" | grep 'SYN' | awk '{for(i=1;i<=NF;i++) if($i ~ /^SRC=/) print $i}' | cut -d'=' -f2 | sort | uniq -c | sort -nr | head -n 5 || echo "No SYN flood indicators found."

# Real-Time Monitoring Instructions
echo ""
echo "Tip: Monitor real-time traffic using: tail -f /var/log/ufw.log"
echo ""
echo "==========================="

# Ensure content is flushed to the file
sync
