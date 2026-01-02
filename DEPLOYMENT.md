# Deployment Guide

Blog CMS is fully serverless-ready with Netlify Functions + Turso SQLite.

## Architecture

```
Client Browser
    ↓
Netlify CDN
    ↓
Netlify Functions (Go binary)
    ↓
Turso SQLite (remote database)
    ↓
Admin Panel + API
```

## Netlify Deployment

### Quick Deploy (Recommended)

1. **Prepare GitHub Repository**
   ```bash
   git add .
   git commit -m "Ready for deployment"
   git push origin main
   ```

2. **Connect to Netlify**
   - Go to [netlify.com](https://netlify.com)
   - Click "New site from Git"
   - Select GitHub repository
   - Build settings auto-detect from `netlify.toml`

3. **Set Environment Variables**
   
   In Netlify dashboard: `Site settings → Environment`:
   
   ```
   DATABASE_URL = libsql://your-db-name-your-org.turso.io?authToken=your-token
   ADMIN_PASSWORD = your-secure-password-here
   JWT_SECRET = your-random-secret-key-here
   ENV = production
   ```

4. **Deploy**
   - Click "Deploy site"
   - Wait for build to complete
   - Test API at `https://your-site.netlify.app/api`

### Manual CLI Deployment

```bash
# Install Netlify CLI
npm install -g netlify-cli

# Login
netlify login

# Deploy
netlify deploy --prod
```

## Database Setup (Turso)

### Create Database

```bash
# Install Turso CLI
curl -sSfL https://get.turso.io | bash

# Create database
turso db create my-blog

# Get connection URL
turso db show my-blog
# Shows: libsql://my-blog-[org].turso.io

# Get auth token
turso auth tokens issue
# Copy the full token
```

### Initialize Schema

With database URL and token set:
```bash
DATABASE_URL="libsql://..." go run cmd/cms/main.go
```

## Environment Variables

### Required
- `DATABASE_URL` - Turso connection string (libsql://...)
- `ADMIN_PASSWORD` - Password for login
- `JWT_SECRET` - Secret key for JWT tokens (use `openssl rand -hex 32`)

### Optional
- `ENV` - "production" or "development" (default: development)
- `PORT` - Server port (default: 8080, ignored on Netlify)

## Build Configuration

The included `netlify.toml` handles everything:

```toml
[build]
  command = "mkdir -p netlify/functions && go build -o netlify/functions/cms ./cmd/functions/main.go"
  functions = "netlify/functions"
  publish = "public"

[[redirects]]
  from = "/api/*"
  to = "/.netlify/functions/cms/:splat"
```

What this does:
- Builds Go binary in `netlify/functions/cms`
- Routes `/api/*` requests to the function
- Serves static files from `public/` directory

## CI/CD Pipeline

### Automatic Deployment on Push

Netlify automatically:
1. Detects push to main branch
2. Runs build command: compiles Go binary
3. Deploys function to Netlify edge network
4. Available at `https://your-site.netlify.app`

No manual steps needed after first setup.

### View Logs

```bash
netlify logs
```

Or in dashboard: `Deploys → Select deploy → Logs`

## Database Backups

Turso automatically backs up your database. To export:

```bash
# Export all data
turso db dump my-blog > backup.sql

# Restore
turso db dump my-blog < backup.sql
```

## DNS & Custom Domain

In Netlify dashboard:

1. Go to `Domain settings`
2. Click "Add custom domain"
3. Point your domain's DNS to Netlify:
   - Use Netlify's nameservers, OR
   - Add CNAME record pointing to your Netlify subdomain

Example with popular providers:
- **Godaddy**: DNS settings → point to Netlify nameservers
- **Route53**: Create CNAME record
- **Cloudflare**: Use Netlify nameservers

## Monitoring

### Health Check
```bash
curl https://your-site.netlify.app/api
# Should return: {"status":"ok","version":"1.0"}
```

### Function Logs
- Netlify Dashboard → Functions
- Shows execution time, errors, cold starts

### Database Monitoring
- Turso Dashboard → Database name
- Shows queries, storage usage, bandwidth

## Scaling

### Request Limits
- Netlify: 100 requests/minute free tier
- Turso: Generous free tier, scales pay-as-you-go

### Cold Start Optimization
- First request after deploy: ~1-2s
- Subsequent requests: <100ms
- Go is fast for serverless (minimal overhead)

## Troubleshooting Deployment

### Build Fails
```
Error: module not found
```
Solution:
```bash
go mod tidy  # Ensure all dependencies are listed
git commit -am "Update dependencies"
git push
```

### Function Returns 500

Check Netlify logs:
```bash
netlify logs
```

Common causes:
- Missing `DATABASE_URL` environment variable
- Invalid Turso token
- Database schema not initialized

Fix:
```bash
# Reinitialize schema
DATABASE_URL="libsql://..." go run cmd/cms/main.go
```

### Cannot Connect to Database

Verify Turso token:
```bash
turso auth tokens list
```

Test connection locally:
```bash
DATABASE_URL="libsql://..." ./cms
```

### API Endpoints Return Errors

Test locally:
```bash
DATABASE_URL="file:./test.db" ADMIN_PASSWORD="test" JWT_SECRET="secret" ./cms
curl http://localhost:8080/api/posts
```

## Performance Tips

### Reduce Cold Starts
- Netlify keeps functions warm for 15 min after last request
- Keep JWT_SECRET small (not a long string)
- Database queries should have proper indexes

### Optimize Database
```bash
turso db stats my-blog
```

Turso has indexes on:
- posts.slug
- posts.status
- posts.published_at
- Add custom indexes if needed

### Monitor Bandwidth
Turso's free tier includes generous bandwidth. Track usage:
```bash
turso db stats my-blog
```

## Rollback Deployment

If something breaks:

1. Netlify dashboard: `Deploys`
2. Find previous working deploy
3. Click → `Restore this deploy`

Done. Previous version is live again.

## DNS Propagation

After changing DNS, wait 24-48 hours for full propagation.

During propagation, both old and new servers may respond.

Check status:
```bash
dig your-domain.com
nslookup your-domain.com
```

## Success Checklist

- [ ] Go binary builds: `go build ./cmd/functions/main.go`
- [ ] Database initializes: `go run cmd/cms/main.go`
- [ ] Local tests pass: `./cms` + curl tests
- [ ] GitHub repo is public/accessible
- [ ] Netlify is connected to GitHub
- [ ] Environment variables are set in Netlify
- [ ] First deploy completes without errors
- [ ] API endpoint responds: `curl https://your-site/api`
- [ ] Login works: `curl -X POST https://your-site/api/auth/login ...`
- [ ] Can create posts via API

## Security Checklist

- [ ] `JWT_SECRET` is strong and random (32+ chars)
- [ ] `ADMIN_PASSWORD` is strong (12+ chars, mixed case, symbols)
- [ ] Don't commit `.env` file (in `.gitignore`)
- [ ] Netlify environment variables are private
- [ ] Turso token is kept secret (rotate periodically)
- [ ] HTTPS only (Netlify enforces automatically)
- [ ] Consider using `PASSWORD_HASH` for production (bcrypt)

## Next Steps

1. [Set up admin UI](./README.md#next-steps) (HTMX frontend)
2. [Configure Hugo static site](./hugo/)
3. [Add daily sync job](./cronjob/)
4. [Deploy to custom domain](https://docs.netlify.com/domains-https/custom-domains/)
