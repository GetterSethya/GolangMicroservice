package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
)

type AppServer struct {
	Cfg    AppConfig
	Server http.Server
}

func NewAppServer(cfg AppConfig) *AppServer {

	router := mux.NewRouter().PathPrefix("/v1/image").Subrouter()
	imageService := NewImageService(cfg.Host, cfg.Port)
	imageService.RegisterRoutes(router)

	return &AppServer{
		Cfg: cfg,
		Server: http.Server{
			Addr:    cfg.Port,
			Handler: router,
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

	log.Println("Image service is running on port:", s.Cfg.Port)

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return s.Server.ListenAndServe()
	})

	g.Go(func() error {
		<-gctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		log.Println("SIGTERM detected, will attempt to graceful shutdown...")
		return s.Server.Shutdown(shutdownCtx)
	})

	if err := g.Wait(); err != nil {
		log.Println(err)
	}
}
