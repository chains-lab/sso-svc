-- name: CreateAccount :one
INSERT INTO accounts (
    email
) VALUES (
    $1
) RETURNING *;

-- name: GetAccountByID :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountByEmail :one
SELECT * FROM accounts
WHERE email = $1 LIMIT 1;