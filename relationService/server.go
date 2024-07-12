package main

import (
	"log"
	"net/http"

	"google.golang.org/grpc/resolver"

	"github.com/GetterSethya/userProto"
	"github.com/gorilla/mux"
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

	// dial grpc user service
	resolver.Register(userRB)
	userServiceGrpcConn, err := generateUserServiceGrpcConn(cfg.UserServiceHostName)
	if err != nil {
		log.Fatalf("Cannot connect to user Grpc server: %v", err)
	}

	userGrpcClient := userProto.NewUserClient(userServiceGrpcConn)
	routes := mux.NewRouter().PathPrefix("/v1/relation").Subrouter()

	relationService := NewRelationService(store, userGrpcClient)
	relationService.RegisterRoutes(routes)

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
	log.Println("relationService is running on port:", PORT)
	s.Server.ListenAndServe()
}
