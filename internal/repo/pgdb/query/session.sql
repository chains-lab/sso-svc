-- name: CreateSession :one
INSERT INTO sessions (id, account_id, hash_token)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM sessions
WHERE id = $1;

-- name: GetAccountSession :one
SELECT *
FROM sessions
WHERE account_id = $1
  AND id = $2;

-- name: GetSessionsByAccountID :many
SELECT *
FROM sessions
WHERE account_id = $1
ORDER BY created_at DESC, id DESC
    LIMIT $2
OFFSET $3;

-- name: CountSessionsByAccountID :one
SELECT count(*)
FROM sessions
WHERE account_id = $1;

-- name: UpdateSessionToken :one
UPDATE sessions
SET hash_token = $2, last_used = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteSessionByID :exec
DELETE FROM sessions
WHERE id = $1;

-- name: DeleteAccountSession :exec
DELETE FROM sessions
WHERE account_id = $1
  AND id = $2;

-- name: DeleteSessionsByAccountID :exec
DELETE FROM sessions
WHERE account_id = $1;