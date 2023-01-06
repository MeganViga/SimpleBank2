package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"fmt"
)
// to debug deadlock in postgres --> https://wiki.postgresql.org/wiki/Lock_Monitoring
func TestTransferTX(t *testing.T){
	store := NewStore(testDB)
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println(">>>before tx: ",account1.Balance,account2.Balance)
	n := 5
	amount := int64(10)
	arg := TransferTXParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: amount,
	}
	errs := make(chan error)
	results := make(chan TransferTXResult)
	for i := 0; i< n;i++{
		txName := fmt.Sprintf("tx %d",i+1)
		go func(){
			ctx := context.WithValue(context.Background(),txKey,txName)
			result, err := store.TransferTX(ctx,arg)
			errs <- err
			results <- result

		}()
	}
	//check results
	existed := make(map[int]bool)
	for i :=0; i < n;i++{
		err := <-errs
		require.NoError(t,err)
		result := <-results

		//check Transfer
		transfer := result.Transfer
		require.NotEmpty(t,transfer)
		require.Equal(t,account1.ID,transfer.FromAccountID)
		require.Equal(t,account2.ID,transfer.ToAccountID)
		require.Equal(t,arg.Amount,transfer.Amount)
		require.NotZero(t,transfer.ID)
		require.NotZero(t,transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(),transfer.ID)
		require.NoError(t,err)

		//check From entry
		fromentry := result.FromEntry
		require.NotEmpty(t,fromentry)
		require.Equal(t,account1.ID,fromentry.AccountID)
		require.Equal(t,-arg.Amount,fromentry.Amount)
		require.NotZero(t,fromentry.ID)
		require.NotZero(t,fromentry.CreatedAt)

		// check to entry
		toentry := result.ToEntry
		require.NotEmpty(t,toentry)
		require.Equal(t,account2.ID,toentry.AccountID)
		require.Equal(t,arg.Amount,toentry.Amount)
		require.NotZero(t,toentry.ID)
		require.NotZero(t,toentry.CreatedAt)
		
		//check accounts
		fromaccount := result.FromAccount
		require.NotEmpty(t,fromaccount)
		require.Equal(t,account1.ID,fromaccount.ID)

		toaccount := result.ToAccount
		require.NotEmpty(t,toaccount)
		require.Equal(t,account2.ID,toaccount.ID)

		//check account balance
		fmt.Println(">>>tx: ",fromaccount.Balance,toaccount.Balance)
		diff1 := account1.Balance - fromaccount.Balance
		diff2 := toaccount.Balance - account2.Balance
		require.Equal(t,diff1,diff2)
		require.True(t,diff1 > 0)
		require.True(t,diff1% amount == 0)

		k := int(diff1/amount)
		require.True(t,k>=1 && k<=n)
		require.NotContains(t,existed,k)
		existed[k] = true


	}

	// check updated account balances
	updatedAccount1, err := store.GetAccount(context.Background(),account1.ID)
	require.NoError(t,err)

	updatedAccount2, err := store.GetAccount(context.Background(),account2.ID)
	require.NoError(t,err)
	fmt.Println(">>>after tx: ",updatedAccount1.Balance,updatedAccount2.Balance)

	require.Equal(t,account1.Balance - int64(n)*amount,updatedAccount1.Balance)
	require.Equal(t,account2.Balance + int64(n)*amount,updatedAccount2.Balance)


}