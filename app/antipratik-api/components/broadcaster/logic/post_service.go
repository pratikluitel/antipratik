package logic

import "context"

// PostService is the interface the broadcaster logic uses to fetch post data
// for rendering email templates. Implemented by an adapter in main.go that
// wraps postslogic.PostLogic — keeping the broadcaster independent of the
// posts component's concrete types.
type PostService interface {
	GetPostsByIDs(ctx context.Context, ids []string) ([]PostSummary, error)
}

// PostSummary holds the subset of post data needed to render email templates.
type PostSummary struct {
	ID                string
	Type              string // essay | photo | music | short | video | link
	Title             string
	Slug              string
	Excerpt           string
	Body              string // essays: full HTML body
	ImageURL           string // photos: first image URL
	ThumbnailMediumURL string // photos/video/link: medium thumbnail URL
	ThumbnailLargeURL  string // photos: large thumbnail URL
	AlbumArtMediumURL  string // music: album art medium thumbnail
	VideoURL           string // video: video URL for click-through
	LinkURL            string // link: external URL for click-through
	Category           string // link: category (e.g. "video") for click-through routing
	Domain             string // link: domain for display
	CreatedAt         string
}
