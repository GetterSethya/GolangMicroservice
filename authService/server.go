package main

import (
	"log"
	"net/http"

	"github.com/GetterSethya/library"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/resolver"
)

type Server struct {
	ListenAddr string
	Cfg        AppConfig
}

func NewServer(listenAddr string, cfg AppConfig) *Server {
	return &Server{
		ListenAddr: listenAddr,
		Cfg:        cfg,
	}
}

func (s *Server) Run() {

	rb := &exampleResolverBuilder{
		UserServiceHostname: s.Cfg.UserServiceHostname,
	}

	// dial grpc user service
	resolver.Register(rb)

	conn, err := generateGrpcConn(s.Cfg.UserServiceHostname)
	if err != nil {
		log.Fatalf("Cannot connect to Grpc server:%v", err)
	}

	grpcClient := library.NewUserClient(conn)
	routes := mux.NewRouter().PathPrefix("/v1/auth").Subrouter()

	userService := NewAuthService(s.Cfg.JwtSecret, grpcClient, s.Cfg.RefreshSecret)
	userService.RegisterRoutes(routes)

	log.Println("authService is running on port:", PORT)
	log.Fatal(http.ListenAndServe(s.ListenAddr, routes))
}
