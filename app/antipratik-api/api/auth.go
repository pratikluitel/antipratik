package api

import (
	"encoding/json"
	"net/http"

	"github.com/pratikluitel/antipratik/logic"
)

// AuthHandlerImpl handles authentication HTTP requests.
type AuthHandlerImpl struct {
	auth logic.AuthLogic
}

// NewAuthHandler creates a new AuthHandlerImpl.
func NewAuthHandler(auth logic.AuthLogic) *AuthHandlerImpl {
	return &AuthHandlerImpl{auth: auth}
}

// Login handles POST /api/auth/login
func (h *AuthHandlerImpl) Login(w http.ResponseWriter, r *http.Request) {
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
		if logic.IsValidationError(err) {
			writeError(w, http.StatusBadRequest, err.Error())
		} else {
			writeError(w, http.StatusUnauthorized, "invalid credentials")
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}
