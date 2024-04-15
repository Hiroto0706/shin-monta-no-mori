-- name: CreateCharacter :one
INSERT INTO characters (name, src)
VALUES ($1, $2)
RETURNING *;
-- name: GetCharacter :one
SELECT *
FROM characters
WHERE id = $1
LIMIT 1;
-- name: ListCharacters :many
SELECT *
FROM characters
ORDER BY id DESC
LIMIT $1 OFFSET $2;
-- name: UpdateCharacter :one
UPDATE characters
SET name = $2,
  src = $3
WHERE id = $1
RETURNING *;
-- name: DeleteCharacter :exec
DELETE FROM characters
WHERE id = $1;