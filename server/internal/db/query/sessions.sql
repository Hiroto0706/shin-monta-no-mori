-- name: CreateSession :one
INSERT INTO sessions (
    id,
    name,
    email,
    refresh_token,
    expires_at,
    user_agent,
    client_ip
  )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
-- name: GetSession :one
SELECT *
FROM sessions
WHERE id = $1
LIMIT 1;