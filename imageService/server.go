package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type AppServer struct {
	Cfg    AppConfig
	Server http.Server
}

func NewAppServer(cfg AppConfig) *AppServer {
	router := mux.NewRouter().PathPrefix("/v1/image").Subrouter()
	imageService := NewImageService(cfg)
	imageService.RegisterRoutes(router)

	return &AppServer{
		Cfg: cfg,
		Server: http.Server{
			Addr:    cfg.InternalPort,
			Handler: router,
		},
	}
}

func (s *AppServer) Run() {
	log.Println("Image service is running on port:", s.Cfg.InternalPort)
	s.Server.ListenAndServe()
}
