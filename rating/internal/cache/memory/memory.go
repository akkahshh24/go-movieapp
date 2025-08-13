package memory

import (
	"context"

	"github.com/akkahshh24/movieapp/rating/internal/cache"
	"github.com/akkahshh24/movieapp/rating/pkg/model"
)

// Cache defines a rating cache.
// It stores aggregated ratings for records in memory.
type Cache struct {
	data map[model.RecordType]map[model.RecordID]float64
}

// New creates a new memory cache.
func New() *Cache {
	return &Cache{data: map[model.RecordType]map[model.RecordID]float64{}}
}

// Get retrieves the aggregated rating for a given record.
func (c *Cache) Get(_ context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	if _, ok := c.data[recordType]; !ok {
		return 0, cache.ErrNotFound
	}

	rating, ok := c.data[recordType][recordID]
	if !ok {
		return 0, cache.ErrNotFound
	}

	return rating, nil
}

// Put adds or updates the aggregated rating for a given record.
func (c *Cache) Put(_ context.Context, recordID model.RecordID, recordType model.RecordType, rating float64) error {
	if _, ok := c.data[recordType]; !ok {
		c.data[recordType] = map[model.RecordID]float64{}
	}

	c.data[recordType][recordID] = rating
	return nil
}
