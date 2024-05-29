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

type AppImage struct {
	Thumbnail string `json:"thumbnail"`
	Original  string `json:"original"`
	Filename  string `json:"filename"`
}

// grpc port
const GRPCPORT = ":4001"

func main() {

	var wg sync.WaitGroup

	cfg := InitConfig()

	httpServer := NewAppServer(cfg)
	grpcServer := NewGrpcServer(cfg, GRPCPORT)

	wg.Add(1)
	go func() {
		defer wg.Done()
		httpServer.Run()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		grpcServer.RunGrpc()
	}()

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

	// shutdown grpc server
	func() {
		grpcServer.Server.GracefulStop()
		log.Println("GRPC server closed")
	}()

	wg.Wait()
}
