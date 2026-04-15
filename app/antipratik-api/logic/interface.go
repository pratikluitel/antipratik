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
	CreateEssay(ctx context.Context, input models.EssayPostInput) (models.EssayPost, error)
	CreateShort(ctx context.Context, input models.ShortPostInput) (models.ShortPost, error)
	CreateMusic(ctx context.Context, preID string, input models.MusicPostInput) (models.MusicPost, error)
	CreatePhoto(ctx context.Context, preID string, input models.PhotoPostInput) (models.PhotoPost, error)
	CreateVideo(ctx context.Context, preID string, input models.VideoPostInput) (models.VideoPost, error)
	CreateLinkPost(ctx context.Context, preID string, input models.LinkPostInput) (models.LinkPost, error)
	UpdateEssay(ctx context.Context, id string, input models.UpdateEssayPost) (models.EssayPost, error)
	UpdateShort(ctx context.Context, id string, input models.UpdateShortPost) (models.ShortPost, error)
	UpdateMusic(ctx context.Context, id string, input models.UpdateMusicPost) (models.MusicPost, error)
	UpdatePhoto(ctx context.Context, id string, input models.PhotoPostInput) (models.PhotoPost, error)
	UpdateVideo(ctx context.Context, id string, input models.UpdateVideoPost) (models.VideoPost, error)
	UpdateLinkPost(ctx context.Context, id string, input models.UpdateLinkPost) (models.LinkPost, error)
	DeletePost(ctx context.Context, id string) error

	// Individual photo image operations
	AddPhotoImage(ctx context.Context, postID string, image models.PhotoImage) (*models.PhotoImage, error)
	GetPhotoImage(ctx context.Context, postID string, imageIDStr string) (*models.PhotoImage, error)
	UpdatePhotoImage(ctx context.Context, postID string, imageIDStr string, input models.UpdatePhotoImage) (*models.PhotoImage, error)
	// DeletePhotoImage returns nil, nil when the image is not found (API maps to 404).
	// Returns a ValidationError if the post has only one image.
	DeletePhotoImage(ctx context.Context, postID string, imageIDStr string) (notFound bool, err error)
}

// LinkLogic defines the business operations on external links.
type LinkLogic interface {
	// GetLinks returns all external links.
	GetLinks(ctx context.Context) ([]models.ExternalLink, error)

	// GetFeaturedLinks returns up to 4 featured external links.
	GetFeaturedLinks(ctx context.Context) ([]models.ExternalLink, error)

	// Write operations
	CreateLink(ctx context.Context, input models.CreateExternalLink) (string, error)
	UpdateLink(ctx context.Context, id string, input models.UpdateExternalLink) (models.ExternalLink, error)
	DeleteLink(ctx context.Context, id string) error
}

// NewsletterLogic defines newsletter subscription operations.
type NewsletterLogic interface {
	// Subscribe validates the email and persists it as a subscriber.
	Subscribe(ctx context.Context, email string) error
}

// AuthLogic defines authentication operations.
type AuthLogic interface {
	Login(ctx context.Context, username, password string) (string, error)
	ValidateToken(ctx context.Context, token string) error
}

// SetupLogic defines application bootstrapping operations that run at startup.
// It ensures main.go never calls the store layer directly.
type SetupLogic interface {
	// UpsertAdminUser creates or updates the admin user with the given password.
	UpsertAdminUser(ctx context.Context, password string) error
	// GetOrCreateJWTSecret returns the persisted JWT signing secret,
	// generating and storing a new one if none exists yet.
	GetOrCreateJWTSecret(ctx context.Context) (string, error)
}
