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
