package api

import (
	"fmt"
	"net/http"

	db "github.com/MeganViga/SimpleBank2/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountRequest struct{
	Owner string `json:"owner" binding:"required"`
	Balance int64  `json:"balance" binding:"required"`
	Curremcy string  `json:"currency" binding:"required"`
}
func (s *Server)createUser(ctx *gin.Context){
	var req createAccountRequest
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