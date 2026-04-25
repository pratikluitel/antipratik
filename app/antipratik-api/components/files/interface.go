package files

import (
	"context"
	"io"
	"net/http"
)

// UploadLogic handles file storage and thumbnail generation for uploaded assets.
type UploadLogic interface {
	// UploadPhotoFiles stores multiple photo images (with 4 thumbnail variants each) for the
	// given post. File IDs follow <postID>-<uuid>.<ext>; thumbnail IDs follow
	// <postID>-<uuid>-<size>.<ext> where size is tiny/small/medium/large.
	UploadPhotoFiles(ctx context.Context, postID string, files []FileInput) ([]PhotoImageResult, error)

	// UploadMusicFiles stores audio and/or album art for the given post.
	// Either file may be nil (to skip that upload), but at least one must be non-nil.
	// Audio is stored at music/<postID>.<ext>; album art at photos/<postID>-albumart.<ext>.
	UploadMusicFiles(ctx context.Context, postID string, audioFile *FileInput, albumArtFile *FileInput) (MusicFilesResult, error)

	// UploadThumbnail stores a single thumbnail image plus a 20px-wide tiny variant.
	// suffix is appended to the postID in the stored file name (e.g. "thumb").
	UploadThumbnail(ctx context.Context, postID string, suffix string, file FileInput) (ThumbnailResult, error)

	// UploadVideoFile stores a video file for the given post.
	// Accepted extensions: mp4, webm, mov. Stored at videos/<postID>.<ext>.
	// Returns a relative URL of the form /files/videos/<postID>.<ext>.
	UploadVideoFile(ctx context.Context, postID string, file FileInput) (VideoFileResult, error)
}

// UploaderService exposes file upload capabilities to other components.
// Inject this interface rather than importing files/logic directly.
type UploaderService interface {
	UploadPhotoFiles(ctx context.Context, postID string, files []FileInput) ([]PhotoImageResult, error)
	UploadMusicFiles(ctx context.Context, postID string, audioFile *FileInput, albumArtFile *FileInput) (MusicFilesResult, error)
	UploadThumbnail(ctx context.Context, postID string, suffix string, file FileInput) (ThumbnailResult, error)
	UploadVideoFile(ctx context.Context, postID string, file FileInput) (VideoFileResult, error)
}

// StorageService exposes file retrieval and deletion to other components.
// Inject this interface rather than importing files/store directly.
type StorageService interface {
	Delete(ctx context.Context, key string) error
}

// FileStore is the interface for storing and retrieving uploaded files.
// All keys are slash-separated paths, e.g. "photos/abc.jpg" or "thumbnails/abc-small.jpg".
type FileStore interface {
	// Put stores the content from r under key with the given MIME content type.
	Put(ctx context.Context, key string, r io.Reader, contentType string) error
	// Get retrieves the content stored at key.
	// Returns a seekable body (caller must close), the content type, and any error.
	// The returned body implements io.ReadSeekCloser so callers can serve Range requests.
	Get(ctx context.Context, key string) (io.ReadSeekCloser, string, error)
	// GetRange retrieves a byte range of the file stored at key.
	// rangeHeader is the raw RFC 7233 Range request header value (e.g. "bytes=0-1023");
	// pass "" to fetch the full object.
	// Returns: body (caller must close), content-type, Content-Range response header
	// ("bytes start-end/total"; empty when serving the full object), content-length of
	// the returned body, and any error.
	GetRange(ctx context.Context, key, rangeHeader string) (io.ReadCloser, string, string, int64, error)
	// Delete removes the file stored at key. It is not an error if the key does not exist.
	Delete(ctx context.Context, key string) error
}

type FilesAPI interface {
	ServeFile(w http.ResponseWriter, r *http.Request)
	ServeThumbnail(w http.ResponseWriter, r *http.Request)
}
