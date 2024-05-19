package main

import (
	"time"
)

// http port
const PORT = ":3002"

// grpc port
const GRPCPORT = ":4002"

// rabbitmq port
const RABBITMQPORT = ":5672"

func main() {

	sqliteStorage := NewSqliteStorage()
	sqliteStorage.Init()

	// set db conn limit
	sqliteStorage.db.SetMaxOpenConns(25)
	sqliteStorage.db.SetMaxIdleConns(25)
	sqliteStorage.db.SetConnMaxLifetime(5 * time.Minute)

	//gRPC server :4002
	go func() {
		server := NewGrpcServer(GRPCPORT, sqliteStorage)
		server.RunGrpc()
	}()

	//http server :3002
	go func() {
		server := NewServer(PORT, sqliteStorage)
		server.Run()
	}()

	select {}
}
