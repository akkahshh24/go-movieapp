package memory

import (
	"context"

	"github.com/akkahshh24/movieapp/rating/internal/repository"
	"github.com/akkahshh24/movieapp/rating/pkg/model"
)

// Repository defines a rating repository.
type Repository struct {
	data map[model.RecordType]map[model.RecordID][]model.Rating
	// For example, {movie: {movie_id: {4, 5, 5}}}
}

// New creates a new memory repository.
func New() *Repository {
	return &Repository{map[model.RecordType]map[model.RecordID][]model.Rating{}}
}

// Get retrieves all ratings for a given record.
func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	// Check if the record type exists in the repository.
	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrNotFound
	}

	// Check if the record ID exists for the given record type.
	if ratings, ok := r.data[recordType][recordID]; !ok || len(ratings) == 0 {
		return nil, repository.ErrNotFound
	}
	return r.data[recordType][recordID], nil
}

// Put adds a rating for a given record.
func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	// Initialize the record type map if it doesn't exist.
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[model.RecordID][]model.Rating{}
	}

	// Initialize the record ID slice if it doesn't exist.
	r.data[recordType][recordID] = append(r.data[recordType][recordID], *rating)
	return nil
}
