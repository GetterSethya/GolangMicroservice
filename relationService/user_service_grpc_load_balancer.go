package main

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type UserServiceResolverBuilder struct {
	UserServiceHostname string
}

type UserServiceResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (rb *UserServiceResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	addrs := []string{
		rb.UserServiceHostname + "1" + GRPC_USER_SERVICE_PORT,
		rb.UserServiceHostname + "2" + GRPC_USER_SERVICE_PORT,
	}

	r := &UserServiceResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			USER_SERVICE_NAME: addrs,
		},
	}
	r.start()
	return r, nil
}

func (*UserServiceResolverBuilder) Scheme() string { return USER_SCHEME }
func (r *UserServiceResolver) start() {
	addrsStrs := r.addrsStore[r.target.Endpoint()]
	addrs := make([]resolver.Address, len(addrsStrs))
	for i, s := range addrsStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*UserServiceResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*UserServiceResolver) Close()                                  {}

func init() {
	resolver.Register(&UserServiceResolverBuilder{})
}

func generateUserServiceGrpcConn(hostname string) (*grpc.ClientConn, error) {
	if hostname == "localhost" {
		// connect kaya biasa
		return grpc.NewClient("localhost"+GRPC_USER_SERVICE_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// round robin
		return grpc.NewClient(
			fmt.Sprintf("%s:///%s", USER_SCHEME, USER_SERVICE_NAME),
			grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
}
