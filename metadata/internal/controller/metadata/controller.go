package metadata

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/akkahshh24/movieapp/metadata/internal/repository"
	"github.com/akkahshh24/movieapp/metadata/pkg/model"
)

// ErrNotFound is returned when a requested record is not found.
var ErrNotFound = errors.New("not found")

//go:generate mockgen -source=controller.go -destination=../../../../gen/mock/metadata/repository/repository.go -package=repository
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
		log.Println("Returning metadata from cache for " + id)
		return cacheRes, nil
	}

	// Get the metadata from the repository.
	res, err := c.repo.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	// Update the cache with the retrieved metadata.
	if err := c.cache.Put(ctx, id, res); err != nil {
		log.Println("Error updating cache: " + err.Error())
	}

	return res, err
}

// Put writes movie metadata to repository.
func (c *Controller) Put(ctx context.Context, m *model.Metadata) error {
	if err := c.repo.Put(ctx, m.ID, m); err != nil {
		return fmt.Errorf("failed to put metadata: %w", err)
	}

	// Update the cache with the new metadata.
	if err := c.cache.Put(ctx, m.ID, m); err != nil {
		return fmt.Errorf("failed to update cache: %w", err)
	}

	return nil
}
