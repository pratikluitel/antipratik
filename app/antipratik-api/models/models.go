// Package models defines the shared data structures used across all layers.
// JSON tags use camelCase to match the TypeScript types in the frontend.
package models

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
	ID        string   `json:"id"`
	Type      string   `json:"type"`
	CreatedAt string   `json:"createdAt"`
	Tags      []string `json:"tags"`
	Title     string   `json:"title"`
	AlbumArt  string   `json:"albumArt"`
	AudioURL  string   `json:"audioUrl"`
	Duration  int      `json:"duration"`
	Album     *string  `json:"album,omitempty"`
}

func (m MusicPost) postType() string { return "music" }

// PhotoImage is a single image within a PhotoPost.
type PhotoImage struct {
	URL     string  `json:"url"`
	Alt     string  `json:"alt"`
	Caption *string `json:"caption,omitempty"`
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
	ID           string   `json:"id"`
	Type         string   `json:"type"`
	CreatedAt    string   `json:"createdAt"`
	Tags         []string `json:"tags"`
	Title        string   `json:"title"`
	ThumbnailURL string   `json:"thumbnailUrl"`
	VideoURL     string   `json:"videoUrl"`
	Duration     int      `json:"duration"`
	Playlist     *string  `json:"playlist,omitempty"`
}

func (v VideoPost) postType() string { return "video" }

// LinkPost is a curated external link shared as a post.
type LinkPost struct {
	ID           string   `json:"id"`
	Type         string   `json:"type"`
	CreatedAt    string   `json:"createdAt"`
	Tags         []string `json:"tags"`
	Title        string   `json:"title"`
	URL          string   `json:"url"`
	Domain       string   `json:"domain"`
	Description  *string  `json:"description,omitempty"`
	ThumbnailURL *string  `json:"thumbnailUrl,omitempty"`
	Category     *string  `json:"category,omitempty"`
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
