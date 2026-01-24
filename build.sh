#!/bin/bash
set -e

# Install Hugo
echo "Installing Hugo..."
curl -L https://github.com/gohugoio/hugo/releases/download/v0.136.5/hugo_0.136.5_linux-amd64.tar.gz | tar xz

# Build
echo "Building site..."
./hugo --minify

# Copy output
echo "Deploying..."
cp -r public/* .
