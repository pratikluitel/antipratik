package logic

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	commonerrors "github.com/pratikluitel/antipratik/common/errors"
	"github.com/pratikluitel/antipratik/components/posts"
)

// linkLogic implements LinkLogic.
type linkLogic struct {
	store posts.LinkStore
}

// NewlinkLogic creates a new linkLogic backed by the given store.
func NewLinkLogic(s posts.LinkStore) posts.LinkLogic {
	return &linkLogic{store: s}
}

// GetLinks returns all external links.
func (s *linkLogic) GetLinks(ctx context.Context) ([]posts.ExternalLink, error) {
	links, err := s.store.GetLinks(ctx)
	if err != nil {
		return nil, fmt.Errorf("linkLogic.GetLinks: %w", err)
	}
	return links, nil
}

// GetFeaturedLinks returns up to 4 featured external links.
func (s *linkLogic) GetFeaturedLinks(ctx context.Context) ([]posts.ExternalLink, error) {
	links, err := s.store.GetFeaturedLinks(ctx)
	if err != nil {
		return nil, fmt.Errorf("linkLogic.GetFeaturedLinks: %w", err)
	}
	return links, nil
}

func validateExternalLink(input posts.CreateExternalLink) error {
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

func (s *linkLogic) CreateLink(ctx context.Context, input posts.CreateExternalLink) (string, error) {
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
		return "", fmt.Errorf("linkLogic.CreateLink: %w", err)
	}
	return id, nil
}

func (s *linkLogic) UpdateLink(ctx context.Context, id string, input posts.UpdateExternalLink) (posts.ExternalLink, error) {
	if err := commonerrors.RequireNonEmpty("id", id); err != nil {
		return posts.ExternalLink{}, err
	}
	cur, err := s.store.GetLinkByID(ctx, id)
	if err != nil {
		return posts.ExternalLink{}, fmt.Errorf("linkLogic.UpdateLink: %w", err)
	}

	merged := posts.CreateExternalLink{
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
		return posts.ExternalLink{}, err
	}
	domain, err := extractDomain(merged.URL)
	if err != nil {
		return posts.ExternalLink{}, err
	}
	merged.Domain = domain

	if err = s.store.UpdateLink(ctx, id, merged); err != nil {
		return posts.ExternalLink{}, fmt.Errorf("linkLogic.UpdateLink: %w", err)
	}
	return posts.ExternalLink{
		ID: id, Title: merged.Title, URL: merged.URL, Domain: merged.Domain,
		Description: merged.Description, Featured: merged.Featured, Category: merged.Category,
	}, nil
}

func (s *linkLogic) DeleteLink(ctx context.Context, id string) error {
	return s.store.DeleteLink(ctx, id)
}
