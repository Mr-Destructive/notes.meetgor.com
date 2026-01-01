# Blog System

A modern blog platform combining Turso database, password-protected editor, and static site generation.

## Features

- **Multiple Post Types**: Articles, reviews, thoughts, links, TILs, quotes, lists, notes, snippets, essays, tutorials, interviews, experiments
- **Password-Protected Admin Panel**: Rich markdown editor with form templates for each post type
- **Database**: Turso-hosted SQLite with automatic backup and version history
- **Static Generation**: GitHub Actions cronjob exports posts to markdown and generates static HTML with Hugo
- **Extensible Schema**: Type-specific metadata fields for customization

## Tech Stack

- **Backend**: Hono (TypeScript)
- **Database**: Turso (SQLite)
- **Frontend**: Plain HTML/JS (lightweight)
- **Export**: Node.js + Turso client
- **Static Site**: Hugo
- **CI/CD**: GitHub Actions

## Setup

### 1. Create Turso Database

```bash
# Install Turso CLI
curl -sSfL https://get.turso.io | bash

# Create database
turso db create my-blog

# Get connection details
turso db show my-blog
```

### 2. Backend Setup

```bash
cd backend
cp ../.env.example .env

# Edit .env with:
# - TURSO_CONNECTION_URL
# - TURSO_AUTH_TOKEN
# - ADMIN_PASSWORD (set a strong password)
# - JWT_SECRET

npm install
npm run db:init  # Initialize schema

npm run dev  # Start dev server on port 3000
```

### 3. Generate Password Hash

```bash
node -e "const bcrypt = require('bcryptjs'); bcrypt.hash('your-password', 10).then(hash => console.log(hash));"
```

Add the hash to `.env` as `PASSWORD_HASH`.

### 4. Frontend

Open `frontend/login.html` in browser (or serve with `npx http-server`).

- Login with your password
- Use the rich editor to create posts
- Select post type to get type-specific form fields
- Save as draft or publish directly

### 5. GitHub Actions Setup

```bash
# Add secrets to GitHub repository
gh secret set TURSO_CONNECTION_URL
gh secret set TURSO_AUTH_TOKEN
```

The cronjob will:
- Run every 6 hours (configurable)
- Fetch published posts from Turso
- Generate markdown files with Hugo frontmatter
- Commit to repo
- Build and deploy static site to GitHub Pages

## API Endpoints

```
POST   /api/auth/login           # Login with password
POST   /api/posts                # Create post
GET    /api/posts                # List posts
GET    /api/posts/:id            # Get post by ID or slug
PUT    /api/posts/:id            # Update post
DELETE /api/posts/:id            # Delete post
GET    /api/posts/:id/revisions  # Get revision history
```

Query parameters:
- `type`: Filter by post type (article, review, etc.)
- `status`: Filter by status (draft, published, archived)
- `limit`: Results per page (default: 50)
- `offset`: Pagination offset

## Post Types & Templates

Each post type has predefined form fields:

- **article** - Full articles with categories
- **review** - Book/movie/product reviews with ratings
- **thought** - Quick reflections
- **link** - Curated links with commentary
- **til** - Today I Learned (difficulty levels)
- **quote** - Quotations with attribution
- **list** - Ordered/unordered curated lists
- **note** - Quick notes (can be private)
- **snippet** - Code snippets with syntax highlighting
- **essay** - Long-form personal essays
- **tutorial** - Step-by-step guides with difficulty
- **interview** - Q&A format interviews
- **experiment** - Technical experiments

## Database Schema

```sql
posts (id, type_id, title, slug, content, excerpt, status, tags, metadata, ...)
post_types (id, name, slug, description)
revisions (id, post_id, content, created_at)
```

Type-specific data stored in `metadata` JSON column for flexibility.

## Extending

### Add New Post Type

1. Add to `POST_TEMPLATES` in `backend/src/types.ts`
2. Define type-specific metadata interface
3. Create form fields in `frontend/editor.js` POST_TEMPLATES
4. Metadata automatically exports to Hugo frontmatter

### Custom Hugo Theme

```
hugo/
├── themes/
│   └── my-theme/
├── content/     # Auto-generated markdown
└── config.toml
```

Run `hugo -d ../public` to build static site.

## Deployment

### Backend
- **Vercel**: `npm run build && npm run start`
- **Railway/Fly.io**: Docker + Node.js
- **Heroku**: Buildpack for Node.js

### Static Site
- **GitHub Pages**: Auto-deployed via Actions
- **Netlify**: Connect repo, set build command to `hugo -d public`
- **Vercel**: Static export with `hugo` build

## Cronjob Customization

Edit `.github/workflows/sync-posts.yml`:

```yaml
- cron: '0 */6 * * *'  # Run every 6 hours
- cron: '0 * * * *'    # Run every hour
- cron: '0 0 * * *'    # Run daily
```

## Security Notes

- Passwords are bcrypt-hashed
- JWTs valid for 7 days (configurable)
- All connections to Turso are HTTPS
- GitHub Actions secrets never exposed
- Admin panel is password-protected

## Local Development

```bash
# Backend
cd backend
npm run dev

# Frontend (separate terminal)
cd frontend
npx http-server

# Database (local SQLite for testing)
DATABASE_URL=file:./dev.db npm run db:init
```

Use local SQLite during development, switch to Turso for production.

## Troubleshooting

**Auth fails**: Check PASSWORD_HASH matches your password
**Posts not exporting**: Verify TURSO_CONNECTION_URL and token
**Cronjob doesn't run**: Check GitHub Actions secrets are set
**Hugo build fails**: Verify `hugo/config.toml` and theme setup

---

Made with ❤️ for bloggers
