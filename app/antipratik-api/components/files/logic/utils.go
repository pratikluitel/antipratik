package logic

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"

	"golang.org/x/image/draw"
	"golang.org/x/image/webp"
)

func decodeImage(r multipart.File, ext string) (image.Image, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var img image.Image
	switch ext {
	case ".webp":
		img, err = webp.Decode(bytes.NewReader(data))
	case ".heic", ".heif":
		img, err = decodeHEIF(data)
	default:
		img, _, err = image.Decode(bytes.NewReader(data))
	}
	if err != nil {
		return nil, err
	}

	if ext == ".jpg" || ext == ".jpeg" {
		orientation := getEXIFOrientation(bytes.NewReader(data))
		if orientation > 1 {
			img = applyOrientation(img, orientation)
		}
	}

	return img, nil
}

func getEXIFOrientation(r io.Reader) uint32 {
	data, err := io.ReadAll(r)
	if err != nil {
		return 1
	}
	orientation, err := parseEXIFOrientation(data)
	if err != nil {
		return orientation
	}
	return orientation
}

func parseEXIFOrientation(data []byte) (uint32, error) {
	if len(data) < 4 {
		return 1, fmt.Errorf("too short")
	}

	var tiffData []byte

	switch {
	case data[0] == 0xFF && data[1] == 0xD8:
		tiffData = extractJPEGExifBlock(data)
		if tiffData == nil {
			return 1, fmt.Errorf("no EXIF block in JPEG")
		}
	case (data[0] == 0x49 && data[1] == 0x49) || (data[0] == 0x4D && data[1] == 0x4D):
		tiffData = data
	default:
		return 1, fmt.Errorf("not a JPEG or TIFF")
	}

	return readTIFFOrientation(tiffData)
}

func extractJPEGExifBlock(data []byte) []byte {
	i := 2
	for i+4 <= len(data) {
		if data[i] != 0xFF {
			return nil
		}
		marker := data[i+1]
		if i+4 > len(data) {
			return nil
		}
		segLen := int(data[i+2])<<8 | int(data[i+3])
		if segLen < 2 || i+2+segLen > len(data) {
			return nil
		}
		if marker == 0xE1 {
			payload := data[i+4 : i+2+segLen]
			if len(payload) > 6 && string(payload[:6]) == "Exif\x00\x00" {
				return payload[6:]
			}
		}
		if marker == 0xDA {
			return nil
		}
		i += 2 + segLen
	}
	return nil
}

func readTIFFOrientation(data []byte) (uint32, error) {
	if len(data) < 8 {
		return 1, fmt.Errorf("TIFF data too short")
	}

	var bo binary.ByteOrder
	switch {
	case data[0] == 0x49 && data[1] == 0x49:
		bo = binary.LittleEndian
	case data[0] == 0x4D && data[1] == 0x4D:
		bo = binary.BigEndian
	default:
		return 1, fmt.Errorf("invalid TIFF byte order marker")
	}

	if bo.Uint16(data[2:4]) != 42 {
		return 1, fmt.Errorf("invalid TIFF magic number")
	}

	ifdOffset := bo.Uint32(data[4:8])
	if uint64(ifdOffset)+2 > uint64(len(data)) {
		return 1, fmt.Errorf("IFD offset out of bounds")
	}

	numEntries := bo.Uint16(data[ifdOffset : ifdOffset+2])
	const entrySize = 12
	const orientationTag = 0x0112

	for i := 0; i < int(numEntries); i++ {
		entryOffset := int(ifdOffset) + 2 + i*entrySize
		if entryOffset+entrySize > len(data) {
			break
		}
		entry := data[entryOffset : entryOffset+entrySize]

		tagID := bo.Uint16(entry[0:2])
		if tagID != orientationTag {
			continue
		}

		dataType := bo.Uint16(entry[2:4])
		count := bo.Uint32(entry[4:8])
		if dataType != 3 || count != 1 {
			return 1, fmt.Errorf("unexpected orientation tag format")
		}

		val := bo.Uint16(entry[8:10])
		if val < 1 || val > 8 {
			return 1, fmt.Errorf("orientation value out of EXIF spec range")
		}
		return uint32(val), nil
	}

	return 1, fmt.Errorf("orientation tag not found")
}

func applyOrientation(img image.Image, orientation uint32) image.Image {
	switch orientation {
	case 2:
		return flipH(img)
	case 3:
		return rotate180(img)
	case 4:
		return flipV(img)
	case 5:
		return flipH(rotate90CCW(img))
	case 6:
		return rotate90CW(img)
	case 7:
		return flipH(rotate90CW(img))
	case 8:
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
			dst.Set(height-1-(y-bounds.Min.Y), x-bounds.Min.X, img.At(x, y))
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
			dst.Set(y-bounds.Min.Y, width-1-(x-bounds.Min.X), img.At(x, y))
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
			dst.Set(bounds.Min.X+width-1-(x-bounds.Min.X), bounds.Min.Y+height-1-(y-bounds.Min.Y), img.At(x, y))
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
			dst.Set(bounds.Min.X+width-1-(x-bounds.Min.X), y, img.At(x, y))
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
			dst.Set(x, bounds.Min.Y+height-1-(y-bounds.Min.Y), img.At(x, y))
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
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
	}
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

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
	case ".heic":
		return "image/heic"
	case ".heif":
		return "image/heif"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".ogg":
		return "audio/ogg"
	case ".m4a":
		return "audio/mp4"
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	default:
		return "application/octet-stream"
	}
}

// storageExt returns the file extension used when storing the encoded image.
// WebP inputs are encoded as JPEG, so they are stored with a .jpg extension.
func storageExt(ext string) string {
	if ext == ".webp" || ext == ".heic" || ext == ".heif" {
		return ".jpg"
	}
	return ext
}
