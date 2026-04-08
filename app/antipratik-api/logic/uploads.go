// Package logic — upload service: validates, stores files, and generates thumbnails.
package logic

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/pratikluitel/antipratik/store"
)

// FileInput bundles a multipart file and its header for upload operations.
type FileInput struct {
	File   multipart.File
	Header *multipart.FileHeader
}

// PhotoImageResult holds the URL fields for a single uploaded photo image.
type PhotoImageResult struct {
	OriginalURL       string
	ThumbnailSmallURL string
	ThumbnailMedURL   string
	ThumbnailLargeURL string
}

// MusicFilesResult holds the URL fields from a music post file upload.
type MusicFilesResult struct {
	AudioURL    string
	AlbumArtURL string // empty string if no album art was uploaded
}

// UploadLogic handles file storage and thumbnail generation for uploaded assets.
type UploadLogic interface {
	// UploadPhotoFiles stores multiple photo images (with 3 thumbnail variants each) for the
	// given post. File IDs follow <postID>-<index>.<ext>; thumbnail IDs follow
	// <postID>-<index>-<size>.<ext> where size is small/medium/large.
	UploadPhotoFiles(ctx context.Context, postID string, files []FileInput) ([]PhotoImageResult, error)

	// UploadMusicFiles stores audio and/or album art for the given post.
	// Either file may be nil (to skip that upload), but at least one must be non-nil.
	// Audio is stored at music/<postID>.<ext>; album art at photos/<postID>-albumart.<ext>.
	UploadMusicFiles(ctx context.Context, postID string, audioFile *FileInput, albumArtFile *FileInput) (MusicFilesResult, error)

	// UploadThumbnail stores a single thumbnail image and returns its serving URL.
	// suffix is appended to the postID in the stored file name (e.g. "thumb").
	UploadThumbnail(ctx context.Context, postID string, suffix string, file FileInput) (string, error)
}

// UploadService implements UploadLogic.
type UploadService struct {
	files store.FileStore
}

// NewUploadService constructs an UploadService.
func NewUploadService(files store.FileStore) *UploadService {
	return &UploadService{files: files}
}

var allowedPhotoExts = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
var allowedMusicExts = map[string]bool{".mp3": true, ".wav": true, ".ogg": true, ".m4a": true}

// storageExt returns the file extension used when storing the encoded image.
// WebP inputs are encoded as JPEG, so they are stored with a .jpg extension.
func storageExt(ext string) string {
	if ext == ".webp" {
		return ".jpg"
	}
	return ext
}

// UploadPhotoFiles implements UploadLogic.
func (s *UploadService) UploadPhotoFiles(ctx context.Context, postID string, files []FileInput) ([]PhotoImageResult, error) {
	if err := requireNonEmpty("postId", postID); err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, validationErr("at least one image file is required")
	}

	results := make([]PhotoImageResult, 0, len(files))
	for i, fi := range files {
		ext := strings.ToLower(filepath.Ext(fi.Header.Filename))
		if !allowedPhotoExts[ext] {
			return nil, validationErr(fmt.Sprintf("images[%d]: file must be one of jpg, jpeg, png, webp — got %q", i, ext))
		}

		src, err := decodeImage(fi.File, ext)
		if err != nil {
			return nil, validationErr(fmt.Sprintf("images[%d]: could not decode image: %v", i, err))
		}

		sExt := storageExt(ext)
		ct := contentTypeForExt(sExt)
		fileID := fmt.Sprintf("%s-%d%s", postID, i, sExt)
		origKey := "photos/" + fileID

		origBuf, err := encodeImage(src, ext)
		if err != nil {
			return nil, fmt.Errorf("UploadPhotoFiles encode original[%d]: %w", i, err)
		}
		if err := s.files.Put(ctx, origKey, bytes.NewReader(origBuf), ct); err != nil {
			return nil, fmt.Errorf("UploadPhotoFiles store original[%d]: %w", i, err)
		}

		sizes := []struct {
			name     string
			maxWidth int
		}{
			{"small", 300},
			{"medium", 600},
			{"large", 1200},
		}
		thumbURLs := make([]string, len(sizes))
		for j, sz := range sizes {
			thumb := resizeImage(src, sz.maxWidth)
			buf, err := encodeImage(thumb, ext)
			if err != nil {
				return nil, fmt.Errorf("UploadPhotoFiles encode thumbnail[%d][%s]: %w", i, sz.name, err)
			}
			thumbID := fmt.Sprintf("%s-%d-%s%s", postID, i, sz.name, sExt)
			if err := s.files.Put(ctx, "thumbnails/"+thumbID, bytes.NewReader(buf), ct); err != nil {
				return nil, fmt.Errorf("UploadPhotoFiles store thumbnail[%d][%s]: %w", i, sz.name, err)
			}
			thumbURLs[j] = "/thumbnails/" + thumbID
		}

		results = append(results, PhotoImageResult{
			OriginalURL:       "/files/" + fileID,
			ThumbnailSmallURL: thumbURLs[0],
			ThumbnailMedURL:   thumbURLs[1],
			ThumbnailLargeURL: thumbURLs[2],
		})
	}
	return results, nil
}

// UploadMusicFiles implements UploadLogic.
func (s *UploadService) UploadMusicFiles(ctx context.Context, postID string, audioFile *FileInput, albumArtFile *FileInput) (MusicFilesResult, error) {
	if err := requireNonEmpty("postId", postID); err != nil {
		return MusicFilesResult{}, err
	}

	var result MusicFilesResult

	// Upload audio (optional)
	if audioFile != nil {
		audioExt := strings.ToLower(filepath.Ext(audioFile.Header.Filename))
		if !allowedMusicExts[audioExt] {
			return MusicFilesResult{}, validationErr(fmt.Sprintf("audioFile must be one of mp3, wav, ogg, m4a — got %q", audioExt))
		}
		audioFileID := postID + audioExt
		if err := s.files.Put(ctx, "music/"+audioFileID, audioFile.File, contentTypeForExt(audioExt)); err != nil {
			return MusicFilesResult{}, fmt.Errorf("UploadMusicFiles store audio: %w", err)
		}
		result.AudioURL = "/files/" + audioFileID
	}

	// Upload album art (optional)
	if albumArtFile != nil {
		artExt := strings.ToLower(filepath.Ext(albumArtFile.Header.Filename))
		if !allowedPhotoExts[artExt] {
			return MusicFilesResult{}, validationErr(fmt.Sprintf("albumArtFile must be one of jpg, jpeg, png, webp — got %q", artExt))
		}
		artSExt := storageExt(artExt)
		artFileID := postID + "-albumart" + artSExt
		if err := s.files.Put(ctx, "photos/"+artFileID, albumArtFile.File, contentTypeForExt(artSExt)); err != nil {
			return MusicFilesResult{}, fmt.Errorf("UploadMusicFiles store album art: %w", err)
		}
		result.AlbumArtURL = "/files/" + artFileID
	}

	return result, nil
}

// UploadThumbnail implements UploadLogic.
func (s *UploadService) UploadThumbnail(ctx context.Context, postID string, suffix string, file FileInput) (string, error) {
	if err := requireNonEmpty("postId", postID); err != nil {
		return "", err
	}
	ext := strings.ToLower(filepath.Ext(file.Header.Filename))
	if !allowedPhotoExts[ext] {
		return "", validationErr(fmt.Sprintf("thumbnailFile must be one of jpg, jpeg, png, webp — got %q", ext))
	}
	sExt := storageExt(ext)
	fileID := postID + "-" + suffix + sExt
	if err := s.files.Put(ctx, "photos/"+fileID, file.File, contentTypeForExt(sExt)); err != nil {
		return "", fmt.Errorf("UploadThumbnail store: %w", err)
	}
	return "/files/" + fileID, nil
}
