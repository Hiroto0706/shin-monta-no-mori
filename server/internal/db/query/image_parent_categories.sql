-- name: CreateImageParentCategoryRelations :one
INSERT INTO image_parent_categories_relations (image_id, parent_category_id)
VALUES ($1, $2)
RETURNING *;
-- name: ListImageParentCategoryRelationsByImageID :many
SELECT *
FROM image_parent_categories_relations
WHERE image_id = $1
ORDER BY image_id DESC;
-- name: ListImageParentCategoryRelationsByParentCategoryID :many
SELECT *
FROM image_parent_categories_relations
WHERE parent_category_id = $1
ORDER BY parent_category_id DESC;
-- name: UpdateImageParentCategoryRelations :one
UPDATE image_parent_categories_relations
SET image_id = $2,
  parent_category_id = $3
WHERE id = $1
RETURNING *;
-- name: DeleteImageParentCategoryRelations :exec
DELETE FROM image_parent_categories_relations
WHERE id = $1;
-- name: DeleteAllImageParentCategoryRelationsByImageID :exec
DELETE FROM image_parent_categories_relations
WHERE image_id = $1;
-- name: DeleteAllImageParentCategoryRelationsByParentCategoryID :exec
DELETE FROM image_parent_categories_relations
WHERE parent_category_id = $1;