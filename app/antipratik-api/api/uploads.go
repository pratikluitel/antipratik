package api

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/pratikluitel/antipratik/logic"
	"github.com/pratikluitel/antipratik/store"
)

func streamFile(w http.ResponseWriter, r *http.Request, body io.ReadCloser, ct string) {
	defer body.Close()
	w.Header().Set("Content-Type", ct)
	if _, err := io.Copy(w, body); err != nil && r.Context().Err() == nil {
		log.Printf("streamFile copy error: %v", err)
	}
}

// UploadHandler handles file upload and file serving endpoints.
type UploadHandler interface {
	UploadPhoto(w http.ResponseWriter, r *http.Request)
	UploadMusic(w http.ResponseWriter, r *http.Request)
	ServeFile(w http.ResponseWriter, r *http.Request)
	ServeThumbnail(w http.ResponseWriter, r *http.Request)
}

// UploadHandlerImpl is the concrete implementation of UploadHandler.
type UploadHandlerImpl struct {
	logic     logic.UploadLogic
	fileStore store.FileStore
}

// NewUploadHandler returns a new UploadHandlerImpl.
func NewUploadHandler(l logic.UploadLogic, fs store.FileStore) *UploadHandlerImpl {
	return &UploadHandlerImpl{logic: l, fileStore: fs}
}

// UploadPhoto handles POST /uploads/photos.
// Expects multipart/form-data with fields: postId (string) and file (binary).
func (h *UploadHandlerImpl) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	postID := r.FormValue("postId")
	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file field is required")
		return
	}
	defer file.Close()

	resp, err := h.logic.UploadPhoto(r.Context(), postID, file, header)
	if err != nil {
		handleLogicError(w, "UploadPhoto", err)
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

// UploadMusic handles POST /uploads/music.
// Expects multipart/form-data with fields: postId (string) and file (binary).
func (h *UploadHandlerImpl) UploadMusic(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	postID := r.FormValue("postId")
	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file field is required")
		return
	}
	defer file.Close()

	resp, err := h.logic.UploadMusic(r.Context(), postID, file, header)
	if err != nil {
		handleLogicError(w, "UploadMusic", err)
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

// ServeFile handles GET /files/{fileId}.
// Tries photos/ prefix first, then music/.
func (h *UploadHandlerImpl) ServeFile(w http.ResponseWriter, r *http.Request) {
	fileID := r.PathValue("fileId")
	if fileID == "" {
		writeError(w, http.StatusBadRequest, "fileId is required")
		return
	}

	// Try photos first, then music.
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
func (h *UploadHandlerImpl) ServeThumbnail(w http.ResponseWriter, r *http.Request) {
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
