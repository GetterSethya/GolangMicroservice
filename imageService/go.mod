module github.com/GetterSethya/imageService

go 1.21.1

require (
	github.com/GetterSethya/imageProto v1.0.0
	github.com/GetterSethya/library v1.0.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	golang.org/x/image v0.15.0
	golang.org/x/sync v0.7.0
	google.golang.org/grpc v1.64.0
)

require (
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)

replace github.com/GetterSethya/library => ../library

replace github.com/GetterSethya/imageProto => ../imageProto
