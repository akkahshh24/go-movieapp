package metadata

import (
	"context"
	"errors"
	"fmt"

	"github.com/akkahshh24/movieapp/metadata/internal/repository"
	"github.com/akkahshh24/movieapp/metadata/pkg/model"
)

// ErrNotFound is returned when a requested record is not found.
var ErrNotFound = errors.New("not found")

type metadataRepository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
	Put(ctx context.Context, id string, metadata *model.Metadata) error
}

// Controller defines a metadata service controller.
type Controller struct {
	repo  metadataRepository
	cache metadataRepository
}

// New creates a metadata service controller.
func New(repo metadataRepository, cache metadataRepository) *Controller {
	return &Controller{repo, cache}
}

// Get returns movie metadata by id.
func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	// Get the metadata from the cache first.
	cacheRes, err := c.cache.Get(ctx, id)
	if err == nil {
		fmt.Println("Returning metadata from a cache for " + id)
		return cacheRes, nil
	}

	// Get the metadata from the repository.
	res, err := c.repo.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}

	// Update the cache with the retrieved metadata.
	if err := c.cache.Put(ctx, id, res); err != nil {
		fmt.Println("Error updating cache: " + err.Error())
	}

	return res, err
}
