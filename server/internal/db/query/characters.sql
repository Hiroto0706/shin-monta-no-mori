-- name: CreateCharacter :one
INSERT INTO characters (name, src, filename, priority_level)
VALUES ($1, $2, $3, $4)
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
-- name: ListAllCharacters :many
SELECT *
FROM characters
ORDER BY priority_level DESC,
  id DESC;
-- name: UpdateCharacter :one
UPDATE characters
SET name = $2,
  src = $3,
  filename = $4,
  updated_at = $5,
  priority_level = $6
WHERE id = $1
RETURNING *;
-- name: DeleteCharacter :exec
DELETE FROM characters
WHERE id = $1;
-- name: SearchCharacters :many
SELECT DISTINCT *
FROM characters
WHERE name LIKE '%' || COALESCE(sqlc.arg(query)) || '%'
  OR filename LIKE '%' || COALESCE(sqlc.arg(query)) || '%'
ORDER BY priority_level DESC,
  id DESC
LIMIT $1 OFFSET $2;
-- name: CountCharacters :one
SELECT count(*)
FROM characters;
-- name: CountSearchCharacters :one
SELECT DISTINCT count(*)
FROM characters
WHERE name LIKE '%' || COALESCE(sqlc.arg(query)) || '%'
  OR filename LIKE '%' || COALESCE(sqlc.arg(query)) || '%';