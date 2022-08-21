package api

import (
	db "github/dutt23/bank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

// Returns a new instance of server
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Add routes
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	server.router = router
	return server
}

// Start the input server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
