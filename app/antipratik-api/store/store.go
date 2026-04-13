// Package store defines the database interfaces and their SQLite implementations.
package store

import (
	"context"
	"time"

	"github.com/pratikluitel/antipratik/models"
)

// PostStore handles all post-related database operations.
type PostStore interface {
	// GetPosts returns posts matching the optional type and tag filters,
	// sorted newest first. Empty slices mean "no filter" (return all).
	GetPosts(ctx context.Context, types, tags []string) ([]models.Post, error)

	// GetPostBySlug returns the essay with the given slug, or nil if not found.
	GetPostBySlug(ctx context.Context, slug string) (*models.EssayPost, error)

	// GetPostByID returns any post type by its ID, or an error if not found.
	GetPostByID(ctx context.Context, id string) (models.Post, error)

	// Write operations
	CreatePost(ctx context.Context, postType string, id string, createdAt string) error
	CreateEssayData(ctx context.Context, id string, input models.CreateEssayPost) error
	CreateShortData(ctx context.Context, id string, input models.CreateShortPost) error
	CreateMusicData(ctx context.Context, id string, input models.CreateMusicPost) error
	CreatePhotoData(ctx context.Context, id string, input models.CreatePhotoPost) error
	CreateVideoData(ctx context.Context, id string, input models.CreateVideoPost) error
	CreateLinkPostData(ctx context.Context, id string, input models.CreateLinkPost) error
	UpdateEssay(ctx context.Context, id string, input models.CreateEssayPost) error
	UpdateShort(ctx context.Context, id string, input models.CreateShortPost) error
	UpdateMusic(ctx context.Context, id string, input models.CreateMusicPost) error
	UpdatePhoto(ctx context.Context, id string, input models.CreatePhotoPost) error
	UpdateVideo(ctx context.Context, id string, input models.CreateVideoPost) error
	UpdateLinkPost(ctx context.Context, id string, input models.CreateLinkPost) error
	DeletePost(ctx context.Context, id string) error
}

// LinkStore handles all external link database operations.
type LinkStore interface {
	// GetLinks returns all external links.
	GetLinks(ctx context.Context) ([]models.ExternalLink, error)

	// GetFeaturedLinks returns up to 4 featured external links.
	GetFeaturedLinks(ctx context.Context) ([]models.ExternalLink, error)

	// GetLinkByID returns an external link by ID, or an error if not found.
	GetLinkByID(ctx context.Context, id string) (models.ExternalLink, error)

	// Write operations
	CreateLink(ctx context.Context, id string, input models.CreateExternalLink) error
	UpdateLink(ctx context.Context, id string, input models.CreateExternalLink) error
	DeleteLink(ctx context.Context, id string) error
}

// NewsletterStore handles newsletter subscriber database operations.
type NewsletterStore interface {
	// Subscribe inserts an email into newsletter_subscribers.
	// Returns ErrDuplicate if the email already exists.
	Subscribe(ctx context.Context, email string) error
}

// UserStore handles user authentication database operations.
type UserStore interface {
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByToken(ctx context.Context, token string) (*User, error)
	UpsertToken(ctx context.Context, username string, token string, expiresAt time.Time) error
}

// User represents an authenticated user.
type User struct {
	ID             string
	Username       string
	PasswordHash   string
	CurrentToken   *string
	TokenExpiresAt *time.Time
}
