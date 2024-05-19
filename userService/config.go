package main

import "os"

type AppConfig struct {
	RabbitMQHostname string
}

func InitConfig() AppConfig {

	rabbitMQHostname := os.Getenv("RABBITMQ_HOSTNAME")
	if rabbitMQHostname == "" {
		rabbitMQHostname = "localhost"
	}

	return AppConfig{
		RabbitMQHostname: rabbitMQHostname,
	}
}
