package files

import "mime/multipart"

// FileInput bundles a multipart file and its header for upload operations.
type FileInput struct {
	File   multipart.File
	Header *multipart.FileHeader
}

// PhotoImageResult holds the URL fields for a single uploaded photo image.
type PhotoImageResult struct {
	OriginalURL       string
	ThumbnailTinyURL  string
	ThumbnailSmallURL string
	ThumbnailMedURL   string
	ThumbnailLargeURL string
}

// ThumbnailResult holds the URL fields for an uploaded thumbnail (all 4 sizes).
type ThumbnailResult struct {
	URL      string
	TinyURL  string
	SmallURL string
	MedURL   string
	LargeURL string
}

// MusicFilesResult holds the URL fields from a music post file upload.
type MusicFilesResult struct {
	AudioURL         string
	AlbumArtURL      string
	AlbumArtTinyURL  string
	AlbumArtSmallURL string
	AlbumArtMedURL   string
	AlbumArtLargeURL string
}

// VideoFileResult holds the relative URL for an uploaded video file.
type VideoFileResult struct {
	VideoURL string // relative: /files/videos/<postID>.<ext>
}
