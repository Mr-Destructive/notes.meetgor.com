# Code Overview

## Go Source Files

### Entry Points

#### `cmd/cms/main.go`
- Database initialization script
- Creates schema from embedded SQL
- Seeds 12 post types
- Run once: `go run cmd/cms/main.go`
- ~40 lines

#### `cmd/functions/main.go`  
- Main API server
- HTTP handler for all endpoints
- Route parsing (/api/posts, /api/series, etc.)
- Netlify Functions compatible
- ~400 lines

### Internal Packages

#### `internal/db/client.go`
- Database connection management
- SQLite + Turso support
- Schema initialization
- ~50 lines

#### `internal/db/schema.sql`
- Embedded database schema
- 6 tables (posts, post_types, series, post_series, revisions, settings)
- 7 indexes
- 12 default post types
- ~100 lines

#### `internal/db/posts.go`
- Post CRUD operations
- CreatePost, GetPost, ListPosts, UpdatePost, DeletePost
- GetPostTypes, GetTags
- Revision management
- Filter and pagination support
- ~300 lines

#### `internal/db/series.go`
- Series CRUD operations
- CreateSeries, GetSeries, ListSeries, UpdateSeries, DeleteSeries
- Add/remove posts from series
- Get series posts (ordered)
- Get post's series
- ~250 lines

#### `internal/models/models.go`
- Post struct with all fields
- PostCreate, PostUpdate request types
- PostType, Series, Revision models
- JSON marshaling helpers
- ListOptions and TagCount types
- Database row scanning helpers
- ~200 lines

#### `internal/handler/auth.go`
- Password authentication
- JWT token generation
- Token verification
- Logout (clear cookies)
- Bcrypt password checking
- ~150 lines

#### `internal/handler/posts.go`
- Placeholder for complex post logic
- Reserved for validation, markdown processing
- ~10 lines (stub)

#### `internal/handler/series.go`
- Placeholder for series logic
- Reserved for validation, ordering
- ~10 lines (stub)

#### `internal/util/util.go`
- ID generation (hex encoded)
- Slug generation from strings
- Password hashing (bcrypt)
- Reading time estimation
- String utilities (truncate, etc.)
- ~80 lines

## File Statistics

| File | Lines | Purpose |
|------|-------|---------|
| cmd/functions/main.go | 400 | API server + routing |
| cmd/cms/main.go | 40 | DB initialization |
| internal/db/posts.go | 300 | Post operations |
| internal/db/series.go | 250 | Series operations |
| internal/models/models.go | 200 | Data structures |
| internal/handler/auth.go | 150 | Authentication |
| internal/db/client.go | 50 | DB connection |
| internal/util/util.go | 80 | Utilities |
| internal/db/schema.sql | 100 | Database schema |
| **TOTAL** | **~1,600** | **All Go code** |

## Architecture Diagram

```
Request (HTTP)
    ↓
cmd/functions/main.go
    ├─ Route parsing
    ├─ Auth middleware
    └─ Handler dispatch
         ↓
    internal/handler/
    ├─ auth.go (login, logout, verify)
    ├─ posts.go (CRUD posts)
    └─ series.go (CRUD series)
         ↓
    internal/db/
    ├─ client.go (DB connection)
    ├─ posts.go (Query execution)
    ├─ series.go (Query execution)
    └─ schema.sql (Table definitions)
         ↓
    internal/models/ (Data types)
         ↓
    SQLite/Turso Database
         ↓
JSON Response
```

## Key Functions

### Authentication (internal/handler/auth.go)
```go
HandleLogin(w, r, db)          // POST /auth/login
HandleLogout(w, r)             // POST /auth/logout
HandleVerify(w, r)             // GET /auth/verify
VerifyToken(next http.Handler) // Middleware
```

### Posts (internal/db/posts.go)
```go
CreatePost(ctx, post)          // POST /posts
GetPost(ctx, idOrSlug)         // GET /posts/:id
ListPosts(ctx, opts)           // GET /posts?...
UpdatePost(ctx, id, update)    // PUT /posts/:id
DeletePost(ctx, id)            // DELETE /posts/:id
GetPostTypes(ctx)              // GET /types
GetTags(ctx)                   // GET /tags
```

### Series (internal/db/series.go)
```go
CreateSeries(ctx, series)      // POST /series
GetSeries(ctx, idOrSlug)       // GET /series/:id
ListSeries(ctx, limit, offset) // GET /series
GetSeriesPosts(ctx, seriesID)  // GET /series/:id/posts
AddPostToSeries(ctx, ...)      // Internal
```

### Utilities (internal/util/util.go)
```go
GenerateID()                   // Create unique IDs
GenerateSlug(s)               // Create URL-friendly slugs
HashPassword(password)        // Bcrypt hashing
CheckPassword(hash, password) // Bcrypt verification
ReadingTime(content)          // Estimate reading time
```

## Database Interactions

### Direct Query Execution

#### GetPost by ID or Slug
```go
query := "SELECT ... FROM posts WHERE id = ? OR slug = ? LIMIT 1"
d.conn.QueryRowContext(ctx, query, idOrSlug, idOrSlug)
```

#### ListPosts with Filters
```go
// Dynamic WHERE clause based on filters
whereClause := ""
if opts.Status != "" {
    whereClause += " AND status = ?"
}
// ... build args array
query := "SELECT ... FROM posts WHERE 1=1" + whereClause
d.conn.QueryContext(ctx, query, args...)
```

#### CreatePost with JSON
```go
tagsJSON, _ := json.Marshal(post.Tags)
query := "INSERT INTO posts (...) VALUES (...)"
d.conn.ExecContext(ctx, query, ..., string(tagsJSON), ...)
```

## Error Handling

### HTTP Errors
- 200: Success
- 201: Created
- 400: Bad request
- 401: Unauthorized
- 404: Not found
- 405: Method not allowed
- 500: Server error
- 501: Not implemented

### Database Errors
- Wrapped with context: `fmt.Errorf("failed to create post: %w", err)`
- Returned as JSON: `{"error":"message"}`

## Testing Patterns

### Manual Testing
```bash
# Start server
./cms

# In another terminal
curl -X GET http://localhost:8080/api/posts
curl -X POST http://localhost:8080/api/auth/login ...
curl -X POST http://localhost:8080/api/posts ...
```

### What to Test
1. Health: `GET /api` → `{"status":"ok"}`
2. Types: `GET /api/types` → Array of 12 types
3. Auth: `POST /auth/login` → JWT token
4. CRUD: Create → Read → Update → Delete
5. Filters: `?status=draft&type=article`

## Dependencies

### Go Modules
- `github.com/mattn/go-sqlite3` - SQLite driver
- `github.com/joho/godotenv` - .env file support
- `golang.org/x/crypto/bcrypt` - Password hashing
- `github.com/golang-jwt/jwt/v5` - JWT tokens
- `encoding/json` - JSON marshaling (stdlib)
- `database/sql` - SQL interface (stdlib)
- `net/http` - HTTP server (stdlib)

## Build Steps

```bash
# 1. Compile
go build -o cms ./cmd/functions/main.go

# 2. Run
./cms

# 3. Or compile and run directly
go run cmd/functions/main.go
```

Result: Single 12MB binary with no external dependencies at runtime.

## Initialization Flow

```
go run cmd/cms/main.go
    ↓
Load .env file
    ↓
Create database connection (SQLite or Turso)
    ↓
Read schema.sql (embedded in binary)
    ↓
Execute schema: CREATE TABLE IF NOT EXISTS ...
    ↓
Insert default post types (12 total)
    ↓
Success: "✓ Found 12 post types"
```

## Runtime Flow

```
Request comes in (HTTP)
    ↓
cmd/functions/main.go Handler()
    ↓
Parse route: /api/{resource}/{id?}/{action?}
    ↓
Route to handleAuth, handlePosts, handleSeries, etc.
    ↓
Handler calls database methods: db.GetPost(), db.CreatePost(), etc.
    ↓
db.{operation} executes SQL query
    ↓
Scan results into struct
    ↓
Return JSON response
```

## Code Quality

- ✓ No magic numbers
- ✓ Clear function names
- ✓ Proper error handling with context
- ✓ Consistent formatting (gofmt)
- ✓ Comments on exported functions
- ✓ Database queries use parameterized statements (SQL injection safe)
- ✓ Struct tags for JSON marshaling
- ✓ Middleware pattern for auth

## Next Steps (Phase 2+)

### Phase 2: Frontend
Add to `internal/handler/`:
- `render.go` - HTMX template rendering
- Add to `web/`:
- `templates/*.html` - HTMX forms and components
- `static/css/style.css` - Minimal CSS

### Phase 3: Link Previews
Add to `internal/editor/`:
- `previews.go` - OpenGraph extraction
- HTTP client for fetching metadata

### Phase 4: Hugo Export
Add:
- `internal/export/` - Markdown generation
- `cronjob/sync.go` - GitHub Actions integration

### Phase 5: Public Site
- Hugo templates in `hugo/`
- Static site generation
- Search indexing

---

## Code Reading Guide

**Start here** → cmd/functions/main.go (see main route structure)
↓
**Then read** → internal/db/posts.go (see data access patterns)
↓
**Then read** → internal/handler/auth.go (see request handling)
↓
**Then read** → internal/models/models.go (see data structures)

This gives you the complete request → database → response flow.
