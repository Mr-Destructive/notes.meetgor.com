# Blog CMS Architecture

Pure Go + HTMX implementation with SQLite/Turso backend.

## Project Structure

```
.
├── cmd/
│   └── cms/
│       └── main.go              # Entry point
├── internal/
│   ├── db/
│   │   ├── client.go            # Database connection
│   │   ├── posts.go             # Post CRUD operations
│   │   ├── series.go            # Series/collection operations
│   │   └── schema.sql           # Database schema (embedded)
│   ├── models/
│   │   └── models.go            # Data structures
│   ├── handler/
│   │   ├── auth.go              # Authentication handlers
│   │   ├── posts.go             # Post handlers
│   │   ├── series.go            # Series handlers
│   │   └── render.go            # HTMX template rendering
│   ├── editor/
│   │   ├── markdown.go          # Markdown utilities
│   │   └── previews.go          # Link/image preview extraction
│   └── util/
│       └── util.go              # Helper utilities
├── web/
│   ├── templates/
│   │   ├── layout.html          # Base layout
│   │   ├── login.html           # Login page
│   │   ├── dashboard.html       # Post list
│   │   ├── editor.html          # Post editor
│   │   └── components/          # HTMX partials
│   └── static/
│       ├── css/
│       │   └── style.css        # Minimal CSS
│       └── js/
│           └── editor.js        # Client-side markdown editor
├── go.mod
├── go.sum
├── .env.example
└── README.md
```

## Database Schema

### Main Tables

- **post_types**: Article, review, thought, link, TIL, quote, list, note, snippet, essay, tutorial, interview
- **posts**: Blog posts with type_id, slug, content, tags (JSON), metadata (JSON), status (draft/published/archived)
- **series**: Collections/series for grouping posts
- **post_series**: Many-to-many relationship between posts and series
- **revisions**: Version history for posts

## Key Features Completed

✓ Database schema with all post types  
✓ SQLite/Turso client in Go  
✓ Full CRUD models for posts and series  
✓ ID generation and slug utilities  
✓ Password hashing with bcrypt  

## Next Steps

1. **HTMX Frontend** - Server-rendered templates with interactive CRUD
2. **Authentication** - Password-based login with JWT tokens
3. **Post Editor** - Markdown editor with live preview
4. **Series Management** - Add/edit/delete series with post ordering
5. **Export/Sync** - Daily export to Hugo markdown files
6. **Link Previews** - Extract images, YouTube, Twitter metadata

## Running Locally

```bash
# Set up environment
export DATABASE_URL="file:./blog.db"
export ADMIN_PASSWORD="yourpassword"
export JWT_SECRET="your-secret-key"

# Initialize database
go run cmd/cms/main.go

# Run server (next)
go run cmd/cms/main.go --serve
```

## Tech Stack

- **Backend**: Pure Go (net/http)
- **Templates**: HTML + HTMX for interactivity
- **Database**: SQLite with go-sqlite3 driver (swappable to Turso)
- **Auth**: bcrypt + JWT
- **Frontend**: Minimal CSS, no frameworks
- **Editor**: Simple markdown textarea with preview

## Configuration

All configuration via environment variables (see `.env.example`):
- `DATABASE_URL` - SQLite path or Turso connection string
- `ADMIN_PASSWORD` - Password for login
- `JWT_SECRET` - Secret for session tokens
- `PORT` - Server port (default 8080)
