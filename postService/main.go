package main

import (
	"sync"
	"time"

)

// http port
const PORT = ":3004"

// grpc port
const GRPC_USER_SERVICE_PORT = ":4002"
const GRPC_NUM_INSTANCE = 2

const exampleScheme = "example"
const exampleServiceName = "user-service"

// rabbitmq port
const RABBITMQ_PORT = ":5672"

func main() {

	var wg sync.WaitGroup

	cfg := InitConfig()

	sqliteStorage := NewSqliteStorage()
	sqliteStorage.Init()

	// set db conn limit
	sqliteStorage.db.SetMaxOpenConns(25)
	sqliteStorage.db.SetMaxIdleConns(25)
	sqliteStorage.db.SetConnMaxLifetime(5 * time.Minute)

	server := NewServer(PORT, sqliteStorage, cfg)

	wg.Add(1)
	go func() {
		server.Run()
		defer wg.Done()
	}()

	rabbitMq := NewRabbitMQ(cfg, sqliteStorage)
	rabbitMq.Run()

	//biar main func tidak exit duluan
	wg.Wait()
}
