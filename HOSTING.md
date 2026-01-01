# Hosting Options

## Database (Turso)
**Free and managed** - Already set up, nothing to do.
- Turso handles everything
- Free tier: up to 9GB storage
- No deployment needed

---

## Backend API

### **Vercel** (Recommended - Easiest)
- **Cost**: Free tier (generous limits)
- **Setup**: 5 minutes
- **Pros**: One-click GitHub integration, auto-deploys on push
- **Cons**: Serverless (cold starts possible)

```bash
# Install Vercel CLI
npm install -g vercel

# Deploy from root directory
vercel

# Follow prompts, link to GitHub repo
```

Add `vercel.json`:
```json
{
  "buildCommand": "cd backend && npm install",
  "outputDirectory": "backend",
  "env": {
    "TURSO_CONNECTION_URL": "@turso_url",
    "TURSO_AUTH_TOKEN": "@turso_token",
    "ADMIN_PASSWORD": "@admin_password",
    "PASSWORD_HASH": "@password_hash",
    "JWT_SECRET": "@jwt_secret"
  }
}
```

Then set environment variables in Vercel dashboard.

---

### **Railway** (Second Choice)
- **Cost**: Free $5/month credit, then ~$0.50/hour
- **Setup**: 10 minutes
- **Pros**: Simple, always-on, great for small projects
- **Cons**: Paid after credit runs out

```bash
# Install Railway CLI
npm install -g @railway/cli

# Login
railway login

# Create project
railway init

# Deploy
railway up
```

---

### **Fly.io**
- **Cost**: Free tier (3 shared-cpu-1x VMs)
- **Setup**: 15 minutes
- **Pros**: Good performance, generous free tier
- **Cons**: Requires Docker knowledge

```bash
# Install Fly CLI
curl -L https://fly.io/install.sh | sh

# Login
fly auth login

# Deploy
fly launch
fly deploy
```

---

### **Render**
- **Cost**: Free tier works but sleeps, paid from $7/month
- **Setup**: 10 minutes
- **Pros**: Straightforward, GitHub integration

1. Go to render.com
2. Connect GitHub repo
3. Create new Web Service
4. Set environment variables
5. Deploy

---

### **Heroku** (Deprecating)
Not recommended - they removed free tier.

---

## Frontend

### **Same as Backend** (Simplest)
Deploy frontend folder from same host:
- Vercel: Put `frontend/` files in public folder
- Railway: Use Node to serve static files
- Fly.io: Same
- Render: Same

### **Separate Static Hosting** (Optional)
- **Netlify** - Free tier, easy deploy
- **GitHub Pages** - Free, but only static
- **Cloudflare Pages** - Free, fast

Deploy `frontend/` folder separately if you prefer separation.

---

## Static Blog (Hugo Output)

Already handled by GitHub Actions!

The cronjob automatically:
1. Exports posts from Turso
2. Generates Hugo site
3. Deploys to GitHub Pages (free)

Or deploy to:
- **Vercel** - Free static hosting
- **Netlify** - Free static hosting
- **Cloudflare Pages** - Free, very fast

---

## Complete Setup Example: Vercel + GitHub Pages

**Cheapest & Simplest** (Everything free):

1. **Backend**: Vercel (free tier)
2. **Frontend**: Same Vercel deployment
3. **Database**: Turso (free tier)
4. **Static Blog**: GitHub Pages (free)

### Step-by-step:

**1. Push to GitHub**
```bash
git init
git add .
git commit -m "Initial commit"
git remote add origin https://github.com/your-username/blog
git push -u origin main
```

**2. Deploy Backend + Frontend to Vercel**
```bash
npm install -g vercel
vercel link  # Link to GitHub repo
vercel env add TURSO_CONNECTION_URL
vercel env add TURSO_AUTH_TOKEN
vercel env add ADMIN_PASSWORD
vercel env add PASSWORD_HASH
vercel env add JWT_SECRET
vercel deploy
```

Or use Vercel dashboard:
- Go to vercel.com
- Click "Add New" → "Project"
- Import GitHub repo
- Set environment variables
- Deploy

**3. Configure Workflow Secrets** (for cronjob)
```bash
gh secret set TURSO_CONNECTION_URL -b "libsql://..."
gh secret set TURSO_AUTH_TOKEN -b "..."
```

**4. Enable GitHub Pages**
- Go to repo Settings → Pages
- Set source to "Deploy from a branch"
- Branch: `gh-pages` (cronjob creates this)

That's it! Everything runs for free.

---

## Cost Breakdown

| Component | Free Tier | Notes |
|-----------|-----------|-------|
| **Turso DB** | Yes (9GB) | No upgrades needed for years |
| **Backend** | Vercel/Railway | ~10k req/day free |
| **Frontend** | Included | Served from backend |
| **Static Blog** | GitHub Pages | Unlimited, 1GB limit |
| **GitHub Actions** | Yes | 2000 mins/month free |
| **Domain** | ~$10/year | Optional (namecheap) |

**Total Monthly Cost: $0-10** (if you buy custom domain)

---

## Recommended Stack for You

```
┌─────────────────────────────────────────────┐
│                                             │
│  Browser                                    │
│    ↓                                        │
│  Frontend (HTML/JS) ─→ Vercel              │
│    ↓                                        │
│  API (Node.js) ─────→ Vercel               │
│    ↓                                        │
│  Database (Turso)                          │
│    ↓                                        │
│  Static Blog ────────→ GitHub Pages        │
│                                             │
└─────────────────────────────────────────────┘

All free. Deploy takes ~5 minutes.
```

---

## Migration Path

**Month 1**: Free tier (Vercel)
**Month 2+**: Still free (all tiers have generous limits)
**If needed**: Pay upgrade (usually <$20/month)

---

## Quick Deployment Checklist

- [ ] Push repo to GitHub
- [ ] Create Vercel account (free)
- [ ] Connect GitHub repo to Vercel
- [ ] Set environment variables in Vercel dashboard
- [ ] Trigger GitHub Actions manually
- [ ] Check GitHub Pages for static blog
- [ ] Set custom domain (optional, $10/year)

Done. Your blog is live.

