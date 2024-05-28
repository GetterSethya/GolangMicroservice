package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GetterSethya/userProto"
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

type AppServer struct {
	Cfg      AppConfig
	Server   http.Server
	UserServiceGrpcConn *grpc.ClientConn
}

func NewServer(listenAddr string, cfg AppConfig) *AppServer {

	rb := &exampleResolverBuilder{
		UserServiceHostname: cfg.UserServiceHostname,
	}

	// dial grpc user service
	resolver.Register(rb)

	conn, err := generateGrpcConn(cfg.UserServiceHostname)
	if err != nil {
		log.Fatalf("Cannot connect to Grpc server:%v", err)
	}

	grpcClient := userProto.NewUserClient(conn)
	routes := mux.NewRouter().PathPrefix("/v1/auth").Subrouter()

	userService := NewAuthService(cfg.JwtSecret, grpcClient, cfg.RefreshSecret)
	userService.RegisterRoutes(routes)
	return &AppServer{
		Cfg:      cfg,
		UserServiceGrpcConn: conn,
		Server: http.Server{
			Addr:    listenAddr,
			Handler: routes,
		},
	}
}

func (s *AppServer) Run() {

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()
	log.Println("authService is running on port:", PORT)

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return s.Server.ListenAndServe()
	})

	g.Go(func() error {
		<-gctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer shutdownCancel()

		if err := s.UserServiceGrpcConn.Close(); err != nil {
			log.Println("Erron when closing grpc connections:", err)
		}

		log.Println("SIGTERM detected, will attempt to graceful shutdown...")
		return s.Server.Shutdown(shutdownCtx)
	})

	if err := g.Wait(); err != nil {
		log.Println(err)
	}

}
