package api

import (
	"encoding/json"
	"net/http"

	"github.com/pratikluitel/antipratik/common/logging"
	"github.com/pratikluitel/antipratik/common/requests"
	"github.com/pratikluitel/antipratik/components/posts"
)

// linkHandler implements LinkHandler.
type linkHandler struct {
	logic posts.LinkLogic
	log   logging.Logger
}

// NewLinkHandler creates a new linkHandler using the given logic layer.
func NewLinkHandler(l posts.LinkLogic, log logging.Logger) posts.LinkHandler {
	return &linkHandler{logic: l, log: log}
}

// GetLinks handles GET /api/links
func (h *linkHandler) GetLinks(w http.ResponseWriter, r *http.Request) {
	links, err := h.logic.GetLinks(r.Context())
	if err != nil {
		handleLogicError(w, h.log, "GetLinks", err)
		return
	}
	requests.WriteJSON(w, http.StatusOK, links)
}

// GetFeaturedLinks handles GET /api/links/featured
func (h *linkHandler) GetFeaturedLinks(w http.ResponseWriter, r *http.Request) {
	links, err := h.logic.GetFeaturedLinks(r.Context())
	if err != nil {
		handleLogicError(w, h.log, "GetFeaturedLinks", err)
		return
	}
	requests.WriteJSON(w, http.StatusOK, links)
}

func (h *linkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	var input posts.CreateExternalLink
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := h.logic.CreateLink(r.Context(), input)
	if err != nil {
		handleLogicError(w, h.log, "CreateLink", err)
		return
	}
	requests.WriteJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *linkHandler) UpdateLink(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input posts.UpdateExternalLink
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	link, err := h.logic.UpdateLink(r.Context(), id, input)
	if err != nil {
		handleLogicError(w, h.log, "UpdateLink", err)
		return
	}
	requests.WriteJSON(w, http.StatusOK, link)
}

func (h *linkHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.logic.DeleteLink(r.Context(), id); err != nil {
		handleLogicError(w, h.log, "DeleteLink", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
