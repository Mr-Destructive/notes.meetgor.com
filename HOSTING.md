# Hosting & Deployment Guide

Complete guide on how to host the Blog CMS and the data flow.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        YOUR USERS                               │
│                  (Browser/Mobile Client)                        │
└──────────────────────────────┬──────────────────────────────────┘
                               │ HTTPS
                               ↓
┌──────────────────────────────────────────────────────────────────┐
│                     NETLIFY CDN EDGE                             │
│              (Global Content Delivery Network)                   │
└──────────────────────────────┬──────────────────────────────────┘
                               │
                               ↓
┌──────────────────────────────────────────────────────────────────┐
│                  NETLIFY FUNCTIONS (USA Region)                 │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │        Go Binary (cms) - 12MB                            │ │
│  │                                                            │ │
│  │  Handler: routes /api/* requests                          │ │
│  │  - Parse route                                            │ │
│  │  - Validate auth                                          │ │
│  │  - Call database                                          │ │
│  │  - Return JSON                                            │ │
│  └────────────────────────────────────────────────────────────┘ │
│                               │                                  │
│                  (Cold start: ~1-2 seconds)                    │
│                  (Warm requests: <100ms)                       │
└──────────────────────────────┬──────────────────────────────────┘
                               │ HTTPS
                               ↓
┌──────────────────────────────────────────────────────────────────┐
│                      TURSO DATABASE                              │
│              (Managed SQLite - Global Replication)              │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │              SQLite Database                            │ │
│  │                                                            │ │
│  │  - posts (blog posts)                                     │ │
│  │  - post_types (12 types)                                 │ │
│  │  - series (collections)                                  │ │
│  │  - post_series (many-to-many)                            │ │
│  │  - revisions (version history)                           │ │
│  │  - settings (config)                                     │ │
│  │                                                            │ │
│  │  Storage: Up to 9GB free tier                            │ │
│  │  Queries: <10ms (replicated globally)                    │ │
│  └────────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────┘
```

## Data Flow

### 1. Create a Blog Post

```
User (Browser)
    ↓
POST /api/posts
    ↓
    ├─ Parse JSON body
    ├─ Validate (non-empty title, slug)
    ├─ Generate unique ID
    └─ Insert into database
         ↓
    Turso SQLite
         ↓
    Return JSON (post object with ID)
         ↓
    Browser shows post created
```

**Time**: ~10-50ms (database insert)

### 2. List Published Posts

```
Browser
    ↓
GET /api/posts?status=published
    ↓
    ├─ Parse query parameters
    ├─ Count matching posts (for pagination)
    └─ Query database with filters
         ↓
    SELECT * FROM posts 
    WHERE status = 'published'
    ORDER BY published_at DESC
         ↓
    Return: JSON array + total count
         ↓
    Browser renders list
```

**Time**: ~5-15ms (read-only query)

### 3. Update Blog Post

```
Browser (authenticated)
    ↓
PUT /api/posts/{id}
    ↓
    ├─ Verify JWT token
    ├─ Fetch current post
    ├─ Create revision (backup old content)
    ├─ Update fields
    └─ Save to database
         ↓
    Turso SQLite
         ↓
    Return: Updated post object
         ↓
    Browser shows "Saved"
```

**Time**: ~15-30ms (2 queries)

### 4. Delete Blog Post

```
Browser (authenticated)
    ↓
DELETE /api/posts/{id}
    ↓
    ├─ Verify JWT token
    ├─ Delete from posts table
    ├─ Cascade delete revisions
    └─ Commit transaction
         ↓
    Turso SQLite
         ↓
    Return: 200 OK
```

**Time**: ~10ms

### 5. Create Series/Collection

```
Browser (authenticated)
    ↓
POST /api/series
    ↓
    ├─ Parse JSON (name, slug)
    ├─ Generate ID
    └─ Insert series
         ↓
    Turso SQLite
         ↓
    Return: Series object
```

**Time**: ~10ms

### 6. Add Post to Series

```
Browser
    ↓
POST /api/series/{id}/add-post
    ↓
    ├─ Verify series exists
    ├─ Verify post exists
    └─ Insert into post_series table
         ↓
    Turso SQLite
         ↓
    Return: Updated mapping
```

**Time**: ~15ms

---

## Deployment Options

### Option 1: Netlify (Recommended for this project)

#### Pros
✓ Netlify Functions = serverless Go binary  
✓ Auto-deploys on git push  
✓ Global edge network for fast delivery  
✓ Automatic HTTPS  
✓ No servers to manage  
✓ Free tier: 100 requests/minute  

#### Cons
✗ Cold starts: 1-2 seconds (not ideal for high traffic)  
✗ Regional latency if user far from deployment region  

#### Cost
- **Free**: 100 req/min, 125K function invocations/month
- **Pro**: $19/month = unlimited (recommended at scale)
- **Turso DB**: $0-29/month depending on usage

#### Setup
```bash
# 1. Push to GitHub
git push origin main

# 2. In Netlify UI:
# - Connect GitHub repo
# - Set env vars
# - Auto-deploys on push

# 3. Test
curl https://your-site.netlify.app/api
```

---

### Option 2: Railway.app (Good alternative)

#### Pros
✓ Better cold start times (~500ms)  
✓ Persistent containers  
✓ Good for Go apps  
✓ Simple deployment  

#### Cons
✗ Different deployment mechanism  
✗ Always-running container (not truly serverless)  

#### Cost
- First $5/month free
- Then pay-as-you-go ($0.0000231/second)
- ~$10-15/month typical

#### Deployment
```bash
# Install Railway CLI
npm install -g @railway/cli

# Login
railway login

# Deploy
railway up
```

---

### Option 3: Fly.io (High performance)

#### Pros
✓ Near-instant cold starts  
✓ Geographic distribution  
✓ Great Go support  
✓ Good uptime  

#### Cons
✗ Different platform philosophy  
✗ Requires some DevOps knowledge  

#### Cost
- First 3 shared-cpu-1x VMs: Free
- Standard: $0.0000231/second per VM
- ~$5-10/month typical

#### Deployment
```bash
flyctl launch
flyctl deploy
```

---

### Option 4: Self-Hosted (Full Control)

Deploy the Go binary to your own VPS.

#### Pros
✓ Full control  
✓ No vendor lock-in  
✓ Potentially cheaper at scale  

#### Cons
✗ You manage uptime  
✗ You handle security patches  
✗ You scale manually  

#### Cost
- **Linode/DigitalOcean**: $5-10/month
- **VPS**: Generally $3-30/month

#### Setup
```bash
# 1. Build binary
go build -o cms ./cmd/functions/main.go

# 2. SSH to VPS
ssh user@your-vps.com

# 3. Upload binary & .env
scp cms user@your-vps.com:/home/app/
scp .env user@your-vps.com:/home/app/

# 4. Run with process manager
# Using systemd:
cat > /etc/systemd/system/cms.service << 'EOF'
[Unit]
Description=Blog CMS
After=network.target

[Service]
Type=simple
User=app
WorkingDirectory=/home/app
ExecStart=/home/app/cms
Restart=always
Environment="DATABASE_URL=..."
Environment="ADMIN_PASSWORD=..."
Environment="JWT_SECRET=..."

[Install]
WantedBy=multi-user.target
EOF

systemctl start cms
systemctl enable cms

# 5. Use nginx as reverse proxy
# Listen 80/443 → Forward to localhost:8080
```

---

## Recommended Setup: Netlify + Turso

### Architecture
```
GitHub Repo
    ↓
Connected to Netlify
    ↓
On push:
    1. Netlify runs: go build -o netlify/functions/cms ./cmd/functions/main.go
    2. Creates serverless function
    3. Deploys globally
    ↓
Turso Database
    ↓
Available at: https://your-site.netlify.app
```

### Setup Steps

#### 1. Create Turso Database
```bash
# Install Turso CLI
curl -sSfL https://get.turso.io | bash

# Create database
turso db create blog-cms

# Get URL
turso db show blog-cms
# Output: libsql://blog-cms-xyz.turso.io

# Create auth token
turso auth tokens issue
# Output: long token string
```

#### 2. Prepare GitHub
```bash
git add .
git commit -m "Ready for Netlify deployment"
git push origin main
```

#### 3. Deploy on Netlify
```bash
# Visit: https://app.netlify.com/start
# - Connect GitHub repo
# - Auto-detects netlify.toml
# - Asks for env variables:

DATABASE_URL = libsql://blog-cms-xyz.turso.io?authToken=token
ADMIN_PASSWORD = your-secure-password
JWT_SECRET = random-secret-key
ENV = production

# Click Deploy
# Netlify auto-builds and deploys
```

#### 4. Test Deployment
```bash
# API health
curl https://your-site.netlify.app/api
# {"status":"ok","version":"1.0"}

# Login
curl -X POST https://your-site.netlify.app/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"password":"your-password"}'

# Create post
curl -X POST https://your-site.netlify.app/api/posts \
  -H "Authorization: Bearer $TOKEN" \
  -d '{...}'
```

---

## Performance & Optimization

### Request Timeline (Warm)
```
User Request
    ├─ Netlify routing: 1-5ms
    ├─ Go handler: 1-10ms
    ├─ Database query: 5-50ms (depending on complexity)
    └─ Response serialization: <1ms
    ─────────────────────────
    Total: 10-70ms (typical)
```

### Request Timeline (Cold Start)
```
First request after deployment:
    ├─ Boot Go runtime: 500-1500ms
    ├─ Establish DB connection: 100-500ms
    ├─ Handle request: 10-20ms
    └─ Return response: 1-2ms
    ─────────────────────────
    Total: 1-2 seconds
```

### Optimization Tips

**1. Minimize Cold Starts**
- Use persistent connections (connection pool)
- Lazy-load heavy dependencies
- Monitor with: Netlify Analytics

**2. Speed Up Queries**
- Add indexes (already done in schema)
- Use pagination (limit + offset)
- Cache frequently accessed data

**3. Reduce Latency**
- Use Turso geographic replication
- Serve static files from CDN
- Compress responses (gzip)

**4. Monitor Performance**
- Netlify Dashboard → Analytics
- Turso Dashboard → Query stats
- Enable request logging

---

## Database Strategy

### Development
```
Local SQLite:
├─ File: ./blog.db
├─ No setup needed
├─ Perfect for testing
└─ Fast development cycle
```

### Production
```
Turso Database:
├─ Managed SQLite
├─ Global replication
├─ Automatic backups
├─ HTTPS encrypted
└─ 9GB free tier
```

### Backup Strategy
```
Daily exports:
├─ GitHub Actions runs daily
├─ Exports all posts as JSON
├─ Commits to separate branch
├─ Can restore from git history
```

---

## Environment Variables

### Required
```
DATABASE_URL=libsql://[db]-[org].turso.io?authToken=[token]
ADMIN_PASSWORD=your-secure-password
JWT_SECRET=randomly-generated-secret-key-32chars
```

### Optional
```
PORT=8080                    (Netlify ignores this)
ENV=production|development
```

### How to Secure
✓ Never commit .env file  
✓ Use Netlify environment variables  
✓ Rotate JWT_SECRET periodically  
✓ Use strong ADMIN_PASSWORD (20+ chars)  
✓ Use Netlify Secrets UI (not public)  

---

## Troubleshooting Deployment

### "Database connection error"
- [ ] Check DATABASE_URL is correct
- [ ] Verify auth token hasn't expired
- [ ] Test locally: `DATABASE_URL=... ./cms`

### "Function timeout"
- [ ] Check database query performance
- [ ] Verify network latency to database
- [ ] Add pagination to list endpoints

### "Cold start too slow"
- This is normal on Netlify
- Upgrade to Fly.io or Railway for faster starts
- Or use Railway/Fly.io instead

### "403 Unauthorized"
- [ ] Check ADMIN_PASSWORD matches
- [ ] Check JWT_SECRET hasn't changed
- [ ] Verify token hasn't expired (7 days)

### "Deploy fails"
- [ ] Run `go build` locally to verify
- [ ] Check `go mod tidy` - update dependencies
- [ ] Review Netlify build logs

---

## Scaling

### Current Limits
- **Requests/min**: 100 (free tier) → unlimited (pro)
- **Database size**: 9GB (free) → unlimited
- **Concurrent connections**: ~100
- **Request timeout**: 30 seconds

### At 1,000 requests/month
- Cost: Free tier sufficient
- Database: No issues

### At 100,000 requests/month
- Netlify: Need Pro ($19/month)
- Database: Still free tier
- Total cost: ~$20-30/month

### At 10M requests/month
- Netlify: ~$50-100/month
- Database: ~$50-100/month
- Total: $100-200/month
- _Consider migrating to Railway/Fly.io_

---

## Monitoring & Logging

### Netlify Dashboard
```
https://app.netlify.com → Select site
├─ Analytics: Request counts, performance
├─ Deploys: Deployment history
├─ Functions: Execution logs
└─ Logs: Real-time request logs
```

### Turso Dashboard
```
https://app.turso.io
├─ Query stats: Which queries are slow
├─ Storage: Data usage
├─ Billing: Current usage
└─ Backups: Automatic daily backups
```

### Manual Logging
Add to any handler:
```go
log.Printf("DEBUG: query took %d ms", elapsed)
```

Netlify shows logs in: Netlify Dashboard → Functions → select function → logs

---

## Security Checklist

- [ ] ADMIN_PASSWORD is strong (20+ chars, mixed case, symbols)
- [ ] JWT_SECRET is random (run: `openssl rand -hex 32`)
- [ ] .env file is in .gitignore
- [ ] Netlify env vars are private (not shared)
- [ ] HTTPS enforced (Netlify does automatically)
- [ ] Turso token rotated every 90 days
- [ ] Backups enabled (automatic in Turso)
- [ ] Regular security updates (Go 1.23+)

---

## Recommended Path

1. **Start**: Deploy to Netlify (free, easy)
2. **Grow**: Monitor performance
3. **Scale**: Upgrade Netlify plan if needed
4. **Optimize**: Migrate to Railway/Fly.io if cold starts matter

This setup is production-ready for 99% of use cases.

---

## Questions?

- Netlify Docs: https://docs.netlify.com/functions/overview
- Turso Docs: https://docs.turso.io
- Go Deployment: https://golang.org/doc/deploy.html
