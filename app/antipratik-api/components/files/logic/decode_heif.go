//go:build cgo

package logic

import (
	"fmt"
	"image"

	"github.com/strukturag/libheif/go/heif"
)

func decodeHEIF(data []byte) (image.Image, error) {
	hctx, hErr := heif.NewContext()
	if hErr != nil {
		return nil, fmt.Errorf("creating heif context: %w", hErr)
	}
	if hErr = hctx.ReadFromMemory(data); hErr != nil {
		return nil, fmt.Errorf("reading heif data: %w", hErr)
	}
	handle, hErr := hctx.GetPrimaryImageHandle()
	if hErr != nil {
		return nil, fmt.Errorf("getting heif primary image: %w", hErr)
	}
	decoded, hErr := handle.DecodeImage(heif.ColorspaceUndefined, heif.ChromaUndefined, nil)
	if hErr != nil {
		return nil, fmt.Errorf("decoding heif image: %w", hErr)
	}
	return decoded.GetImage()
}
