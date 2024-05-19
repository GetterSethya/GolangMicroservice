package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GetterSethya/library"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	ListenAddr string
	Store      *SqliteStorage
	Cfg        AppConfig
}

func NewServer(listenAddr string, store *SqliteStorage, cfg AppConfig) *Server {
	return &Server{
		ListenAddr: listenAddr,
		Store:      store,
		Cfg:        cfg,
	}
}

func (s *Server) Run() {

	routes := mux.NewRouter().PathPrefix("/v1/user").Subrouter()

	//rabbitmq conn
	amqpConnString := fmt.Sprintf("amqp://guest:guest@%s%s/", s.Cfg.RabbitMQHostname, RABBITMQPORT)
	conn, err := amqp.Dial(amqpConnString)
	if err != nil {
		log.Println("Error when creating rabbitMq connection:", err)
	}
	defer conn.Close()

	rabbitMQ := library.NewRabbitMq(conn)
	userService := NewUserService(s.Store, rabbitMQ)
	userService.RegisterRoutes(routes)

	log.Println("userService is running on port:", PORT)
	log.Fatal(http.ListenAndServe(s.ListenAddr, routes))
}
