package api

import (
	"errors"
	"net/http"

	"github.com/pratikluitel/antipratik/common/logging"
	"github.com/pratikluitel/antipratik/components/files"
	"github.com/pratikluitel/antipratik/components/files/store"
)

// fileServingHandler serves uploaded files and thumbnails from storage.
type fileServingHandler struct {
	fileStore files.FileStore
	log       logging.Logger
}

// NewFileServingHandler returns a new fileServingHandler.
func NewFileServingHandler(fs files.FileStore, log logging.Logger) files.FilesAPI {
	return &fileServingHandler{fileStore: fs, log: log}
}

// ServeFile handles GET /files/{fileId}.
// Forwards the Range request header to the store so R2 returns only the requested
// bytes — no full-file buffering for seeks on large video files.
func (h *fileServingHandler) ServeFile(w http.ResponseWriter, r *http.Request) {
	fileID := r.PathValue("fileId")
	if fileID == "" {
		writeError(w, http.StatusBadRequest, "fileId is required")
		return
	}

	rangeHeader := r.Header.Get("Range")

	for _, prefix := range []string{"photos/", "music/", "videos/"} {
		body, ct, contentRange, contentLength, err := h.fileStore.GetRange(r.Context(), prefix+fileID, rangeHeader)
		if err == nil {
			streamFileRange(w, body, ct, contentRange, contentLength)
			return
		}
		if !errors.Is(err, store.ErrFileNotFound) {
			h.log.Error("ServeFile error", "key", prefix+fileID, "err", err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}
	writeError(w, http.StatusNotFound, "file not found")
}

// ServeThumbnail handles GET /thumbnails/{thumbnailId}.
func (h *fileServingHandler) ServeThumbnail(w http.ResponseWriter, r *http.Request) {
	thumbnailID := r.PathValue("thumbnailId")
	if thumbnailID == "" {
		writeError(w, http.StatusBadRequest, "thumbnailId is required")
		return
	}

	body, ct, err := h.fileStore.Get(r.Context(), "thumbnails/"+thumbnailID)
	if err != nil {
		if errors.Is(err, store.ErrFileNotFound) {
			writeError(w, http.StatusNotFound, "thumbnail not found")
			return
		}
		h.log.Error("ServeThumbnail error", "key", "thumbnails/"+thumbnailID, "err", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	streamFile(w, r, body, ct)
}
