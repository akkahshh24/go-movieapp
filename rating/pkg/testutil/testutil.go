package testutil

import (
	"github.com/akkahshh24/movieapp/gen"
	cachememory "github.com/akkahshh24/movieapp/rating/internal/cache/memory"
	"github.com/akkahshh24/movieapp/rating/internal/controller/rating"
	grpchandler "github.com/akkahshh24/movieapp/rating/internal/handler/grpc"
	repomemory "github.com/akkahshh24/movieapp/rating/internal/repository/memory"
)

// NewTestRatingGRPCServer creates a new rating gRPC server to be used in tests.
func NewTestRatingGRPCServer() gen.RatingServiceServer {
	repo := repomemory.New()
	cache := cachememory.New()
	ctrl := rating.New(repo, cache, nil)
	return grpchandler.New(ctrl)
}
