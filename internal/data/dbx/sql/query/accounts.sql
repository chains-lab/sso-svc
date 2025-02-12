-- name: CreateAccount :one
INSERT INTO accounts (
    email,
    role
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetAccountByID :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountByEmail :one
SELECT * FROM accounts
WHERE email = $1 LIMIT 1;

-- name: UpdateAccountRole :one
UPDATE accounts SET
    role = $2
WHERE id = $1
RETURNING *;