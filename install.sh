#!/bin/bash

# Check if the ".devops" directory already exists in the home directory
if [ ! -d "$HOME/.devops" ]; then
  # If it doesn't exist, create it
  mkdir "$HOME/.devops"
fi

# Download the tar.gz file
wget -O "$HOME/.devops/devops.tar.gz" "https://storage.googleapis.com/devops-cli-artifacts/releases/devops/0.1.0/devops_0.1.0_Linux_x86_64.tar.gz"

# Extract the tar.gz file
tar -xzf "$HOME/.devops/devops.tar.gz" -C "$HOME/.devops"

# Add the binary to the system path
echo 'export PATH=$PATH:$HOME/.devops' >> ~/.bashrc
source ~/.bashrc

rm "$HOME/.devops/devops.tar.gz"
"$HOME/.devops"

# Write a bash script to the following
# 1) create a ".devops" directory inside the home directory of user. If it doesn't already exists
# 2) Download a tar.gz file from this link "https://github.com/sharadregoti/devops/releases/download/v0.1.0/devops_0.1.0_Linux_x86_64.tar.gz"
# 3) Untar the above downloaded file which contains an executable binary.
# 4) Add the binary to system path
# 5) Delete the downloaded file & extracted content
