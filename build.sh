#!/bin/bash
set -e

# Build Hugo site
hugo --minify

# Copy public to output
cp -r public/* .
