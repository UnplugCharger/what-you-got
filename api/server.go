package api

import (
	"github.com/gin-gonic/gin"

	db "hackathon/db/sqlc"
	"hackathon/utils"
)

type Server struct {
	store  db.Store
	router *gin.Engine
	config utils.Config
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {

	server := &Server{store: store, config: config}

	server.setupRouter()

	return server, nil

}

func (server *Server) setupRouter() {
	router := gin.Default()

	//router.POST("/users", server.createUser)
	//router.POST("/users/login", server.loginUser)
	//router.POST("/tokens/renew_access", server.renewAccessToken)
	//
	//authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	//authRoutes.POST("/accounts", server.createAccount)
	//authRoutes.GET("/accounts/:id", server.getAccount)
	//authRoutes.GET("/accounts", server.listAccounts)
	//
	//authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
