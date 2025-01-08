-- name: CreateSession :one
INSERT INTO sessions (id, user_id, token, device_name, client, IP)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1;

-- name: GetUserSession :one
SELECT * FROM sessions
WHERE id = $1 AND user_id = $2;

-- name: GetSessionsByUserID :many
SELECT * FROM sessions
WHERE user_id = $1;

-- name: GetSessionToken :one
SELECT token FROM sessions
WHERE id = $1 AND user_id = $2;

-- name: UpdateSessionToken :exec
UPDATE sessions
SET
    token = $3,
    last_used = now()
WHERE id = $1 AND user_id = $2;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;

-- name: DeleteUserSessions :exec
DELETE FROM sessions
WHERE user_id = $1;

-- name: DeleteUserSession :exec
DELETE FROM sessions
WHERE id = $1 AND user_id = $2;