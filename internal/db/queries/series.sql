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
