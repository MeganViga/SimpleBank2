package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Provide all queries to execute DB queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}
var txKey = struct{}{}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("tx err: %v, rbErr: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTXParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}
type TransferTXResult struct {
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	Transfer    Transfer `json:"transfer"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}
// TransferTX perform a money transfer from one account to another
//Creates a transfer record, from entry, to entry,update account's balance within a single database transaction
func (store *Store) TransferTX(ctx context.Context, arg TransferTXParams)(TransferTXResult,error) {
	var result TransferTXResult

	err := store.execTx(ctx, func(q *Queries)error{
		var err error
		txName := ctx.Value(txKey)
		fmt.Println(txName,"create transfer")
		result.Transfer, err = q.CreateTransfer(ctx,CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,

		})
		if err != nil{
			return err
		}
		fmt.Println(txName,"create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil{
			return err
		}
		fmt.Println(txName,"create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}
		//TODO: Update Account Balance
		// fmt.Println(txName,"get account 1")
		// account1 , err := q.GetAccountForUpdate(ctx,arg.FromAccountID)
		// if err != nil{
		// 	return err
		// }

		//below if condition to avoid deadlock
		if arg.FromAccountID < arg.ToAccountID{
			result.FromAccount,result.ToAccount ,err =  addMoney(ctx,q,arg.FromAccountID,-arg.Amount,arg.ToAccountID,arg.Amount)
		// 	fmt.Println(txName,"update account 1")
		// result.FromAccount, err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
		// 	ID: arg.FromAccountID,
		// 	Amount: -arg.Amount,
		// })
		if err != nil{
			return err
		}
		// // fmt.Println(txName,"get account 1")
		// // account2 , err := q.GetAccountForUpdate(ctx,arg.ToAccountID)
		// // if err != nil{
		// // 	return err
		// // }
		// fmt.Println(txName,"update account 2")
		// result.ToAccount, err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
		// 	ID: arg.ToAccountID,
		// 	Amount: arg.Amount,
		// })
		// if err != nil{
		// 	return err
		// }
		}else{
			result.ToAccount,result.FromAccount ,err =  addMoney(ctx,q,arg.ToAccountID,arg.Amount,arg.FromAccountID,-arg.Amount)

			// result.ToAccount, err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
			// 	ID: arg.ToAccountID,
			// 	Amount: arg.Amount,
			// })
			if err != nil{
				return err
			}

			// result.FromAccount, err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
			// 	ID: arg.FromAccountID,
			// 	Amount: -arg.Amount,
			// })
			// if err != nil{
			// 	return err
			// }

		}
		
		return nil 
	})
	return result, err
}

func addMoney(ctx context.Context,q *Queries,account1ID int64, amount1 int64,account2ID int64,amount2 int64)(account1 Account,account2 Account, err error){
	account1, err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
		ID: account1ID,
		Amount: amount1,
	})
	if err != nil{
		return 
	}

	account2, err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
		ID: account2ID,
		Amount: amount2,
	})
	if err != nil{
		return
	}
	return
}
