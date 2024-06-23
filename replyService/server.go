package main

import (
	"log"
	"net/http"

	"github.com/GetterSethya/postProto"
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

	postRB := &PostServiceResolverBuilder{
		PostServiceHostname: cfg.PostServiceHostName,
	}

	resolver.Register(userRB)
	userServiceGrpcConn, err := generateUserServiceGrpcConn(cfg.UserServiceHostName)
	if err != nil {
		log.Fatalf("Cannot connect to user Grpc server: %v", err)
	} else {
		log.Println("create connection to user Grpc server:", userServiceGrpcConn.Target())
	}

	resolver.Register(postRB)
	postServiceGrpcConn, err := generatePostServiceGrpcConn(cfg.PostServiceHostName)
	if err != nil {
		log.Fatalf("Cannot connect to post Grpc server: %v", err)
	} else {
		log.Println("create connection to post Grpc server:", postServiceGrpcConn.Target(), postServiceGrpcConn.GetState().String())
	}

	userServiceGrpcClient := userProto.NewUserClient(userServiceGrpcConn)
	postServiceGrpcClient := postProto.NewPostClient(postServiceGrpcConn)
	routes := mux.NewRouter().PathPrefix("/v1/reply").Subrouter()

	replyService := NewReplyService(store, userServiceGrpcClient, postServiceGrpcClient)
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
