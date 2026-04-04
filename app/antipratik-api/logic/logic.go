// Package logic contains the business logic layer.
// It sits between the HTTP api layer and the store layer.
package logic

import (
	"context"

	"github.com/pratikluitel/antipratik/models"
)

// PostLogic defines the business operations on posts.
type PostLogic interface {
	// GetPosts returns posts matching the given filter, newest first.
	GetPosts(ctx context.Context, filter models.FilterState) ([]models.Post, error)

	// GetPost returns the essay with the given slug, or nil if not found.
	GetPost(ctx context.Context, slug string) (*models.EssayPost, error)
}

// LinkLogic defines the business operations on external links.
type LinkLogic interface {
	// GetLinks returns all external links.
	GetLinks(ctx context.Context) ([]models.ExternalLink, error)

	// GetFeaturedLinks returns up to 4 featured external links.
	GetFeaturedLinks(ctx context.Context) ([]models.ExternalLink, error)
}
