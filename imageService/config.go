package main

import (
	"log"
	"os"
)

type AppConfig struct {
	Host string
	Port string
}

func InitConfig() AppConfig {
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	if port == "" {
		log.Println("PORT environment variable is missing, fallback to :3001")
		port = ":3001"
	}

	if host == "" {
		log.Println("HOST environment variable is missing, fallback to localhost")
		host = "localhost"
	}

	return AppConfig{
		Host: host,
		Port: port,
	}
}
