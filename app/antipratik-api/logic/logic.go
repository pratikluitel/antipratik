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

	// Write operations
	CreateEssay(ctx context.Context, input models.CreateEssayPost) (string, error)
	CreateShort(ctx context.Context, input models.CreateShortPost) (string, error)
	CreateMusic(ctx context.Context, input models.CreateMusicPost) (string, error)
	CreatePhoto(ctx context.Context, input models.CreatePhotoPost) (string, error)
	CreateVideo(ctx context.Context, input models.CreateVideoPost) (string, error)
	CreateLinkPost(ctx context.Context, input models.CreateLinkPost) (string, error)
	UpdateEssay(ctx context.Context, id string, input models.CreateEssayPost) error
	UpdateShort(ctx context.Context, id string, input models.CreateShortPost) error
	UpdateMusic(ctx context.Context, id string, input models.CreateMusicPost) error
	UpdatePhoto(ctx context.Context, id string, input models.CreatePhotoPost) error
	UpdateVideo(ctx context.Context, id string, input models.CreateVideoPost) error
	UpdateLinkPost(ctx context.Context, id string, input models.CreateLinkPost) error
	DeletePost(ctx context.Context, id string) error
}

// LinkLogic defines the business operations on external links.
type LinkLogic interface {
	// GetLinks returns all external links.
	GetLinks(ctx context.Context) ([]models.ExternalLink, error)

	// GetFeaturedLinks returns up to 4 featured external links.
	GetFeaturedLinks(ctx context.Context) ([]models.ExternalLink, error)

	// Write operations
	CreateLink(ctx context.Context, input models.CreateExternalLink) (string, error)
	UpdateLink(ctx context.Context, id string, input models.CreateExternalLink) error
	DeleteLink(ctx context.Context, id string) error
}

// AuthLogic defines authentication operations.
type AuthLogic interface {
	Login(ctx context.Context, username, password string) (string, error)
	ValidateToken(ctx context.Context, token string) error
}
