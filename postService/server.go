package main

import (
	"google.golang.org/grpc/resolver"
	"log"
	"net/http"

	"github.com/GetterSethya/library"
	"github.com/gorilla/mux"
)

type Server struct {
	ListenAddr string
	Store      *SqliteStorage
	Cfg        AppConfig
}

func NewServer(listenAddr string, store *SqliteStorage, cfg AppConfig) *Server {
	return &Server{
		ListenAddr: listenAddr,
		Store:      store,
		Cfg:        cfg,
	}
}

func (s *Server) Run() {

	rb := &exampleResolverBuilder{
		UserServiceHostname: s.Cfg.UserServiceHostName,
	}

	// dial grpc user service
	resolver.Register(rb)
	conn, err := generateGrpcConn(s.Cfg.UserServiceHostName)
	if err != nil {
		log.Fatalf("Cannot connect to Grpc server:%v", err)
	}

	c := library.NewUserClient(conn)
	routes := mux.NewRouter().PathPrefix("/v1/post").Subrouter()

	userService := NewUserService(s.Store, c)
	userService.RegisterRoutes(routes)

	log.Println("postService is running on port:", PORT)
	log.Fatal(http.ListenAndServe(s.ListenAddr, routes))
}
