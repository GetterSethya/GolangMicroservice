package main

import (
	"time"
)

const PORT = ":3002"

func main() {

	sqliteStorage := NewSqliteStorage()
	sqliteStorage.Init()

	// set db conn limit
	sqliteStorage.db.SetMaxOpenConns(25)
	sqliteStorage.db.SetMaxIdleConns(25)
	sqliteStorage.db.SetConnMaxLifetime(5 * time.Minute)

	server := NewServer(PORT, sqliteStorage)
    server.Run()
}

