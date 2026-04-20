// Package services exposes files component capabilities as injectable interfaces
// for use by other components.
package services

import (
	"context"
	"io"

	"github.com/pratikluitel/antipratik/components/files/store"
)

// StorageService exposes file retrieval and deletion to other components.
// Inject this interface rather than importing files/store directly.
type StorageService interface {
	Get(ctx context.Context, key string) (io.ReadSeekCloser, string, error)
	Delete(ctx context.Context, key string) error
}

type storageService struct {
	store store.FileStore
}

// NewStorageService returns a StorageService backed by the given FileStore.
func NewStorageService(s store.FileStore) StorageService {
	return &storageService{store: s}
}

func (s *storageService) Get(ctx context.Context, key string) (io.ReadSeekCloser, string, error) {
	return s.store.Get(ctx, key)
}

func (s *storageService) Delete(ctx context.Context, key string) error {
	return s.store.Delete(ctx, key)
}
