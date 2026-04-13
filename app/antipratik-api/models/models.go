// Package models defines the shared data structures used across all layers.
// JSON tags use camelCase to match the TypeScript types in the frontend.
package models

import "mime/multipart"

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
	Tags               []string `json:"tags"`
	Title              string   `json:"title"`
	Slug               string   `json:"slug"`
	Excerpt            string   `json:"excerpt"`
	Body               string   `json:"body"`
	ReadingTimeMinutes int      `json:"readingTimeMinutes"`
}

func (e EssayPost) postType() string { return "essay" }

// ShortPost is a brief text post; tags serve as hashtags.
type ShortPost struct {
	ID        string   `json:"id"`
	Type      string   `json:"type"`
	CreatedAt string   `json:"createdAt"`
	Tags      []string `json:"tags"`
	Body      string   `json:"body"`
}

func (s ShortPost) postType() string { return "short" }

// MusicPost is a music track with album art and audio URL.
type MusicPost struct {
	ID              string   `json:"id"`
	Type            string   `json:"type"`
	CreatedAt       string   `json:"createdAt"`
	Tags            []string `json:"tags"`
	Title           string   `json:"title"`
	AlbumArt        string   `json:"albumArt"`
	AlbumArtTinyURL *string  `json:"albumArtTinyUrl,omitempty"`
	AudioURL        string   `json:"audioUrl"`
	Duration        int      `json:"duration"`
	Album           *string  `json:"album,omitempty"`
}

func (m MusicPost) postType() string { return "music" }

// PhotoImage is a single image within a PhotoPost.
type PhotoImage struct {
	URL               string  `json:"url"`
	Alt               string  `json:"alt"`
	Caption           *string `json:"caption,omitempty"`
	ThumbnailTinyURL  *string `json:"thumbnailTinyUrl,omitempty"`
	ThumbnailSmallURL *string `json:"thumbnailSmallUrl,omitempty"`
	ThumbnailMedURL   *string `json:"thumbnailMediumUrl,omitempty"`
	ThumbnailLargeURL *string `json:"thumbnailLargeUrl,omitempty"`
}

// PhotoPost contains one or more images and an optional location.
type PhotoPost struct {
	ID        string       `json:"id"`
	Type      string       `json:"type"`
	CreatedAt string       `json:"createdAt"`
	Tags      []string     `json:"tags"`
	Images    []PhotoImage `json:"images"`
	Location  *string      `json:"location,omitempty"`
}

func (p PhotoPost) postType() string { return "photo" }

// VideoPost is a video with thumbnail, URL, and duration.
type VideoPost struct {
	ID               string   `json:"id"`
	Type             string   `json:"type"`
	CreatedAt        string   `json:"createdAt"`
	Tags             []string `json:"tags"`
	Title            string   `json:"title"`
	ThumbnailURL     string   `json:"thumbnailUrl"`
	ThumbnailTinyURL *string  `json:"thumbnailTinyUrl,omitempty"`
	VideoURL         string   `json:"videoUrl"`
	Duration         int      `json:"duration"`
	Playlist         *string  `json:"playlist,omitempty"`
}

func (v VideoPost) postType() string { return "video" }

// LinkPost is a curated external link shared as a post.
type LinkPost struct {
	ID               string   `json:"id"`
	Type             string   `json:"type"`
	CreatedAt        string   `json:"createdAt"`
	Tags             []string `json:"tags"`
	Title            string   `json:"title"`
	URL              string   `json:"url"`
	Domain           string   `json:"domain"`
	Description      *string  `json:"description,omitempty"`
	ThumbnailURL     *string  `json:"thumbnailUrl,omitempty"`
	ThumbnailTinyURL *string  `json:"thumbnailTinyUrl,omitempty"`
	Category         *string  `json:"category,omitempty"`
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
	Featured    bool   `json:"featured"`
	Category    string `json:"category"`
}

// ── Filter ────────────────────────────────────────────────────────────────────

// FilterState represents the active filters for the posts feed.
type FilterState struct {
	ActiveTypes []string
	ActiveTags  []string
}

// ── Create/Update input types ─────────────────────────────────────────────────

type CreateEssayPost struct {
	Title              string   `json:"title"`
	Slug               string   `json:"slug"`
	Excerpt            string   `json:"excerpt"`
	Body               string   `json:"body"`
	ReadingTimeMinutes int      `json:"-"`
	Tags               []string `json:"tags"`
}

type CreateShortPost struct {
	Body string   `json:"body"`
	Tags []string `json:"tags"`
}

type CreateMusicPost struct {
	Title           string   `json:"title"`
	AlbumArt        string   `json:"albumArt"`
	AlbumArtTinyURL string   `json:"albumArtTinyUrl,omitempty"`
	AudioURL        string   `json:"audioURL"`
	Duration        int      `json:"duration"`
	Album           *string  `json:"album,omitempty"`
	Tags            []string `json:"tags"`
}

type CreatePhotoPost struct {
	Images   []PhotoImage `json:"images"`
	Location *string      `json:"location,omitempty"`
	Tags     []string     `json:"tags"`
}

type CreateVideoPost struct {
	Title            string   `json:"title"`
	ThumbnailURL     string   `json:"thumbnailURL"`
	ThumbnailTinyURL string   `json:"thumbnailTinyUrl,omitempty"`
	VideoURL         string   `json:"videoURL"`
	Duration         int      `json:"duration"`
	Playlist         *string  `json:"playlist,omitempty"`
	Tags             []string `json:"tags"`
}

type CreateLinkPost struct {
	Title            string   `json:"title"`
	URL              string   `json:"url"`
	Domain           string   `json:"-"`
	Description      *string  `json:"description,omitempty"`
	ThumbnailURL     *string  `json:"thumbnailURL,omitempty"`
	ThumbnailTinyURL *string  `json:"thumbnailTinyUrl,omitempty"`
	Category         *string  `json:"category,omitempty"`
	Tags             []string `json:"tags"`
}

type CreateExternalLink struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Domain      string `json:"-"`
	Description string `json:"description"`
	Featured    bool   `json:"featured"`
	Category    string `json:"category"`
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
	Title           *string
	AudioURL        *string
	AlbumArt        *string
	AlbumArtTinyURL *string
	Duration        *int
	Album           *string
	Tags            []string
}

// UpdatePhotoPost is the partial-update input for photo posts (multipart/form-data).
type UpdatePhotoPost struct {
	Images   []PhotoImage
	Location *string
	Tags     []string
}

// UpdateVideoPost is the partial-update input for video posts (multipart/form-data).
type UpdateVideoPost struct {
	Title            *string
	ThumbnailURL     *string
	ThumbnailTinyURL *string
	VideoURL         *string
	Duration         *int
	Playlist         *string
	Tags             []string
}

// UpdateLinkPost is the partial-update input for link posts (multipart/form-data).
type UpdateLinkPost struct {
	Title            *string
	URL              *string
	ThumbnailURL     *string
	ThumbnailTinyURL *string
	Description      *string
	Category         *string
	Tags             []string
}

// UpdateExternalLink is the partial-update input for external links (JSON body).
type UpdateExternalLink struct {
	Title       *string `json:"title"`
	URL         *string `json:"url"`
	Description *string `json:"description"`
	Featured    *bool   `json:"featured"`
	Category    *string `json:"category"`
}

// FileInput bundles a multipart file and its header for upload operations.
type FileInput struct {
	File   multipart.File
	Header *multipart.FileHeader
}
