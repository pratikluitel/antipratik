package api

import (
	"encoding/json"
	"net/http"

	"github.com/pratikluitel/antipratik/common/logging"
	"github.com/pratikluitel/antipratik/components/posts/logic"
	"github.com/pratikluitel/antipratik/components/posts/models"
)

// LinkHandlerImpl implements LinkHandler.
type LinkHandlerImpl struct {
	logic logic.LinkLogic
	log   logging.Logger
}

// NewLinkHandler creates a new LinkHandlerImpl using the given logic layer.
func NewLinkHandler(l logic.LinkLogic, log logging.Logger) *LinkHandlerImpl {
	return &LinkHandlerImpl{logic: l, log: log}
}

// GetLinks handles GET /api/links
func (h *LinkHandlerImpl) GetLinks(w http.ResponseWriter, r *http.Request) {
	links, err := h.logic.GetLinks(r.Context())
	if err != nil {
		handleLogicError(w, h.log, "GetLinks", err)
		return
	}
	writeJSON(w, http.StatusOK, links)
}

// GetFeaturedLinks handles GET /api/links/featured
func (h *LinkHandlerImpl) GetFeaturedLinks(w http.ResponseWriter, r *http.Request) {
	links, err := h.logic.GetFeaturedLinks(r.Context())
	if err != nil {
		handleLogicError(w, h.log, "GetFeaturedLinks", err)
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
		handleLogicError(w, h.log, "CreateLink", err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *LinkHandlerImpl) UpdateLink(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input models.UpdateExternalLink
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	link, err := h.logic.UpdateLink(r.Context(), id, input)
	if err != nil {
		handleLogicError(w, h.log, "UpdateLink", err)
		return
	}
	writeJSON(w, http.StatusOK, link)
}

func (h *LinkHandlerImpl) DeleteLink(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.logic.DeleteLink(r.Context(), id); err != nil {
		handleLogicError(w, h.log, "DeleteLink", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
