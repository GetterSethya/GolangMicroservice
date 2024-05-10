package main

import (
	"fmt"
	"log"
	"os"

	"github.com/GetterSethya/library"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

const PORT = ":3003"
const GRPC_USER_SERVICE_PORT = ":4002"
const GRPC_NUM_INSTANCE = 2
const exampleScheme = "example"
const exampleServiceName = "user-service"

type exampleResolverBuilder struct {
	UserServiceHostname string
}

type exampleResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

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

func (rb *exampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	addrs := []string{
		rb.UserServiceHostname + "1" + GRPC_USER_SERVICE_PORT,
		rb.UserServiceHostname + "2" + GRPC_USER_SERVICE_PORT,
	}

	r := &exampleResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			exampleServiceName: addrs,
		},
	}
	r.start()
	return r, nil
}

func (*exampleResolverBuilder) Scheme() string { return exampleScheme }
func (r *exampleResolver) start() {
	addrsStrs := r.addrsStore[r.target.Endpoint()]
	addrs := make([]resolver.Address, len(addrsStrs))
	for i, s := range addrsStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*exampleResolver) Close()                                  {}

func init() {
	resolver.Register(&exampleResolverBuilder{})
}

func generateGrpcConn(hostname string) (*grpc.ClientConn, error) {
	if hostname == "localhost" {
		// connect kaya biasa
		return grpc.Dial("localhost"+GRPC_USER_SERVICE_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// round robin
		return grpc.Dial(
			fmt.Sprintf("%s:///%s", exampleScheme, exampleServiceName),
			grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
}
