// Package store contains the files component data layer (storage backends).
package store

import (
	"context"
	"io"
)

// FileStore is the interface for storing and retrieving uploaded files.
// All keys are slash-separated paths, e.g. "photos/abc.jpg" or "thumbnails/abc-small.jpg".
type FileStore interface {
	// Put stores the content from r under key with the given MIME content type.
	Put(ctx context.Context, key string, r io.Reader, contentType string) error
	// Get retrieves the content stored at key.
	// Returns a seekable body (caller must close), the content type, and any error.
	// The returned body implements io.ReadSeekCloser so callers can serve Range requests.
	Get(ctx context.Context, key string) (io.ReadSeekCloser, string, error)
	// Delete removes the file stored at key. It is not an error if the key does not exist.
	Delete(ctx context.Context, key string) error
}
