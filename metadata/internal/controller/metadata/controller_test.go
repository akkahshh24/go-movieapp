package metadata

import (
	"context"
	"errors"
	"testing"

	gen "github.com/akkahshh24/movieapp/gen/mock/metadata/repository"
	"github.com/akkahshh24/movieapp/metadata/internal/repository"
	"github.com/akkahshh24/movieapp/metadata/pkg/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name        string
		expCacheRes *model.Metadata
		expCacheErr error
		expRepoRes  *model.Metadata
		expRepoErr  error
		wantRes     *model.Metadata
		wantErr     error
	}{
		{
			name:        "cache hit",
			expCacheRes: &model.Metadata{},
			expCacheErr: nil,
			wantRes:     &model.Metadata{},
			wantErr:     nil,
		},
		{
			name:        "not found",
			expCacheRes: nil,
			expCacheErr: repository.ErrNotFound,
			expRepoErr:  repository.ErrNotFound,
			wantErr:     ErrNotFound,
		},
		{
			name:        "unexpected error",
			expCacheRes: nil,
			expCacheErr: repository.ErrNotFound,
			expRepoRes:  nil,
			expRepoErr:  errors.New("unexpected error"),
			wantErr:     errors.New("unexpected error"),
		},
		{
			name:        "success",
			expCacheRes: nil,
			expCacheErr: repository.ErrNotFound,
			expRepoRes:  &model.Metadata{},
			expRepoErr:  nil,
			wantRes:     &model.Metadata{},
			wantErr:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := gen.NewMockmetadataRepository(ctrl)
			cacheMock := gen.NewMockmetadataRepository(ctrl)
			c := New(repoMock, cacheMock)

			ctx := context.Background()
			id := "id"

			// Cache expectation
			cacheMock.EXPECT().Get(ctx, id).Return(tt.expCacheRes, tt.expCacheErr)

			// If cache hit, repo shouldn't be called
			if tt.expCacheErr != nil {
				repoMock.EXPECT().Get(ctx, id).Return(tt.expRepoRes, tt.expRepoErr)

				// If repo succeeds, cache should be updated
				if tt.expRepoErr == nil {
					cacheMock.EXPECT().Put(ctx, id, tt.expRepoRes).Return(nil)
				}
			}

			res, err := c.Get(ctx, id)
			assert.Equal(t, tt.wantRes, res, tt.name)
			assert.Equal(t, tt.wantErr, err, tt.name)
		})
	}
}
