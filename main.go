package main

import (
	"context"
	"database/sql"
	"github/dutt23/bank/api"
	db "github/dutt23/bank/db/sqlc"
	"github/dutt23/bank/gapi"
	"github/dutt23/bank/pb"
	"github/dutt23/bank/util"
	"github/dutt23/bank/worker"
	"net"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"

	_ "github/dutt23/bank/docs/statik"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {

	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal().AnErr("Cannot load config files", err)
	}

	if config.ENV == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal().AnErr("cannot connect to database", err)
	}

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	runMigrations(config.MigrationURL, config.DBSource)
	go runTaskProcessors(redisOpt, store)
	// runGinServer(config, store)
	go runGatewayServer(config, store, taskDistributor)
	gprcServer(config, store, taskDistributor)
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal().AnErr("Cannot start server", err)
	}
	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal().AnErr("Cannot start sever ", err)
	}
}

func gprcServer(config util.Config, store db.Store, distributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, distributor)
	if err != nil {
		log.Fatal().AnErr("Cannot start sever ", err)
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)

	if err != nil {
		log.Fatal().AnErr("Cannot create listener", err)
	}

	log.Printf("starting grpc server at %s.", listener.Addr().String())

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal().AnErr("Could not start server", err)
	}
}

func runGatewayServer(config util.Config, store db.Store, distributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, distributor)
	if err != nil {
		log.Fatal().AnErr("Cannot start sever ", err)
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
		log.Fatal().AnErr("Cannot register handler server", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// fs := http.FileServer(http.Dir("./docs/swagger"))

	statikFs, err := fs.New()

	if err != nil {
		log.Fatal().AnErr("Cannot serve static file system", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFs))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.ServerAddress)

	if err != nil {
		log.Fatal().AnErr("Cannot create listener", err)
	}

	log.Printf("starting Http gateway server server at %s.", listener.Addr().String())

	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)

	if err != nil {
		log.Fatal().AnErr("cannot start http gateway server", err)
	}
}

func runMigrations(migrationURL, dbSource string) {
	m, err := migrate.New(migrationURL, dbSource)

	if err != nil {
		log.Fatal().AnErr("Cannot create new migrate instance", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().AnErr("Cannot run migrate up on the instance", err)
	}

	log.Info().Msg("database migration successful")
}

func runTaskProcessors(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Info().Msg("start task processor")

	if err := taskProcessor.Start(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start  task processor")
	}
}
