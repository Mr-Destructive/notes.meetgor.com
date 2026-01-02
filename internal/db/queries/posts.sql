-- Ensure tables exist
-- name: InitPostTables :exec
CREATE TABLE IF NOT EXISTS post_types (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  description TEXT
);
CREATE TABLE IF NOT EXISTS posts (
  id TEXT PRIMARY KEY,
  type_id TEXT NOT NULL,
  title TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  content TEXT NOT NULL,
  excerpt TEXT,
  status TEXT DEFAULT 'draft',
  is_featured BOOLEAN DEFAULT 0,
  tags TEXT,
  metadata TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  published_at DATETIME,
  FOREIGN KEY(type_id) REFERENCES post_types(id)
);
CREATE TABLE IF NOT EXISTS revisions (
  id TEXT PRIMARY KEY,
  post_id TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_posts_type ON posts(type_id);
CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
CREATE INDEX IF NOT EXISTS idx_posts_published_at ON posts(published_at);
CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);
INSERT OR IGNORE INTO post_types (id, name, slug, description) VALUES
  ('article', 'Article', 'article', 'Full-length articles'),
  ('review', 'Review', 'review', 'Book, movie, or product reviews'),
  ('thought', 'Thought', 'thought', 'Quick thoughts and reflections'),
  ('link', 'Link', 'link', 'Curated links with commentary'),
  ('til', 'TIL', 'til', 'Today I Learned'),
  ('quote', 'Quote', 'quote', 'Quotations and excerpts'),
  ('list', 'List', 'list', 'Curated lists'),
  ('note', 'Note', 'note', 'Quick notes'),
  ('snippet', 'Snippet', 'snippet', 'Code snippets'),
  ('essay', 'Essay', 'essay', 'Long-form essays'),
  ('tutorial', 'Tutorial', 'tutorial', 'Step-by-step guides'),
  ('interview', 'Interview', 'interview', 'Q&A interviews');

-- name: CreatePost :one
INSERT INTO posts (id, type_id, title, slug, content, excerpt, status, is_featured, tags, metadata, created_at, updated_at, published_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetPost :one
SELECT * FROM posts
WHERE id = ? OR slug = ?
LIMIT 1;

-- name: ListPosts :many
SELECT * FROM posts
WHERE 1=1
  AND (sqlc.arg('status') IS NULL OR status = sqlc.arg('status'))
  AND (sqlc.arg('type_id') IS NULL OR type_id = sqlc.arg('type_id'))
ORDER BY created_at DESC, published_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: UpdatePost :one
UPDATE posts SET
  title = COALESCE(sqlc.narg('title'), title),
  slug = COALESCE(sqlc.narg('slug'), slug),
  content = COALESCE(sqlc.narg('content'), content),
  excerpt = COALESCE(sqlc.narg('excerpt'), excerpt),
  status = COALESCE(sqlc.narg('status'), status),
  is_featured = COALESCE(sqlc.narg('is_featured'), is_featured),
  tags = COALESCE(sqlc.narg('tags'), tags),
  metadata = COALESCE(sqlc.narg('metadata'), metadata),
  published_at = COALESCE(sqlc.narg('published_at'), published_at),
  updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = ?;

-- name: CountPosts :one
SELECT COUNT(*) as count FROM posts
WHERE 1=1
  AND (sqlc.arg('status') IS NULL OR status = sqlc.arg('status'))
  AND (sqlc.arg('type_id') IS NULL OR type_id = sqlc.arg('type_id'));

-- name: GetPostTypes :many
SELECT * FROM post_types ORDER BY name;

-- name: CreateRevision :exec
INSERT INTO revisions (id, post_id, content, created_at)
VALUES (?, ?, ?, ?);

-- name: GetRevisions :many
SELECT * FROM revisions
WHERE post_id = ?
ORDER BY created_at DESC;
