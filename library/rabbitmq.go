package library

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMq struct {
	Conn *amqp.Connection
}

func NewRabbitMq(conn *amqp.Connection) *RabbitMq {
	return &RabbitMq{
		Conn: conn,
	}
}

