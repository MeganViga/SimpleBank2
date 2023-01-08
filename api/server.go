package api

import (

	"github.com/MeganViga/SimpleBank2/db/sqlc"

	"github.com/gin-gonic/gin"
)


type Server struct{
	store db.Store
	router *gin.Engine
}

func NewServer(store db.Store)*Server{
	server :=  &Server{
		store :store,
	}
	router := gin.Default()
	router.POST("/users",server.createUser)
	server.router = router
	return server
}

func (s *Server)StartServer(address string)error{
	return s.router.Run(address)
}