package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	// http port
	PORT = ":3006"

	// grpc client port
	GRPC_USER_SERVICE_PORT         = ":4002"
	GRPC_USER_SERVICE_NUM_INSTANCE = 2

	USER_SCHEME       = "user"
	USER_SERVICE_NAME = "user-service"

	// rabbitmq port
	RABBITMQ_PORT = ":5672"
)

func main() {
	var wg sync.WaitGroup

	cfg := InitConfig()

	sqliteStorage := NewSqliteStorage()
	sqliteStorage.Init()

	// set db conn limit
	sqliteStorage.db.SetMaxOpenConns(25)
	sqliteStorage.db.SetMaxIdleConns(25)
	sqliteStorage.db.SetConnMaxLifetime(5 * time.Minute)

	// http server
	httpServer := NewServer(PORT, sqliteStorage, cfg)

	wg.Add(1)
	go func() {
		defer wg.Done()
		httpServer.Run()
	}()

	// rabbitmq consumer
	rabbitMq := NewRabbitMQ(cfg, sqliteStorage)
	go rabbitMq.Run()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	<-sigs
	log.Println("SIGTERM detected, will attempt to graceful shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// shutdown http.server
	if err := httpServer.Server.Shutdown(shutdownCtx); err != nil {
		log.Println("Error when trying to shutdown http server:", err)
	} else {
		log.Println("http server closed")
	}

	// shutdown rabbitmq
	rabbitMq.Close()

	wg.Wait()
}
