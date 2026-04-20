// Package api contains the broadcaster HTTP layer.
package api

import (
	"encoding/json"
	"net/http"

	"github.com/pratikluitel/antipratik/common/logging"
	"github.com/pratikluitel/antipratik/components/broadcaster/logic"
)

// NewsletterHandlerImpl handles newsletter subscription requests.
type NewsletterHandlerImpl struct {
	logic logic.NewsletterLogic
	log   logging.Logger
}

// NewNewsletterHandler creates a new NewsletterHandlerImpl.
func NewNewsletterHandler(l logic.NewsletterLogic, log logging.Logger) *NewsletterHandlerImpl {
	return &NewsletterHandlerImpl{logic: l, log: log}
}

type subscribeRequest struct {
	Email string `json:"email"`
}

// Subscribe handles POST /api/subscribe.
func (h *NewsletterHandlerImpl) Subscribe(w http.ResponseWriter, r *http.Request) {
	var input subscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.logic.Subscribe(r.Context(), input.Email); err != nil {
		handleLogicError(w, h.log, "Subscribe", err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{})
}
