# GitHub Workflows

## sync-posts.yml (Local Build)

Runs the CMS binary, exports posts, and builds Hugo site all in one GitHub Actions job.

**When to use**: When you want the workflow to build and export locally in GitHub Actions

**Cron Schedule**: Every 6 hours (0, 6, 12, 18 UTC)

**Required Secrets**:
- `DATABASE_URL` - SQLite connection or Turso URL (e.g., `libsql://db-name.turso.io?authToken=token`)
- `ADMIN_PASSWORD` - Admin password for authentication
- `JWT_SECRET` - JWT secret for token generation

**Process**:
1. Checkout repository
2. Build Go CMS binary
3. Start CMS server
4. Authenticate and trigger export
5. Commit changes (if any)
6. Build Hugo site
7. Deploy to GitHub Pages

**Advantages**:
- No external dependencies
- Self-contained workflow
- Works offline with local database

---

## sync-posts-remote.yml (Remote API)

Calls a remote CMS API to fetch and export posts.

**When to use**: When your CMS is already deployed (Netlify, Vercel, custom server)

**Cron Schedule**: Every 6 hours (0, 6, 12, 18 UTC)

**Required Secrets**:
- `CMS_URL` - Full URL to your deployed CMS (e.g., `https://cms.example.com`)
- `ADMIN_PASSWORD` - Admin password

**Process**:
1. Checkout repository
2. Call remote CMS export API
3. Download exported posts
4. Parse and generate markdown files
5. Commit changes (if any)
6. Build Hugo site
7. Deploy to GitHub Pages

**Advantages**:
- Lighter workflow (no build step)
- Uses existing deployed CMS
- Faster execution

---

## Setup Instructions

### Step 1: Create Secrets

In GitHub repo → Settings → Secrets and variables → Actions:

**For local build workflow:**
```
DATABASE_URL = file:./blog.db
ADMIN_PASSWORD = your-secure-password
JWT_SECRET = your-secret-key
```

**For remote API workflow:**
```
CMS_URL = https://your-cms-domain.com
ADMIN_PASSWORD = your-secure-password
```

### Step 2: Enable GitHub Pages

1. Go to Settings → Pages
2. Set Source to "GitHub Actions"
3. Save

### Step 3: Verify Workflow

Check GitHub Actions tab for successful runs.

---

## Disabling Auto-Sync

Comment out the schedule in the workflow YAML:

```yaml
on:
  # schedule:
  #   - cron: '0 */6 * * *'
  workflow_dispatch:  # Manual only
```

Then manually trigger via: Actions tab → Workflow → "Run workflow"

---

## Customizing Cron Schedule

Edit the cron expression in either workflow:

```yaml
on:
  schedule:
    - cron: '0 0 * * *'  # Daily at midnight
    # - cron: '0 9,15 * * *'  # Twice daily at 9 AM and 3 PM
    # - cron: '0 * * * *'  # Every hour
```

[Cron syntax reference](https://crontab.guru/)

---

## Monitoring

View workflow runs:
1. GitHub repo → Actions tab
2. Click on "Sync Posts and Build Hugo Site"
3. See run history with status

Each run shows:
- ✓ Steps completed
- ✗ Failed steps with error logs
- Commit history with bot changes

---

## Troubleshooting

**Workflow fails with authentication error**:
- Check `ADMIN_PASSWORD` secret is correct
- Verify `JWT_SECRET` is set (if using local build)

**No posts exported**:
- Ensure posts have `status: "published"`
- Check API is working: `curl http://localhost:8080/api/exports`
- View workflow logs for details

**GitHub Pages not deploying**:
- Verify Settings → Pages → Source is "GitHub Actions"
- Check branch protection rules don't block bot commits
- Ensure Hugo theme exists in repository

**Commits not pushing**:
- Check default branch is `main`
- Verify repo permissions (token may need adjustment)
- Enable branch protection for merge checks if needed

---

## Notes

- `[skip ci]` in commit message prevents re-triggering other workflows
- Bot uses email `bot@blog.local` to identify sync commits
- Exports only **published** posts (draft/archived excluded)
- Hugo build happens even if no changes detected (ensures fresh output)
