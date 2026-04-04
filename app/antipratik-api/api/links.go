package api

import (
	"log"
	"net/http"

	"github.com/pratikluitel/antipratik/logic"
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
