package api

import (
	"github.com/MeganViga/SimpleBank2/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()
	router.POST("/users", server.createAccount)
	router.GET("getuser/:id", server.getAccount)
	router.POST("/transfers", server.createTransfer)
	server.router = router
	return server
}

func (s *Server) StartServer(address string) error {
	return s.router.Run(address)
}
