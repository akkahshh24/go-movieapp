package mysql

import (
	"context"
	"database/sql"

	"github.com/akkahshh24/movieapp/metadata/internal/repository"
	"github.com/akkahshh24/movieapp/metadata/pkg/model"
	_ "github.com/go-sql-driver/mysql"
)

// Repository defines a MySQL-based movie matadata repository.
type Repository struct {
	db *sql.DB
}

// New creates a new MySQL-based repository.
func New(dsn string) (*Repository, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &Repository{db}, nil
}

// Get retrieves movie metadata by movie id.
func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	var title, description, director string
	row := r.db.QueryRowContext(ctx, "SELECT title, description, director FROM movies WHERE id = ?", id)
	if err := row.Scan(&title, &description, &director); err != nil {
		// If no rows are found, return a not found error.
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		// For other errors, return the error.
		return nil, err
	}

	return &model.Metadata{
		ID:          id,
		Title:       title,
		Description: description,
		Director:    director,
	}, nil
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO movies (id, title, description, director) VALUES (?, ?, ?, ?)",
		id, metadata.Title, metadata.Description, metadata.Director)
	return err
}
