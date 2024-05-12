-- name: CreateImageCharacterRelations :one
INSERT INTO image_characters_relations (image_id, character_id)
VALUES ($1, $2)
RETURNING *;
-- name: ListImageCharacterRelationsByImageID :many
SELECT *
FROM image_characters_relations
WHERE image_id = $1
ORDER BY image_id DESC;
-- name: ListImageCharacterRelationsByParentCategoryID :many
SELECT *
FROM image_characters_relations
WHERE character_id = $1
ORDER BY character_id DESC;
-- name: UpdateImageCharacterRelations :one
UPDATE image_characters_relations
SET image_id = $2,
  character_id = $3
WHERE id = $1
RETURNING *;
-- name: DeleteImageCharacterRelations :exec
DELETE FROM image_characters_relations
WHERE id = $1;
-- name: DeleteAllImageCharacterRelationsByImageID :exec
DELETE FROM image_characters_relations
WHERE image_id = $1;
-- name: DeleteAllImageCharacterRelationsByCharacterID :exec
DELETE FROM image_characters_relations
WHERE character_id = $1;