package main

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type exampleResolverBuilder struct {
	UserServiceHostname string
}

type exampleResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
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
