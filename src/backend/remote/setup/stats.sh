#!/bin/bash
set -euo pipefail  # Enhanced error handling: -e (exit on error), -u (treat unset variables as errors), -o pipefail (catch errors in pipelines)

# ======================================================================== #
# VARIABLES
# ======================================================================== #

LOG_FILE="/var/log/ufw.log"
AUTH_LOG="/var/log/auth.log"
OUTPUT_FILE="ufw_stats_output.txt"

# ======================================================================== #
# UFW Activity Statistics
# ======================================================================== #

# Check if required log files exist
check_log_files() {
    if [ ! -f "$LOG_FILE" ]; then
        echo "UFW log file not found at $LOG_FILE. Please ensure UFW logging is enabled."
        exit 1
    fi

    if [ ! -f "$AUTH_LOG" ]; then
        echo "Auth log file not found at $AUTH_LOG. Ensure SSH logs are enabled."
        exit 1
    fi
}

# Total Blocked Connections
display_blocked_connections() {
    echo "Blocked Connections: $(grep -c 'BLOCK' "$LOG_FILE" || echo "0")"
}

# Total Allowed Connections
display_allowed_connections() {
    echo "Allowed Connections: $(grep -c 'ALLOW' "$LOG_FILE" || echo "0")"
}

# Top 10 Blocked IPs
display_top_blocked_ips() {
    echo "Top 10 Blocked IPs:"
    grep 'BLOCK' "$LOG_FILE" | awk '{for(i=1;i<=NF;i++) if($i ~ /^SRC=/) print $i}' | cut -d'=' -f2 | sort | uniq -c | sort -nr | head || echo "No blocked IPs found."
}

# Top 10 Accessed Ports
display_top_accessed_ports() {
    echo "Top 10 Accessed Ports:"
    grep 'BLOCK' "$LOG_FILE" | awk '{for(i=1;i<=NF;i++) if($i ~ /^DPT=/) print $i}' | cut -d'=' -f2 | sort | uniq -c | sort -nr | head || echo "No accessed ports found."
}

# Top 10 Offending IPs
display_top_offending_ips() {
  echo "Top 10 Offending IPs:"
  grep 'BLOCK' "$LOG_FILE" | awk '{for(i=1;i<=NF;i++) if($i ~ /^SRC=/) print $i}' | cut -d'=' -f2 | sort | uniq -c | sort -nr | head || echo "No offending IPs found."
}

# Top 10 Targeted Ports
display_top_targeted_ports() {
  echo "Top 10 Targeted Ports:"
  grep 'BLOCK' "$LOG_FILE" | grep -oP 'DPT=\K[0-9]+' | sort | uniq -c | sort -nr | head || echo "No targeted ports found."
}

# Suspicious Small Packets (Likely Port Scans or Probes)
display_suspicious_packets() {
    echo "Suspicious Small Packets (Likely Probes):"
    grep 'BLOCK' "$LOG_FILE" | awk '{for(i=1;i<=NF;i++) if($i ~ /^LEN=/) print $i}' | cut -d'=' -f2 | sort -n | uniq -c | head || echo "No suspicious small packets found."
}

# SYN Flood Indicators (Multiple SYN Packets from Same IP)
display_syn_flood_indicators() {
    echo "Possible SYN Flood IPs (Top 5):"
    grep 'BLOCK' "$LOG_FILE" | grep 'SYN' | awk '{for(i=1;i<=NF;i++) if($i ~ /^SRC=/) print $i}' | cut -d'=' -f2 | sort | uniq -c | sort -nr | head -n 5 || echo "No SYN flood indicators found."
}

# ======================================================================== #
# SSH Activity Statistics
# ======================================================================== #

# Top 10 IPs for Failed SSH Login Attempts
display_failed_ssh_attempts() {
    echo "Top 10 Failed SSH Login Attempts (IPs):"
    grep -E "Failed password|invalid user|Connection closed.*preauth|Disconnected.*preauth" "$AUTH_LOG" | \
        awk '{for(i=1;i<=NF;i++) if($i ~ /[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+/) print $i}' | \
        sort | uniq -c | sort -nr | head || echo "No failed SSH logins found."
}

# Top 10 IPs for Successful SSH Logins
display_successful_ssh_attempts() {
    echo "Top 10 Successful SSH Logins (IPs):"
    grep 'Accepted' "$AUTH_LOG" | awk '{for(i=1;i<=NF;i++) if($i ~ /[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+/) print $i}' \
        | sort | uniq -c | sort -nr | head || echo "No successful SSH logins found."
}

# SSH Summary (Failed and Successful Attempts)
display_ssh_summary() {
    FAILED=$(grep -c 'Failed password' "$AUTH_LOG" || echo "0")
    SUCCESS=$(grep -c 'Accepted' "$AUTH_LOG" || echo "0")
    echo "Summary of SSH Attempts:"
    echo " - Failed SSH Attempts: $FAILED"
    echo " - Successful SSH Attempts: $SUCCESS"
}

# ======================================================================== #
# Monitoring
# ======================================================================== #

# Real-Time Monitoring Instructions
display_monitoring_tips() {
    echo "Tip: Monitor SSH logins using: tail -f /var/log/auth.log"
    echo "Tip: Monitor real-time traffic using: tail -f /var/log/ufw.log"
}

# Write Statistics to Output File
write_statistics() {
    exec > "$OUTPUT_FILE" 2>&1
    echo "==========================="
    echo " UFW and SSH Activity Statistics"
    echo "==========================="

    display_blocked_connections
    echo ""
    display_allowed_connections
    echo ""
    display_top_blocked_ips
    echo ""
    display_top_accessed_ports
    echo ""
    display_top_offending_ips
    echo ""
    display_top_targeted_ports
    echo ""
    display_suspicious_packets
    echo ""
    display_syn_flood_indicators
    echo ""
    echo "--------------- SSH section ---------------"
    display_failed_ssh_attempts
    echo ""
    display_successful_ssh_attempts
    echo ""
    display_ssh_summary
    echo ""

    echo "==========================="
}

# ======================================================================== #
# MAIN SCRIPT
# ======================================================================== #

main() {
    check_log_files
    write_statistics
    echo "Statistics have been written to ${OUTPUT_FILE}."
}

main "$@"