// Package models defines the shared data structures used across all layers.
// JSON tags use camelCase to match the TypeScript types in the frontend.
package posts

// PostType is the discriminator string stored in the posts table.
type PostType = string

// Named post type constants — use these instead of bare string literals.
const (
	PostTypeEssay PostType = "essay"
	PostTypeShort PostType = "short"
	PostTypeMusic PostType = "music"
	PostTypePhoto PostType = "photo"
	PostTypeVideo PostType = "video"
	PostTypeLink  PostType = "link"
)

// Post is an interface satisfied by all content types.
// The unexported method prevents external implementations.
type Post interface {
	postType() string
}

// ── Base fields shared by all post types ─────────────────────────────────────

// EssayPost is a long-form written piece with a slug and reading time.
type EssayPost struct {
	ID                 string   `json:"id"`
	Type               string   `json:"type"`
	CreatedAt          string   `json:"createdAt"`
	Title              string   `json:"title"`
	Slug               string   `json:"slug"`
	Excerpt            string   `json:"excerpt"`
	Body               string   `json:"body"`
	Tags               []string `json:"tags"`
	ReadingTimeMinutes int      `json:"readingTimeMinutes"`
}

func (e EssayPost) postType() string { return "essay" }

// ShortPost is a brief text post; tags serve as hashtags.
type ShortPost struct {
	ID        string   `json:"id"`
	Type      string   `json:"type"`
	CreatedAt string   `json:"createdAt"`
	Body      string   `json:"body"`
	Tags      []string `json:"tags"`
}

func (s ShortPost) postType() string { return "short" }

// MusicPost is a music track with album art and audio URL.
type MusicPost struct {
	AlbumArtTinyURL  *string  `json:"albumArtTinyUrl,omitempty"`
	AlbumArtSmallURL *string  `json:"albumArtSmallUrl,omitempty"`
	AlbumArtMedURL   *string  `json:"albumArtMediumUrl,omitempty"`
	AlbumArtLargeURL *string  `json:"albumArtLargeUrl,omitempty"`
	Album            *string  `json:"album,omitempty"`
	ID               string   `json:"id"`
	Type             string   `json:"type"`
	CreatedAt        string   `json:"createdAt"`
	Title            string   `json:"title"`
	AlbumArt         string   `json:"albumArt"`
	AudioURL         string   `json:"audioUrl"`
	Tags             []string `json:"tags"`
	Duration         int      `json:"duration"`
}

func (m MusicPost) postType() string { return "music" }

// PhotoImage is a single image within a PhotoPost.
type PhotoImage struct {
	Caption           *string `json:"caption,omitempty"`
	ThumbnailTinyURL  *string `json:"thumbnailTinyUrl,omitempty"`
	ThumbnailSmallURL *string `json:"thumbnailSmallUrl,omitempty"`
	ThumbnailMedURL   *string `json:"thumbnailMediumUrl,omitempty"`
	ThumbnailLargeURL *string `json:"thumbnailLargeUrl,omitempty"`
	URL               string  `json:"url"`
	Alt               string  `json:"alt"`
	ID                int     `json:"id"`
}

// UpdatePhotoImage is a partial-update input for a single photo image.
// Nil fields mean "leave unchanged".
type UpdatePhotoImage struct {
	Caption *string
	Alt     *string
}

// PhotoPost contains one or more images and an optional location.
type PhotoPost struct {
	Location  *string      `json:"location,omitempty"`
	ID        string       `json:"id"`
	Type      string       `json:"type"`
	CreatedAt string       `json:"createdAt"`
	Images    []PhotoImage `json:"images"`
	Tags      []string     `json:"tags"`
}

func (p PhotoPost) postType() string { return "photo" }

// VideoPost is a video with thumbnail, URL, and duration.
type VideoPost struct {
	ThumbnailTinyURL  *string  `json:"thumbnailTinyUrl,omitempty"`
	ThumbnailSmallURL *string  `json:"thumbnailSmallUrl,omitempty"`
	ThumbnailMedURL   *string  `json:"thumbnailMediumUrl,omitempty"`
	ThumbnailLargeURL *string  `json:"thumbnailLargeUrl,omitempty"`
	Playlist          *string  `json:"playlist,omitempty"`
	ID                string   `json:"id"`
	Type              string   `json:"type"`
	CreatedAt         string   `json:"createdAt"`
	Title             string   `json:"title"`
	ThumbnailURL      string   `json:"thumbnailUrl"`
	VideoURL          string   `json:"videoUrl"`
	Tags              []string `json:"tags"`
	Duration          int      `json:"duration"`
}

func (v VideoPost) postType() string { return "video" }

// LinkPost is a curated external link shared as a post.
type LinkPost struct {
	Description       *string  `json:"description,omitempty"`
	ThumbnailURL      *string  `json:"thumbnailUrl,omitempty"`
	ThumbnailTinyURL  *string  `json:"thumbnailTinyUrl,omitempty"`
	ThumbnailSmallURL *string  `json:"thumbnailSmallUrl,omitempty"`
	ThumbnailMedURL   *string  `json:"thumbnailMediumUrl,omitempty"`
	ThumbnailLargeURL *string  `json:"thumbnailLargeUrl,omitempty"`
	Category          *string  `json:"category,omitempty"`
	ID                string   `json:"id"`
	Type              string   `json:"type"`
	CreatedAt         string   `json:"createdAt"`
	Title             string   `json:"title"`
	URL               string   `json:"url"`
	Domain            string   `json:"domain"`
	Tags              []string `json:"tags"`
}

func (l LinkPost) postType() string { return "link" }

// ── External links ────────────────────────────────────────────────────────────

// ExternalLink is a curated link to an external platform (not a post).
// Icons are derived from Category on the frontend — no iconUrl field.
type ExternalLink struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Domain      string `json:"domain"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Featured    bool   `json:"featured"`
}

// ── Filter ────────────────────────────────────────────────────────────────────

// FilterState represents the active filters for the posts feed.
type FilterState struct {
	ActiveTypes []string
	ActiveTags  []string
}

// ── Create/Update input types ─────────────────────────────────────────────────

type EssayPostInput struct {
	Title              string   `json:"title"`
	Slug               string   `json:"slug"`
	Excerpt            string   `json:"excerpt"`
	Body               string   `json:"body"`
	Tags               []string `json:"tags"`
	ReadingTimeMinutes int      `json:"-"`
}

type ShortPostInput struct {
	Body string   `json:"body"`
	Tags []string `json:"tags"`
}

type MusicPostInput struct {
	AlbumArtTinyURL  *string  `json:"albumArtTinyUrl,omitempty"`
	AlbumArtSmallURL *string  `json:"albumArtSmallUrl,omitempty"`
	AlbumArtMedURL   *string  `json:"albumArtMediumUrl,omitempty"`
	AlbumArtLargeURL *string  `json:"albumArtLargeUrl,omitempty"`
	Album            *string  `json:"album,omitempty"`
	Title            string   `json:"title"`
	AlbumArt         string   `json:"albumArt"`
	AudioURL         string   `json:"audioURL"`
	Tags             []string `json:"tags"`
	Duration         int      `json:"duration"`
}

// PhotoPostInput is the input type for both creating and updating photo posts.
// There is no separate Update type because all fields are optional on update
// (images replace in full when provided; location and tags use nil = no-change semantics).
type PhotoPostInput struct {
	Images   []PhotoImage `json:"images"`
	Location *string      `json:"location,omitempty"`
	Tags     []string     `json:"tags"`
}

type VideoPostInput struct {
	ThumbnailTinyURL  *string  `json:"thumbnailTinyUrl,omitempty"`
	ThumbnailSmallURL *string  `json:"thumbnailSmallUrl,omitempty"`
	ThumbnailMedURL   *string  `json:"thumbnailMediumUrl,omitempty"`
	ThumbnailLargeURL *string  `json:"thumbnailLargeUrl,omitempty"`
	Playlist          *string  `json:"playlist,omitempty"`
	Title             string   `json:"title"`
	ThumbnailURL      string   `json:"thumbnailURL"`
	VideoURL          string   `json:"videoURL"`
	Tags              []string `json:"tags"`
	Duration          int      `json:"duration"`
}

type LinkPostInput struct {
	Title             string   `json:"title"`
	URL               string   `json:"url"`
	Domain            string   `json:"-"`
	Description       *string  `json:"description,omitempty"`
	ThumbnailURL      *string  `json:"thumbnailURL,omitempty"`
	ThumbnailTinyURL  *string  `json:"thumbnailTinyUrl,omitempty"`
	ThumbnailSmallURL *string  `json:"thumbnailSmallUrl,omitempty"`
	ThumbnailMedURL   *string  `json:"thumbnailMediumUrl,omitempty"`
	ThumbnailLargeURL *string  `json:"thumbnailLargeUrl,omitempty"`
	Category          *string  `json:"category,omitempty"`
	Tags              []string `json:"tags"`
}

type CreateExternalLink struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Domain      string `json:"-"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Featured    bool   `json:"featured"`
}

// ── Partial Update input types ────────────────────────────────────────────────

// UpdateEssayPost is the partial-update input for essay posts (JSON body).
type UpdateEssayPost struct {
	Title   *string  `json:"title"`
	Slug    *string  `json:"slug"`
	Excerpt *string  `json:"excerpt"`
	Body    *string  `json:"body"`
	Tags    []string `json:"tags"`
}

// UpdateShortPost is the partial-update input for short posts (JSON body).
type UpdateShortPost struct {
	Body *string  `json:"body"`
	Tags []string `json:"tags"`
}

// UpdateMusicPost is the partial-update input for music posts (multipart/form-data).
type UpdateMusicPost struct {
	Title            *string
	AudioURL         *string
	AlbumArt         *string
	AlbumArtTinyURL  *string
	AlbumArtSmallURL *string
	AlbumArtMedURL   *string
	AlbumArtLargeURL *string
	Album            *string
	Tags             []string
}

// UpdateVideoPost is the partial-update input for video posts (multipart/form-data).
type UpdateVideoPost struct {
	Title             *string
	ThumbnailURL      *string
	ThumbnailTinyURL  *string
	ThumbnailSmallURL *string
	ThumbnailMedURL   *string
	ThumbnailLargeURL *string
	VideoURL          *string
	Duration          *int
	Playlist          *string
	Tags              []string
}

// UpdateLinkPost is the partial-update input for link posts (multipart/form-data).
type UpdateLinkPost struct {
	Title             *string
	URL               *string
	ThumbnailURL      *string
	ThumbnailTinyURL  *string
	ThumbnailSmallURL *string
	ThumbnailMedURL   *string
	ThumbnailLargeURL *string
	Description       *string
	Category          *string
	Tags              []string
}

// UpdateExternalLink is the partial-update input for external links (JSON body).
type UpdateExternalLink struct {
	Title       *string `json:"title"`
	URL         *string `json:"url"`
	Description *string `json:"description"`
	Featured    *bool   `json:"featured"`
	Category    *string `json:"category"`
}
