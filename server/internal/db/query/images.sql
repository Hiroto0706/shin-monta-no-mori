-- name: CreateImage :one
INSERT INTO images (title, original_src, simple_src)
VALUES ($1, $2, $3)
RETURNING *;
-- name: GetImage :one
SELECT *
FROM images
WHERE id = $1
LIMIT 1;
-- name: ListImage :many
SELECT *
FROM images
ORDER BY id DESC
LIMIT $1 OFFSET $2;
-- name: UpdateImage :one
UPDATE images
SET title = $2,
  original_src = $3,
  simple_src = $4
WHERE id = $1
RETURNING *;
-- name: DeleteImage :exec
DELETE FROM images
WHERE id = $1;
-- name: SearchImages :many
SELECT *
FROM images
WHERE title LIKE '%' || COALESCE(sqlc.arg(title)) || '%'
ORDER BY id DESC
LIMIT $1 OFFSET $2;