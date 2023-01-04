package main

import (
	"context"
	"database/sql"
	"github/dutt23/bank/api"
	db "github/dutt23/bank/db/sqlc"
	"github/dutt23/bank/gapi"
	"github/dutt23/bank/pb"
	"github/dutt23/bank/util"
	"log"
	"net"
	"net/http"

	_ "github/dutt23/bank/docs/statik"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

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
	go runGatewayServer(config, store)
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

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot start sever ", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterBankHandlerServer(ctx, grpcMux, server)

	if err != nil {
		log.Fatal("Cannot register handler server", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// fs := http.FileServer(http.Dir("./docs/swagger"))

	statikFs, err := fs.New()

	if err != nil {
		log.Fatal("Cannot serve static file system", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFs))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot create listener", err)
	}

	log.Printf("starting Http gateway server server at %s.", listener.Addr().String())

	err = http.Serve(listener, mux)

	if err != nil {
		log.Fatal("cannot start http gateway server", err)
	}
}
