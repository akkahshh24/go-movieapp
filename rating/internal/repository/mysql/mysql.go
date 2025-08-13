package mysql

import (
	"context"
	"database/sql"

	"github.com/akkahshh24/movieapp/rating/pkg/model"
	_ "github.com/go-sql-driver/mysql"
)

const (
	driverName     = "mysql"
	dataSourceName = "root:password@/movieapp"
)

// Repository defines a MySQL-based rating repository.
type Repository struct {
	db *sql.DB
}

// New creates a new MySQL-based rating repository.
func New() (*Repository, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Repository{db}, nil
}

// Get retrieves all ratings for a given record.
func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT user_id, value FROM ratings WHERE record_id = ? AND record_type = ?", recordID, recordType)
	if err != nil {
		return nil, err
	}
	// Ensure rows are closed after processing.
	defer rows.Close()

	var res []model.Rating
	for rows.Next() {
		var userID string
		var value int32
		if err := rows.Scan(&userID, &value); err != nil {
			return nil, err
		}

		// Append the rating to the result slice.
		res = append(res, model.Rating{
			UserID: model.UserID(userID),
			Value:  model.RatingValue(value),
		})
	}

	return res, nil
}

// Put adds a rating for a given record.
func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO ratings (record_id, record_type, user_id, value) VALUES (?, ?, ?, ?)",
		recordID, recordType, rating.UserID, rating.Value)
	return err
}
