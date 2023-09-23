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

	router.POST("/users", server.createUser)
	router.POST("/receipts/process/:userid", server.processReceipt)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
