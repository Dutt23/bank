package main

import (
	"database/sql"
	"github/dutt23/bank/api"
	db "github/dutt23/bank/db/sqlc"
	"github/dutt23/bank/util"
	"log"

	_ "github.com/lib/pq"
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
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start sever ", err)
	}
}
