-- name: CreateOperator :one
INSERT INTO operators (name, hashed_password, email)
VALUES ($1, $2, $3)
RETURNING *;
-- name: GetOperatorByEmail :one
SELECT *
FROM operators
WHERE email = $1
LIMIT 1;
-- name: UpdateOperator :one
UPDATE operators
SET name = $2,
  hashed_password = $3,
  email = $4
WHERE id = $1
RETURNING *;
-- -- name: ListOperator :many
-- SELECT *
-- FROM operators
-- ORDER BY id DESC
-- LIMIT $1 OFFSET $2;
-- -- name: DeleteOperator :exec
-- DELETE FROM operators
-- WHERE id = $1;