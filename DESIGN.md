# Blog System Design

## Database Schema (SQLite/Turso)

```sql
-- Post types: articles, reviews, thoughts, links, tils, quotes, lists, notes, snippets, essays
CREATE TABLE post_types (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  description TEXT
);

CREATE TABLE posts (
  id TEXT PRIMARY KEY,
  type_id TEXT NOT NULL REFERENCES post_types(id),
  title TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  content TEXT NOT NULL,
  excerpt TEXT,
  status TEXT DEFAULT 'draft', -- draft, published, archived
  is_featured BOOLEAN DEFAULT 0,
  tags TEXT, -- JSON array ["tag1", "tag2"]
  metadata TEXT, -- JSON for type-specific data
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  published_at DATETIME,
  FOREIGN KEY(type_id) REFERENCES post_types(id)
);

CREATE TABLE revisions (
  id TEXT PRIMARY KEY,
  post_id TEXT NOT NULL REFERENCES posts(id),
  content TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE settings (
  key TEXT PRIMARY KEY,
  value TEXT
);

-- Indexes
CREATE INDEX idx_posts_type ON posts(type_id);
CREATE INDEX idx_posts_status ON posts(status);
CREATE INDEX idx_posts_published_at ON posts(published_at);
CREATE INDEX idx_posts_slug ON posts(slug);
```

## Post Type Templates & Metadata

```json
{
  "article": {
    "fields": ["title", "content", "excerpt", "tags", "reading_time"],
    "metadata": { "reading_time": "number", "category": "string" }
  },
  "review": {
    "fields": ["title", "subject", "rating", "content", "tags"],
    "metadata": { "type": "book|movie|product", "rating": "1-5", "author": "string", "subject_link": "url" }
  },
  "thought": {
    "fields": ["title", "content", "tags"],
    "metadata": { "mood": "string" }
  },
  "link": {
    "fields": ["title", "url", "excerpt", "tags"],
    "metadata": { "source_url": "url", "domain": "string" }
  },
  "til": {
    "fields": ["title", "content", "tags"],
    "metadata": { "category": "string", "difficulty": "beginner|intermediate|advanced" }
  },
  "quote": {
    "fields": ["text", "author", "tags", "source"],
    "metadata": { "author": "string", "source": "url", "context": "string" }
  },
  "list": {
    "fields": ["title", "items", "tags"],
    "metadata": { "items": "json_array", "list_type": "ordered|unordered" }
  },
  "note": {
    "fields": ["title", "content", "tags"],
    "metadata": { "is_private": "boolean" }
  },
  "snippet": {
    "fields": ["title", "code", "language", "tags"],
    "metadata": { "language": "js|python|sql|etc", "syntax_highlight": "url" }
  },
  "essay": {
    "fields": ["title", "content", "excerpt", "tags"],
    "metadata": { "reading_time": "number", "word_count": "number" }
  }
}
```

## Project Structure

```
blog-system/
├── backend/
│   ├── src/
│   │   ├── db/
│   │   │   ├── schema.sql
│   │   │   ├── migrations.ts
│   │   │   └── client.ts
│   │   ├── api/
│   │   │   ├── auth.ts (login endpoint)
│   │   │   ├── posts.ts (CRUD)
│   │   │   ├── types.ts (post type definitions)
│   │   │   └── middleware.ts (auth validation)
│   │   ├── editors/
│   │   │   ├── form-templates.ts
│   │   │   └── validators.ts
│   │   ├── export/
│   │   │   ├── markdown-generator.ts
│   │   │   └── hugo-formatter.ts
│   │   └── index.ts
│   ├── package.json
│   └── .env.example
│
├── frontend/
│   ├── src/
│   │   ├── pages/
│   │   │   ├── login.html
│   │   │   ├── dashboard.html
│   │   │   └── editor.html
│   │   ├── js/
│   │   │   ├── editor.js (Monaco/TinyMCE integration)
│   │   │   ├── auth.js
│   │   │   └── api-client.js
│   │   └── css/
│   │
│   └── package.json
│
├── cronjob/
│   ├── export.ts (GitHub Actions script)
│   ├── config.ts
│   └── templates/
│       └── post.md.template
│
├── .github/workflows/
│   ├── sync-posts.yml (daily/hourly)
│   └── build-static.yml (on sync completion)
│
└── hugo/
    ├── content/
    ├── layouts/
    └── config.toml
```

## API Endpoints

```
POST   /api/auth/login          - Authenticate with password
POST   /api/posts               - Create post
GET    /api/posts              - List posts (with filters)
GET    /api/posts/:id          - Get post
PUT    /api/posts/:id          - Update post
DELETE /api/posts/:id          - Delete post
GET    /api/posts/:id/revisions - Get revision history
GET    /api/types              - Get post type templates
GET    /api/export             - Export posts as markdown/json
```

## Tech Stack Options

**Backend:**
- Node.js + Hono/Express (lightweight API)
- Turso client (better-sqlite3 for dev)
- Monaco Editor API or TinyMCE for rich editing
- JWT for session management

**Frontend:**
- Plain HTML/JS + Fetch API (lightweight)
- Or: Astro/SvelteKit for more control

**Export:**
- GitHub Actions every 6 hours
- Fetch from Turso → Generate markdown files
- Hugo builds static site
- Deploy to Vercel/Netlify/GitHub Pages

## Features to Consider

Additional post types:
- **snippets** - code snippets with syntax highlighting
- **essays** - long-form writing with reading time
- **notes** - quick notes, can be private
- **interviews** - Q&A format
- **experiments** - technical experiments/learnings
- **tutorials** - step-by-step guides
- **compilations** - roundups of related content

## Data Flow

1. **Create/Edit**: User logs in → Editor UI → API → Turso DB → DB updated
2. **Export**: GitHub Actions cron triggers → Fetch from Turso → Generate markdown + frontmatter → Commit to repo
3. **Build**: Hugo processes markdown → Static HTML → Deploy

## Password Strategy

```ts
// Simple password hash (use bcrypt/argon2 in production)
// Single password or per-user could be added later
const PASSWORD_HASH = await bcrypt.hash(process.env.ADMIN_PASSWORD, 10);

// Verify on login
const token = jwt.sign({ admin: true }, process.env.JWT_SECRET, { expiresIn: '7d' });
```

## Deployment

- **Database**: Turso (managed SQLite)
- **API**: Vercel / Railway / Fly.io
- **Frontend**: Same as API or separate
- **Static Site**: GitHub Pages / Netlify / Vercel (static export)
- **Secrets**: GitHub Actions repository secrets for passwords/keys
