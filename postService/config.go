package main

import (
	"log"
	"os"
)

type AppConfig struct {
	UserServiceHostName  string
	ImageServiceHostName string
	RabbitMQHostname     string
}

func InitConfig() AppConfig {

	userServiceHostName := os.Getenv("USER_SERVICE_HOSTNAME")
	if userServiceHostName == "" {
		log.Println("USER_SERVICE_HOSTNAME key is not found, fallback to 'localhost'")
		userServiceHostName = "localhost"
	}

	imageServiceHostName := os.Getenv("IMAGE_SERVICE_HOSTNAME")
	if imageServiceHostName == "" {
		log.Println("IMAGE_SERVICE_HOSTNAME key is not found, fallback to 'localhost'")
		imageServiceHostName = "localhost"
	}

	rabbitMQHostname := os.Getenv("RABBITMQ_HOSTNAME")
	if rabbitMQHostname == "" {
		log.Println("RABBITMQ_HOSTNAME is not found, fallback to localhost")
		rabbitMQHostname = "localhost"
	}

	return AppConfig{
		UserServiceHostName:  userServiceHostName,
		ImageServiceHostName: imageServiceHostName,
		RabbitMQHostname:     rabbitMQHostname,
	}
}
