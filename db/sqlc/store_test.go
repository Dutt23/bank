package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTxn(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> Before : ", account1.Owner, account2.Owner)
	fmt.Println(">> Before : ", account1.Balance, account2.Balance)
	n := 4
	amount := int64(10)

	errs := make(chan error, n)
	results := make(chan TransferTxResult, n)
	defer close(errs)
	defer close(results)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.NotEmpty(t, result.ToEntry)
		require.NotEmpty(t, result.FromEntry)
		require.NotEmpty(t, result.ToAccountID)
		require.NotEmpty(t, result.FromAccountID)

		tr, err := store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		require.Equal(t, tr.ID, transfer.ID)

		toE := result.ToEntry
		fromE := result.FromEntry
		toEntry, err := store.GetEntry(context.Background(), toE.ID)
		require.NoError(t, err)
		require.Equal(t, toEntry.ID, toE.ID)
		require.Equal(t, toEntry.AccountID, toE.AccountID)

		fromEntry, err := store.GetEntry(context.Background(), fromE.ID)
		require.NoError(t, err)
		require.Equal(t, fromEntry.ID, fromE.ID)
		require.Equal(t, fromEntry.AccountID, fromE.AccountID)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		// check balances
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%10 == 0) // 1 * amount, 2 * amount, 3 * amount ...

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAcc1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAcc2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> Txx : ", updatedAcc1.Balance, updatedAcc2.Balance)
	require.Equal(t, account1.Balance-(int64(n)*amount), updatedAcc1.Balance)
	require.Equal(t, account2.Balance+(int64(n)*amount), updatedAcc2.Balance)
}

func TestTransferTxnDeadLock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> Before : ", account1.Balance, account2.Balance)
	n := 10
	amount := int64(10)

	errs := make(chan error, n)
	defer close(errs)
	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		ToAccountID := account2.ID
		if i%2 == 1 {
			fromAccountID = account2.ID
			ToAccountID = account1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: fromAccountID,
				ToAccountId:   ToAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedAcc1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAcc2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> Txx : ", updatedAcc1.Balance, updatedAcc2.Balance)
	require.Equal(t, account1.Balance, updatedAcc1.Balance)
	require.Equal(t, account2.Balance, updatedAcc2.Balance)
}
