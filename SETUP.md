# Step-by-Step Setup Guide

## Prerequisites

- Node.js 18+ (check: `node --version`)
- npm (comes with Node.js)
- Git (for GitHub Actions)
- Turso account (free tier at turso.io)
- GitHub account (for Actions)

## Step 1: Create Turso Database

```bash
# Install Turso CLI
curl -sSfL https://get.turso.io | bash
export PATH="$PATH:$HOME/.turso"

# Login (opens browser)
turso auth login

# Create database
turso db create blog-db

# Get connection details
turso db show blog-db --json

# You'll need:
# - URL (libsql://...)
# - TOKEN (from turso db tokens create blog-db)

turso db tokens create blog-db
```

Save the connection URL and token somewhere safe.

---

## Step 2: Setup Backend

```bash
cd backend

# Create .env file
cat > .env << 'EOF'
# Database
TURSO_CONNECTION_URL=libsql://your-db.turso.io
TURSO_AUTH_TOKEN=your-token-here

# Auth
ADMIN_PASSWORD=your-strong-password-here
JWT_SECRET=your-random-jwt-secret-here

# Server
API_PORT=3000
API_URL=http://localhost:3000

# For dev, optional local SQLite:
DATABASE_URL=file:./dev.db
EOF
```

Replace the placeholder values with your actual Turso details.

### Generate Password Hash

```bash
# Install dependencies first
npm install

# Generate bcrypt hash of your password
node << 'EOF'
const bcrypt = require('bcryptjs');
const password = 'your-strong-password-here'; // Use same as ADMIN_PASSWORD

bcrypt.hash(password, 10).then(hash => {
  console.log('Add this to .env as PASSWORD_HASH:');
  console.log(`PASSWORD_HASH=${hash}`);
});
EOF
```

Copy the hash and add it to `.env`:
```
PASSWORD_HASH=your-generated-hash-here
```

### Initialize Database

```bash
# Create and populate database schema
npm run db:init

# You should see: "Database initialized successfully"
```

### Start Backend Server

```bash
npm run dev

# You should see: "Server running on http://localhost:3000"
```

Keep this terminal open.

---

## Step 3: Setup Frontend

Open new terminal:

```bash
cd frontend

# Start a simple HTTP server
npx http-server -p 8080

# Open browser: http://localhost:8080/login.html
```

### First Login

1. Open http://localhost:8080/login.html
2. Enter your `ADMIN_PASSWORD`
3. Click "Sign In"
4. You should be redirected to editor.html

If it fails, check:
- Is backend running on port 3000?
- Is password correct?
- Check browser console (F12) for errors

---

## Step 4: Create Your First Post

1. Stay on editor.html
2. Fill in form:
   - **Title**: "Hello World"
   - **Type**: Article
   - **Content**: "My first blog post"
   - Add a tag (e.g., "intro")
3. Click "Save as Draft"
4. Or click "Publish" directly

Check backend console - you should see database queries.

---

## Step 5: Verify Database

Terminal:
```bash
# Connect to Turso database
turso db shell blog-db

# Check posts table
SELECT * FROM posts;

# Exit with Ctrl+D
```

You should see your post in the database.

---

## Step 6: Test API Directly (Optional)

```bash
# Get auth token
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"password":"your-password"}'

# Response: {"token":"eyJ...","expires_in":604800}

# Save token
TOKEN="eyJ..."

# Get all posts
curl http://localhost:3000/api/posts \
  -H "Authorization: Bearer $TOKEN"

# Create post
curl -X POST http://localhost:3000/api/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "type_id":"thought",
    "title":"Test Post",
    "content":"Just testing the API",
    "tags":["test"]
  }'
```

---

## Step 7: GitHub Actions Setup (Optional for Now)

Skip if you just want local setup. Setup later when ready to deploy.

```bash
# Push to GitHub first
git init
git add .
git commit -m "Initial commit"
git remote add origin https://github.com/your-username/blog
git push -u origin main

# Add secrets to GitHub
gh secret set TURSO_CONNECTION_URL -b "libsql://your-db.turso.io"
gh secret set TURSO_AUTH_TOKEN -b "your-token-here"

# Trigger workflow manually or wait for schedule
gh workflow run sync-posts.yml
```

---

## Step 8: Hugo Setup (Optional)

```bash
# Install Hugo
# macOS: brew install hugo
# Linux: sudo apt install hugo
# Windows: choco install hugo-extended

# Test Hugo
hugo version

# Initialize if not exists
hugo new site hugo

# Add a theme
cd hugo
git init
git submodule add https://github.com/theNewDynamic/gohugo-theme-ananke.git themes/ananke
echo "theme = 'ananke'" >> hugo.toml

# Test build
hugo -d ../public

# Serve locally
hugo server -D
```

---

## Troubleshooting

### "Connection refused" on login
- Check backend is running: `npm run dev`
- Verify port 3000 is not in use: `lsof -i :3000`
- Check firewall

### "Invalid credentials"
- Verify password matches `ADMIN_PASSWORD` in `.env`
- Verify `PASSWORD_HASH` is set correctly
- Regenerate hash if unsure

### "Database error"
- Check `TURSO_CONNECTION_URL` is correct
- Verify `TURSO_AUTH_TOKEN` is valid
- Try: `turso db shell blog-db` to test connection
- Check `.env` file exists in backend folder

### Posts not saving
- Check backend console for errors
- Open browser DevTools (F12) → Console
- Check network tab for 401/500 errors
- Verify JWT token is being sent

### "CORS errors"
- Backend should handle CORS automatically with Hono
- If issues, check request headers in DevTools

### Can't run `npm run db:init`
- Check you're in `backend/` folder
- Run `npm install` first
- Check Node.js version: `node --version` (should be 18+)

---

## Next Steps

Once setup works:

1. **Create more posts** - Try different post types in editor
2. **Add Hugo theme** - Customize `hugo/` for static site
3. **Configure GitHub Actions** - Setup cronjob and auto-deploy
4. **Deploy backend** - Push to Vercel, Railway, or Fly.io
5. **Deploy static site** - Setup GitHub Pages or Netlify

---

## File Structure After Setup

```
blog/author/
├── .env                    ← Created in step 2
├── backend/
│   ├── .env               ← Your secrets
│   ├── dev.db             ← Local SQLite (optional)
│   ├── node_modules/      ← npm packages
│   ├── src/
│   │   ├── db/
│   │   │   └── schema.sql
│   │   ├── api/
│   │   ├── types.ts
│   │   └── index.ts       ← Main server
│   └── package.json
├── frontend/
│   ├── login.html
│   ├── editor.html
│   ├── editor.js
│   └── login.js
├── hugo/                  ← Add after step 8
│   ├── config.toml
│   ├── content/           ← Auto-populated by cronjob
│   └── themes/
└── cronjob/
    └── export.ts

```

---

## Commands Reference

```bash
# Backend
cd backend && npm run dev        # Start dev server
cd backend && npm run db:init    # Initialize database
cd backend && npm run build      # Build for production

# Frontend
cd frontend && npx http-server   # Serve on localhost:8080

# Turso
turso db show blog-db            # Show database info
turso db shell blog-db           # Open database shell
turso db tokens create blog-db   # Create new token
turso db delete blog-db          # Delete database

# GitHub
git push origin main             # Push to GitHub
gh workflow run sync-posts.yml   # Manually trigger action
gh secret set KEY VALUE          # Set repository secret
```

---

Ready? Start with **Step 1** below.
