module github.com/GetterSethya/relationService

go 1.21.1

require (
	github.com/GetterSethya/library v0.0.0-00010101000000-000000000000
	github.com/GetterSethya/userProto v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/rabbitmq/amqp091-go v1.10.0
	google.golang.org/grpc v1.65.0
)

require (
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace github.com/GetterSethya/library => ../library

replace github.com/GetterSethya/userProto => ../userProto
