package main

import (
	"log"
	"os"

	"goFinalProject/pkg/db"
	"goFinalProject/pkg/server"
)

const defaultDBFile = "scheduler.db"

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defaultDBFile
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
