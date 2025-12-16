-- name: CreateAccount :one
INSERT INTO accounts (username, role, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateAccountEmail :one
INSERT INTO account_emails (account_id, email, verified)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateAccountPassword :one
INSERT INTO account_passwords (account_id, hash)
VALUES ($1, $2)
RETURNING *;

-- name: GetAccountByID :one
SELECT * FROM accounts
WHERE id = $1;

-- name: GetAccountByUsername :one
SELECT * FROM accounts
WHERE username = $1;

-- name: GetAccountByEmail :one
SELECT a.* FROM accounts a
JOIN account_emails ae ON a.id = ae.account_id
WHERE ae.email = $1;

-- name: GetAccountEmail :one
SELECT * FROM account_emails
WHERE account_id = $1;

-- name: GetAccountPassword :one
SELECT * FROM account_passwords
WHERE account_id = $1;

-- name: UpdateVerifiedEmail :one
UPDATE account_emails
SET verified = $2, updated_at = NOW()
WHERE account_id = $1
RETURNING *;

-- name: UpdateAccountUsername :one
UPDATE accounts
SET username = $2, updated_at = NOW(), username_updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateAccountPassword :one
UPDATE account_passwords
SET hash = $2, updated_at = NOW()
WHERE account_id = $1
RETURNING *;

-- name: UpdateAccountStatus :one
UPDATE accounts
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;