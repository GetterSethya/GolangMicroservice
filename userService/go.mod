module github.com/GetterSethya/userService

go 1.21.1

require (
	github.com/GetterSethya/library v1.0.0
	github.com/GetterSethya/userProto v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.1
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/rabbitmq/amqp091-go v1.10.0
	golang.org/x/crypto v0.22.0
	golang.org/x/sync v0.7.0
	google.golang.org/grpc v1.63.2
)

require (
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240227224415-6ceb2ff114de // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)

replace github.com/GetterSethya/library => ../library

replace github.com/GetterSethya/userProto => ../userProto
