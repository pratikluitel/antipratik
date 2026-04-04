package logic

import (
	"context"
	"fmt"

	"github.com/pratikluitel/antipratik/models"
	"github.com/pratikluitel/antipratik/store"
)

// LinkService implements LinkLogic.
type LinkService struct {
	store store.LinkStore
}

// NewLinkService creates a new LinkService backed by the given store.
func NewLinkService(s store.LinkStore) *LinkService {
	return &LinkService{store: s}
}

// GetLinks returns all external links.
func (s *LinkService) GetLinks(ctx context.Context) ([]models.ExternalLink, error) {
	links, err := s.store.GetLinks(ctx)
	if err != nil {
		return nil, fmt.Errorf("LinkService.GetLinks: %w", err)
	}
	return links, nil
}

// GetFeaturedLinks returns up to 4 featured external links.
func (s *LinkService) GetFeaturedLinks(ctx context.Context) ([]models.ExternalLink, error) {
	links, err := s.store.GetFeaturedLinks(ctx)
	if err != nil {
		return nil, fmt.Errorf("LinkService.GetFeaturedLinks: %w", err)
	}
	return links, nil
}
