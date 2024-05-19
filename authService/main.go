package main

const PORT = ":3003"
const GRPC_USER_SERVICE_PORT = ":4002"
const GRPC_NUM_INSTANCE = 2
const exampleScheme = "example"
const exampleServiceName = "user-service"

func main() {

	cfg := InitConfig()

	// start http server
	server := NewServer(PORT, cfg)
	server.Run()

}
