//go:build !cgo

package logic

import (
	"fmt"
	"image"
)

func decodeHEIF(_ []byte) (image.Image, error) {
	return nil, fmt.Errorf("HEIC/HEIF decoding requires CGO (build without -tags nocgo or enable CGO)")
}
