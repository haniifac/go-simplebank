package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/haniifac/simplebank/db/sqlc"
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

	// router
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts", server.listAccounts)
	router.GET("/accounts/:id", server.getAccount)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	err := server.router.Run(address)
	return err
}

func errResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
