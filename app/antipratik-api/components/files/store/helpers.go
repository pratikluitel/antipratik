package store

import (
	"fmt"
	"path/filepath"

	"github.com/pratikluitel/antipratik/components/files"
)

// resolveRange converts a ParsedRange to absolute start/end byte positions against a known total size.
func resolveRange(r *files.ParsedRange, total int64) (start, end int64, ok bool) {
	if r.Start == nil {
		// Suffix range: bytes=-N
		n := *r.End
		if n <= 0 {
			return 0, 0, false
		}
		s := total - n
		if s < 0 {
			s = 0
		}
		return s, total - 1, true
	}
	s := *r.Start
	if s < 0 || s >= total {
		return 0, 0, false
	}
	if r.End == nil {
		return s, total - 1, true
	}
	e := *r.End
	if e < s {
		return 0, 0, false
	}
	if e >= total {
		e = total - 1
	}
	return s, e, true
}

// contentTypeFromKey returns a MIME type based on the file extension in key.
func contentTypeFromKey(key string) string {
	switch filepath.Ext(key) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
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

// parsedRangeToS3Header converts a ParsedRange back to an RFC 7233 Range header string for the S3 API.
func parsedRangeToS3Header(r *files.ParsedRange) string {
	if r.Start == nil {
		return fmt.Sprintf("bytes=-%d", *r.End)
	}
	if r.End == nil {
		return fmt.Sprintf("bytes=%d-", *r.Start)
	}
	return fmt.Sprintf("bytes=%d-%d", *r.Start, *r.End)
}
