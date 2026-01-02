# Complete Deployment Guide

## Overview

This guide covers deployment of both the admin backend and static blog site.

## Part 1: Admin Backend Deployment

### Option A: Netlify Functions (Recommended)

The backend is built for Netlify Functions serverless deployment.

#### Setup

1. **Push to GitHub**
```bash
git add .
git commit -m "Ready for Netlify deployment"
git push origin main
```

2. **Connect to Netlify**
   - Visit [netlify.com](https://netlify.com)
   - Click "New site from Git"
   - Select GitHub repository
   - Configure build:
     - Build command: `go build ./cmd/functions/main.go`
     - Publish directory: `./public`
     - Functions directory: `./`

3. **Set Environment Variables**

In Netlify settings > Environment:

```
DATABASE_URL=libsql://your-db.turso.io?authToken=YOUR_TOKEN
ADMIN_PASSWORD=your-secure-password
JWT_SECRET=your-jwt-secret-key
ENV=production
```

4. **Database: Turso (Recommended)**

Turso provides SQLite as a service - perfect for this backend.

```bash
# Install turso CLI
brew install tursodatabase/tap/turso

# Create database
turso db create myblog

# Get connection string
turso db show --url myblog

# Get auth token
turso db tokens create myblog
```

Set `DATABASE_URL` to the connection string.

#### Deploy

Push to main branch - Netlify builds and deploys automatically.

Site URL: `https://your-site.netlify.app/api`

### Option B: Docker (Self-Hosted)

```dockerfile
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o cms ./cmd/functions/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/cms /usr/local/bin/
EXPOSE 8080
CMD ["cms"]
```

Deploy with:
```bash
docker build -t blog-cms .
docker run -e DATABASE_URL=... -e ADMIN_PASSWORD=... -p 8080:8080 blog-cms
```

### Option C: Railway, Render, Heroku

All support Go applications. Configuration is similar to Netlify.

## Part 2: Static Blog Deployment

### Option A: GitHub Pages (Recommended)

1. **Export blog to GitHub repository**

```bash
# Export markdown from admin panel
curl -X POST http://localhost:8080/api/exports

# or trigger from admin UI
# Navigate to Admin > Export > [Export to Markdown]
```

2. **Create GitHub repository for blog**

```bash
cd exports
git init
git add .
git commit -m "Initial blog export"
git remote add origin https://github.com/username/blog.git
git push -u origin main
```

3. **Add Hugo theme**

```bash
git submodule add https://github.com/theNewDynamic/gohugo-theme-ananke themes/ananke
# Update hugo.toml
```

4. **Enable GitHub Pages**
   - Go to repository Settings
   - Pages > Source: GitHub Actions
   - Workflow automatically runs on push

5. **Configure domain** (Optional)
   - Settings > Pages > Custom domain
   - Update `hugo.toml` baseURL

### Option B: Vercel

```bash
cd exports

# Login to Vercel
vercel login

# Deploy
vercel --prod
```

Vercel automatically detects Hugo and builds on each push.

### Option C: Netlify

```bash
cd exports
npm install -g netlify-cli
netlify deploy --prod
```

### Option D: Self-Hosted with GitHub Actions

Create your own deployment:

```yaml
# .github/workflows/deploy.yml
name: Deploy Blog

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Hugo
        uses: peaceiris/actions-hugo@v2
        with:
          hugo-version: 'latest'
      
      - name: Build
        run: hugo --minify
      
      - name: Deploy via SSH
        env:
          DEPLOY_KEY: ${{ secrets.DEPLOY_KEY }}
          DEPLOY_HOST: ${{ secrets.DEPLOY_HOST }}
          DEPLOY_USER: ${{ secrets.DEPLOY_USER }}
        run: |
          mkdir -p ~/.ssh
          echo "$DEPLOY_KEY" > ~/.ssh/id_ed25519
          chmod 600 ~/.ssh/id_ed25519
          ssh-keyscan -H $DEPLOY_HOST >> ~/.ssh/known_hosts
          rsync -avz --delete public/ $DEPLOY_USER@$DEPLOY_HOST:/var/www/blog/
```

## Part 3: Full Stack Deployment Example

### Complete Setup: Netlify + GitHub Pages

**Admin Backend**: Netlify Functions  
**Blog Frontend**: GitHub Pages  
**Database**: Turso

#### Step-by-Step

1. **Prepare Turso Database**

```bash
turso db create myblog
export DATABASE_URL=$(turso db show --url myblog)
export DB_TOKEN=$(turso db tokens create myblog)
```

2. **Deploy Backend to Netlify**

- Push main branch to GitHub
- Connect repo to Netlify
- Set environment variables:
  - `DATABASE_URL`
  - `ADMIN_PASSWORD`
  - `JWT_SECRET`
- Auto-deploys on push

3. **Access Admin Panel**

Navigate to: `https://your-site.netlify.app/`
- Dashboard loads
- Login with admin password
- Create/edit posts

4. **Export Blog**

From admin panel:
- Click "Export" in sidebar
- Click "Export to Markdown"
- Files generated in `./exports/`

5. **Deploy Blog**

```bash
cd exports
git init
git add .
git commit -m "Blog posts"
git remote add origin https://github.com/username/blog-site.git
git push -u origin main
```

6. **Enable GitHub Pages**

- Repo Settings > Pages
- Source: GitHub Actions
- Custom domain (optional)

7. **Set Up Auto-Export**

Create scheduled export:

```bash
# In your blog repository
curl -X POST https://your-site.netlify.app/api/exports \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

Use GitHub Actions cron to export daily:

```yaml
# .github/workflows/sync-blog.yml
name: Sync Blog

on:
  schedule:
    - cron: '0 0 * * *'  # Daily at midnight
  workflow_dispatch:

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Export blog
        run: |
          curl -X POST ${{ secrets.BLOG_API }}/api/exports \
            -H "Authorization: Bearer ${{ secrets.BLOG_TOKEN }}" \
            -o export.json
```

## Monitoring & Maintenance

### Health Checks

```bash
# Admin API health
curl https://your-site.netlify.app/api

# Blog site
curl https://blog-site.github.io/

# Post listing
curl https://your-site.netlify.app/api/posts
```

### Database Backups

For Turso:
```bash
turso db dump myblog > backup.sql
```

For self-hosted:
```bash
sqlite3 blog.db ".dump" > backup.sql
```

### Monitoring

Set up monitoring for:
- Netlify Functions runtime errors
- Database connection status
- GitHub Actions workflow failures
- Blog build times

## Troubleshooting

### Netlify Deploy Fails

1. Check build logs: Netlify UI > Deploys
2. Verify environment variables set
3. Test locally: `PORT=8080 ./main`
4. Check Go version compatibility

### GitHub Pages Not Building

1. Verify workflow exists: `.github/workflows/deploy.yml`
2. Check workflow logs
3. Ensure theme exists
4. Validate `hugo.toml` syntax

### Database Connection Error

1. Verify `DATABASE_URL` format
2. Test connection locally
3. Check Turso status page
4. Verify auth token validity

### Admin Not Loading

1. Check browser console for errors
2. Verify CSS loads: `https://site/css/admin.css`
3. Check CORS headers
4. Verify JWT token in cookies

## Production Checklist

- [ ] Environment variables set (admin, JWT secret)
- [ ] Database backed up
- [ ] Custom domain configured
- [ ] HTTPS enabled (auto on Netlify/GitHub Pages)
- [ ] Admin password changed from default
- [ ] Database connection tested
- [ ] Export workflow verified
- [ ] Blog theme selected and deployed
- [ ] Monitoring/alerts configured
- [ ] Backup schedule created

## Cost Estimate (2024 Pricing)

- Netlify Functions: Free tier (300,000 requests/month)
- Turso Database: Free tier (9GB storage, unlimited reads)
- GitHub Pages: Free
- Custom Domain: ~$12/year
- **Total: Free or minimal cost**

## Summary

Your blog system is now ready for production:

1. **Admin Backend** - Netlify Functions (serverless)
2. **Database** - Turso (SQLite in cloud)
3. **Static Blog** - GitHub Pages (free hosting)
4. **Deployment** - GitHub Actions (automated)
5. **Domain** - Custom domain (optional)

All components are designed for minimal cost, maximum reliability, and zero maintenance.
