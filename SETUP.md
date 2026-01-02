# Setup Guide

Complete setup for Blog CMS + Static Site with pure Go + HTMX backend and Netlify deployment.

## Development Setup

### 1. Prerequisites

- Go 1.23+
- Git
- curl/httpie (for testing)

### 2. Clone & Configure

```bash
git clone <repo>
cd blog
cp .env.example .env
```

Edit `.env` with your preferences:
```bash
# For development (simple password)
DATABASE_URL=file:./blog.db
ADMIN_PASSWORD=your-secure-password
JWT_SECRET=your-random-jwt-secret
ENV=development
PORT=8080
```

### 3. Initialize Database

First time only - creates tables and seeds post types:

```bash
go run cmd/cms/main.go
```

You should see output like:
```
✓ Database initialized successfully
✓ Schema created
✓ Found 12 post types:
  - Article (article)
  - Review (review)
  ...
```

### 4. Start Local Server

```bash
go build -o cms ./cmd/functions/main.go
./cms
```

Or with environment variables:
```bash
DATABASE_URL="file:./blog.db" ADMIN_PASSWORD="test" JWT_SECRET="secret" ./cms
```

Server runs on `http://localhost:8080`

### 5. Test the API

```bash
# Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"password":"your-password"}' | jq -r '.token')

# Create a post
curl -X POST http://localhost:8080/api/posts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type_id": "article",
    "title": "Hello World",
    "slug": "hello-world",
    "content": "# Hello\n\nThis is my first post!",
    "excerpt": "My first post",
    "tags": ["hello"],
    "status": "draft"
  }'

# List posts
curl http://localhost:8080/api/posts | jq .

# Get a single post
curl http://localhost:8080/api/posts/{post-id} | jq .
```

Full API documentation: see [API.md](./API.md)

## Production Setup (Netlify)

### Prerequisites

- GitHub repository
- Netlify account
- Turso database account

### 1. Create Turso Database

```bash
# Install Turso CLI if needed
curl -sSfL https://get.turso.io | bash

# Create database
turso db create my-blog-cms

# Get connection URL
turso db show my-blog-cms
# Copy the libsql:// URL

# Get auth token
turso auth tokens issue
```

### 2. Add Netlify Configuration

File: `netlify.toml` (already included)
```toml
[build]
  command = "mkdir -p netlify/functions && go build -o netlify/functions/cms ./cmd/functions/main.go"
  functions = "netlify/functions"
  publish = "public"

[[redirects]]
  from = "/api/*"
  to = "/.netlify/functions/cms/:splat"
  status = 200
```

### 3. Push to GitHub

```bash
git add .
git commit -m "Initial blog cms setup"
git push origin main
```

### 4. Deploy to Netlify

Option A: Through Netlify UI
1. Go to [netlify.com](https://netlify.com)
2. Click "New site from Git"
3. Select your repository
4. Set environment variables (see below)
5. Deploy

Option B: Via CLI
```bash
netlify login
netlify deploy --prod
```

### 5. Configure Environment Variables

In Netlify dashboard, go to `Site settings → Environment` and add:

```
DATABASE_URL=libsql://your-db-name-org.turso.io?authToken=your-token
ADMIN_PASSWORD=your-secure-password
JWT_SECRET=your-random-jwt-secret
ENV=production
```

### 6. Verify Deployment

```bash
# Test API endpoint
curl https://your-site.netlify.app/api/types | jq .

# Login
TOKEN=$(curl -s -X POST https://your-site.netlify.app/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"password":"your-password"}' | jq -r '.token')

# Create post
curl -X POST https://your-site.netlify.app/api/posts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{...}'
```

## Database Migration

### From Local SQLite to Turso

1. Export local data (if needed):
```bash
sqlite3 blog.db "SELECT * FROM posts;" > posts_backup.sql
```

2. Update `DATABASE_URL` in `.env`:
```bash
# Before:
DATABASE_URL=file:./blog.db

# After:
DATABASE_URL=libsql://[your-turso-db].turso.io?authToken=[token]
```

3. Initialize Turso database with schema:
```bash
go run cmd/cms/main.go
```

4. Migrate data (if needed - use your backup SQL)

## Troubleshooting

### "database file is locked"
Multiple processes accessing the same SQLite file. Close other connections or restart.

### "no such table: posts"
Run database initialization:
```bash
go run cmd/cms/main.go
```

### "Invalid credentials"
Check:
- `ADMIN_PASSWORD` env variable is set
- Password matches what you're sending
- Or use `PASSWORD_HASH` with bcrypt hash instead

### Netlify deployment fails
Check:
- Go version: `netlify env:list`
- Ensure `DATABASE_URL` is set in Netlify environment
- Check build logs in Netlify dashboard
- Verify Turso token is valid: `turso auth tokens list`

### Posts not appearing after creating
Check post `status` field. Only `published` posts show in public queries by default:
```bash
# See all posts including drafts
curl http://localhost:8080/api/posts?status=draft
```

## Next Steps

- [ ] Build HTMX frontend (admin dashboard, editor)
- [ ] Add link preview extraction
- [ ] Create Hugo export job
- [ ] Set up automatic daily sync to static site
- [ ] Add full-text search
- [ ] Build public site with filtering/tagging
- [ ] Add analytics

## File Structure

```
.
├── cmd/
│   ├── cms/main.go                  # DB initialization
│   └── functions/main.go            # API server (Netlify)
├── internal/
│   ├── db/                          # Database layer
│   ├── models/                      # Data structures
│   ├── handler/                     # API handlers
│   ├── editor/                      # [coming: markdown utils]
│   └── util/                        # Helpers
├── netlify.toml                     # Netlify config
├── .env.example                     # Configuration template
├── go.mod / go.sum                  # Go dependencies
├── API.md                           # API documentation
├── ARCHITECTURE.md                  # System design
├── README.md                        # Overview
└── SETUP.md                         # This file
```

## Support

See [API.md](./API.md) for complete endpoint reference.

See [ARCHITECTURE.md](./ARCHITECTURE.md) for system design details.
