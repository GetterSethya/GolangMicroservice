package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GetterSethya/library"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

type AppServer struct {
	Store    *SqliteStorage
	Cfg      AppConfig
	Server   http.Server
	AmqpConn *amqp.Connection
}

func NewServer(listenAddr string, store *SqliteStorage, cfg AppConfig) *AppServer {

	routes := mux.NewRouter().PathPrefix("/v1/user").Subrouter()

	//rabbitmq conn
	amqpConnString := fmt.Sprintf("amqp://guest:guest@%s%s/", cfg.RabbitMQHostname, RABBITMQPORT)
	conn, err := amqp.Dial(amqpConnString)
	if err != nil {
		log.Println("Error when creating rabbitMq connection:", err)
	}
	// defer conn.Close()

	rabbitMQ := library.NewRabbitMq(conn)
	userService := NewUserService(store, rabbitMQ)
	userService.RegisterRoutes(routes)

	return &AppServer{
		Store: store,
		Cfg:   cfg,
		Server: http.Server{
			Addr:    listenAddr,
			Handler: routes,
		},
		AmqpConn: conn,
	}
}


func (s *AppServer) Run() {

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()
	log.Println("userService is running on port:", PORT)

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return s.Server.ListenAndServe()
	})

	g.Go(func() error {
		<-gctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer shutdownCancel()

		if err := s.AmqpConn.Close(); err != nil {
			log.Println("Error when closing rabbitmq connection:", err)
		}

		log.Println("SIGTERM detected, will attempt to graceful shutdown...")
		return s.Server.Shutdown(shutdownCtx)
	})

	if err := g.Wait(); err != nil {
		log.Println(err)
	}

}
