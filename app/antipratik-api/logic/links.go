package logic

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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

func validateExternalLink(input models.CreateExternalLink) error {
	if err := requireNonEmpty("title", input.Title); err != nil {
		return err
	}
	if err := requireNonEmpty("url", input.URL); err != nil {
		return err
	}
	if err := requireNonEmpty("description", input.Description); err != nil {
		return err
	}
	if err := requireNonEmpty("category", input.Category); err != nil {
		return err
	}
	return nil
}

func (s *LinkService) CreateLink(ctx context.Context, input models.CreateExternalLink) (string, error) {
	if err := validateExternalLink(input); err != nil {
		return "", err
	}
	domain, err := extractDomain(input.URL)
	if err != nil {
		return "", err
	}
	input.Domain = domain
	id := uuid.New().String()
	if err := s.store.CreateLink(ctx, id, input); err != nil {
		return "", fmt.Errorf("LinkService.CreateLink: %w", err)
	}
	return id, nil
}

func (s *LinkService) UpdateLink(ctx context.Context, id string, input models.CreateExternalLink) error {
	if err := requireNonEmpty("id", id); err != nil {
		return err
	}
	if err := validateExternalLink(input); err != nil {
		return err
	}
	domain, err := extractDomain(input.URL)
	if err != nil {
		return err
	}
	input.Domain = domain
	return s.store.UpdateLink(ctx, id, input)
}

func (s *LinkService) DeleteLink(ctx context.Context, id string) error {
	return s.store.DeleteLink(ctx, id)
}
