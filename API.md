# Blog CMS API

Pure Go serverless functions API, Netlify Function compatible.

## Base URL

- **Local**: `http://localhost:8080/api`
- **Netlify**: `https://your-site.netlify.app/api`
- **Custom Domain**: `https://yourdomain.com/api`

## Authentication

All protected endpoints require a JWT token via:
- **Header**: `Authorization: Bearer <token>`
- **Cookie**: `auth_token=<token>`

Tokens expire in 7 days.

## Endpoints

### Authentication

#### POST `/auth/login`
Login with password.

**Request:**
```json
{
  "password": "your-admin-password"
}
```

**Response:**
```json
{
  "token": "eyJ0eXAiOiJKV1QiLCJhbGc...",
  "expires_at": "2026-01-09T20:56:30Z"
}
```

**Errors:**
- `400` - Invalid request
- `401` - Invalid credentials

---

#### POST `/auth/logout`
Logout and clear token.

**Response:**
```json
{
  "status": "logged out"
}
```

---

#### GET `/auth/verify`
Check if token is valid.

**Response:**
```json
{
  "valid": true
}
```

**Errors:**
- `401` - No token or invalid token

---

### Posts

#### GET `/posts`
List all posts with filters.

**Query Parameters:**
- `limit` (int, default: 50) - Results per page
- `offset` (int, default: 0) - Pagination offset
- `type` (string) - Filter by post type (article, review, etc.)
- `status` (string) - Filter by status (draft, published, archived)
- `tag` (string) - Filter by tag
- `series` (string) - Filter by series ID

**Response:**
```json
{
  "posts": [
    {
      "id": "abc123",
      "type_id": "article",
      "title": "My Post",
      "slug": "my-post",
      "content": "...",
      "excerpt": "Short description",
      "status": "published",
      "is_featured": false,
      "tags": ["golang", "web"],
      "metadata": {},
      "created_at": "2026-01-02T20:56:30Z",
      "updated_at": "2026-01-02T20:56:30Z",
      "published_at": "2026-01-02T20:56:30Z"
    }
  ],
  "total": 42
}
```

**Errors:**
- `500` - Server error

---

#### GET `/posts/:id`
Get a single post by ID or slug.

**Response:**
```json
{
  "id": "abc123",
  "type_id": "article",
  "title": "My Post",
  "slug": "my-post",
  "content": "...",
  "excerpt": "Short description",
  "status": "published",
  "is_featured": false,
  "tags": ["golang", "web"],
  "metadata": {},
  "created_at": "2026-01-02T20:56:30Z",
  "updated_at": "2026-01-02T20:56:30Z",
  "published_at": "2026-01-02T20:56:30Z"
}
```

**Errors:**
- `404` - Post not found

---

#### POST `/posts`
Create a new post. **Requires auth**.

**Request:**
```json
{
  "type_id": "article",
  "title": "My New Article",
  "slug": "my-new-article",
  "content": "# Markdown content...",
  "excerpt": "Brief summary",
  "tags": ["golang", "tutorial"],
  "metadata": {
    "reading_time": 5,
    "category": "Development"
  },
  "is_featured": false,
  "status": "draft",
  "published_at": null
}
```

**Response:**
```json
{
  "id": "abc123",
  "type_id": "article",
  ...
}
```

**Errors:**
- `400` - Invalid request
- `401` - Unauthorized
- `500` - Server error

---

#### PUT `/posts/:id`
Update a post. **Requires auth**.

**Request:** (all fields optional)
```json
{
  "title": "Updated Title",
  "content": "...",
  "status": "published",
  "published_at": "2026-01-02T20:56:30Z"
}
```

**Response:**
```json
{
  "id": "abc123",
  ...
}
```

**Errors:**
- `400` - Invalid request
- `401` - Unauthorized
- `500` - Server error

---

#### DELETE `/posts/:id`
Delete a post. **Requires auth**.

**Response:**
```json
{
  "status": "deleted"
}
```

**Errors:**
- `400` - Post ID required
- `401` - Unauthorized
- `500` - Server error

---

### Series

#### GET `/series`
List all series.

**Query Parameters:**
- `limit` (int, default: 50)
- `offset` (int, default: 0)

**Response:**
```json
[
  {
    "id": "series1",
    "name": "Go Tutorial Series",
    "slug": "go-tutorial-series",
    "description": "Learn Go from basics",
    "created_at": "2026-01-02T20:56:30Z"
  }
]
```

---

#### GET `/series/:id`
Get a single series.

**Response:**
```json
{
  "id": "series1",
  "name": "Go Tutorial Series",
  "slug": "go-tutorial-series",
  "description": "Learn Go from basics",
  "created_at": "2026-01-02T20:56:30Z"
}
```

**Errors:**
- `404` - Series not found

---

#### GET `/series/:id/posts`
Get all posts in a series (ordered).

**Response:**
```json
[
  {
    "id": "post1",
    "title": "Part 1: Basics",
    ...
  },
  {
    "id": "post2",
    "title": "Part 2: Advanced",
    ...
  }
]
```

---

#### POST `/series`
Create a new series. **Requires auth**.

**Request:**
```json
{
  "name": "Go Tutorial Series",
  "slug": "go-tutorial-series",
  "description": "Learn Go from basics"
}
```

**Response:**
```json
{
  "id": "series1",
  "name": "Go Tutorial Series",
  ...
}
```

---

#### DELETE `/series/:id`
Delete a series. **Requires auth**.

**Response:**
```json
{
  "status": "deleted"
}
```

---

### Post Types

#### GET `/types`
Get all available post types.

**Response:**
```json
[
  {
    "id": "article",
    "name": "Article",
    "slug": "article",
    "description": "Full-length articles"
  },
  {
    "id": "review",
    "name": "Review",
    "slug": "review",
    "description": "Book, movie, or product reviews"
  },
  ...
]
```

---

### Tags

#### GET `/tags`
Get all unique tags with post counts.

**Response:**
```json
[
  {
    "tag": "golang",
    "count": 12
  },
  {
    "tag": "web",
    "count": 8
  }
]
```

---

### Exports

#### GET `/exports`
Export posts as JSON or markdown.

**Query Parameters:**
- `format` (string, default: "json") - "json" or "markdown"
- `status` (string, default: "published") - Filter by status

**Response (JSON):**
```json
[
  {
    "id": "post1",
    "title": "My Post",
    ...
  }
]
```

**Errors:**
- `501` - Format not supported

---

## Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad request
- `401` - Unauthorized
- `404` - Not found
- `405` - Method not allowed
- `500` - Server error
- `501` - Not implemented

## Error Response

All errors return JSON:
```json
{
  "error": "Error message"
}
```

## Examples

### Create and publish a post
```bash
# Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"password":"test"}' | jq -r '.token')

# Create post
curl -X POST http://localhost:8080/api/posts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type_id": "article",
    "title": "Hello World",
    "slug": "hello-world",
    "content": "# Hello\n\nThis is my first post!",
    "excerpt": "My first post",
    "tags": ["hello"],
    "status": "published",
    "published_at": "2026-01-02T20:56:30Z"
  }'
```

### List published articles
```bash
curl http://localhost:8080/api/posts?status=published&type=article&limit=10
```

### Get posts in a series
```bash
curl http://localhost:8080/api/series/go-tutorial-series/posts
```

## Deployment

### Netlify
Automatically deployed from `cmd/functions/main.go`. Just push to main branch.

### Local Development
```bash
export DATABASE_URL="file:./blog.db"
export ADMIN_PASSWORD="your-password"
export JWT_SECRET="your-secret-key"
go run cmd/functions/main.go
```

Server runs on `http://localhost:8080` by default.
