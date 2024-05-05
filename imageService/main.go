package main

import (
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
)

type AppImage struct {
	Thumbnail string `json:"thumbnail"`
	Original  string `json:"original"`
}

func main() {
	PORT := os.Getenv("PORT")
	HOST := os.Getenv("HOST")

	if PORT == "" {
		PORT = ":3001"
	}

	if HOST == "" {
		HOST = "localhost"
	}

	server := NewServer(PORT)

	router := mux.NewRouter().PathPrefix("/v1/image").Subrouter()

	imageService := NewImageService(HOST, PORT)
	imageService.RegisterRoutes(router)

	log.Println("Image service is running on port:", PORT)
	log.Fatal(http.ListenAndServe(server.ListenAddr, router))
}

