package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAcc int64 `json:"from_acc"`
	ToAcc   int64 `json:"to_acc"`
	Amount  int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfers Transfers `json:"transfer"`
	FromAcc   Accounts  `json:"from_acc"`
	ToAcc     Accounts  `json:"to_acc"`
	FromEntry Entries   `json:"from_entry"`
	ToEntry   Entries   `json:"to_entry"`
}

// var txKey = struct{}{}

//TransferTx perform money transfer from one account to the other
//It creates a transfer record, add bank entries, and update accounts'
// balance within a single database transaction

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// txName := ctx.Value(txKey)

		// fmt.Println(txName, "create transfer")
		result.Transfers, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAcc: arg.FromAcc,
			ToAcc:   arg.ToAcc,
			Amount:  arg.Amount,
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAcc,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		// fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAcc,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		//TODO: update accounts' balance
		// get account -> update

		if arg.FromAcc < arg.ToAcc {
			result.FromAcc, result.ToAcc, err = addMoney(ctx, q, arg.FromAcc, -arg.Amount, arg.ToAcc, arg.Amount)

			if err != nil {
				return err
			}

		} else {
			result.ToAcc, result.FromAcc, err = addMoney(ctx, q, arg.ToAcc, arg.Amount, arg.FromAcc, -arg.Amount)

			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Accounts, account2 Accounts, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})

	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	return
}
