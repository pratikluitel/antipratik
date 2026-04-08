package logic

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"

	"github.com/rwcarlsen/goexif/exif"
	"golang.org/x/image/draw"
	"golang.org/x/image/webp"
)

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
