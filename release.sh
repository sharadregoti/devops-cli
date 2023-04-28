#!/bin/bash

# echo "Did you uncomment the production URL (y/n)?"
# read answer

# if [ "$answer" != "${answer#[Nn]}" ] ;then
#     echo "Exitting..."
#     exit 0
# fi

# echo "Did increment version number in code (y/n)?"
# read answer

# if [ "$answer" != "${answer#[Nn]}" ] ;then
#     echo "Exitting..."
#     exit 0
# fi

# echo "Did create a new git tag as per current version number (y/n)?"
# read answer

# if [ "$answer" != "${answer#[Nn]}" ] ;then
#     echo "Exitting..."
#     exit 0
# fi

echo "Building frontend..."
cd devops-frontend
# Build frontend
vite build
# Create tar file
tar -czvf ui.tar.gz dist/

# Build Core Binary
echo "Building core binary..."
cd ../
goreleaser release --clean --skip-validate --snapshot

# Build kubernetes plugins
cd plugins/kubernetes
goreleaser release --clean --skip-validate

cd ../

# Build kubernetes plugins
cd helm
goreleaser release --clean --skip-validate
