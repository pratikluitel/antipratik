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
