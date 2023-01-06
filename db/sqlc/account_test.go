package db

import (
	"context"

	"testing"

	"github.com/MeganViga/SimpleBank2/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T)Account{
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  int64(util.RandomMoney()),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	return account
}
func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}
