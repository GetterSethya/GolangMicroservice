package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GetterSethya/library"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	ListenAddr string
	Store      *SqliteStorage
}

func NewServer(listenAddr string, store *SqliteStorage) *Server {
	return &Server{
		ListenAddr: listenAddr,
		Store:      store,
	}
}

func (s *Server) Run() {

	routes := mux.NewRouter().PathPrefix("/v1/user").Subrouter()

	rabbitMQHostname := os.Getenv("RABBITMQ_HOSTNAME")
	if rabbitMQHostname == "" {
		rabbitMQHostname = "localhost"
	}

	//rabbitmq conn
	// amqp://guest:guest@localhost:5672/
	amqpConnString := fmt.Sprintf("amqp://guest:guest@%s%s/", rabbitMQHostname, RABBITMQPORT)
	log.Println(amqpConnString)
	conn, err := amqp.Dial(amqpConnString)
	if err != nil {
		log.Println("Error when creating rabbitMq connection:", err)
	}

	rabbitMQ := library.NewRabbitMq(conn)
	userService := NewUserService(s.Store, rabbitMQ)
	userService.RegisterRoutes(routes)

	log.Println("userService is running on port:", PORT)
	log.Fatal(http.ListenAndServe(s.ListenAddr, routes))
}
