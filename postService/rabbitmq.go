package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Cfg   AppConfig
	Store *SqliteStorage
}

func NewRabbitMQ(cfg AppConfig, store *SqliteStorage) *RabbitMQ {
	return &RabbitMQ{
		Cfg:   cfg,
		Store: store,
	}
}

func (r *RabbitMQ) Run() {

	connString := fmt.Sprintf("amqp://guest:guest@%s%s/", r.Cfg.RabbitMQHostname, RABBITMQ_PORT)
	rabbitMQConn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatalf("Error when creating connection to rabbit mq: %+v", err)
	}

	defer rabbitMQConn.Close()

	consumer := NewConsumer(rabbitMQConn, r.Store)
	consumer.Consume(r.Cfg.RabbitMQHostname)
}
