-- name: CreateParentCategory :one
INSERT INTO parent_categories (name, src, filename)
VALUES ($1, $2, $3)
RETURNING *;
-- name: GetParentCategory :one
SELECT *
FROM parent_categories
WHERE id = $1
LIMIT 1;
-- name: ListParentCategories :many
SELECT *
FROM parent_categories
ORDER BY id DESC;
-- name: UpdateParentCategory :one
UPDATE parent_categories
SET name = $2,
  src = $3,
  filename = $4
WHERE id = $1
RETURNING *;
-- name: DeleteParentCategory :exec
DELETE FROM parent_categories
WHERE id = $1;
-- name: SearchParentCategories :many
SELECT DISTINCT *
FROM parent_categories
WHERE name LIKE '%' || COALESCE(sqlc.arg(query)) || '%'
  OR filename LIKE '%' || COALESCE(sqlc.arg(query)) || '%'
ORDER BY id DESC;