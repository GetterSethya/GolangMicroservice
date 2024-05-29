package main

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type ImageServiceResolverBuilder struct {
	ImageServiceHostname string
}

type ImageServiceResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (rb *ImageServiceResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	addrs := []string{
		rb.ImageServiceHostname + "1" + GRPC_IMAGE_SERVICE_PORT,
		rb.ImageServiceHostname + "2" + GRPC_IMAGE_SERVICE_PORT,
	}

	r := &ImageServiceResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			IMAGE_SERVICE_NAME: addrs,
		},
	}
	r.start()
	return r, nil
}

func (*ImageServiceResolverBuilder) Scheme() string { return IMAGE_SCHEME }
func (r *ImageServiceResolver) start() {
	addrsStrs := r.addrsStore[r.target.Endpoint()]
	addrs := make([]resolver.Address, len(addrsStrs))
	for i, s := range addrsStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*ImageServiceResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*ImageServiceResolver) Close()                                  {}

func init() {
	resolver.Register(&ImageServiceResolverBuilder{})
}

func generateImageServiceGrpcConn(hostname string) (*grpc.ClientConn, error) {
	if hostname == "localhost" {
		// connect kaya biasa
		return grpc.NewClient("localhost"+GRPC_IMAGE_SERVICE_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// round robin
		return grpc.NewClient(
			fmt.Sprintf("%s:///%s", IMAGE_SCHEME, IMAGE_SERVICE_NAME),
			grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
}
