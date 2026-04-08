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
	// Essay and Short use auto-generated IDs.
	// Music, Photo, Video, and LinkPost accept a preID: if non-empty it is used as the post ID
	// (so the API layer can generate the ID before uploading files and keep them in sync);
	// if empty, a new UUID is generated.
	CreateEssay(ctx context.Context, input models.CreateEssayPost) (models.EssayPost, error)
	CreateShort(ctx context.Context, input models.CreateShortPost) (models.ShortPost, error)
	CreateMusic(ctx context.Context, preID string, input models.CreateMusicPost) (models.MusicPost, error)
	CreatePhoto(ctx context.Context, preID string, input models.CreatePhotoPost) (models.PhotoPost, error)
	CreateVideo(ctx context.Context, preID string, input models.CreateVideoPost) (models.VideoPost, error)
	CreateLinkPost(ctx context.Context, preID string, input models.CreateLinkPost) (models.LinkPost, error)
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
