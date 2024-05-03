-- name: CreateImage :one
INSERT INTO images (
    title,
    original_src,
    simple_src,
    original_filename,
    simple_filename
  )
VALUES ($1, $2, $3, $4, $5)
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
  simple_src = $4,
  original_filename = $5,
  simple_filename = $6
WHERE id = $1
RETURNING *;
-- name: DeleteImage :exec
DELETE FROM images
WHERE id = $1;
-- name: SearchImages :many
SELECT DISTINCT *
FROM images
WHERE title LIKE '%' || COALESCE(sqlc.arg(query)) || '%'
  OR original_filename LIKE '%' || COALESCE(sqlc.arg(query)) || '%'
ORDER BY id DESC
LIMIT $1 OFFSET $2;