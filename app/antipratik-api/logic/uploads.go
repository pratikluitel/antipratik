// Package logic — upload service: validates, stores files, and generates thumbnails.
package logic

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/pratikluitel/antipratik/store"
	"github.com/rwcarlsen/goexif/exif"
	"golang.org/x/image/draw"
	"golang.org/x/image/webp"
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

	// UploadMusicFiles stores an audio file and optional album art for the given post.
	// audioFile is required; albumArtFile may be nil when no album art is provided.
	UploadMusicFiles(ctx context.Context, postID string, audioFile FileInput, albumArtFile *FileInput) (MusicFilesResult, error)

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

		ct := contentTypeForExt(ext)
		fileID := fmt.Sprintf("%s-%d%s", postID, i, ext)
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
			thumbID := fmt.Sprintf("%s-%d-%s%s", postID, i, sz.name, ext)
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
func (s *UploadService) UploadMusicFiles(ctx context.Context, postID string, audioFile FileInput, albumArtFile *FileInput) (MusicFilesResult, error) {
	if err := requireNonEmpty("postId", postID); err != nil {
		return MusicFilesResult{}, err
	}

	// Upload audio.
	audioExt := strings.ToLower(filepath.Ext(audioFile.Header.Filename))
	if !allowedMusicExts[audioExt] {
		return MusicFilesResult{}, validationErr(fmt.Sprintf("audioFile must be one of mp3, wav, ogg, m4a — got %q", audioExt))
	}
	audioFileID := postID + audioExt
	if err := s.files.Put(ctx, "music/"+audioFileID, audioFile.File, contentTypeForExt(audioExt)); err != nil {
		return MusicFilesResult{}, fmt.Errorf("UploadMusicFiles store audio: %w", err)
	}

	result := MusicFilesResult{
		AudioURL: "/files/" + audioFileID,
	}

	// Upload album art (optional).
	if albumArtFile != nil {
		artExt := strings.ToLower(filepath.Ext(albumArtFile.Header.Filename))
		if !allowedPhotoExts[artExt] {
			return MusicFilesResult{}, validationErr(fmt.Sprintf("albumArtFile must be one of jpg, jpeg, png, webp — got %q", artExt))
		}
		artFileID := postID + "-albumart" + artExt
		if err := s.files.Put(ctx, "photos/"+artFileID, albumArtFile.File, contentTypeForExt(artExt)); err != nil {
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
	fileID := postID + "-" + suffix + ext
	if err := s.files.Put(ctx, "photos/"+fileID, file.File, contentTypeForExt(ext)); err != nil {
		return "", fmt.Errorf("UploadThumbnail store: %w", err)
	}
	return "/files/" + fileID, nil
}

// ── Image helpers ─────────────────────────────────────────────────────────────

func decodeImage(r multipart.File, ext string) (image.Image, error) {
	// Read the file into a buffer so we can parse EXIF data
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var img image.Image
	switch ext {
	case ".webp":
		img, err = webp.Decode(bytes.NewReader(data))
	default:
		img, _, err = image.Decode(bytes.NewReader(data))
	}
	if err != nil {
		return nil, err
	}

	// Try to apply EXIF orientation
	if ext == ".jpg" || ext == ".jpeg" {
		orientation := getEXIFOrientation(bytes.NewReader(data))
		if orientation > 1 {
			img = applyOrientation(img, orientation)
		}
	}

	return img, nil
}

// getEXIFOrientation extracts the orientation tag from EXIF metadata.
// Returns 1 (normal) if no EXIF orientation is found or if there's an error.
func getEXIFOrientation(r io.Reader) uint32 {
	exifData, err := exif.Decode(r)
	if err != nil {
		return 1 // Default to normal orientation
	}
	orientation, err := exifData.Get(exif.Orientation)
	if err != nil {
		return 1 // Default to normal orientation
	}
	val, err := orientation.Int(0)
	if err != nil {
		return 1
	}
	return uint32(val)
}

// applyOrientation transforms an image based on EXIF orientation tag.
// Orientation values:
// 1: Normal
// 2: Flip horizontal
// 3: Rotate 180°
// 4: Flip vertical
// 5: Rotate 90° CCW + flip horizontal
// 6: Rotate 90° CW
// 7: Rotate 90° CW + flip horizontal
// 8: Rotate 90° CCW
func applyOrientation(img image.Image, orientation uint32) image.Image {
	switch orientation {
	case 2:
		// Flip horizontal
		return flipH(img)
	case 3:
		// Rotate 180
		return rotate180(img)
	case 4:
		// Flip vertical
		return flipV(img)
	case 5:
		// Rotate 90 CCW + flip
		return flipH(rotate90CCW(img))
	case 6:
		// Rotate 90 CW
		return rotate90CW(img)
	case 7:
		// Rotate 90 CW + flip
		return flipH(rotate90CW(img))
	case 8:
		// Rotate 90 CCW
		return rotate90CCW(img)
	default:
		return img
	}
}

func rotate90CW(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, height, width))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Map (x, y) to (height - 1 - (y - min.Y), x - min.X)
			dstX := height - 1 - (y - bounds.Min.Y)
			dstY := x - bounds.Min.X
			dst.Set(dstX, dstY, img.At(x, y))
		}
	}
	return dst
}

func rotate90CCW(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, height, width))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Map (x, y) to (y - min.Y, width - 1 - (x - min.X))
			dstX := y - bounds.Min.Y
			dstY := width - 1 - (x - bounds.Min.X)
			dst.Set(dstX, dstY, img.At(x, y))
		}
	}
	return dst
}

func rotate180(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Map (x, y) to (width - 1 - (x - min.X), height - 1 - (y - min.Y))
			dstX := bounds.Min.X + width - 1 - (x - bounds.Min.X)
			dstY := bounds.Min.Y + height - 1 - (y - bounds.Min.Y)
			dst.Set(dstX, dstY, img.At(x, y))
		}
	}
	return dst
}

func flipH(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dstX := bounds.Min.X + width - 1 - (x - bounds.Min.X)
			dst.Set(dstX, y, img.At(x, y))
		}
	}
	return dst
}

func flipV(img image.Image) image.Image {
	bounds := img.Bounds()
	height := bounds.Dy()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dstY := bounds.Min.Y + height - 1 - (y - bounds.Min.Y)
			dst.Set(x, dstY, img.At(x, y))
		}
	}
	return dst
}

func encodeImage(img image.Image, ext string) ([]byte, error) {
	var buf bytes.Buffer
	var err error
	switch ext {
	case ".png":
		err = png.Encode(&buf, img)
	default:
		// jpg, jpeg, webp → encode as JPEG
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
	}
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// resizeImage scales src so its width is at most maxWidth, preserving aspect ratio.
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
