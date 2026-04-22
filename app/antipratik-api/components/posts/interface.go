package posts

import (
	"context"
	"net/http"
)

// PostLogic defines the business operations on
type PostLogic interface {
	// GetPosts returns posts matching the given filter, newest first.
	GetPosts(ctx context.Context, filter FilterState) ([]Post, error)

	// GetTags returns all tag names sorted alphabetically.
	GetTags(ctx context.Context) ([]string, error)

	// GetPost returns the essay with the given slug, or nil if not found.
	GetPost(ctx context.Context, slug string) (*EssayPost, error)

	// GetPostsByIDs returns posts for each given ID, preserving the order of ids.
	// IDs not found are silently skipped.
	GetPostsByIDs(ctx context.Context, ids []string) ([]Post, error)

	// Write operations
	// Essay and Short use auto-generated IDs.
	// Music, Photo, Video, and LinkPost accept a preID: if non-empty it is used as the post ID
	// (so the API layer can generate the ID before uploading files and keep them in sync);
	// if empty, a new UUID is generated.
	CreateEssay(ctx context.Context, input EssayPostInput) (EssayPost, error)
	CreateShort(ctx context.Context, input ShortPostInput) (ShortPost, error)
	CreateMusic(ctx context.Context, preID string, input MusicPostInput) (MusicPost, error)
	CreatePhoto(ctx context.Context, preID string, input PhotoPostInput) (PhotoPost, error)
	CreateVideo(ctx context.Context, preID string, input VideoPostInput) (VideoPost, error)
	CreateLinkPost(ctx context.Context, preID string, input LinkPostInput) (LinkPost, error)
	UpdateEssay(ctx context.Context, id string, input UpdateEssayPost) (EssayPost, error)
	UpdateShort(ctx context.Context, id string, input UpdateShortPost) (ShortPost, error)
	UpdateMusic(ctx context.Context, id string, input UpdateMusicPost) (MusicPost, error)
	UpdatePhoto(ctx context.Context, id string, input PhotoPostInput) (PhotoPost, error)
	UpdateVideo(ctx context.Context, id string, input UpdateVideoPost) (VideoPost, error)
	UpdateLinkPost(ctx context.Context, id string, input UpdateLinkPost) (LinkPost, error)
	DeletePost(ctx context.Context, id string) error

	// Individual photo image operations
	AddPhotoImage(ctx context.Context, postID string, image PhotoImage) (*PhotoImage, error)
	GetPhotoImage(ctx context.Context, postID string, imageIDStr string) (*PhotoImage, error)
	UpdatePhotoImage(ctx context.Context, postID string, imageIDStr string, input UpdatePhotoImage) (*PhotoImage, error)
	// DeletePhotoImage returns nil, nil when the image is not found (API maps to 404).
	// Returns a ValidationError if the post has only one image.
	DeletePhotoImage(ctx context.Context, postID string, imageIDStr string) (notFound bool, err error)
}

// LinkLogic defines the business operations on external links.
type LinkLogic interface {
	// GetLinks returns all external links.
	GetLinks(ctx context.Context) ([]ExternalLink, error)

	// GetFeaturedLinks returns up to 4 featured external links.
	GetFeaturedLinks(ctx context.Context) ([]ExternalLink, error)

	// Write operations
	CreateLink(ctx context.Context, input CreateExternalLink) (string, error)
	UpdateLink(ctx context.Context, id string, input UpdateExternalLink) (ExternalLink, error)
	DeleteLink(ctx context.Context, id string) error
}

// PostStore handles all post-related database operations.
type PostStore interface {
	// GetPosts returns posts matching the optional type and tag filters,
	// sorted newest first. Empty slices mean "no filter" (return all).
	GetPosts(ctx context.Context, types, tags []string) ([]Post, error)

	// GetPostBySlug returns the essay with the given slug, or nil if not found.
	GetPostBySlug(ctx context.Context, slug string) (*EssayPost, error)

	// GetPostByID returns any post type by its ID, or an error if not found.
	GetPostByID(ctx context.Context, id string) (Post, error)

	// GetPostsByIDs returns posts for each given ID, preserving the order of ids.
	// IDs that are not found are silently skipped.
	GetPostsByIDs(ctx context.Context, ids []string) ([]Post, error)

	// Write operations
	CreatePost(ctx context.Context, postType string, id string, createdAt string) error
	CreateEssayData(ctx context.Context, id string, input EssayPostInput) error
	CreateShortData(ctx context.Context, id string, input ShortPostInput) error
	CreateMusicData(ctx context.Context, id string, input MusicPostInput) error
	CreatePhotoData(ctx context.Context, id string, input PhotoPostInput) error
	CreateVideoData(ctx context.Context, id string, input VideoPostInput) error
	CreateLinkPostData(ctx context.Context, id string, input LinkPostInput) error
	UpdateEssay(ctx context.Context, id string, input EssayPostInput) error
	UpdateShort(ctx context.Context, id string, input ShortPostInput) error
	UpdateMusic(ctx context.Context, id string, input MusicPostInput) error
	UpdatePhoto(ctx context.Context, id string, input PhotoPostInput) error
	UpdateVideo(ctx context.Context, id string, input VideoPostInput) error
	UpdateLinkPost(ctx context.Context, id string, input LinkPostInput) error
	DeletePost(ctx context.Context, id string) error

	// GetAllTags returns all tag names sorted alphabetically.
	GetAllTags(ctx context.Context) ([]string, error)

	// Individual photo image operations
	AddPhotoImage(ctx context.Context, postID string, image PhotoImage) (*PhotoImage, error)
	GetPhotoImage(ctx context.Context, postID string, imageID int) (*PhotoImage, error)
	UpdatePhotoImage(ctx context.Context, postID string, imageID int, input UpdatePhotoImage) (*PhotoImage, error)
	DeletePhotoImage(ctx context.Context, postID string, imageID int) error
}

// LinkStore handles all external link database operations.
type LinkStore interface {
	// GetLinks returns all external links.
	GetLinks(ctx context.Context) ([]ExternalLink, error)

	// GetFeaturedLinks returns up to 4 featured external links.
	GetFeaturedLinks(ctx context.Context) ([]ExternalLink, error)

	// GetLinkByID returns an external link by ID, or an error if not found.
	GetLinkByID(ctx context.Context, id string) (ExternalLink, error)

	// Write operations
	CreateLink(ctx context.Context, id string, input CreateExternalLink) error
	UpdateLink(ctx context.Context, id string, input CreateExternalLink) error
	DeleteLink(ctx context.Context, id string) error
}

// PostHandler handles HTTP requests for post resources.
type PostHandler interface {
	GetPosts(w http.ResponseWriter, r *http.Request)
	GetPost(w http.ResponseWriter, r *http.Request)
	GetTags(w http.ResponseWriter, r *http.Request)
	CreateEssay(w http.ResponseWriter, r *http.Request)
	CreateShort(w http.ResponseWriter, r *http.Request)
	CreateMusic(w http.ResponseWriter, r *http.Request)
	CreatePhoto(w http.ResponseWriter, r *http.Request)
	CreateVideo(w http.ResponseWriter, r *http.Request)
	CreateLinkPost(w http.ResponseWriter, r *http.Request)
	UpdateEssay(w http.ResponseWriter, r *http.Request)
	UpdateShort(w http.ResponseWriter, r *http.Request)
	UpdateMusic(w http.ResponseWriter, r *http.Request)
	UpdatePhoto(w http.ResponseWriter, r *http.Request)
	UpdateVideo(w http.ResponseWriter, r *http.Request)
	UpdateLinkPost(w http.ResponseWriter, r *http.Request)
	DeletePost(w http.ResponseWriter, r *http.Request)
	AddPhotoImage(w http.ResponseWriter, r *http.Request)
	GetPhotoImage(w http.ResponseWriter, r *http.Request)
	UpdatePhotoImage(w http.ResponseWriter, r *http.Request)
	DeletePhotoImage(w http.ResponseWriter, r *http.Request)
}

// LinkHandler handles HTTP requests for external link resources.
type LinkHandler interface {
	GetLinks(w http.ResponseWriter, r *http.Request)
	GetFeaturedLinks(w http.ResponseWriter, r *http.Request)
	CreateLink(w http.ResponseWriter, r *http.Request)
	UpdateLink(w http.ResponseWriter, r *http.Request)
	DeleteLink(w http.ResponseWriter, r *http.Request)
}

// PostsService exposes post retrieval to other components via posts/services.
type PostsService interface {
	GetPostsByIDs(ctx context.Context, ids []string) ([]Post, error)
}
