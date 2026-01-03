# Vercel Deployment Setup

## Overview

The blog will be automatically deployed to Vercel (notes.meetgor.com) whenever you push changes.

## Setup Steps

### 1. Create Vercel Project

You have two options:

**Option A: Automatic (Recommended)**
1. Go to https://vercel.com
2. Import this GitHub repository
3. Select the `exports` directory as the root
4. Hugo will auto-detect and build

**Option B: Manual**
1. Create new project on Vercel dashboard
2. Connect to GitHub repo
3. Set environment:
   - **Root Directory:** `exports`
   - **Build Command:** `hugo --minify`
   - **Output Directory:** `public`
   - **Framework Preset:** Other

### 2. Add GitHub Secrets

These secrets enable the GitHub Actions workflow to deploy:

1. Go to GitHub repo → Settings → Secrets and variables → Actions
2. Add these secrets:

| Secret | Value | Where to find |
|--------|-------|---------------|
| `VERCEL_TOKEN` | Vercel authentication token | Vercel Settings → Tokens |
| `VERCEL_ORG_ID` | Your Vercel team/org ID | Vercel Settings → General |
| `VERCEL_PROJECT_ID` | Project ID for this blog | Vercel project → Settings → General |

**Get VERCEL_TOKEN:**
1. Go to https://vercel.com/account/tokens
2. Click "Create Token"
3. Name it "GitHub Actions"
4. Copy and add to GitHub Secrets

**Get VERCEL_ORG_ID:**
1. Go to https://vercel.com/dashboard/settings
2. Copy "Team ID" from General section

**Get VERCEL_PROJECT_ID:**
1. Go to project Settings → General
2. Copy "Project ID"

### 3. Configure Domain

Point `notes.meetgor.com` to Vercel:

**Option A: Using Vercel's Nameservers (Recommended)**
1. In Vercel project → Settings → Domains
2. Add `notes.meetgor.com`
3. Follow Vercel's DNS setup instructions
4. Update your registrar to use Vercel's nameservers

**Option B: Using CNAME (without changing nameservers)**
1. Add domain to Vercel project
2. Get the CNAME target (usually `cname.vercel.com`)
3. In your DNS provider, create:
   ```
   notes.meetgor.com CNAME cname.vercel.com
   ```

**Option C: Using A Records**
1. Get Vercel's IP addresses from project settings
2. Create A records in your DNS:
   ```
   notes.meetgor.com A 76.76.19.132
   notes.meetgor.com A 76.76.47.132
   ```

### 4. Verify Configuration

Vercel should now be configured. To test:

1. Make a commit to main branch
2. GitHub Actions will run `deploy-vercel.yml`
3. Check deployment at https://notes.meetgor.com/

## How It Works

```
You commit to main
        ↓
GitHub Actions triggers
        ↓
Hugo builds site (exports/public/)
        ↓
Vercel action deploys to Vercel
        ↓
Site live at notes.meetgor.com
```

## Multiple Domains/Subdomains

You can route multiple repos to different subdomains:

**Your Current Setup:**
```
dev.meetgor.com    → mr-destructive.github.io (GitHub Pages)
notes.meetgor.com  → this repo (Vercel)
```

To add more subdomains:

1. Create new Vercel project for each subdomain
2. Add DNS CNAME record for subdomain
3. Connect GitHub repo to Vercel project
4. Each repo deploys independently

**DNS Setup Example:**
```
dev.meetgor.com      CNAME cname.vercel.com    (or GitHub Pages)
notes.meetgor.com    CNAME cname.vercel.com    (Vercel)
blog.meetgor.com     CNAME cname.vercel.com    (another Vercel project)
www.meetgor.com      CNAME meetgor.com         (main domain)
```

## Environment Variables (Optional)

If you need environment variables in the build, add them in Vercel:

1. Vercel project → Settings → Environment Variables
2. Add variables needed for build
3. They'll be available during `hugo` build

Example:
```
SITE_TITLE=Notes
SITE_URL=https://notes.meetgor.com
```

## Deployments

### Automatic Deployments
- Any push to `main` branch
- Changes to `exports/**` files
- Changes to `.github/workflows/deploy-vercel.yml`

### Manual Deployment
1. Go to GitHub repo
2. Click "Actions" tab
3. Select "Deploy to Vercel" workflow
4. Click "Run workflow"

### Preview Deployments
Every pull request creates a preview deployment on Vercel automatically.

## Troubleshooting

### Deployment Failed

Check GitHub Actions logs:
1. Go to repo → Actions tab
2. Find failed workflow
3. Click to see error details
4. Common issues:
   - Hugo not installed (fix: use peaceiris/actions-hugo)
   - Missing `vercel.json` (we created it)
   - Wrong secrets (verify VERCEL_TOKEN, etc)

### Domain Not Working

1. Verify DNS is pointing to Vercel
2. Check Vercel project Settings → Domains
3. Wait 24-48 hours for DNS propagation
4. Use `nslookup notes.meetgor.com` to verify DNS

### Build Fails in Vercel

1. Check Vercel project Deployments tab
2. Click failed deployment to see logs
3. Verify Hugo can build locally:
   ```bash
   cd exports
   hugo --minify
   ```

## GitHub Pages vs Vercel

**GitHub Pages:**
- ✅ Free
- ✅ Built-in to GitHub
- ❌ One free repo per account (mr-destructive.github.io)
- ❌ Can't use multiple subdomains easily

**Vercel:**
- ✅ Free tier with good limits
- ✅ Multiple projects/subdomains
- ✅ Better performance (CDN, edge functions)
- ✅ Preview deployments
- ✅ Automatic previews for PRs
- ✅ Better analytics and monitoring

For your use case, **Vercel is better** because:
1. You already have dev.meetgor.com on GitHub Pages
2. You want notes.meetgor.com separate
3. Vercel handles multiple subdomains easily
4. Better performance and features

## Summary

Your setup:
- ✅ `vercel.json` created (Vercel knows how to build)
- ✅ GitHub Actions workflow created (auto-deploy on push)
- ⏳ Add GitHub Secrets (VERCEL_TOKEN, VERCEL_ORG_ID, VERCEL_PROJECT_ID)
- ⏳ Configure domain in Vercel (add notes.meetgor.com)
- ⏳ Update DNS to point to Vercel

Once DNS propagates, visiting notes.meetgor.com will show your blog!
