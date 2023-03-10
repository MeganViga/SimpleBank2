package api

import (
	"fmt"
	"net/http"

	db "github.com/MeganViga/SimpleBank2/db/sqlc"
	"github.com/gin-gonic/gin"
)

type transferRequestParams struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,oneof=USD INR SGP"`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	var req transferRequestParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if !s.validAccount(ctx, req.FromAccountID,req.Currency){
		return
	}
	if !s.validAccount(ctx, req.ToAccountID,req.Currency){
		return
	}
	fmt.Println("Request Func", req)
	arg := db.TransferTXParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := s.store.TransferTX(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, result)

}

func(s *Server)validAccount(ctx *gin.Context,accountId int64, currency string)bool{
	account, err := s.store.GetAccount(ctx,accountId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return false
	}
	if account.Currency != currency{
		err := fmt.Errorf("account [%d]currency mismatch %s vs %s",accountId,account.Currency,currency)
		//fmt.Printf(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return false
	}
	return true
}
