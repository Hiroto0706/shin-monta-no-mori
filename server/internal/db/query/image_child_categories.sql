-- name: CreateImageChildCategoryRelations :one
INSERT INTO image_child_categories_relations (image_id, child_category_id)
VALUES ($1, $2)
RETURNING *;
-- name: ListImageChildCategoryRelationsByImageID :many
SELECT *
FROM image_child_categories_relations
WHERE image_id = $1
ORDER BY image_id DESC;
-- name: ListImageChildCategoryRelationsByChildCategoryID :many
SELECT *
FROM image_child_categories_relations
WHERE child_category_id = $1
ORDER BY child_category_id DESC;
-- name: ListImageChildCategoryRelationsByChildCategoryIDWithPagination :many
SELECT *
FROM image_child_categories_relations
WHERE child_category_id = $3
ORDER BY child_category_id DESC
LIMIT $1 OFFSET $2;
-- name: UpdateImageChildCategoryRelations :one
UPDATE image_child_categories_relations
SET image_id = $2,
  child_category_id = $3
WHERE id = $1
RETURNING *;
-- name: DeleteImageChildCategoryRelations :exec
DELETE FROM image_child_categories_relations
WHERE id = $1;
-- name: DeleteAllImageChildCategoryRelationsByImageID :exec
DELETE FROM image_child_categories_relations
WHERE image_id = $1;
-- name: DeleteAllImageChildCategoryRelationsByChildCategoryID :exec
DELETE FROM image_child_categories_relations
WHERE child_category_id = $1;