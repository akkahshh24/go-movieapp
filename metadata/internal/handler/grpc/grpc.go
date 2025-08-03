package grpc

import (
	"context"
	"errors"

	"github.com/akkahshh24/movieapp/gen"
	"github.com/akkahshh24/movieapp/metadata/internal/controller/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a metadata gRPC handler.
type Handler struct {
	// Embed the generated server interface to implement it.
	// This allows us to use the generated methods directly.
	gen.UnimplementedMetadataServiceServer
	ctrl *metadata.Controller
}

// New creates a new metadata gRPC handler.
func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// GetMetadata returns the metadata of the requested movie ID.
func (h *Handler) GetMetadata(ctx context.Context, req *gen.GetMetadataRequest) (*gen.GetMetadataResponse, error) {
	// Validate the request
	if req == nil || req.MovieId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	// Call the controller to get the metadata
	movieMetaData, err := h.ctrl.Get(ctx, req.MovieId)
	if err != nil && errors.Is(err, metadata.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// Convert the metadata to the proto response format.
	return &gen.GetMetadataResponse{Metadata: movieMetaData.ToProto()}, nil
}
