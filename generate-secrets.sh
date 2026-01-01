#!/bin/bash

# Color codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}Blog System - Secret Generator${NC}"
echo ""

# Generate JWT secret
JWT_SECRET=$(openssl rand -base64 32)
echo -e "${GREEN}JWT_SECRET:${NC}"
echo "$JWT_SECRET"
echo ""

# Password hashing
echo -e "${BLUE}To generate PASSWORD_HASH:${NC}"
echo ""
echo "1. Set a strong password (min 12 chars recommended)"
echo ""
echo "2. Run this command (replace 'your-password' with your actual password):"
echo ""
echo -e "${YELLOW}node -e \"const bcrypt = require('bcryptjs'); bcrypt.hash('your-password', 10).then(h => console.log(h))\"${NC}"
echo ""
echo "3. Copy the output hash"
echo ""

# Summary
echo -e "${BLUE}To complete setup:${NC}"
echo ""
echo "1. Create Turso database (if you haven't):"
echo "   ${YELLOW}turso auth login${NC}"
echo "   ${YELLOW}turso db create blog-db${NC}"
echo "   ${YELLOW}turso db tokens create blog-db${NC}"
echo ""
echo "2. Edit ${YELLOW}backend/.env${NC} with:"
echo ""
echo "   TURSO_CONNECTION_URL=<from turso db show blog-db>"
echo "   TURSO_AUTH_TOKEN=<from turso db tokens create>"
echo "   ADMIN_PASSWORD=<your-strong-password>"
echo "   PASSWORD_HASH=<from bcrypt command above>"
echo "   JWT_SECRET=$JWT_SECRET"
echo ""
echo -e "${GREEN}Keep this terminal output for reference!${NC}"
