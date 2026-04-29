package api

import (
	"net/http"

	commonerrors "github.com/pratikluitel/antipratik/common/errors"
	"github.com/pratikluitel/antipratik/common/logging"
	"github.com/pratikluitel/antipratik/common/requests"
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
	req := &auth.LoginRequest{}
	if err := requests.DecodeJSONBody(w, r, req); err != nil {
		return
	}

	token, err := h.auth.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		if commonerrors.Is(err) {
			requests.WriteError(w, http.StatusBadRequest, err.Error())
		} else {
			requests.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		}
		return
	}

	requests.WriteJSON(w, http.StatusOK, auth.Token{Token: token})
}
