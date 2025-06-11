package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	Querier
}

// Store provides all functions to execute SQL queries and transactions
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
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rollbackErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	ToEntry     Entry    `json:"to_entry"`
	FromEntry   Entry    `json:"from_entry"`
}

// TransferTx performs money transfer from one account to another within a single transaction operation
// Steps: 1) create transfer record, 2) add account entries, 3) Update each account's balance
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Check for overdraft balance
		fromAccount, _, err := lockTwoAccounts(ctx, q, arg.FromAccountID, arg.ToAccountID)

		// fromAccount, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}
		if fromAccount.Balance < arg.Amount {
			return fmt.Errorf("insufficient funds")
		}

		// Step 1)
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// Step 2)
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// always start AddAccountBalance in the same order (smaller id first) to avoid exclusive lock deadlock
		// if account1 row is locked and account2 row is locked from two transfers or more, deadlock will occur.
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

func addMoney(q *Queries, accountId1 int64, amount1 int64, accountId2 int64, amount2 int64) (Account, Account, error) {
	acc1, err := q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		ID:     accountId1,
		Amount: amount1,
	})
	if err != nil {
		return Account{}, Account{}, err
	}

	acc2, err := q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		ID:     accountId2,
		Amount: amount2,
	})
	if err != nil {
		return Account{}, Account{}, err
	}

	return acc1, acc2, nil
}

func lockTwoAccounts(ctx context.Context, q *Queries, id1, id2 int64) (Account, Account, error) {
	if id1 < id2 {
		acc1, err := q.GetAccountForUpdate(ctx, id1)
		if err != nil {
			return Account{}, Account{}, err
		}
		acc2, err := q.GetAccountForUpdate(ctx, id2)
		if err != nil {
			return Account{}, Account{}, err
		}
		return acc1, acc2, nil
	}

	acc2, err := q.GetAccountForUpdate(ctx, id2)
	if err != nil {
		return Account{}, Account{}, err
	}
	acc1, err := q.GetAccountForUpdate(ctx, id1)
	if err != nil {
		return Account{}, Account{}, err
	}
	return acc1, acc2, nil
}
