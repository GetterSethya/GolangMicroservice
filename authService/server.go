package main

import (
	"log"
	"net/http"

	"github.com/GetterSethya/library"
	"github.com/gorilla/mux"
)

type Server struct {
	ListenAddr    string
	JwtSecret     string
	RefreshSecret string
}

func NewServer(listenAddr string, jwtSecret string, refreshSecret string) *Server {
	return &Server{
		ListenAddr:    listenAddr,
		JwtSecret:     jwtSecret,
		RefreshSecret: refreshSecret,
	}
}

func (s *Server) Run(grpcClient library.UserClient) {

	routes := mux.NewRouter().PathPrefix("/v1/auth").Subrouter()

	userService := NewAuthService(s.JwtSecret, grpcClient, s.RefreshSecret)
	userService.RegisterRoutes(routes)

	log.Println("authService is running on port:", PORT)
	log.Fatal(http.ListenAndServe(s.ListenAddr, routes))
}
