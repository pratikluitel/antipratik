// Package logic — upload service: validates, stores files, and generates thumbnails.
package logic

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pratikluitel/antipratik/models"
	"github.com/pratikluitel/antipratik/store"
)

// maxConcurrentUploads caps the number of photo images processed in parallel
// to avoid overwhelming the file store or exhausting memory on the server.
const maxConcurrentUploads = 4

// Thumbnail width constants — change these to adjust all generated variants.
const (
	thumbWidthTiny   = 20
	thumbWidthSmall  = 300
	thumbWidthMedium = 600
	thumbWidthLarge  = 1200
)

// Storage key prefixes and URL prefixes — never hardcode these inline.
const (
	storePrefixPhotos     = "photos/"
	storePrefixMusic      = "music/"
	storePrefixThumbnails = "thumbnails/"
	urlPrefixFiles        = "/files/"
	urlPrefixThumbnails   = "/thumbnails/"
)

// PhotoImageResult holds the URL fields for a single uploaded photo image.
type PhotoImageResult struct {
	OriginalURL       string
	ThumbnailTinyURL  string
	ThumbnailSmallURL string
	ThumbnailMedURL   string
	ThumbnailLargeURL string
}

// ThumbnailResult holds the URL and tiny-placeholder URL for an uploaded thumbnail.
type ThumbnailResult struct {
	URL     string
	TinyURL string
}

// MusicFilesResult holds the URL fields from a music post file upload.
type MusicFilesResult struct {
	AudioURL        string
	AlbumArtURL     string // empty string if no album art was uploaded
	AlbumArtTinyURL string // empty string if no album art was uploaded
}

// UploadLogic handles file storage and thumbnail generation for uploaded assets.
type UploadLogic interface {
	// UploadPhotoFiles stores multiple photo images (with 3 thumbnail variants each) for the
	// given post. File IDs follow <postID>-<index>.<ext>; thumbnail IDs follow
	// <postID>-<index>-<size>.<ext> where size is small/medium/large.
	UploadPhotoFiles(ctx context.Context, postID string, files []models.FileInput) ([]PhotoImageResult, error)

	// UploadMusicFiles stores audio and/or album art for the given post.
	// Either file may be nil (to skip that upload), but at least one must be non-nil.
	// Audio is stored at music/<postID>.<ext>; album art at photos/<postID>-albumart.<ext>.
	UploadMusicFiles(ctx context.Context, postID string, audioFile *models.FileInput, albumArtFile *models.FileInput) (MusicFilesResult, error)

	// UploadThumbnail stores a single thumbnail image plus a 20px-wide tiny variant.
	// suffix is appended to the postID in the stored file name (e.g. "thumb").
	UploadThumbnail(ctx context.Context, postID string, suffix string, file models.FileInput) (ThumbnailResult, error)
}

// UploadService implements UploadLogic.
type UploadService struct {
	files store.FileStore
}

// NewUploadService constructs an UploadService.
func NewUploadService(files store.FileStore) *UploadService {
	return &UploadService{files: files}
}

var allowedPhotoExts = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".heic": true, ".heif": true}
var allowedMusicExts = map[string]bool{".mp3": true, ".wav": true, ".ogg": true, ".m4a": true}

// storageExt returns the file extension used when storing the encoded image.
// WebP inputs are encoded as JPEG, so they are stored with a .jpg extension.
func storageExt(ext string) string {
	if ext == ".webp" || ext == ".heic" || ext == ".heif" {
		return ".jpg"
	}
	return ext
}

// photoImageWork is the result (or error) for a single image in a concurrent upload batch.
type photoImageWork struct {
	index  int
	result PhotoImageResult
	err    error
}

// UploadPhotoFiles implements UploadLogic.
// Images are processed concurrently up to maxConcurrentUploads at a time so
// that large photo batches don't serialize unnecessarily, while avoiding
// unbounded goroutine or memory growth.
func (s *UploadService) UploadPhotoFiles(ctx context.Context, postID string, files []models.FileInput) ([]PhotoImageResult, error) {
	if err := requireNonEmpty("postId", postID); err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, validationErr("at least one image file is required")
	}

	results := make([]PhotoImageResult, len(files))
	sem := make(chan struct{}, maxConcurrentUploads)
	work := make(chan photoImageWork, len(files))

	var wg sync.WaitGroup
	for i, fi := range files {
		wg.Add(1)
		go func(i int, fi models.FileInput) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			result, err := s.uploadOnePhoto(ctx, postID, i, fi)
			work <- photoImageWork{index: i, result: result, err: err}
		}(i, fi)
	}

	go func() {
		wg.Wait()
		close(work)
	}()

	for w := range work {
		if w.err != nil {
			return nil, w.err
		}
		results[w.index] = w.result
	}
	return results, nil
}

var thumbnailSizes = []struct {
	name     string
	maxWidth int
}{
	{"tiny", thumbWidthTiny},
	{"small", thumbWidthSmall},
	{"medium", thumbWidthMedium},
	{"large", thumbWidthLarge},
}

// uploadOnePhoto encodes, resizes, and stores a single photo image and its thumbnails.
func (s *UploadService) uploadOnePhoto(ctx context.Context, postID string, i int, fi models.FileInput) (PhotoImageResult, error) {
	ext := strings.ToLower(filepath.Ext(fi.Header.Filename))
	if !allowedPhotoExts[ext] {
		return PhotoImageResult{}, validationErr(fmt.Sprintf("images[%d]: file must be one of jpg, jpeg, png, webp — got %q", i, ext))
	}

	src, err := decodeImage(fi.File, ext)
	if err != nil {
		return PhotoImageResult{}, validationErr(fmt.Sprintf("images[%d]: could not decode image: %v", i, err))
	}

	sExt := storageExt(ext)
	ct := contentTypeForExt(sExt)
	fileID := fmt.Sprintf("%s-%d%s", postID, i, sExt)

	origBuf, err := encodeImage(src, ext)
	if err != nil {
		return PhotoImageResult{}, fmt.Errorf("UploadPhotoFiles encode original[%d]: %w", i, err)
	}
	if err := s.files.Put(ctx, storePrefixPhotos+fileID, bytes.NewReader(origBuf), ct); err != nil {
		return PhotoImageResult{}, fmt.Errorf("UploadPhotoFiles store original[%d]: %w", i, err)
	}

	thumbURLs := make([]string, len(thumbnailSizes))
	for j, sz := range thumbnailSizes {
		thumb := resizeImage(src, sz.maxWidth)
		buf, err := encodeImage(thumb, ext)
		if err != nil {
			return PhotoImageResult{}, fmt.Errorf("UploadPhotoFiles encode thumbnail[%d][%s]: %w", i, sz.name, err)
		}
		thumbID := fmt.Sprintf("%s-%d-%s%s", postID, i, sz.name, sExt)
		if err := s.files.Put(ctx, storePrefixThumbnails+thumbID, bytes.NewReader(buf), ct); err != nil {
			return PhotoImageResult{}, fmt.Errorf("UploadPhotoFiles store thumbnail[%d][%s]: %w", i, sz.name, err)
		}
		thumbURLs[j] = urlPrefixThumbnails + thumbID
	}

	return PhotoImageResult{
		OriginalURL:       urlPrefixFiles + fileID,
		ThumbnailTinyURL:  thumbURLs[0],
		ThumbnailSmallURL: thumbURLs[1],
		ThumbnailMedURL:   thumbURLs[2],
		ThumbnailLargeURL: thumbURLs[3],
	}, nil
}

// storeImageWithTiny decodes, encodes, and stores an image at photos/<fileID>
// and a 20px-wide tiny variant at thumbnails/<fileID>-tiny<ext>.
// It returns the /files/ and /thumbnails/ URLs for both.
func (s *UploadService) storeImageWithTiny(ctx context.Context, fileID string, file models.FileInput, ext string) (photoURL, tinyURL string, err error) {
	sExt := storageExt(ext)
	ct := contentTypeForExt(sExt)
	storeID := fileID + sExt

	src, err := decodeImage(file.File, ext)
	if err != nil {
		return "", "", validationErr(fmt.Sprintf("%s: could not decode image: %v", fileID, err))
	}

	buf, err := encodeImage(src, ext)
	if err != nil {
		return "", "", fmt.Errorf("storeImageWithTiny encode %s: %w", fileID, err)
	}
	if err := s.files.Put(ctx, storePrefixPhotos+storeID, bytes.NewReader(buf), ct); err != nil {
		return "", "", fmt.Errorf("storeImageWithTiny store %s: %w", fileID, err)
	}

	tiny := resizeImage(src, thumbWidthTiny)
	tinyBuf, err := encodeImage(tiny, ext)
	if err != nil {
		return "", "", fmt.Errorf("storeImageWithTiny encode tiny %s: %w", fileID, err)
	}
	tinyStoreID := fileID + "-tiny" + sExt
	if err := s.files.Put(ctx, storePrefixThumbnails+tinyStoreID, bytes.NewReader(tinyBuf), ct); err != nil {
		return "", "", fmt.Errorf("storeImageWithTiny store tiny %s: %w", fileID, err)
	}

	return urlPrefixFiles + storeID, urlPrefixThumbnails + tinyStoreID, nil
}

// UploadMusicFiles implements UploadLogic.
func (s *UploadService) UploadMusicFiles(ctx context.Context, postID string, audioFile *models.FileInput, albumArtFile *models.FileInput) (MusicFilesResult, error) {
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
		if err := s.files.Put(ctx, storePrefixMusic+audioFileID, audioFile.File, contentTypeForExt(audioExt)); err != nil {
			return MusicFilesResult{}, fmt.Errorf("UploadMusicFiles store audio: %w", err)
		}
		result.AudioURL = urlPrefixFiles + audioFileID
	}

	// Upload album art (optional)
	if albumArtFile != nil {
		artExt := strings.ToLower(filepath.Ext(albumArtFile.Header.Filename))
		if !allowedPhotoExts[artExt] {
			return MusicFilesResult{}, validationErr(fmt.Sprintf("albumArtFile must be one of jpg, jpeg, png, webp, heic, heif — got %q", artExt))
		}
		artURL, artTinyURL, err := s.storeImageWithTiny(ctx, postID+"-albumart", *albumArtFile, artExt)
		if err != nil {
			return MusicFilesResult{}, err
		}
		result.AlbumArtURL = artURL
		result.AlbumArtTinyURL = artTinyURL
	}

	return result, nil
}

// UploadThumbnail implements UploadLogic.
func (s *UploadService) UploadThumbnail(ctx context.Context, postID string, suffix string, file models.FileInput) (ThumbnailResult, error) {
	if err := requireNonEmpty("postId", postID); err != nil {
		return ThumbnailResult{}, err
	}
	ext := strings.ToLower(filepath.Ext(file.Header.Filename))
	if !allowedPhotoExts[ext] {
		return ThumbnailResult{}, validationErr(fmt.Sprintf("thumbnailFile must be one of jpg, jpeg, png, webp — got %q", ext))
	}
	photoURL, tinyURL, err := s.storeImageWithTiny(ctx, postID+"-"+suffix, file, ext)
	if err != nil {
		return ThumbnailResult{}, err
	}
	return ThumbnailResult{URL: photoURL, TinyURL: tinyURL}, nil
}
