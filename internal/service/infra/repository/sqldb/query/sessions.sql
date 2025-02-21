-- name: CreateSession :one
INSERT INTO sessions (id, account_id, token, client, IP)
VALUES ($1, $2, $3, $4 , $5)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1;

-- name: GetSessionsByAccountID :many
SELECT * FROM sessions
WHERE account_id = $1;

-- name: UpdateSessionToken :one
UPDATE sessions
SET
    token = $3,
    IP = $4,
    last_used = now()
WHERE id = $1 AND account_id = $2
RETURNING *;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;

-- name: DeleteAccountSessions :exec
DELETE FROM sessions
WHERE account_id = $1;