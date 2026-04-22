package api

import (
	"net/http"

	"github.com/pratikluitel/antipratik/components/auth"
)

// JWTAuthMiddleware validates the Bearer token on protected routes.
func JWTAuthMiddleware(authLogic auth.AuthLogic) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
				writeError(w, http.StatusUnauthorized, "missing or invalid authorization header")
				return
			}
			token := authHeader[7:]
			if err := authLogic.ValidateToken(r.Context(), token); err != nil {
				writeError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
