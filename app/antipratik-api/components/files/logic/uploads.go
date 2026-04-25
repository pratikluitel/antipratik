// Package logic — upload service: validates, stores files, and generates thumbnails.
package logic

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"

	commonerrors "github.com/pratikluitel/antipratik/common/errors"
	"github.com/pratikluitel/antipratik/components/files"
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

var thumbnailSizes = []struct {
	name     string
	maxWidth int
}{
	{"tiny", thumbWidthTiny},
	{"small", thumbWidthSmall},
	{"medium", thumbWidthMedium},
	{"large", thumbWidthLarge},
}

// Storage key prefixes and URL prefixes — never hardcode these inline.
const (
	storePrefixPhotos     = "photos/"
	storePrefixMusic      = "music/"
	storePrefixVideos     = "videos/"
	storePrefixThumbnails = "thumbnails/"
	urlPrefixFiles        = "/files/"
	urlPrefixThumbnails   = "/thumbnails/"
)

var allowedPhotoExts = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".heic": true, ".heif": true}
var allowedMusicExts = map[string]bool{".mp3": true, ".wav": true, ".ogg": true, ".m4a": true}
var allowedVideoExts = map[string]bool{".mp4": true, ".webm": true, ".mov": true}

// uploadLogic implements UploadLogic.
type uploadLogic struct {
	files files.FileStore
}

// NewUploadLogic constructs an uploadLogic.
func NewUploadLogic(files files.FileStore) files.UploadLogic {
	return &uploadLogic{files: files}
}

// photoImageWork is the result (or error) for a single image in a concurrent upload batch.
type photoImageWork struct {
	err    error
	result files.PhotoImageResult
	index  int
}

// UploadPhotoFiles implements UploadLogic.
// Images are processed concurrently up to maxConcurrentUploads at a time so
// that large photo batches don't serialize unnecessarily, while avoiding
// unbounded goroutine or memory growth.
func (s *uploadLogic) UploadPhotoFiles(ctx context.Context, postID string, fls []files.FileInput) ([]files.PhotoImageResult, error) {
	if err := commonerrors.RequireNonEmpty("postId", postID); err != nil {
		return nil, err
	}
	if len(fls) == 0 {
		return nil, commonerrors.New("at least one image file is required")
	}

	results := make([]files.PhotoImageResult, len(fls))
	sem := make(chan struct{}, maxConcurrentUploads)
	work := make(chan photoImageWork, len(fls))

	var wg sync.WaitGroup
	for i, fi := range fls {
		wg.Add(1)
		go func(i int, fi files.FileInput) {
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

// uploadOnePhoto encodes, resizes, and stores a single photo image and its thumbnails.
// File IDs use a UUID to guarantee uniqueness: <postID>-<uuid>.<ext>.
func (s *uploadLogic) uploadOnePhoto(ctx context.Context, postID string, i int, fi files.FileInput) (files.PhotoImageResult, error) {
	ext := strings.ToLower(filepath.Ext(fi.Header.Filename))
	if !allowedPhotoExts[ext] {
		return files.PhotoImageResult{}, commonerrors.New(fmt.Sprintf("images[%d]: file must be one of jpg, jpeg, png, webp — got %q", i, ext))
	}

	src, err := decodeImage(fi.File, ext)
	if err != nil {
		return files.PhotoImageResult{}, commonerrors.New(fmt.Sprintf("images[%d]: could not decode image: %v", i, err))
	}

	sExt := storageExt(ext)
	ct := contentTypeForExt(sExt)
	imgUUID := uuid.New().String()
	fileID := fmt.Sprintf("%s-%s%s", postID, imgUUID, sExt)

	origBuf, err := encodeImage(src, ext)
	if err != nil {
		return files.PhotoImageResult{}, fmt.Errorf("UploadPhotoFiles encode original[%d]: %w", i, err)
	}
	if err := s.files.Put(ctx, storePrefixPhotos+fileID, bytes.NewReader(origBuf), ct); err != nil {
		return files.PhotoImageResult{}, fmt.Errorf("UploadPhotoFiles store original[%d]: %w", i, err)
	}

	thumbURLs := make([]string, len(thumbnailSizes))
	for j, sz := range thumbnailSizes {
		thumb := resizeImage(src, sz.maxWidth)
		buf, err := encodeImage(thumb, ext)
		if err != nil {
			return files.PhotoImageResult{}, fmt.Errorf("UploadPhotoFiles encode thumbnail[%d][%s]: %w", i, sz.name, err)
		}
		thumbID := fmt.Sprintf("%s-%s-%s%s", postID, imgUUID, sz.name, sExt)
		if err := s.files.Put(ctx, storePrefixThumbnails+thumbID, bytes.NewReader(buf), ct); err != nil {
			return files.PhotoImageResult{}, fmt.Errorf("UploadPhotoFiles store thumbnail[%d][%s]: %w", i, sz.name, err)
		}
		thumbURLs[j] = urlPrefixThumbnails + thumbID
	}

	return files.PhotoImageResult{
		OriginalURL:       urlPrefixFiles + fileID,
		ThumbnailTinyURL:  thumbURLs[0],
		ThumbnailSmallURL: thumbURLs[1],
		ThumbnailMedURL:   thumbURLs[2],
		ThumbnailLargeURL: thumbURLs[3],
	}, nil
}

// allSizesResult holds the original URL and all 4 thumbnail URLs for a stored image.
type allSizesResult struct {
	OriginalURL string
	TinyURL     string
	SmallURL    string
	MedURL      string
	LargeURL    string
}

// storeImageAllSizes decodes, encodes, and stores an image at photos/<fileID> plus
// all 4 thumbnail variants (tiny/small/medium/large) at thumbnails/<fileID>-<size><ext>.
func (s *uploadLogic) storeImageAllSizes(ctx context.Context, fileID string, file files.FileInput, ext string) (allSizesResult, error) {
	sExt := storageExt(ext)
	ct := contentTypeForExt(sExt)
	storeID := fileID + sExt

	src, err := decodeImage(file.File, ext)
	if err != nil {
		return allSizesResult{}, commonerrors.New(fmt.Sprintf("%s: could not decode image: %v", fileID, err))
	}

	buf, err := encodeImage(src, ext)
	if err != nil {
		return allSizesResult{}, fmt.Errorf("storeImageAllSizes encode %s: %w", fileID, err)
	}
	if err = s.files.Put(ctx, storePrefixPhotos+storeID, bytes.NewReader(buf), ct); err != nil {
		return allSizesResult{}, fmt.Errorf("storeImageAllSizes store %s: %w", fileID, err)
	}

	result := allSizesResult{OriginalURL: urlPrefixFiles + storeID}
	for _, sz := range thumbnailSizes {
		thumb := resizeImage(src, sz.maxWidth)
		thumbBuf, err := encodeImage(thumb, ext)
		if err != nil {
			return allSizesResult{}, fmt.Errorf("storeImageAllSizes encode %s-%s: %w", fileID, sz.name, err)
		}
		thumbStoreID := fileID + "-" + sz.name + sExt
		if err := s.files.Put(ctx, storePrefixThumbnails+thumbStoreID, bytes.NewReader(thumbBuf), ct); err != nil {
			return allSizesResult{}, fmt.Errorf("storeImageAllSizes store %s-%s: %w", fileID, sz.name, err)
		}
		thumbURL := urlPrefixThumbnails + thumbStoreID
		switch sz.name {
		case "tiny":
			result.TinyURL = thumbURL
		case "small":
			result.SmallURL = thumbURL
		case "medium":
			result.MedURL = thumbURL
		case "large":
			result.LargeURL = thumbURL
		}
	}
	return result, nil
}

// UploadMusicFiles implements UploadLogic.
func (s *uploadLogic) UploadMusicFiles(ctx context.Context, postID string, audioFile *files.FileInput, albumArtFile *files.FileInput) (files.MusicFilesResult, error) {
	if err := commonerrors.RequireNonEmpty("postId", postID); err != nil {
		return files.MusicFilesResult{}, err
	}

	var result files.MusicFilesResult

	if audioFile != nil {
		audioExt := strings.ToLower(filepath.Ext(audioFile.Header.Filename))
		if !allowedMusicExts[audioExt] {
			return files.MusicFilesResult{}, commonerrors.New(fmt.Sprintf("audioFile must be one of mp3, wav, ogg, m4a — got %q", audioExt))
		}
		audioFileID := postID + audioExt
		if err := s.files.Put(ctx, storePrefixMusic+audioFileID, audioFile.File, contentTypeForExt(audioExt)); err != nil {
			return files.MusicFilesResult{}, fmt.Errorf("UploadMusicFiles store audio: %w", err)
		}
		result.AudioURL = urlPrefixFiles + audioFileID
	}

	if albumArtFile != nil {
		artExt := strings.ToLower(filepath.Ext(albumArtFile.Header.Filename))
		if !allowedPhotoExts[artExt] {
			return files.MusicFilesResult{}, commonerrors.New(fmt.Sprintf("albumArtFile must be one of jpg, jpeg, png, webp, heic, heif — got %q", artExt))
		}
		art, err := s.storeImageAllSizes(ctx, postID+"-albumart", *albumArtFile, artExt)
		if err != nil {
			return files.MusicFilesResult{}, err
		}
		result.AlbumArtURL = art.OriginalURL
		result.AlbumArtTinyURL = art.TinyURL
		result.AlbumArtSmallURL = art.SmallURL
		result.AlbumArtMedURL = art.MedURL
		result.AlbumArtLargeURL = art.LargeURL
	}

	return result, nil
}

// UploadThumbnail implements UploadLogic.
func (s *uploadLogic) UploadThumbnail(ctx context.Context, postID string, suffix string, file files.FileInput) (files.ThumbnailResult, error) {
	if err := commonerrors.RequireNonEmpty("postId", postID); err != nil {
		return files.ThumbnailResult{}, err
	}
	ext := strings.ToLower(filepath.Ext(file.Header.Filename))
	if !allowedPhotoExts[ext] {
		return files.ThumbnailResult{}, commonerrors.New(fmt.Sprintf("thumbnailFile must be one of jpg, jpeg, png, webp — got %q", ext))
	}
	r, err := s.storeImageAllSizes(ctx, postID+"-"+suffix, file, ext)
	if err != nil {
		return files.ThumbnailResult{}, err
	}
	return files.ThumbnailResult{URL: r.OriginalURL, TinyURL: r.TinyURL, SmallURL: r.SmallURL, MedURL: r.MedURL, LargeURL: r.LargeURL}, nil
}

// UploadVideoFile implements UploadLogic.
// Accepted extensions: mp4, webm, mov. Stored at videos/<postID>.<ext>.
func (s *uploadLogic) UploadVideoFile(ctx context.Context, postID string, file files.FileInput) (files.VideoFileResult, error) {
	if err := commonerrors.RequireNonEmpty("postId", postID); err != nil {
		return files.VideoFileResult{}, err
	}
	ext := strings.ToLower(filepath.Ext(file.Header.Filename))
	if !allowedVideoExts[ext] {
		return files.VideoFileResult{}, commonerrors.New(fmt.Sprintf("videoFile must be one of mp4, webm, mov — got %q", ext))
	}
	videoFileID := postID + ext
	if err := s.files.Put(ctx, storePrefixVideos+videoFileID, file.File, contentTypeForExt(ext)); err != nil {
		return files.VideoFileResult{}, fmt.Errorf("UploadVideoFile store: %w", err)
	}
	return files.VideoFileResult{VideoURL: urlPrefixFiles + videoFileID}, nil
}
