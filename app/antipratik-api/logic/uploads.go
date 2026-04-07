// Package logic — upload service: validates, stores files, and generates thumbnails.
package logic

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/pratikluitel/antipratik/models"
	"github.com/pratikluitel/antipratik/store"
	"golang.org/x/image/draw"
	"golang.org/x/image/webp"
)

// UploadLogic defines the upload use-cases.
type UploadLogic interface {
	UploadPhoto(ctx context.Context, postID string, file multipart.File, header *multipart.FileHeader) (models.UploadPhotoResponse, error)
	UploadMusic(ctx context.Context, postID string, file multipart.File, header *multipart.FileHeader) (models.UploadMusicResponse, error)
}

// UploadService implements UploadLogic.
type UploadService struct {
	files   store.FileStore
	baseURL string // e.g. "https://api.example.com" — prepended to /files/ and /thumbnails/ paths
}

// NewUploadService constructs an UploadService.
// baseURL is the server base URL used to construct public file URLs.
func NewUploadService(files store.FileStore, baseURL string) *UploadService {
	return &UploadService{files: files, baseURL: strings.TrimRight(baseURL, "/")}
}

var allowedPhotoExts = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
var allowedMusicExts = map[string]bool{".mp3": true, ".wav": true, ".ogg": true, ".m4a": true}

// UploadPhoto validates the photo, generates 3 thumbnails, stores all files, and returns URLs.
func (s *UploadService) UploadPhoto(ctx context.Context, postID string, file multipart.File, header *multipart.FileHeader) (models.UploadPhotoResponse, error) {
	if err := requireNonEmpty("postId", postID); err != nil {
		return models.UploadPhotoResponse{}, err
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedPhotoExts[ext] {
		return models.UploadPhotoResponse{}, validationErr(fmt.Sprintf("photo file must be one of: jpg, jpeg, png, webp — got %q", ext))
	}

	// Decode image.
	src, err := decodeImage(file, ext)
	if err != nil {
		return models.UploadPhotoResponse{}, validationErr(fmt.Sprintf("could not decode image: %v", err))
	}

	fileID := postID + ext
	originalKey := "photos/" + fileID
	ct := contentTypeForExt(ext)

	// Re-encode original (from decoded image so we read from src, not the already-read file).
	origBuf, err := encodeImage(src, ext)
	if err != nil {
		return models.UploadPhotoResponse{}, fmt.Errorf("UploadPhoto encode original: %w", err)
	}
	if err := s.files.Put(ctx, originalKey, bytes.NewReader(origBuf), ct); err != nil {
		return models.UploadPhotoResponse{}, fmt.Errorf("UploadPhoto store original: %w", err)
	}

	// Generate and store thumbnails.
	sizes := []struct {
		name     string
		maxWidth int
	}{
		{"small", 300},
		{"medium", 600},
		{"large", 1200},
	}
	thumbnailIDs := make([]string, len(sizes))
	for i, sz := range sizes {
		thumb := resizeImage(src, sz.maxWidth)
		buf, err := encodeImage(thumb, ext)
		if err != nil {
			return models.UploadPhotoResponse{}, fmt.Errorf("UploadPhoto encode thumbnail %s: %w", sz.name, err)
		}
		thumbID := postID + "-" + sz.name + ext
		if err := s.files.Put(ctx, "thumbnails/"+thumbID, bytes.NewReader(buf), ct); err != nil {
			return models.UploadPhotoResponse{}, fmt.Errorf("UploadPhoto store thumbnail %s: %w", sz.name, err)
		}
		thumbnailIDs[i] = thumbID
	}

	return models.UploadPhotoResponse{
		FileID:             fileID,
		OriginalURL:        s.baseURL + "/files/" + fileID,
		ThumbnailSmallURL:  s.baseURL + "/thumbnails/" + thumbnailIDs[0],
		ThumbnailMediumURL: s.baseURL + "/thumbnails/" + thumbnailIDs[1],
		ThumbnailLargeURL:  s.baseURL + "/thumbnails/" + thumbnailIDs[2],
	}, nil
}

// UploadMusic validates the audio file, stores it, and returns the audio URL.
func (s *UploadService) UploadMusic(ctx context.Context, postID string, file multipart.File, header *multipart.FileHeader) (models.UploadMusicResponse, error) {
	if err := requireNonEmpty("postId", postID); err != nil {
		return models.UploadMusicResponse{}, err
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedMusicExts[ext] {
		return models.UploadMusicResponse{}, validationErr(fmt.Sprintf("music file must be one of: mp3, wav, ogg, m4a — got %q", ext))
	}

	fileID := postID + ext
	if err := s.files.Put(ctx, "music/"+fileID, file, contentTypeForExt(ext)); err != nil {
		return models.UploadMusicResponse{}, fmt.Errorf("UploadMusic store: %w", err)
	}

	return models.UploadMusicResponse{
		FileID:   fileID,
		AudioURL: s.baseURL + "/files/" + fileID,
	}, nil
}

// ── Image helpers ─────────────────────────────────────────────────────────────

func decodeImage(r multipart.File, ext string) (image.Image, error) {
	switch ext {
	case ".webp":
		return webp.Decode(r)
	default:
		img, _, err := image.Decode(r)
		return img, err
	}
}

func encodeImage(img image.Image, ext string) ([]byte, error) {
	var buf bytes.Buffer
	var err error
	switch ext {
	case ".png":
		err = png.Encode(&buf, img)
	default:
		// jpg, jpeg, webp → store as JPEG for thumbnails
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
	}
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// resizeImage scales src so its width is at most maxWidth, preserving aspect ratio.
// If src is already narrower than maxWidth it is returned unchanged.
func resizeImage(src image.Image, maxWidth int) image.Image {
	bounds := src.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	if w <= maxWidth {
		return src
	}
	newH := h * maxWidth / w
	dst := image.NewRGBA(image.Rect(0, 0, maxWidth, newH))
	draw.BiLinear.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)
	return dst
}

func contentTypeForExt(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".ogg":
		return "audio/ogg"
	case ".m4a":
		return "audio/mp4"
	default:
		return "application/octet-stream"
	}
}
