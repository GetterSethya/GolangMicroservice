package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/GetterSethya/userProto"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	ListenAddr string
	Store       *SqliteStorage
	Server      *grpc.Server
	NetListener net.Listener
	userProto.UnimplementedUserServer
}

func NewGrpcServer(listenAddr string, store *SqliteStorage) *GrpcServer {

	listen, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to start grpc userService server:%v", err)
	}

	return &GrpcServer{
		ListenAddr:  listenAddr,
		Store:       store,
		Server:      grpc.NewServer(),
		NetListener: listen,
	}
}

func (s *GrpcServer) RunGrpc() {
	userProto.RegisterUserServer(s.Server, s)
	log.Println("Grpc userService is running on port:", s.ListenAddr)

	if err := s.Server.Serve(s.NetListener); err != nil {
		log.Fatalf("Failed serve userService grpc server:%v", err)
	}
}

func (s *GrpcServer) IncrementFollowerById(ctx context.Context, req *userProto.RelationReq) (*userProto.RelationResp, error) {
	log.Println("hit increment follower by id grpc")

	id := req.GetId()

	resp := &userProto.RelationResp{}

	err := s.Store.IncrementFollowerById(id)
	if err != nil {
		return resp, err
	}

	resp.Message = "Increment success"

	return resp, nil

}

func (s *GrpcServer) DecrementFollowerById(ctx context.Context, req *userProto.RelationReq) (*userProto.RelationResp, error) {

	log.Println("hit decrement follower by id grpc")

	id := req.GetId()

	resp := &userProto.RelationResp{}

	err := s.Store.DecrementFollowerById(id)
	if err != nil {
		return resp, err
	}

	resp.Message = "Decrement success"

	return resp, nil
}

func (s *GrpcServer) IncrementFollowingById(ctx context.Context, req *userProto.RelationReq) (*userProto.RelationResp, error) {

	log.Println("hit increment following by id grpc")

	id := req.GetId()

	resp := &userProto.RelationResp{}

	err := s.Store.IncrementFollowingById(id)
	if err != nil {
		return resp, err
	}

	resp.Message = "Increment success"

	return resp, nil
}

func (s *GrpcServer) DecrementFollowingById(ctx context.Context, req *userProto.RelationReq) (*userProto.RelationResp, error) {

	log.Println("hit decremenet following by id grpc")

	id := req.GetId()

	resp := &userProto.RelationResp{}

	err := s.Store.DecrementFollowingById(id)
	if err != nil {
		return resp, err
	}

	resp.Message = "Increment success"

	return resp, nil
}

func (s *GrpcServer) GetUserPasswordById(ctx context.Context, req *userProto.GetUserByIdReq) (*userProto.UserPasswordResp, error) {

	log.Println("hit get user password by id grpc")

	id := req.GetId()

	user := &User{}

	err := s.Store.GetUserPasswordById(id, user)
	if err != nil {
		return &userProto.UserPasswordResp{}, err
	}

	returnUser := &userProto.UserPasswordResp{
		Id:           user.Id,
		HashPassword: user.HashPassword,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	return returnUser, nil

}
func (s *GrpcServer) GetUserPasswordByUsername(ctx context.Context, req *userProto.GetUserByUsernameReq) (*userProto.UserPasswordResp, error) {
	log.Println("hit get user password by username grpc")

	username := req.GetUsername()

	user := &User{}

	err := s.Store.GetUserPasswordByUsername(username, user)
	if err != nil {
		return &userProto.UserPasswordResp{}, err
	}

	returnUser := &userProto.UserPasswordResp{
		Id:           user.Id,
		HashPassword: user.HashPassword,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	return returnUser, nil

}

func (s *GrpcServer) CreateUser(ctx context.Context, req *userProto.CreateUserReq) (*userProto.CreateUserResp, error) {
	log.Println("hit handle create user grpc")

	resp := &userProto.CreateUserResp{}

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

func (s *GrpcServer) GetUserById(ctx context.Context, req *userProto.GetUserByIdReq) (*userProto.UserResp, error) {

	log.Println("hit get user by id grpc")

	id := req.GetId()

	user := &ReturnUser{}

	err := s.Store.GetUserById(id, user)
	if err != nil {
		return &userProto.UserResp{}, err
	}

	returnUser := &userProto.UserResp{
		Id:        user.Id,
		Username:  user.Username,
		Name:      user.Name,
		Profile:   user.Profile,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return returnUser, nil

}

func (s *GrpcServer) GetUserByUsername(ctx context.Context, req *userProto.GetUserByUsernameReq) (*userProto.UserResp, error) {

	log.Println("hit get user by username grpc")

	username := req.GetUsername()

	user := &ReturnUser{}

	err := s.Store.GetUserByUsername(username, user)
	if err != nil {
		return &userProto.UserResp{}, err
	}

	returnUser := &userProto.UserResp{
		Id:        user.Id,
		Username:  user.Username,
		Name:      user.Name,
		Profile:   user.Profile,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return returnUser, nil
}
