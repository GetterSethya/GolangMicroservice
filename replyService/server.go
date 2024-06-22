package main

import (
	"log"
	"net/http"

	"github.com/GetterSethya/userProto"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/resolver"
)

type AppServer struct {
	Store  *SqliteStorage
	Cfg    AppConfig
	Server http.Server
}

func NewServer(listenAddr string, store *SqliteStorage, cfg AppConfig) *AppServer {

	userRB := &UserServiceResolverBuilder{
		UserServiceHostname: cfg.UserServiceHostName,
	}

	resolver.Register(userRB)
	userServiceGrpcConn, err := generateUserServiceGrpcConn(cfg.UserServiceHostName)
	if err != nil {
		log.Fatalf("Cannot connect to user Grpc server: %v", err)
	}

	userServiceGrpcClient := userProto.NewUserClient(userServiceGrpcConn)
	routes := mux.NewRouter().PathPrefix("/v1/reply").Subrouter()

	replyService := NewReplyService(store, userServiceGrpcClient)
	replyService.RegisterRoutes(routes)

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
	log.Println("reply service is running on port:", PORT)
	s.Server.ListenAndServe()
}
