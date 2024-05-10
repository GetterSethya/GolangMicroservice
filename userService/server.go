package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	ListenAddr string
	Store      *SqliteStorage
}


func NewServer(listenAddr string, store *SqliteStorage) *Server {
	return &Server{
		ListenAddr: listenAddr,
		Store:      store,
	}
}

func (s *Server) Run() {

	routes := mux.NewRouter().PathPrefix("/v1/user").Subrouter()

	userService := NewUserService(s.Store)
	userService.RegisterRoutes(routes)

	log.Println("userService is running on port:", PORT)
	log.Fatal(http.ListenAndServe(s.ListenAddr, routes))
}
