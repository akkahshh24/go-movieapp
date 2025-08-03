package grpc

import (
	"context"
	"errors"

	"github.com/akkahshh24/movieapp/gen"
	"github.com/akkahshh24/movieapp/movie/internal/controller/movie"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a movie gRPC handler.
type Handler struct {
	gen.UnimplementedMovieServiceServer
	ctrl *movie.Controller
}

// New creates a new movie gRPC handler.
func New(ctrl *movie.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// GetMovieDetails returns moviie details by id.
func (h *Handler) GetMovieDetails(ctx context.Context, req *gen.GetMovieDetailsRequest) (*gen.GetMovieDetailsResponse, error) {
	// Validate the request.
	if req == nil || req.MovieId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	// Call the controller to get movie details.
	// This will fetch the movie metadata and rating from the respective gateways.
	m, err := h.ctrl.Get(ctx, req.MovieId)
	if err != nil && errors.Is(err, movie.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// Convert the movie metadata to the gRPC response format.
	return &gen.GetMovieDetailsResponse{
		MovieDetails: &gen.MovieDetails{
			Metadata: m.Metadata.ToProto(),
			Rating:   *m.Rating,
		},
	}, nil
}
