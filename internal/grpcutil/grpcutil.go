package grpcutil

import (
	"context"
	"math/rand"

	"github.com/akkahshh24/movieapp/pkg/discovery"
	"github.com/akkahshh24/movieapp/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServiceConnection attempts to select a random service instance and returns a gRPC connection to it.
func ServiceConnection(ctx context.Context, serviceName string, registry discovery.Registry) (*grpc.ClientConn, error) {
	// Get the service endpoints from the registry.
	// This will return a list of addresses for the service instances.
	addrs, err := registry.ServiceEndpoints(ctx, model.ServiceName(serviceName))
	if err != nil {
		return nil, err
	}

	// Create a gRPC client connection to a random service instance.
	return grpc.NewClient(addrs[rand.Intn(len(addrs))], grpc.WithTransportCredentials(insecure.NewCredentials()))
}
