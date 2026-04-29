package api

import (
	"errors"
	"net/http"

	"github.com/pratikluitel/antipratik/common/logging"
	"github.com/pratikluitel/antipratik/common/requests"
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
// Parses the Range header in the API layer and passes the resolved range to the store.
func (h *fileServingHandler) ServeFile(w http.ResponseWriter, r *http.Request) {
	fileID := r.PathValue("fileId")
	if fileID == "" {
		requests.WriteError(w, http.StatusBadRequest, "fileId is required")
		return
	}

	var parsedRange *files.ParsedRange
	if raw := r.Header.Get("Range"); raw != "" {
		pr, ok := parseByteRange(raw)
		if !ok {
			requests.WriteError(w, http.StatusRequestedRangeNotSatisfiable, "invalid range")
			return
		}
		parsedRange = pr
	}

	for _, prefix := range []string{"photos/", "music/", "videos/"} {
		body, ct, contentRange, contentLength, err := h.fileStore.GetRange(r.Context(), prefix+fileID, parsedRange)
		if err == nil {
			streamFileRange(w, body, ct, contentRange, contentLength)
			return
		}
		if !errors.Is(err, store.ErrFileNotFound) {
			h.log.Error("ServeFile error", "key", prefix+fileID, "err", err)
			requests.WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}
	requests.WriteError(w, http.StatusNotFound, "file not found")
}

// ServeThumbnail handles GET /thumbnails/{thumbnailId}.
func (h *fileServingHandler) ServeThumbnail(w http.ResponseWriter, r *http.Request) {
	thumbnailID := r.PathValue("thumbnailId")
	if thumbnailID == "" {
		requests.WriteError(w, http.StatusBadRequest, "thumbnailId is required")
		return
	}

	body, ct, err := h.fileStore.Get(r.Context(), "thumbnails/"+thumbnailID)
	if err != nil {
		if errors.Is(err, store.ErrFileNotFound) {
			requests.WriteError(w, http.StatusNotFound, "thumbnail not found")
			return
		}
		h.log.Error("ServeThumbnail error", "key", "thumbnails/"+thumbnailID, "err", err)
		requests.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	streamFile(w, r, body, ct)
}
