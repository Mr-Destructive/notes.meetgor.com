#!/bin/bash

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     Blog System - Quick Setup         ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Check Node.js
if ! command -v node &> /dev/null; then
    echo -e "${RED}✗ Node.js is not installed${NC}"
    echo "  Download from: https://nodejs.org/"
    exit 1
fi

echo -e "${GREEN}✓ Node.js $(node --version)${NC}"

# Check npm
if ! command -v npm &> /dev/null; then
    echo -e "${RED}✗ npm is not installed${NC}"
    exit 1
fi

echo -e "${GREEN}✓ npm $(npm --version)${NC}"
echo ""

# Setup Backend
echo -e "${BLUE}Setting up backend...${NC}"
cd backend || exit 1

if [ ! -f .env ]; then
    echo -e "${YELLOW}Creating .env file...${NC}"
    cp .env.example .env
    echo -e "${YELLOW}⚠ Edit backend/.env with your Turso credentials:${NC}"
    echo "  - TURSO_CONNECTION_URL"
    echo "  - TURSO_AUTH_TOKEN"
    echo "  - ADMIN_PASSWORD"
    echo "  - JWT_SECRET (generate: openssl rand -base64 32)"
    echo ""
    read -p "Press Enter when done..."
fi

if [ ! -d node_modules ]; then
    echo -e "${YELLOW}Installing dependencies...${NC}"
    npm install
fi

echo -e "${GREEN}✓ Backend ready${NC}"
cd ..

# Summary
echo ""
echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║         Setup Complete!              ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo ""
echo "1. ${BLUE}Generate password hash:${NC}"
echo "   cd backend && node -e \"const bcrypt = require('bcryptjs'); bcrypt.hash('your-password', 10).then(h => console.log('PASSWORD_HASH=' + h))\""
echo ""
echo "2. ${BLUE}Add hash to backend/.env${NC}"
echo ""
echo "3. ${BLUE}Initialize database:${NC}"
echo "   cd backend && npm run db:init"
echo ""
echo "4. ${BLUE}Start backend (Terminal 1):${NC}"
echo "   cd backend && npm run dev"
echo ""
echo "5. ${BLUE}Start frontend (Terminal 2):${NC}"
echo "   cd frontend && npx http-server -p 8080"
echo ""
echo "6. ${BLUE}Open browser:${NC}"
echo "   http://localhost:8080/login.html"
echo ""
echo -e "${GREEN}Happy blogging!${NC}"
