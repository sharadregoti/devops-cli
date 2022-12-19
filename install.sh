#!/bin/bash

# Check if the ".devops" directory already exists in the home directory
if [ ! -d "$HOME/.devops" ]; then
  # If it doesn't exist, create it
  mkdir "$HOME/.devops"
fi

# Check if the ".devops/plugins" directory already exists in the home directory
if [ ! -d "$HOME/.devops/plugins" ]; then
  # If it doesn't exist, create it
  mkdir "$HOME/.devops/plugins"
fi

# Download the tar.gz file
wget -O "$HOME/.devops/devops.tar.gz" "https://storage.googleapis.com/devops-cli-artifacts/releases/devops/0.1.0/devops_0.1.0_Linux_x86_64.tar.gz"

# Extract the tar.gz file
tar -xzf "$HOME/.devops/devops.tar.gz" -C "$HOME/.devops"

# Check if the ".devops/plugins" directory already exists in the home directory
if [ ! -d "$HOME/.devops/plugins/kubernetes" ]; then
  # If it doesn't exist, create it
  mkdir "$HOME/.devops/plugins/kubernetes"
fi
# Binaries
wget -O "$HOME/.devops/plugins/kubernetes/devops.tar.gz" "https://storage.googleapis.com/devops-cli-artifacts/releases/devops/0.1.0/devops-kubernetes-plugin_0.1.0_Linux_x86_64.tar.gz"
tar -xzf "$HOME/.devops/plugins/kubernetes/devops.tar.gz" -C "$HOME/.devops/plugins/kubernetes"


# Add the binary to the system path
echo 'export PATH=$PATH:$HOME/.devops' >> ~/.bashrc
source ~/.bashrc

rm "$HOME/.devops/devops.tar.gz"
rm "$HOME/.devops/plugins/kubernetes/devops.tar.gz"