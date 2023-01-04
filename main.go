package main

import (
	"database/sql"
	"github/dutt23/bank/api"
	db "github/dutt23/bank/db/sqlc"
	"github/dutt23/bank/gapi"
	"github/dutt23/bank/pb"
	"github/dutt23/bank/util"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// const (
// 	dbDriver      = "postgres"
// 	dbSource      = "postgresql://root:secret@localhost:5430/bank?sslmode=disable"
// 	serverAddress = "localhost:8080"
// )

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatalln("Cannot load config files", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	store := db.NewStore(conn)
	// runGinServer(config, store)
	gprcServer(config, store)
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("Cannot start server", err)
	}
	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start sever ", err)
	}
}

func gprcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot start sever ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)

	if err != nil {
		log.Fatal("Cannot create listener", err)
	}

	log.Printf("starting grpc server at %s.", listener.Addr().String())

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("Could not start server", err)
	}

}
