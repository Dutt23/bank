-- name: CreateAccount :one
INSERT INTO accounts (
  owner, 
  balance, 
  currency
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
where id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
where id = $1 LIMIT 1 
FOR UPDATE;

-- name: GetAccounts :many
SELECT * FROM accounts
where owner = $1
order by id
LIMIT $2
OFFSET $3;

-- name: UpdateAccount :one
UPDATE accounts 
SET balance = $2
where id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE from accounts where id = $1;

-- name: AddAccountBalance :one
UPDATE accounts 
SET balance = balance + sqlc.arg(amount) 
where id = sqlc.arg(id) 
RETURNING *;