package main

import (
	"log"
	"os"
)

type AppConfig struct {
	JwtSecret           string
	RefreshSecret       string
	UserServiceHostname string
}

func InitConfig() AppConfig {

	jwtSecret := os.Getenv("JWT_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")
	userServiceHostName := os.Getenv("USER_SERVICE_HOSTNAME")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET key not found!")
	}

	if refreshSecret == "" {
		log.Fatal("REFRESH_SECRET key not found!")
	}

	if userServiceHostName == "" {
		log.Println("USER_SERVICE_HOSTNAME key is not found, fallback to 'localhost'")
		userServiceHostName = "localhost"
	}

	return AppConfig{
		JwtSecret:           jwtSecret,
		RefreshSecret:       refreshSecret,
		UserServiceHostname: userServiceHostName,
	}
}
