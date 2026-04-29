// Package api contains the broadcaster HTTP layer.
package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pratikluitel/antipratik/common/logging"
	"github.com/pratikluitel/antipratik/common/requests"
	"github.com/pratikluitel/antipratik/components/broadcaster"
	"github.com/pratikluitel/antipratik/components/broadcaster/store"
)

// broadcasterHandler handles all broadcaster HTTP endpoints.
type broadcasterHandler struct {
	logic broadcaster.BroadcasterLogic
	log   logging.Logger
}

// NewBroadcasterHandler creates a new broadcasterHandler.
func NewBroadcasterHandler(l broadcaster.BroadcasterLogic, log logging.Logger) broadcaster.BroadcasterAPI {
	return &broadcasterHandler{logic: l, log: log}
}

// ── Subscriber endpoints ──────────────────────────────────────────────────────

type subscribeRequest struct {
	Type    string `json:"type"`
	Address string `json:"address"`
}

// Subscribe handles POST /api/subscribe.
func (h *broadcasterHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var input subscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.Type == "" {
		input.Type = "email"
	}
	if err := h.logic.Subscribe(r.Context(), input.Type, input.Address); err != nil {
		handleLogicError(w, h.log, "Subscribe", err)
		return
	}
	requests.WriteJSON(w, http.StatusCreated, map[string]any{})
}

// ResendConfirmation handles POST /api/subscribers/resend-confirmation.
func (h *broadcasterHandler) ResendConfirmation(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Type string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Type == "" {
		body.Type = "email"
	}
	n, err := h.logic.SendConfirmationEmails(r.Context(), body.Type)
	if err != nil {
		handleLogicError(w, h.log, "ResendConfirmation", err)
		return
	}
	requests.WriteJSON(w, http.StatusOK, map[string]any{"sent_count": n})
}

// Confirm handles GET /api/confirm?token=<token>.
func (h *broadcasterHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if err := h.logic.ConfirmSubscription(r.Context(), token); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			requests.WriteError(w, http.StatusNotFound, "token not found")
			return
		}
		handleLogicError(w, h.log, "Confirm", err)
		return
	}
	requests.WriteJSON(w, http.StatusOK, map[string]any{"message": "subscription confirmed"})
}

// Unsubscribe handles GET /api/unsubscribe?token=<token>.
func (h *broadcasterHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if err := h.logic.Unsubscribe(r.Context(), token); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			requests.WriteError(w, http.StatusNotFound, "token not found")
			return
		}
		handleLogicError(w, h.log, "Unsubscribe", err)
		return
	}
	requests.WriteJSON(w, http.StatusOK, map[string]any{"message": "unsubscribed"})
}

// GetSubscribers handles GET /api/subscribers.
func (h *broadcasterHandler) GetSubscribers(w http.ResponseWriter, r *http.Request) {
	subType := r.URL.Query().Get("type")
	if subType == "" {
		subType = "email"
	}
	subs, err := h.logic.GetSubscribers(r.Context(), subType)
	if err != nil {
		handleLogicError(w, h.log, "GetSubscribers", err)
		return
	}
	requests.WriteJSON(w, http.StatusOK, subs)
}

// DeleteSubscriber handles DELETE /api/subscribers/{address}.
func (h *broadcasterHandler) DeleteSubscriber(w http.ResponseWriter, r *http.Request) {
	address := r.PathValue("address")
	if address == "" {
		requests.WriteError(w, http.StatusBadRequest, "address is required")
		return
	}
	if err := h.logic.DeleteSubscriber(r.Context(), address); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			requests.WriteError(w, http.StatusNotFound, "subscriber not found")
			return
		}
		handleLogicError(w, h.log, "DeleteSubscriber", err)
		return
	}
	requests.WriteJSON(w, http.StatusNoContent, map[string]any{})
}

// ── Broadcast endpoints ───────────────────────────────────────────────────────

type broadcastRequest struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Data  struct {
		Caption string   `json:"caption"`
		PostIDs []string `json:"postIDs"`
	} `json:"data"`
}

type broadcastUpdateRequest struct {
	Title string `json:"title"`
	Data  struct {
		Caption string   `json:"caption"`
		PostIDs []string `json:"postIDs"`
	} `json:"data"`
}

// CreateBroadcast handles POST /api/broadcasts.
func (h *broadcasterHandler) CreateBroadcast(w http.ResponseWriter, r *http.Request) {
	var req broadcastRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	preview, err := h.logic.CreateBroadcast(r.Context(), broadcaster.BroadcastInput{
		Type:    req.Type,
		Title:   req.Title,
		Caption: req.Data.Caption,
		PostIDs: req.Data.PostIDs,
	})
	if err != nil {
		handleLogicError(w, h.log, "CreateBroadcast", err)
		return
	}
	requests.WriteJSON(w, http.StatusCreated, map[string]any{"id": preview.ID, "html": preview.HTML})
}

// UpdateBroadcast handles PUT /api/broadcasts/{id}.
func (h *broadcasterHandler) UpdateBroadcast(w http.ResponseWriter, r *http.Request) {
	id, ok := parseBroadcastID(w, r)
	if !ok {
		return
	}
	var req broadcastUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	preview, err := h.logic.UpdateBroadcast(r.Context(), id, broadcaster.BroadcastUpdateInput{
		Title:   req.Title,
		Caption: req.Data.Caption,
		PostIDs: req.Data.PostIDs,
	})
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			requests.WriteError(w, http.StatusNotFound, "broadcast not found")
			return
		}
		handleLogicError(w, h.log, "UpdateBroadcast", err)
		return
	}
	requests.WriteJSON(w, http.StatusOK, map[string]any{"id": preview.ID, "html": preview.HTML})
}

// DeleteBroadcast handles DELETE /api/broadcasts/{id}.
func (h *broadcasterHandler) DeleteBroadcast(w http.ResponseWriter, r *http.Request) {
	id, ok := parseBroadcastID(w, r)
	if !ok {
		return
	}
	if err := h.logic.DeleteBroadcast(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			requests.WriteError(w, http.StatusNotFound, "broadcast not found")
			return
		}
		handleLogicError(w, h.log, "DeleteBroadcast", err)
		return
	}
	requests.WriteJSON(w, http.StatusNoContent, map[string]any{})
}

// GetBroadcasts handles GET /api/broadcasts?type=email.
func (h *broadcasterHandler) GetBroadcasts(w http.ResponseWriter, r *http.Request) {
	bType := r.URL.Query().Get("type")
	if bType == "" {
		bType = "email"
	}
	summaries, err := h.logic.GetBroadcasts(r.Context(), bType)
	if err != nil {
		handleLogicError(w, h.log, "GetBroadcasts", err)
		return
	}
	requests.WriteJSON(w, http.StatusOK, summaries)
}

// DispatchBroadcast handles POST /api/broadcasts/{id}/dispatch.
func (h *broadcasterHandler) DispatchBroadcast(w http.ResponseWriter, r *http.Request) {
	id, ok := parseBroadcastID(w, r)
	if !ok {
		return
	}
	n, err := h.logic.DispatchBroadcast(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			requests.WriteError(w, http.StatusNotFound, "broadcast not found")
			return
		}
		handleLogicError(w, h.log, "DispatchBroadcast", err)
		return
	}
	requests.WriteJSON(w, http.StatusOK, map[string]any{"buffered_count": n})
}

// GetBroadcastSendDetails handles GET /api/broadcasts/{id}/sends.
func (h *broadcasterHandler) GetBroadcastSendDetails(w http.ResponseWriter, r *http.Request) {
	id, ok := parseBroadcastID(w, r)
	if !ok {
		return
	}
	details, err := h.logic.GetBroadcastSends(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			requests.WriteError(w, http.StatusNotFound, "broadcast not found")
			return
		}
		handleLogicError(w, h.log, "GetBroadcastSendDetails", err)
		return
	}
	requests.WriteJSON(w, http.StatusOK, details)
}

// ── Contact endpoint ──────────────────────────────────────────────────────────

type contactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// Contact handles POST /api/contact.
func (h *broadcasterHandler) Contact(w http.ResponseWriter, r *http.Request) {
	var req contactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.logic.SendContactMessage(r.Context(), broadcaster.ContactInput{
		Name:    req.Name,
		Email:   req.Email,
		Message: req.Message,
	}); err != nil {
		handleLogicError(w, h.log, "Contact", err)
		return
	}
	requests.WriteJSON(w, http.StatusOK, map[string]any{})
}
