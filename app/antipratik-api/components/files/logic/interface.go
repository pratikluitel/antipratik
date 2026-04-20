// Package logic contains the files component business logic layer.
package logic

import (
	"context"
	"mime/multipart"
)

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

// UploadLogic handles file storage and thumbnail generation for uploaded assets.
type UploadLogic interface {
	// UploadPhotoFiles stores multiple photo images (with 4 thumbnail variants each) for the
	// given post. File IDs follow <postID>-<uuid>.<ext>; thumbnail IDs follow
	// <postID>-<uuid>-<size>.<ext> where size is tiny/small/medium/large.
	UploadPhotoFiles(ctx context.Context, postID string, files []FileInput) ([]PhotoImageResult, error)

	// UploadMusicFiles stores audio and/or album art for the given post.
	// Either file may be nil (to skip that upload), but at least one must be non-nil.
	// Audio is stored at music/<postID>.<ext>; album art at photos/<postID>-albumart.<ext>.
	UploadMusicFiles(ctx context.Context, postID string, audioFile *FileInput, albumArtFile *FileInput) (MusicFilesResult, error)

	// UploadThumbnail stores a single thumbnail image plus a 20px-wide tiny variant.
	// suffix is appended to the postID in the stored file name (e.g. "thumb").
	UploadThumbnail(ctx context.Context, postID string, suffix string, file FileInput) (ThumbnailResult, error)
}
