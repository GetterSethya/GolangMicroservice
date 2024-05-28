package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GetterSethya/library"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
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
	log.Println("userService is running on port:", PORT)
	s.Server.ListenAndServe()
}
