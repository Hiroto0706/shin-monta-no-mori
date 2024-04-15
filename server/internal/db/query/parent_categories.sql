-- name: CreateParentCategories :one
INSERT INTO parent_categories (name, src)
VALUES ($1, $2)
RETURNING *;
-- name: GetParentCategories :one
SELECT *
FROM parent_categories
WHERE id = $1
LIMIT 1;
-- name: ListParentCategories :many
SELECT *
FROM parent_categories
ORDER BY id DESC
LIMIT $1 OFFSET $2;
-- name: UpdateParentCategories :one
UPDATE parent_categories
SET name = $2,
  src = $3
WHERE id = $1
RETURNING *;
-- name: DeleteParentCategories :exec
DELETE FROM parent_categories
WHERE id = $1;