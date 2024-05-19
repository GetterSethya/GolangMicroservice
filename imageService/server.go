package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	ListenAddr string
	Cfg        AppConfig
}

func NewServer(cfg AppConfig) *Server {
	return &Server{
		Cfg: cfg,
	}
}

func (s *Server) Run() {

	router := mux.NewRouter().PathPrefix("/v1/image").Subrouter()

	imageService := NewImageService(s.Cfg.Host, s.Cfg.Port)
	imageService.RegisterRoutes(router)

	log.Println("Image service is running on port:", s.Cfg.Port)
	log.Fatal(http.ListenAndServe(s.ListenAddr, router))
}
