# GitHub Workflow Setup Guide

## Overview

Two automated workflows are available to sync your blog database with Hugo static site generation:

1. **sync-posts.yml** - Builds locally in GitHub Actions (recommended for local databases)
2. **sync-posts-remote.yml** - Calls your deployed CMS API (recommended for production)

Both run every 6 hours by default and auto-deploy to GitHub Pages.

---

## Quick Start (Choose One)

### Option A: Local Build Workflow

Best for: Local SQLite database or testing

**Step 1: Add Secrets to GitHub**

Navigate to: **Settings → Secrets and variables → Actions → New repository secret**

Add these secrets:

| Name | Value | Example |
|------|-------|---------|
| `DATABASE_URL` | SQLite or Turso URL | `file:./blog.db` or `libsql://db.turso.io?authToken=...` |
| `ADMIN_PASSWORD` | Admin password | `your-secure-password` |
| `JWT_SECRET` | JWT signing key | `your-secret-key-32-chars-or-more` |

**Step 2: Ensure GitHub Pages is Enabled**

1. Go to **Settings → Pages**
2. Set **Source** to "GitHub Actions"
3. Save

**Step 3: Enable Workflow**

The workflow `.github/workflows/sync-posts.yml` runs automatically every 6 hours.

To manually trigger: **Actions → Sync Posts and Build Hugo Site → Run workflow**

---

### Option B: Remote API Workflow

Best for: Deployed CMS on Netlify, Vercel, or custom domain

**Step 1: Ensure CMS is Deployed**

Your Go CMS must be accessible at a public URL with these endpoints:
- `POST {URL}/api/auth/login` - Returns JWT token
- `POST {URL}/api/exports` - Triggers export (requires auth)
- `GET {URL}/api/exports` - Returns published posts as JSON

**Step 2: Add Secrets to GitHub**

Navigate to: **Settings → Secrets and variables → Actions → New repository secret**

Add these secrets:

| Name | Value | Example |
|------|-------|---------|
| `CMS_URL` | Full CMS URL | `https://cms.example.com` |
| `ADMIN_PASSWORD` | Admin password | `your-secure-password` |

**Step 3: Ensure GitHub Pages is Enabled**

1. Go to **Settings → Pages**
2. Set **Source** to "GitHub Actions"
3. Save

**Step 4: Enable Workflow**

Edit `.github/workflows/sync-posts-remote.yml` to be active, or create it if using remote API.

To manually trigger: **Actions → Sync Posts (Remote API) → Run workflow**

---

## Directory Structure

The workflow expects this structure:

```
root/
├── .github/workflows/
│   ├── sync-posts.yml           # Local build workflow
│   ├── sync-posts-remote.yml    # Remote API workflow
│   └── README.md
├── exports/
│   ├── hugo.toml                # Hugo config (auto-generated)
│   ├── content/
│   │   └── posts/
│   │       ├── post-one.md
│   │       ├── post-two.md
│   │       └── ...
│   └── themes/
│       └── [your-theme]/        # Hugo theme
├── public/                       # Generated HTML (GitHub Pages)
├── internal/                     # Go source
├── cmd/                          # Go binaries
└── Makefile
```

---

## Workflow Steps Explained

### Local Build Workflow (sync-posts.yml)

```
1. Checkout code
2. Setup Go 1.23
3. Build CMS binary (go build)
4. Start CMS server
5. Login → Get JWT token
6. Call /api/exports to export
7. Verify markdown files created
8. Git commit changes
9. Setup Hugo
10. Build Hugo site (hugo --minify)
11. Deploy to GitHub Pages
```

**Time**: ~3-5 minutes

### Remote API Workflow (sync-posts-remote.yml)

```
1. Checkout code
2. Call remote CMS /api/auth/login
3. Get JWT token
4. Call remote CMS /api/exports
5. Fetch published posts as JSON
6. Parse and generate markdown files
7. Git commit changes
8. Setup Hugo
9. Build Hugo site (hugo --minify)
10. Deploy to GitHub Pages
```

**Time**: ~1-2 minutes

---

## Testing

### Test Local Build Workflow

```bash
# 1. Simulate CMS export locally
./cms &
sleep 2

# 2. Get token
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"password":"your-password"}' | jq -r '.token')

# 3. Trigger export
curl -X POST http://localhost:8080/api/exports \
  -H "Authorization: Bearer $TOKEN"

# 4. Verify exports/content/posts created
ls -la exports/content/posts/
```

### Test Remote API Workflow

```bash
# Replace https://cms.example.com with your URL
CMS_URL="https://cms.example.com"

# 1. Get token
TOKEN=$(curl -s -X POST "${CMS_URL}/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"password":"your-password"}' | jq -r '.token')

echo "Token: $TOKEN"

# 2. Fetch posts
curl -s "${CMS_URL}/api/exports" | jq .

# 3. Verify response has posts
```

---

## Customizing the Schedule

Edit the cron expression in the workflow YAML:

**Current**: Every 6 hours
```yaml
cron: '0 */6 * * *'  # 0:00, 6:00, 12:00, 18:00 UTC
```

**Alternatives**:
```yaml
cron: '0 0 * * *'        # Daily at midnight
cron: '0 * * * *'        # Every hour
cron: '0 9 * * 1'        # Weekly on Monday at 9 AM
cron: '0 0 1 * *'        # Monthly on 1st at midnight
```

[Cron syntax helper](https://crontab.guru/)

---

## Hugo Theme Setup

The workflow builds Hugo but needs a theme.

### Option 1: Add Theme as Submodule

```bash
cd exports
git submodule add https://github.com/theNewDynamic/gohugo-theme-ananke themes/ananke
git commit -m "Add Ananke theme"
git push
```

Update `exports/hugo.toml`:
```toml
theme = "ananke"
```

### Option 2: Use a CDN Theme

Some themes can be included via shortcodes or external sources. Check [Hugo Themes](https://themes.gohugo.io/).

### Option 3: Create Minimal Theme

Create `exports/themes/blog-theme/` with minimal HTML:

```bash
mkdir -p exports/themes/blog-theme/layouts/_default
mkdir -p exports/themes/blog-theme/static/css

# Create baseof.html, list.html, single.html templates
# See STATIC_SITE_GENERATION.md for examples
```

---

## Monitoring Workflow Runs

### View Workflow Status

1. Go to repo → **Actions** tab
2. Click **"Sync Posts and Build Hugo Site"**
3. View recent runs

### Understanding Run Output

Each run shows:
- ✅ Green = step succeeded
- ❌ Red = step failed
- ⏭️ Skipped = conditional step not triggered

Click on a step to see detailed logs.

### Common Success Indicators

```
✓ Export successful
✓ Generated 5 markdown files
✓ commit 1a2b3c4: chore: sync posts from database
✓ Hugo site built
✓ Deploy to GitHub Pages: Success
```

---

## Troubleshooting

### Workflow Won't Run

**Issue**: Workflow tab is empty or inactive

**Fix**: 
1. Ensure workflow file exists in `.github/workflows/`
2. Check branch is `main` (default)
3. Go to **Settings → Actions** → Ensure "Allow all actions and reusable workflows" is selected

### Authentication Fails

**Issue**: "Failed to authenticate" in logs

**Fix**:
1. Verify `ADMIN_PASSWORD` secret is correct
2. For local build: Verify `JWT_SECRET` is set
3. For remote API: Verify `CMS_URL` is accessible (no 404)
4. Test manually: `curl -X POST {CMS_URL}/api/auth/login`

### No Posts Exported

**Issue**: "Generated 0 markdown files"

**Fix**:
1. Verify posts exist: `curl {CMS_URL}/api/posts?status=published`
2. Ensure posts have `status: "published"`
3. Check draft/archived posts aren't included
4. Review CMS logs for export errors

### Hugo Build Fails

**Issue**: "hugo: not found" or theme errors

**Fix**:
1. Ensure theme exists in `exports/themes/{theme-name}/`
2. Check `hugo.toml` theme name matches folder
3. Try building locally: `cd exports && hugo`
4. View error logs in GitHub Actions

### GitHub Pages Won't Deploy

**Issue**: Site doesn't appear or shows 404

**Fix**:
1. Go to **Settings → Pages**
2. Verify Source is "GitHub Actions"
3. Check that `publish_dir: ./public` is being deployed
4. Wait 1-2 minutes for Pages to rebuild
5. View deployment status: **Deployments** tab

### Bot Commits Not Pushing

**Issue**: Workflow completes but git push fails

**Fix**:
1. Check branch protection rules
2. Ensure default branch is `main`
3. Verify token permissions (auto-provided by GitHub)
4. Check if repo is private vs public
5. Review git config settings in workflow

---

## Security Notes

- Secrets are encrypted and never logged
- `[skip ci]` prevents infinite workflow loops
- Use strong `ADMIN_PASSWORD` (30+ characters recommended)
- Rotate `JWT_SECRET` periodically
- For Turso: Use separate read-only tokens if available
- Monitor Actions usage (free tier has limits)

---

## Next Steps

1. ✅ Choose workflow (local or remote)
2. ✅ Add required secrets
3. ✅ Enable GitHub Pages
4. ✅ Add Hugo theme
5. ✅ Manually test workflow (Actions → Run workflow)
6. ✅ Wait for first auto-run (next cron time)
7. ✅ Verify site deploys to GitHub Pages

**Questions?** See `.github/workflows/README.md` for detailed workflow docs.
