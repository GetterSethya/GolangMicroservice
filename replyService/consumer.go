package main

import (
	"encoding/json"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	Conn  *amqp.Connection
	Store *SqliteStorage
}

func NewConsumer(conn *amqp.Connection, store *SqliteStorage) *Consumer {
	return &Consumer{
		Conn:  conn,
		Store: store,
	}
}

func (c *Consumer) Consume(rabbitMQHostname string) {
	var wg sync.WaitGroup

	ch, err := c.Conn.Channel()
	if err != nil {
		log.Println("Error when creating channel:", err)
	}

	defer ch.Close()

	err = ch.ExchangeDeclare(
		"userServiceExchange",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Error when declaring exchange:", err)
	}

	q, err := ch.QueueDeclare(
		"replyService_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Error when declaring queue:", err)
	}

	err = ch.QueueBind(
		q.Name,
		"user.detail.change",
		"userServiceExchange",
		false,
		nil,
	)
	if err != nil {
		log.Println("Error when binding queue:", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Error when consuming queue:", err)
	}

	wg.Add(1)
	go func() {
		for d := range msgs {
			// update data user
			log.Println("New event receive:", string(d.Body))
			c.handleUpdateUserDetail(d.Body)
		}
		defer wg.Done()
	}()
	wg.Wait()
}

func (c *Consumer) handleUpdateUserDetail(data []byte) {
	type userDetailChangeEvent struct {
		Id      string `json:"id"`
		Name    string `json:"name"`
		Profile string `json:"profile"`
	}

	newUserData := &userDetailChangeEvent{}

	err := json.Unmarshal(data, newUserData)
	if err != nil {
		log.Println("Error when unmarshaling user data:", err)
	}

	log.Print("new user data", newUserData.Id, newUserData.Profile, newUserData.Name)

	if err := c.Store.UpdateUserDetail(newUserData.Id, newUserData.Profile, newUserData.Name); err != nil {
		log.Println("Error when updating post name:", err)
	}
}
