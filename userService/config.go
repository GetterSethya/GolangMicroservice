package main

import (
	"log"
	"os"
)

type AppConfig struct {
	RabbitMQHostname     string
	ImageServiceHostName string
}

func InitConfig() AppConfig {
	rabbitMQHostname := os.Getenv("RABBITMQ_HOSTNAME")
	if rabbitMQHostname == "" {
		rabbitMQHostname = "localhost"
	}

	imageServiceHostName := os.Getenv("IMAGE_SERVICE_HOSTNAME")
	if imageServiceHostName == "" {
		log.Println("IMAGE_SERVICE_HOSTNAME env key is missing, fallback to localhost")
		imageServiceHostName = "localhost"
	}

	return AppConfig{
		RabbitMQHostname:     rabbitMQHostname,
		ImageServiceHostName: imageServiceHostName,
	}
}
