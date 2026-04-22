// Package services exposes files component capabilities as injectable interfaces
// for use by other components.
package services

import (
	"context"
	"io"

	"github.com/pratikluitel/antipratik/components/files"
)

type storageService struct {
	store files.FileStore
}

// NewStorageService returns a StorageService backed by the given Filefiles.
func NewStorageService(s files.FileStore) files.StorageService {
	return &storageService{store: s}
}

func (s *storageService) Get(ctx context.Context, key string) (io.ReadSeekCloser, string, error) {
	return s.store.Get(ctx, key)
}

func (s *storageService) Delete(ctx context.Context, key string) error {
	return s.store.Delete(ctx, key)
}
