// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: transfer.sql

package db

import (
	"context"
)

const createTransfer = `-- name: CreateTransfer :one
INSERT INTO transfers (
  from_acc,
  to_acc,
  amount
) VALUES (
  $1, $2, $3
) RETURNING id, from_acc, to_acc, amount, created_at
`

type CreateTransferParams struct {
	FromAcc int64 `json:"from_acc"`
	ToAcc   int64 `json:"to_acc"`
	Amount  int64 `json:"amount"`
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfers, error) {
	row := q.db.QueryRowContext(ctx, createTransfer, arg.FromAcc, arg.ToAcc, arg.Amount)
	var i Transfers
	err := row.Scan(
		&i.ID,
		&i.FromAcc,
		&i.ToAcc,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getTransfer = `-- name: GetTransfer :one
SELECT id, from_acc, to_acc, amount, created_at FROM transfers
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTransfer(ctx context.Context, id int64) (Transfers, error) {
	row := q.db.QueryRowContext(ctx, getTransfer, id)
	var i Transfers
	err := row.Scan(
		&i.ID,
		&i.FromAcc,
		&i.ToAcc,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const listTransfers = `-- name: ListTransfers :many
SELECT id, from_acc, to_acc, amount, created_at FROM transfers
WHERE 
    from_acc = $1 OR
    to_acc = $2
ORDER BY id
LIMIT $3
OFFSET $4
`

type ListTransfersParams struct {
	FromAcc int64 `json:"from_acc"`
	ToAcc   int64 `json:"to_acc"`
	Limit   int32 `json:"limit"`
	Offset  int32 `json:"offset"`
}

func (q *Queries) ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfers, error) {
	rows, err := q.db.QueryContext(ctx, listTransfers,
		arg.FromAcc,
		arg.ToAcc,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transfers{}
	for rows.Next() {
		var i Transfers
		if err := rows.Scan(
			&i.ID,
			&i.FromAcc,
			&i.ToAcc,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
