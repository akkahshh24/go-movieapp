package rating

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/akkahshh24/movieapp/rating/internal/repository"
	"github.com/akkahshh24/movieapp/rating/pkg/model"
)

// ErrNotFound is returned when no ratings are found for a record.
var ErrNotFound = errors.New("ratings not found for a record")

type ratingRepository interface {
	Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error)
	Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

type ratingCache interface {
	Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error)
	Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating float64) error
}

type ratingIngester interface {
	Ingest(ctx context.Context) (chan model.RatingEvent, error)
}

// Controller defines a rating service controller.
type Controller struct {
	repo     ratingRepository
	cache    ratingCache
	ingester ratingIngester
}

// New creates a rating service controller.
func New(repo ratingRepository, cache ratingCache, ingester ratingIngester) *Controller {
	return &Controller{repo, cache, ingester}
}

// GetAggregatedRating returns the aggregated rating for a record or ErrNotFound if there are no ratings for it.
func (c *Controller) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	// Get the aggregated rating from the cache first.
	cacheRes, err := c.cache.Get(ctx, recordID, recordType)
	if err == nil {
		log.Println("Returning aggregated rating from cache for record:", recordID)
		return cacheRes, nil
	}

	ratings, err := c.repo.Get(ctx, recordID, recordType)
	if err != nil && err == repository.ErrNotFound {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}

	sum := float64(0)
	for _, r := range ratings {
		sum += float64(r.Value)
	}

	aggregatedRating := sum / float64(len(ratings))

	// Update the cache with the aggregated rating.
	if err := c.cache.Put(ctx, recordID, recordType, aggregatedRating); err != nil {
		log.Println("Error updating cache with aggregated rating:", err.Error())
	}

	return aggregatedRating, nil
}

// PutRating writes a rating for a given record.
func (c *Controller) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	if err := c.repo.Put(ctx, recordID, recordType, rating); err != nil {
		return fmt.Errorf("put rating: %w", err)
	}

	ratings, err := c.repo.Get(ctx, recordID, recordType)
	if err != nil && err == repository.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		return err
	}

	sum := float64(0)
	for _, r := range ratings {
		sum += float64(r.Value)
	}

	aggregatedRating := sum / float64(len(ratings))

	// Update the cache with the aggregated rating.
	if err := c.cache.Put(ctx, recordID, recordType, aggregatedRating); err != nil {
		fmt.Println("Error updating cache with aggregated rating:", err.Error())
	}

	return nil
}

// StartIngestion starts the ingestion of rating events.
func (s *Controller) StartIngestion(ctx context.Context) error {
	ch, err := s.ingester.Ingest(ctx)
	if err != nil {
		return err
	}
	for e := range ch {
		fmt.Printf("Consumed a message: %v\n", e)
		if err := s.PutRating(ctx, e.RecordID, e.RecordType, &model.Rating{UserID: e.UserID, Value: e.Value}); err != nil {
			return err
		}
	}
	return nil
}
