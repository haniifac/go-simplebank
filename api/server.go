package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/haniifac/simplebank/db/sqlc"
	"github.com/haniifac/simplebank/token"
	"github.com/haniifac/simplebank/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setRouter()
	return server, nil
}

func (server *Server) setRouter() {
	router := gin.Default()

	// router
	authGroup := router.Group("/", server.authMiddleware(server.tokenMaker))

	authGroup.POST("/accounts", server.createAccount)
	authGroup.GET("/accounts", server.listAccounts)
	authGroup.GET("/accounts/:id", server.getAccount)
	authGroup.POST("/transfers", server.createTransfer)

	authGroup.GET("users/:username", server.GetUser)

	router.POST("/users", server.CreateUser)
	router.POST("/users/login", server.loginUser)

	server.router = router
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
