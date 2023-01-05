package gapi

import (
	"fmt"
	db "github/dutt23/bank/db/sqlc"
	"github/dutt23/bank/pb"
	"github/dutt23/bank/token"
	"github/dutt23/bank/util"
	"github/dutt23/bank/worker"
)

type Server struct {
	pb.UnimplementedBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// Returns a new instance of server
func NewServer(config util.Config, store db.Store, distributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %w", err)
	}

	server := &Server{store: store, tokenMaker: tokenMaker, config: config, taskDistributor: distributor}
	return server, nil
}
