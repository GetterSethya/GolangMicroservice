package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type AppServer struct {
	Cfg    AppConfig
	Server http.Server
}

func NewAppServer(cfg AppConfig) *AppServer {

	router := mux.NewRouter().PathPrefix("/v1/image").Subrouter()
	imageService := NewImageService(cfg.Host, cfg.Port)
	imageService.RegisterRoutes(router)

	return &AppServer{
		Cfg: cfg,
		Server: http.Server{
			Addr:    cfg.Port,
			Handler: router,
		},
	}
}

func (s *AppServer) Run() {
	log.Println("Image service is running on port:", s.Cfg.Port)
	s.Server.ListenAndServe()
}
