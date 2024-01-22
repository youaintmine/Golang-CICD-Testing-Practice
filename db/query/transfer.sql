-- name: CreateTransfer :one
INSERT INTO transfers (
  from_acc,
  to_acc,
  amount
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
WHERE 
    from_acc = $1 OR
    to_acc = $2
ORDER BY id
LIMIT $3
OFFSET $4;