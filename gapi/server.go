package gapi

import (
	"fmt"
	db "github/dutt23/bank/db/sqlc"
	"github/dutt23/bank/pb"
	"github/dutt23/bank/token"
	"github/dutt23/bank/util"
)

type Server struct {
	pb.UnimplementedBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// Returns a new instance of server
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %w", err)
	}

	server := &Server{store: store, tokenMaker: tokenMaker, config: config}
	return server, nil
}
