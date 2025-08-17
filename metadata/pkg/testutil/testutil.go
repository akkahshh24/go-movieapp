package testutil

import (
	"github.com/akkahshh24/movieapp/gen"
	"github.com/akkahshh24/movieapp/metadata/internal/controller/metadata"
	grpchandler "github.com/akkahshh24/movieapp/metadata/internal/handler/grpc"
	"github.com/akkahshh24/movieapp/metadata/internal/repository/memory"
)

// NewTestMetadataGRPCServer creates a new metadata gRPC server to be used in tests.
func NewTestMetadataGRPCServer() gen.MetadataServiceServer {
	repo := memory.New()
	cache := memory.New()
	ctrl := metadata.New(repo, cache)
	return grpchandler.New(ctrl)
}
