package main

import (
	"context"
	"log"
	"net"

	"github.com/GetterSethya/postProto"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	ListenAddr  string
	Store       *SqliteStorage
	Server      *grpc.Server
	NetListener net.Listener
	postProto.UnimplementedPostServer
}

func NewGrpcServer(listenAddr string, store *SqliteStorage) *GrpcServer {
	listen, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to start grpc postService server:%v", err)
	}
	log.Println("newGrpcServer post service on:", listenAddr)

	return &GrpcServer{
		ListenAddr:  listenAddr,
		Store:       store,
		Server:      grpc.NewServer(),
		NetListener: listen,
	}
}

func (s *GrpcServer) RunGrpc() {
	postProto.RegisterPostServer(s.Server, s)
	log.Println("Grpc postService is running on port:", s.ListenAddr)

	if err := s.Server.Serve(s.NetListener); err != nil {
		log.Fatalf("Failed serve postService grpc server:%v", err)
	}
}

func (s *GrpcServer) DecrementReplyById(ctx context.Context, req *postProto.ReplyCountReq) (*postProto.ReplyCountResp, error) {
	log.Println("hit decrement totalReplies by id grpc")
	id := req.GetId()
	resp := &postProto.ReplyCountResp{}

	err := s.Store.DecrementTotalReplyById(id)
	if err != nil {
		return resp, err
	}

	resp.Message = "Increment success"
	return resp, nil
}

func (s *GrpcServer) IncrementReplyById(ctx context.Context, req *postProto.ReplyCountReq) (*postProto.ReplyCountResp, error) {
	log.Println("hit increment totalReplies by id grpc")
	id := req.GetId()
	resp := &postProto.ReplyCountResp{}

	err := s.Store.IncrementTotalReplyById(id)
	if err != nil {
		return resp, err
	}

	resp.Message = "Increment success"
	return resp, nil
}
