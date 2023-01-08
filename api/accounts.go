package api

import (
	"fmt"
	"net/http"

	db "github.com/MeganViga/SimpleBank2/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountRequestParams struct{
	Owner string `json:"owner" binding:"required"`
	Balance int64  `json:"balance" binding:"required"`
	Curremcy string  `json:"currency" binding:"required"`
}
func (s *Server)createUser(ctx *gin.Context){
	var req createAccountRequestParams
	if err := ctx.ShouldBindJSON(&req); err != nil{
		ctx.JSON(http.StatusBadRequest,err)
		return
	}
	fmt.Println(req)
	arg := db.CreateAccountParams{
		Owner: req.Owner,
		Balance: req.Balance,
		Currency: req.Curremcy,
	}

	user, err := s.store.CreateAccount(ctx,arg)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError,err)
		return
	}
	ctx.JSON(http.StatusOK,user)

}
type getAccountRequestParams struct{
	Id int64 `uri:"id" binding:"required"`
}
func (s *Server)getAccount(ctx *gin.Context){
	var req getAccountRequestParams
	if err := ctx.ShouldBindUri(&req); err != nil{
		ctx.JSON(http.StatusBadRequest,err)
		return
	}
	//fmt.Println(req.Id)
	account, err := s.store.GetAccount(ctx,req.Id)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError,err)
	}
	ctx.JSON(http.StatusOK,account)
}