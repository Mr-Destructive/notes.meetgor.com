# Blog CMS + Static Site

A lightweight blog system with:
- **Pure Go + HTMX CMS** - Password-protected editor with crisp UI
- **Netlify Functions** - Serverless backend (Turso SQLite compatible)
- **Hugo Static Site** - Auto-synced posts with previews, tags, series
- **Multiple Post Types** - Articles, reviews, TILs, snippets, essays, etc.

## Stack

- **Backend**: Pure Go (net/http)
- **CMS UI**: HTMX + HTML (server-rendered, no JS framework)
- **Database**: SQLite (local) or Turso (production)
- **Serverless**: Netlify Functions
- **Static Site**: Hugo
- **Sync**: Daily cron job
- **Auth**: Password + JWT tokens

## Features

✓ 12 post types (articles, reviews, TILs, quotes, snippets, essays, etc.)  
✓ Post drafts with revision history  
✓ Series/collections for grouping posts  
✓ Tagging system  
✓ Markdown editor with live preview  
✓ Link preview extraction (images, YouTube, Twitter)  
✓ JSON metadata per post type  
✓ Status: draft/published/archived  
✓ Featured posts  
✓ Auto-export to Hugo markdown  

## Quick Start

### 1. Clone & Setup

```bash
git clone <repo>
cd blog
cp .env.example .env
```

### 2. Configure Environment

```bash
# .env
DATABASE_URL=file:./blog.db           # Local SQLite
ADMIN_PASSWORD=your-secure-password   # Set a strong password
JWT_SECRET=your-jwt-secret-key        # Random string
PORT=8080
ENV=development
```

For **Turso** (production):
```bash
DATABASE_URL=libsql://[db]-[org].turso.io?authToken=[token]
```

### 3. Initialize Database

```bash
go run cmd/cms/main.go
```

Output should show schema initialized with 12 post types.

### 4. Start CMS Server

```bash
export $(cat .env | xargs)
go run cmd/functions/main.go
```

Server runs on `http://localhost:8080`

### 5. Access CMS

Open browser to `http://localhost:8080` (coming next: HTMX admin UI)

## Project Structure

```
.
├── cmd/
│   ├── cms/
│   │   └── main.go              # Database initialization
│   └── functions/
│       └── main.go              # Netlify Functions handler
├── internal/
│   ├── db/
│   │   ├── client.go            # Database connection
│   │   ├── posts.go             # Post CRUD
│   │   ├── series.go            # Series CRUD
│   │   └── schema.sql           # Database schema
│   ├── models/
│   │   └── models.go            # Data structures
│   ├── handler/
│   │   └── auth.go              # Authentication
│   ├── editor/
│   │   └── [coming next]        # Markdown & preview utils
│   └── util/
│       └── util.go              # Helpers
├── web/
│   ├── templates/               # [coming next]
│   └── static/                  # [coming next]
├── netlify.toml                 # Netlify config
├── go.mod
├── API.md                       # API documentation
├── ARCHITECTURE.md              # Technical overview
└── README.md
```

## API

See [API.md](./API.md) for full endpoint documentation.

Quick examples:

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"password":"your-password"}'
```

### Create a post
```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type_id": "article",
    "title": "Hello World",
    "slug": "hello-world",
    "content": "# Content here",
    "excerpt": "Summary",
    "status": "draft"
  }'
```

### List posts
```bash
curl http://localhost:8080/api/posts?status=published
```

## Post Types

Predefined templates for different content:

| Type | Fields | Use Case |
|------|--------|----------|
| **article** | title, content, tags, reading_time | Full articles |
| **review** | title, content, rating, subject_link | Book/movie/product reviews |
| **thought** | title, content, tags | Quick reflections |
| **link** | title, url, excerpt, tags | Curated links |
| **til** | title, content, category, difficulty | Today I Learned |
| **quote** | text, author, source | Quotations |
| **list** | title, items, list_type | Curated lists |
| **note** | title, content, is_private | Quick notes |
| **snippet** | title, code, language | Code snippets |
| **essay** | title, content, reading_time | Long-form writing |
| **tutorial** | title, content, difficulty | Step-by-step guides |
| **interview** | title, questions, answers | Q&A interviews |

Each type stores additional metadata in JSON for flexibility.

## Database Schema

### Tables

- **posts** - Main content (id, type_id, title, slug, content, tags, metadata, status, published_at)
- **post_types** - Templates (article, review, etc.)
- **series** - Collections (name, slug, description)
- **post_series** - Many-to-many post↔series mapping
- **revisions** - Version history
- **settings** - Configuration key-value pairs

See `internal/db/schema.sql` for full schema.

## Development vs Production

### Local Development
```bash
# Uses local SQLite
DATABASE_URL=file:./blog.db go run cmd/functions/main.go
```

### Production (Netlify)
```bash
# Uses Turso + Netlify Functions
DATABASE_URL=libsql://... netlify deploy
```

Netlify automatically:
1. Builds Go binary from `cmd/functions/main.go`
2. Deploys as serverless function
3. Routes `/api/*` to the handler
4. Serves static site from `public/`

## Deployment

### To Netlify

1. Connect GitHub repo
2. Set environment variables:
   - `DATABASE_URL=libsql://...turso.io...`
   - `ADMIN_PASSWORD=...`
   - `JWT_SECRET=...`
   - `ENV=production`
3. Netlify builds and deploys automatically

### To Custom Domain
```bash
netlify deploy --prod
```

### Local Build
```bash
go build -o cms ./cmd/functions/main.go
./cms
```

## Next Steps

- [ ] HTMX CMS frontend (dashboard, editor, components)
- [ ] Markdown editor with preview
- [ ] Link preview extraction (images, YouTube, Twitter embeds)
- [ ] Hugo markdown export job
- [ ] Static site with filtering/tagging/series
- [ ] Search functionality
- [ ] Admin-only posts (private notes)
- [ ] Analytics/stats
- [ ] Categories/topic organization
- [ ] Multiple author support

## Architecture Docs

- [API.md](./API.md) - Full API reference
- [ARCHITECTURE.md](./ARCHITECTURE.md) - System design
- [DESIGN.md](./DESIGN.md) - Original specifications

## License

MIT
