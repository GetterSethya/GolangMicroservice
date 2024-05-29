package main

import (
	"log"
	"os"
)

type AppConfig struct {
	Host         string
	Port         string
	InternalPort string
}

func InitConfig() AppConfig {
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	internalPort := os.Getenv("INTERNAL_PORT")

	if internalPort == "" {
		log.Println("INTERNAL_PORT environment variable is missing, fallback to :3001")
		internalPort = ":3001"
	}

	if port == "" {
		log.Println("PORT environment variable is missing, fallback to :3001")
		port = ":3001"
	}

	if host == "" {
		log.Println("HOST environment variable is missing, fallback to localhost")
		host = "localhost"
	}

	return AppConfig{
		Host:         host,
		Port:         port,
		InternalPort: internalPort,
	}
}
