package main

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type PostServiceResolverBuilder struct {
	PostServiceHostname string
}

type PostServiceResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (rb *PostServiceResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	addrs := []string{
		rb.PostServiceHostname + "1" + GRPC_POST_SERVICE_PORT,
		rb.PostServiceHostname + "2" + GRPC_POST_SERVICE_PORT,
	}

	r := &PostServiceResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			POST_SERVICE_NAME: addrs,
		},
	}
	r.start()
	return r, nil
}

func (*PostServiceResolverBuilder) Scheme() string { return POST_SCHEME }
func (r *PostServiceResolver) start() {
	addrsStrs := r.addrsStore[r.target.Endpoint()]
	addrs := make([]resolver.Address, len(addrsStrs))
	for i, s := range addrsStrs {
		addrs[i] = resolver.Address{Addr: s}
	}

	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*PostServiceResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*PostServiceResolver) Close()                                  {}

func init() {
	resolver.Register(&PostServiceResolverBuilder{})
}

func generatePostServiceGrpcConn(hostname string) (*grpc.ClientConn, error) {
	if hostname == "localhost" {
		// connect kaya biasa
		return grpc.NewClient("localhost"+GRPC_POST_SERVICE_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// round robin
		return grpc.NewClient(
			fmt.Sprintf("%s:///%s", POST_SCHEME, POST_SERVICE_NAME),
			grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
}
