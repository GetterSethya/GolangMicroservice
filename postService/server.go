package main

import (
	"log"
	"net/http"

	"google.golang.org/grpc/resolver"

	"github.com/GetterSethya/imageProto"
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

	imageRb := &ImageServiceResolverBuilder{
		ImageServiceHostname: cfg.ImageServiceHostName,
	}

	// dial grpc user service
	resolver.Register(userRB)
	userServiceGrpcConn, err := generateUserServiceGrpcConn(cfg.UserServiceHostName)
	if err != nil {
		log.Fatalf("Cannot connect to user Grpc server: %v", err)
	}

	resolver.Register(imageRb)
	imageServiceGrpcConn, err := generateImageServiceGrpcConn(cfg.ImageServiceHostName)
	if err != nil {
		log.Fatalf("Cannon connect to image Grpc server: %v", err)
	}

	userGrpcClient := userProto.NewUserClient(userServiceGrpcConn)
	imageGrpcClient := imageProto.NewUserClient(imageServiceGrpcConn)
	routes := mux.NewRouter().PathPrefix("/v1/post").Subrouter()

	userService := NewUserService(store, userGrpcClient, imageGrpcClient)
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
