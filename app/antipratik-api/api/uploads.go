package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/pratikluitel/antipratik/store"
)

// FileServingHandler serves uploaded files and thumbnails from storage.
type FileServingHandler struct {
	fileStore store.FileStore
}

// NewFileServingHandler returns a new FileServingHandler.
func NewFileServingHandler(fs store.FileStore) *FileServingHandler {
	return &FileServingHandler{fileStore: fs}
}

// ServeFile handles GET /files/{fileId}.
// Tries photos/ prefix first, then music/.
func (h *FileServingHandler) ServeFile(w http.ResponseWriter, r *http.Request) {
	fileID := r.PathValue("fileId")
	if fileID == "" {
		writeError(w, http.StatusBadRequest, "fileId is required")
		return
	}

	for _, prefix := range []string{"photos/", "music/"} {
		body, ct, err := h.fileStore.Get(r.Context(), prefix+fileID)
		if err == nil {
			streamFile(w, r, body, ct)
			return
		}
		if !errors.Is(err, store.ErrFileNotFound) {
			log.Printf("ServeFile error (key=%s): %v", prefix+fileID, err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}
	writeError(w, http.StatusNotFound, "file not found")
}

// ServeThumbnail handles GET /thumbnails/{thumbnailId}.
func (h *FileServingHandler) ServeThumbnail(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("ServeThumbnail error (key=thumbnails/%s): %v", thumbnailID, err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	streamFile(w, r, body, ct)
}
