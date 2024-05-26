package main

import (
	"log"
	"net/http"
	"google.golang.org/grpc/resolver"

	"github.com/GetterSethya/library"
	"github.com/gorilla/mux"
)

type AppServer struct {
	Store  *SqliteStorage
	Cfg    AppConfig
	Server http.Server
}

func NewServer(listenAddr string, store *SqliteStorage, cfg AppConfig) *AppServer {

	rb := &exampleResolverBuilder{
		UserServiceHostname: cfg.UserServiceHostName,
	}

	// dial grpc user service
	resolver.Register(rb)
	conn, err := generateGrpcConn(cfg.UserServiceHostName)
	if err != nil {
		log.Fatalf("Cannot connect to Grpc server:%v", err)
	}

	c := library.NewUserClient(conn)
	routes := mux.NewRouter().PathPrefix("/v1/post").Subrouter()

	userService := NewUserService(store, c)
	userService.RegisterRoutes(routes)
	return &AppServer{
		Store: store,
		Cfg:   cfg,
		Server: http.Server{
			Addr:    listenAddr,
			Handler: routes,
		},
	}
}

func (s *AppServer) Run() {
	log.Println("postService is running on port:", PORT)
	s.Server.ListenAndServe()
}
