package logic

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	commonerrors "github.com/pratikluitel/antipratik/common/errors"
	"github.com/pratikluitel/antipratik/components/posts/models"
	"github.com/pratikluitel/antipratik/components/posts/store"
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
	if err := commonerrors.RequireNonEmpty("title", input.Title); err != nil {
		return err
	}
	if err := commonerrors.RequireNonEmpty("url", input.URL); err != nil {
		return err
	}
	if err := commonerrors.RequireNonEmpty("description", input.Description); err != nil {
		return err
	}
	if err := commonerrors.RequireNonEmpty("category", input.Category); err != nil {
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

func (s *LinkService) UpdateLink(ctx context.Context, id string, input models.UpdateExternalLink) (models.ExternalLink, error) {
	if err := commonerrors.RequireNonEmpty("id", id); err != nil {
		return models.ExternalLink{}, err
	}
	cur, err := s.store.GetLinkByID(ctx, id)
	if err != nil {
		return models.ExternalLink{}, fmt.Errorf("LinkService.UpdateLink: %w", err)
	}

	merged := models.CreateExternalLink{
		Title: cur.Title, URL: cur.URL, Domain: cur.Domain,
		Description: cur.Description, Featured: cur.Featured, Category: cur.Category,
	}
	if input.Title != nil {
		merged.Title = *input.Title
	}
	if input.URL != nil {
		merged.URL = *input.URL
	}
	if input.Description != nil {
		merged.Description = *input.Description
	}
	if input.Featured != nil {
		merged.Featured = *input.Featured
	}
	if input.Category != nil {
		merged.Category = *input.Category
	}

	if err = validateExternalLink(merged); err != nil {
		return models.ExternalLink{}, err
	}
	domain, err := extractDomain(merged.URL)
	if err != nil {
		return models.ExternalLink{}, err
	}
	merged.Domain = domain

	if err = s.store.UpdateLink(ctx, id, merged); err != nil {
		return models.ExternalLink{}, fmt.Errorf("LinkService.UpdateLink: %w", err)
	}
	return models.ExternalLink{
		ID: id, Title: merged.Title, URL: merged.URL, Domain: merged.Domain,
		Description: merged.Description, Featured: merged.Featured, Category: merged.Category,
	}, nil
}

func (s *LinkService) DeleteLink(ctx context.Context, id string) error {
	return s.store.DeleteLink(ctx, id)
}
