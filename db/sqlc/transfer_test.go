package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        15,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	require.Equal(t, account1.ID, transfer.FromAccountID)
	require.Equal(t, account2.ID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        15,
	}

	tf, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)

	transfer, err := testQueries.GetTransfer(context.Background(), tf.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	require.Equal(t, account1.ID, transfer.FromAccountID)
	require.Equal(t, account2.ID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	argCreate := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        10,
	}

	for i := 0; i < 5; i++ {
		_, err := testQueries.CreateTransfer(context.Background(), argCreate)
		require.NoError(t, err)
	}

	argList := ListTransfersParams{
		Limit:  5,
		Offset: 0,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), argList)
	require.NoError(t, err)
	require.Len(t, transfers, int(argList.Limit))
}

func TestListTransfersByAccountId(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	argCreate := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        10,
	}

	for i := 0; i < 5; i++ {
		_, err := testQueries.CreateTransfer(context.Background(), argCreate)
		require.NoError(t, err)
	}

	argList := ListTransfersByAccountIdParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Limit:         5,
		Offset:        0,
	}

	transfers, err := testQueries.ListTransfersByAccountId(context.Background(), argList)
	require.NoError(t, err)
	require.Len(t, transfers, int(argList.Limit))
}
