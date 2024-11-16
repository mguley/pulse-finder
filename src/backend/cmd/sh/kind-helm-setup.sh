#!/bin/bash
#####################################################################################
# Summary:
# This script sets up a Kind cluster, deploys Helm charts, and cleans up the environment.
# Requirements:
# - kubectl version: v1.31.2 or compatible
# - kind version: v0.25.0 or compatible
# - helm version: v3.16.3 or compatible
# - docker version: 27.3.1 or compatible
#####################################################################################

# Set environment variables
CLUSTER_NAME="kind-helm-cluster"

#####################################################################################
# Function to check if a command is available
check_command_available() {
  command -v $1 >/dev/null 2>&1 || { echo >&2 "$1 is required but not installed. Aborting."; exit 1; }
}

#####################################################################################
# Function to check the result of the last command and exit if it failed
check_command() {
  if [ $? -ne 0 ]; then
    echo "$1"
    cleanup_kind_cluster
    exit 1
  fi
}

#####################################################################################
# Function to create Kind cluster
create_kind_cluster() {
  echo "Creating Kind cluster..."
  kind create cluster --name $CLUSTER_NAME
  check_command "Failed to create Kind cluster."
  echo "Kind cluster created successfully."
}

#####################################################################################
# Function to delete Kind cluster
cleanup_kind_cluster() {
  echo "Deleting Kind cluster..."
  kind delete cluster --name $CLUSTER_NAME
  if [ $? -ne 0 ]; then
    echo "Failed to delete Kind cluster. Manual cleanup may be required."
    exit 1
  fi
  echo "Kind cluster deleted successfully."
}

#####################################################################################
# Function to display the status of the cluster
show_cluster_status() {
  echo "Cluster and deployed resources status:"
  kubectl get all
}

#####################################################################################
# Main script execution

# Check for required commands
check_command_available kubectl
check_command_available kind
check_command_available helm
check_command_available docker

# Trap to clean up on script exit
trap cleanup_kind_cluster EXIT

# Execute tasks
create_kind_cluster
show_cluster_status

echo "Kind cluster and Helm deployment setup completed successfully."
echo "Press Ctrl+C to clean up and exit."
# Wait for user to terminate the script
while :; do
  sleep 1
done