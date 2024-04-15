-- name: CreateChildCategories :one
INSERT INTO child_categories (name, parent_id)
VALUES ($1, $2)
RETURNING *;
-- name: GetChildCategories :one
SELECT *
FROM child_categories
WHERE id = $1
LIMIT 1;
-- name: ListChildCategories :many
SELECT *
FROM child_categories
ORDER BY id DESC
LIMIT $1 OFFSET $2;
-- name: UpdateChildCategories :one
UPDATE child_categories
SET name = $2,
  parent_id = $3
WHERE id = $1
RETURNING *;
-- name: DeleteChildCategories :exec
DELETE FROM child_categories
WHERE id = $1;