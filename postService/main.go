package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc/resolver"
)

// http port
const PORT = ":3004"

// grpc port
const GRPC_USER_SERVICE_PORT = ":4002"
const GRPC_NUM_INSTANCE = 2

const exampleScheme = "example"
const exampleServiceName = "user-service"

const RABBITMQ_PORT = ":5672"

func main() {

	var wg sync.WaitGroup

	userServiceHostName := os.Getenv("USER_SERVICE_HOSTNAME")
	if userServiceHostName == "" {
		log.Println("USER_SERVICE_HOSTNAME key is not found, fallback to 'localhost'")
		userServiceHostName = "localhost"
	}

	rabbitMqHostname := os.Getenv("RABBITMQ_HOSTNAME")
	if rabbitMqHostname == "" {
		log.Println("RABBITMQ_HOSTNAME is not found, fallback to localhost")
		rabbitMqHostname = "localhost"
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

	server := NewServer(PORT, sqliteStorage)

	wg.Add(1)
	go func() {
		server.Run(conn)
		defer wg.Done()
	}()

	//rabbitMq conn
	connString := fmt.Sprintf("amqp://guest:guest@%s%s/", rabbitMqHostname, RABBITMQ_PORT)
	rabbitMQConn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatalf("Error when creating connection to rabbit mq: %+v", err)
	}

	defer rabbitMQConn.Close()

	consumer := NewConsumer(rabbitMQConn, sqliteStorage)
	consumer.Consume(rabbitMqHostname)

	//biar main func tidak exit duluan
	wg.Wait()
}
