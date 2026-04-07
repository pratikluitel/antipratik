package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pratikluitel/antipratik/logic"
	"github.com/pratikluitel/antipratik/models"
)

// PostHandlerImpl implements PostHandler.
type PostHandlerImpl struct {
	logic logic.PostLogic
}

// NewPostHandler creates a new PostHandlerImpl using the given logic layer.
func NewPostHandler(l logic.PostLogic) *PostHandlerImpl {
	return &PostHandlerImpl{logic: l}
}

// GetPosts handles GET /api/posts
// Query params: type (repeatable), tag (repeatable)
func (h *PostHandlerImpl) GetPosts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := models.FilterState{
		ActiveTypes: q["type"],
		ActiveTags:  q["tag"],
	}

	posts, err := h.logic.GetPosts(r.Context(), filter)
	if err != nil {
		log.Printf("GetPosts error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, posts)
}

// GetPost handles GET /api/posts/{slug}
func (h *PostHandlerImpl) GetPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	post, err := h.logic.GetPost(r.Context(), slug)
	if err != nil {
		log.Printf("GetPost error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if post == nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	writeJSON(w, http.StatusOK, post)
}

// ── Write handlers ────────────────────────────────────────────────────────────

func (h *PostHandlerImpl) CreateEssay(w http.ResponseWriter, r *http.Request) {
	var input models.CreateEssayPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := h.logic.CreateEssay(r.Context(), input)
	if err != nil {
		log.Printf("CreateEssay error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *PostHandlerImpl) CreateShort(w http.ResponseWriter, r *http.Request) {
	var input models.CreateShortPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := h.logic.CreateShort(r.Context(), input)
	if err != nil {
		log.Printf("CreateShort error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *PostHandlerImpl) CreateMusic(w http.ResponseWriter, r *http.Request) {
	var input models.CreateMusicPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := h.logic.CreateMusic(r.Context(), input)
	if err != nil {
		log.Printf("CreateMusic error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *PostHandlerImpl) CreatePhoto(w http.ResponseWriter, r *http.Request) {
	var input models.CreatePhotoPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := h.logic.CreatePhoto(r.Context(), input)
	if err != nil {
		log.Printf("CreatePhoto error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *PostHandlerImpl) CreateVideo(w http.ResponseWriter, r *http.Request) {
	var input models.CreateVideoPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := h.logic.CreateVideo(r.Context(), input)
	if err != nil {
		log.Printf("CreateVideo error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *PostHandlerImpl) CreateLinkPost(w http.ResponseWriter, r *http.Request) {
	var input models.CreateLinkPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := h.logic.CreateLinkPost(r.Context(), input)
	if err != nil {
		log.Printf("CreateLinkPost error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *PostHandlerImpl) UpdateEssay(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input models.CreateEssayPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.logic.UpdateEssay(r.Context(), id, input); err != nil {
		log.Printf("UpdateEssay error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": id})
}

func (h *PostHandlerImpl) UpdateShort(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input models.CreateShortPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.logic.UpdateShort(r.Context(), id, input); err != nil {
		log.Printf("UpdateShort error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": id})
}

func (h *PostHandlerImpl) UpdateMusic(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input models.CreateMusicPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.logic.UpdateMusic(r.Context(), id, input); err != nil {
		log.Printf("UpdateMusic error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": id})
}

func (h *PostHandlerImpl) UpdatePhoto(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input models.CreatePhotoPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.logic.UpdatePhoto(r.Context(), id, input); err != nil {
		log.Printf("UpdatePhoto error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": id})
}

func (h *PostHandlerImpl) UpdateVideo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input models.CreateVideoPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.logic.UpdateVideo(r.Context(), id, input); err != nil {
		log.Printf("UpdateVideo error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": id})
}

func (h *PostHandlerImpl) UpdateLinkPost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input models.CreateLinkPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.logic.UpdateLinkPost(r.Context(), id, input); err != nil {
		log.Printf("UpdateLinkPost error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": id})
}

func (h *PostHandlerImpl) DeletePost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.logic.DeletePost(r.Context(), id); err != nil {
		log.Printf("DeletePost error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── Shared response helpers ───────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON encode error: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
