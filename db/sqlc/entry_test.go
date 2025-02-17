package db

import (
	"context"
	"testing"

	"github.com/haniifac/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := CreateEntryParams{
		AccountID: account1.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.NotZero(t, entry.ID)
	require.Equal(t, account1.ID, entry.AccountID)
	require.NotZero(t, entry.CreatedAt)

	fetchEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchEntry)
}

func TestListEntries(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := CreateEntryParams{
		AccountID: account1.ID,
		Amount:    util.RandomMoney(),
	}

	for i := 0; i < 5; i++ {
		entry, err := testQueries.CreateEntry(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, entry)
		arg.Amount = util.RandomMoney()
	}

	argList := ListEntriesParams{
		AccountID: account1.ID,
		Limit:     5,
		Offset:    0,
	}

	entries, err := testQueries.ListEntries(context.Background(), argList)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Len(t, entries, int(argList.Limit))

}
