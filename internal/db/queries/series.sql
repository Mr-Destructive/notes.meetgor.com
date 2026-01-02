-- Ensure tables exist
-- name: InitSeriesTables :exec
CREATE TABLE IF NOT EXISTS series (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  description TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS post_series (
  post_id TEXT NOT NULL,
  series_id TEXT NOT NULL,
  order_in_series INT,
  PRIMARY KEY(post_id, series_id),
  FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
  FOREIGN KEY(series_id) REFERENCES series(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_series_slug ON series(slug);

-- name: CreateSeries :one
INSERT INTO series (id, name, slug, description, created_at)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetSeries :one
SELECT * FROM series
WHERE id = ? OR slug = ?
LIMIT 1;

-- name: ListSeries :many
SELECT * FROM series
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: UpdateSeries :one
UPDATE series SET
  name = ?,
  slug = ?,
  description = ?
WHERE id = ?
RETURNING *;

-- name: DeleteSeries :exec
DELETE FROM series WHERE id = ?;

-- name: AddPostToSeries :exec
INSERT INTO post_series (post_id, series_id, order_in_series)
VALUES (?, ?, ?)
ON CONFLICT(post_id, series_id) DO UPDATE SET order_in_series = ?;

-- name: RemovePostFromSeries :exec
DELETE FROM post_series WHERE post_id = ? AND series_id = ?;

-- name: GetSeriesPosts :many
SELECT p.* FROM posts p
JOIN post_series ps ON p.id = ps.post_id
WHERE ps.series_id = ?
ORDER BY ps.order_in_series ASC;

-- name: GetPostSeries :many
SELECT s.* FROM series s
JOIN post_series ps ON s.id = ps.series_id
WHERE ps.post_id = ?
ORDER BY ps.order_in_series ASC;
