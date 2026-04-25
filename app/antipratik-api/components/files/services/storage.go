// Package services exposes files component capabilities as injectable interfaces
// for use by other components.
package services

import (
	"context"

	"github.com/pratikluitel/antipratik/components/files"
)

type storageService struct {
	store files.FileStore
}

// NewStorageService returns a StorageService backed by the given Filefiles.
func NewStorageService(s files.FileStore) files.StorageService {
	return &storageService{store: s}
}

func (s *storageService) Delete(ctx context.Context, key string) error {
	return s.store.Delete(ctx, key)
}
