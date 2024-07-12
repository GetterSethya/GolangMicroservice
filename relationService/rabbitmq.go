package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Cfg      AppConfig
	Store    *SqliteStorage
	AmqpConn *amqp.Connection
}

func NewRabbitMQ(cfg AppConfig, store *SqliteStorage) *RabbitMQ {

	connString := fmt.Sprintf("amqp://guest:guest@%s%s/", cfg.RabbitMQHostname, RABBITMQ_PORT)
	rabbitMQConn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatalf("Error when creating connection to rabbit mq: %+v", err)
	} else {
		log.Println("Connected to rabbitmq")
	}

	return &RabbitMQ{
		Cfg:      cfg,
		Store:    store,
		AmqpConn: rabbitMQConn,
	}
}

func (r *RabbitMQ) Run() {
	consumer := NewConsumer(r.AmqpConn, r.Store)
	consumer.Consume(r.Cfg.RabbitMQHostname)
}

func (r *RabbitMQ) Close() {

	if err := r.AmqpConn.Close(); err != nil {
		log.Println("Error when closing rabbitMQ", err)
	} else {
		log.Println("rabbitMQ connection closed..")
	}

}
