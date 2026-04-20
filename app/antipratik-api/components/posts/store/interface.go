package store

import (
	"context"
"github.com/pratikluitel/antipratik/components/posts/models"
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
	CreateEssayData(ctx context.Context, id string, input models.EssayPostInput) error
	CreateShortData(ctx context.Context, id string, input models.ShortPostInput) error
	CreateMusicData(ctx context.Context, id string, input models.MusicPostInput) error
	CreatePhotoData(ctx context.Context, id string, input models.PhotoPostInput) error
	CreateVideoData(ctx context.Context, id string, input models.VideoPostInput) error
	CreateLinkPostData(ctx context.Context, id string, input models.LinkPostInput) error
	UpdateEssay(ctx context.Context, id string, input models.EssayPostInput) error
	UpdateShort(ctx context.Context, id string, input models.ShortPostInput) error
	UpdateMusic(ctx context.Context, id string, input models.MusicPostInput) error
	UpdatePhoto(ctx context.Context, id string, input models.PhotoPostInput) error
	UpdateVideo(ctx context.Context, id string, input models.VideoPostInput) error
	UpdateLinkPost(ctx context.Context, id string, input models.LinkPostInput) error
	DeletePost(ctx context.Context, id string) error

	// GetAllTags returns all tag names sorted alphabetically.
	GetAllTags(ctx context.Context) ([]string, error)

	// Individual photo image operations
	AddPhotoImage(ctx context.Context, postID string, image models.PhotoImage) (*models.PhotoImage, error)
	GetPhotoImage(ctx context.Context, postID string, imageID int) (*models.PhotoImage, error)
	UpdatePhotoImage(ctx context.Context, postID string, imageID int, input models.UpdatePhotoImage) (*models.PhotoImage, error)
	DeletePhotoImage(ctx context.Context, postID string, imageID int) error
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

