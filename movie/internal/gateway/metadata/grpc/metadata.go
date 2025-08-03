package grpc

import (
	"context"

	"github.com/akkahshh24/movieapp/gen"
	"github.com/akkahshh24/movieapp/internal/grpcutil"
	"github.com/akkahshh24/movieapp/metadata/pkg/model"
	"github.com/akkahshh24/movieapp/pkg/discovery"
)

// Gateway defines a movie metadata gRPC gateway.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for a movie metadata service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

// Get returns movie metadata by a movie id.
func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	// Create a gRPC connection to the metadata service.
	// This will select a random service instance from the registry.
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}
	// Close the connection when done.
	// This is important to avoid resource leaks.
	defer conn.Close()

	// Create a gRPC client for the metadata service.
	// This client will be used to call the GetMetadata method.
	client := gen.NewMetadataServiceClient(conn)
	resp, err := client.GetMetadata(ctx, &gen.GetMetadataRequest{MovieId: id})
	if err != nil {
		return nil, err
	}

	// Convert the response to the model.Metadata type.
	return model.ProtoToMetadata(resp.Metadata), nil
}
