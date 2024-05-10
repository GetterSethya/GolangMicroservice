package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/GetterSethya/library"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	ListenAddr string
	library.UnimplementedUserServer
	Store *SqliteStorage
}

func (s *GrpcServer) GetUserPasswordById(ctx context.Context, req *library.GetUserByIdReq) (*library.UserPasswordResp, error) {

	log.Println("hit get user password by id grpc")

	id := req.GetId()

	user := &User{}

	err := s.Store.GetUserPasswordById(id, user)
	if err != nil {
		return &library.UserPasswordResp{}, err
	}

	returnUser := &library.UserPasswordResp{
		Id:           user.Id,
		HashPassword: user.HashPassword,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	return returnUser, nil

}
func (s *GrpcServer) GetUserPasswordByUsername(ctx context.Context, req *library.GetUserByUsernameReq) (*library.UserPasswordResp, error) {
	log.Println("hit get user password by username grpc")

	username := req.GetUsername()

	user := &User{}

	err := s.Store.GetUserPasswordByUsername(username, user)
	if err != nil {
		return &library.UserPasswordResp{}, err
	}

	returnUser := &library.UserPasswordResp{
		Id:           user.Id,
		HashPassword: user.HashPassword,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	return returnUser, nil

}

func (s *GrpcServer) CreateUser(ctx context.Context, req *library.CreateUserReq) (*library.CreateUserResp, error) {
	log.Println("hit handle create user grpc")

	resp := &library.CreateUserResp{}

	unixEpoch := time.Now().Unix()

	if err := s.Store.CreateUser(
		req.GetId(),
		req.GetUsername(),
		req.GetName(),
		req.GetHashPassword(),
		defaultProfile,
		unixEpoch,
		unixEpoch,
	); err != nil {
		log.Println("Error when inserting user:", err)
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
			return resp, fmt.Errorf("Username already used")
		}
		return resp, fmt.Errorf("Something went wrong")
	}

	resp.Message = "User created!"

	return resp, nil

}

func (s *GrpcServer) GetUserById(ctx context.Context, req *library.GetUserByIdReq) (*library.UserResp, error) {

	log.Println("hit get user by id grpc")

	id := req.GetId()

	user := &ReturnUser{}

	err := s.Store.GetUserById(id, user)
	if err != nil {
		return &library.UserResp{}, err
	}

	returnUser := &library.UserResp{
		Id:        user.Id,
		Username:  user.Username,
		Name:      user.Name,
		Profile:   user.Profile,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return returnUser, nil

}

func (s *GrpcServer) GetUserByUsername(ctx context.Context, req *library.GetUserByUsernameReq) (*library.UserResp, error) {

	log.Println("hit get user by username grpc")

	username := req.GetUsername()

	user := &ReturnUser{}

	err := s.Store.GetUserByUsername(username, user)
	if err != nil {
		return &library.UserResp{}, err
	}

	returnUser := &library.UserResp{
		Id:        user.Id,
		Username:  user.Username,
		Name:      user.Name,
		Profile:   user.Profile,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return returnUser, nil
}

func NewGrpcServer(listenAddr string, store *SqliteStorage) *GrpcServer {

	return &GrpcServer{
		ListenAddr: listenAddr,
		Store:      store,
	}
}

func (s *GrpcServer) RunGrpc() {
	listen, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		log.Fatalf("Failed to start grpc userService server:%v", err)
	}

	server := grpc.NewServer()
	library.RegisterUserServer(server, s)
	log.Println("Grpc userService is running on port:", s.ListenAddr)
	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed serve userService grpc server:%v", err)
	}
}
