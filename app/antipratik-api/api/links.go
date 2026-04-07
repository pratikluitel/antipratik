package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pratikluitel/antipratik/logic"
	"github.com/pratikluitel/antipratik/models"
)

// LinkHandlerImpl implements LinkHandler.
type LinkHandlerImpl struct {
	logic logic.LinkLogic
}

// NewLinkHandler creates a new LinkHandlerImpl using the given logic layer.
func NewLinkHandler(l logic.LinkLogic) *LinkHandlerImpl {
	return &LinkHandlerImpl{logic: l}
}

// GetLinks handles GET /api/links
func (h *LinkHandlerImpl) GetLinks(w http.ResponseWriter, r *http.Request) {
	links, err := h.logic.GetLinks(r.Context())
	if err != nil {
		log.Printf("GetLinks error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, links)
}

// GetFeaturedLinks handles GET /api/links/featured
func (h *LinkHandlerImpl) GetFeaturedLinks(w http.ResponseWriter, r *http.Request) {
	links, err := h.logic.GetFeaturedLinks(r.Context())
	if err != nil {
		log.Printf("GetFeaturedLinks error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, links)
}

func (h *LinkHandlerImpl) CreateLink(w http.ResponseWriter, r *http.Request) {
	var input models.CreateExternalLink
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := h.logic.CreateLink(r.Context(), input)
	if err != nil {
		log.Printf("CreateLink error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *LinkHandlerImpl) UpdateLink(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input models.CreateExternalLink
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.logic.UpdateLink(r.Context(), id, input); err != nil {
		log.Printf("UpdateLink error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": id})
}

func (h *LinkHandlerImpl) DeleteLink(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.logic.DeleteLink(r.Context(), id); err != nil {
		log.Printf("DeleteLink error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
