package main

import (
	"log"
	"os"
	"time"

	"github.com/GetterSethya/library"
	"google.golang.org/grpc/resolver"
)

// http port
const PORT = ":3004"
const GRPC_USER_SERVICE_PORT = ":4002"
const GRPC_NUM_INSTANCE = 2
const exampleScheme = "example"
const exampleServiceName = "user-service"

func main() {

	userServiceHostName := os.Getenv("USER_SERVICE_HOSTNAME")
	if userServiceHostName == "" {
		log.Println("USER_SERVICE_HOSTNAME key is not found, fallback to 'localhost'")
		userServiceHostName = "localhost"
	}

	sqliteStorage := NewSqliteStorage()
	sqliteStorage.Init()

	// set db conn limit
	sqliteStorage.db.SetMaxOpenConns(25)
	sqliteStorage.db.SetMaxIdleConns(25)
	sqliteStorage.db.SetConnMaxLifetime(5 * time.Minute)

	rb := &exampleResolverBuilder{
		UserServiceHostname: userServiceHostName,
	}
	// dial grpc user service
	resolver.Register(rb)
	conn, err := generateGrpcConn(userServiceHostName)
	if err != nil {
		log.Fatalf("Cannot connect to Grpc server:%v", err)
	}

	c := library.NewUserClient(conn)

	server := NewServer(PORT, sqliteStorage)
	server.Run(c)

}
