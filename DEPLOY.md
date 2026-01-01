# Deployment Guide

## Setup

### 1. Push to GitHub
```bash
git init
git add .
git commit -m "Initial commit"
git remote add origin https://github.com/your-username/blog
git push -u origin main
```

### 2. Deploy Backend (CMS) to Vercel

#### Option A: Using Vercel Dashboard (Easiest)
1. Go to [vercel.com](https://vercel.com)
2. Click "Add New" → "Project"
3. Import your GitHub repo
4. Framework: "Other" (already has `vercel.json`)
5. Click "Deploy"
6. After deployment, set environment variables in Vercel dashboard:
   - `TURSO_CONNECTION_URL` - your Turso database URL
   - `TURSO_AUTH_TOKEN` - your Turso auth token
   - `ADMIN_PASSWORD` - secure password for admin login
   - `JWT_SECRET` - random secret (generate with `openssl rand -hex 32`)

#### Option B: Using Vercel CLI
```bash
npm install -g vercel
vercel link
vercel env add TURSO_CONNECTION_URL
vercel env add TURSO_AUTH_TOKEN
vercel env add ADMIN_PASSWORD
vercel env add JWT_SECRET
vercel deploy
```

### 3. Set GitHub Actions Secrets
For the cronjob to work, add these secrets to your repo:

```bash
gh secret set TURSO_CONNECTION_URL -b "libsql://..."
gh secret set TURSO_AUTH_TOKEN -b "..."
```

Or in GitHub UI:
1. Go to repo Settings → Secrets and variables → Actions
2. Add `TURSO_CONNECTION_URL`
3. Add `TURSO_AUTH_TOKEN`

### 4. Configure GitHub Pages
1. Go to repo Settings → Pages
2. Source: "Deploy from a branch"
3. Branch: `gh-pages` (the workflow creates this)
4. Save

## Architecture

```
┌──────────────────┐
│   Your Browser   │
└────────┬─────────┘
         │
         ↓
    ┌────────────────────┐
    │  Vercel (Backend)  │
    │ ┌──────────────────┤
    │ │ Node.js API      │  /api/posts, /api/auth
    │ │ + Frontend UI    │  /login.html, /editor.html
    │ └──────────────────┤
    └────────┬───────────┘
             │
             ↓
    ┌────────────────────┐
    │  Turso Database    │  SQLite in the cloud
    └────────────────────┘
             ▲
             │
    ┌────────┴───────────┐
    │ GitHub Actions     │
    │ (every 6 hours)    │
    │ 1. Fetch posts     │
    │ 2. Generate Hugo   │
    │ 3. Deploy to       │
    │    GitHub Pages    │
    └────────────────────┘
             │
             ↓
    ┌────────────────────┐
    │  GitHub Pages      │  https://your-site.com
    │  (Static Blog)     │
    └────────────────────┘
```

## What goes where

| Component | Location | Deployed To |
|-----------|----------|-------------|
| CMS Backend API | `/backend` | Vercel |
| CMS Frontend UI | `/frontend` | Vercel (served by backend) |
| Post Database | Turso | Cloud (free tier) |
| Static Blog | `/hugo/content` | GitHub Pages (free) |
| Export Script | `/cronjob` | GitHub Actions (free) |

## Workflow

1. **Create/Edit Posts**: Visit your Vercel URL → Login → Write post
2. **Publish**: Click "Publish" in editor → Saved to Turso
3. **Auto-Export** (every 6 hours): GitHub Actions fetches published posts from Turso → Generates Hugo markdown → Commits to repo
4. **Live Blog**: Hugo builds static HTML → Deploys to GitHub Pages

## Costs

- **Turso**: Free (9GB storage)
- **Vercel**: Free (generous API limits)
- **GitHub Pages**: Free (1GB limit)
- **GitHub Actions**: Free (2000 mins/month)
- **Domain**: ~$10/year (optional)

**Total: $0/month** (or ~$1/month if you add a custom domain)

## Testing Locally

```bash
# Start backend
cd backend
npm install
npm run dev

# In another terminal, trigger cronjob
cd cronjob
bun install
bun run export
```

## Troubleshooting

**Cronjob not running?**
- Check GitHub Actions tab in repo
- Verify `TURSO_CONNECTION_URL` and `TURSO_AUTH_TOKEN` are set
- Manually trigger: Actions tab → "Sync Posts and Build Static Site" → "Run workflow"

**Frontend not loading on Vercel?**
- Check that backend built successfully
- Verify `frontend/` folder exists with HTML files
- Check Vercel logs for errors

**Hugo build failing?**
- Check GitHub Actions logs
- Ensure `hugo/config.toml` exists
- Verify markdown files are in `hugo/content/`
