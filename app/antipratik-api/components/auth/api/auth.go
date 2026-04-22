package api

import (
	"encoding/json"
	"net/http"

	commonerrors "github.com/pratikluitel/antipratik/common/errors"
	"github.com/pratikluitel/antipratik/common/logging"
	"github.com/pratikluitel/antipratik/components/auth"
)

type authHandler struct {
	auth auth.AuthLogic
	log  logging.Logger
}

// NewAuthHandler creates a new authHandler.
func NewAuthHandler(a auth.AuthLogic, log logging.Logger) auth.AuthAPI {
	return &authHandler{auth: a, log: log}
}

// Login handles POST /api/auth/login
func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	token, err := h.auth.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		if commonerrors.Is(err) {
			writeError(w, http.StatusBadRequest, err.Error())
		} else {
			writeError(w, http.StatusUnauthorized, "invalid credentials")
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}
