package main

import (
	"database/sql"
	"log"

	"github.com/dpsigor/cheatsheet-golang-postgres/api"
	db "github.com/dpsigor/cheatsheet-golang-postgres/db/sqlc"
	"github.com/dpsigor/cheatsheet-golang-postgres/util"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(cfg, store)
	if err != nil {
		log.Fatalf("failed to instantiate server: %v", err)
	}

	err = server.Start(cfg.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
