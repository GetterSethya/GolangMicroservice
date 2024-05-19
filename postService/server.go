package main

import (
	"log"
	"net/http"

	"github.com/GetterSethya/library"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
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

func (s *Server) Run(grpcClient *grpc.ClientConn) {

	c := library.NewUserClient(grpcClient)
	routes := mux.NewRouter().PathPrefix("/v1/post").Subrouter()

	userService := NewUserService(s.Store, c)
	userService.RegisterRoutes(routes)

	log.Println("postService is running on port:", PORT)
	log.Fatal(http.ListenAndServe(s.ListenAddr, routes))
}
