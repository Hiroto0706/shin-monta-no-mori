-- name: CreateChildCategory :one
INSERT INTO child_categories (name, parent_id)
VALUES ($1, $2)
RETURNING *;
-- name: GetChildCategory :one
SELECT *
FROM child_categories
WHERE id = $1
LIMIT 1;
-- name: GetChildCategoriesByParentID :many
SELECT *
FROM child_categories
WHERE parent_id = $1;
-- name: GetChildCategoriesByImageID :many
SELECT *
FROM child_categories
WHERE image_id = $1;
-- name: ListChildCategories :many
SELECT *
FROM child_categories
ORDER BY id DESC
LIMIT $1 OFFSET $2;
-- name: UpdateChildCategory :one
UPDATE child_categories
SET name = $2,
  parent_id = $3
WHERE id = $1
RETURNING *;
-- name: DeleteChildCategory :exec
DELETE FROM child_categories
WHERE id = $1;