package main

import (
	"log"
	"os"

	"github.com/GetterSethya/library"
	"google.golang.org/grpc/resolver"
)

const PORT = ":3003"
const GRPC_USER_SERVICE_PORT = ":4002"
const GRPC_NUM_INSTANCE = 2
const exampleScheme = "example"
const exampleServiceName = "user-service"


func main() {

	jwtSecret := os.Getenv("JWT_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")
	userServiceHostName := os.Getenv("USER_SERVICE_HOSTNAME")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET key not found!")
	}

	if refreshSecret == "" {
		log.Fatal("REFRESH_SECRET key not found!")
	}

	if userServiceHostName == "" {
		log.Println("USER_SERVICE_HOSTNAME key is not found, fallback to 'localhost'")
		userServiceHostName = "localhost"
	}

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

	// start http server
	server := NewServer(PORT, jwtSecret, refreshSecret)
	server.Run(c)

}
