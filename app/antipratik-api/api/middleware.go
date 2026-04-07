package api

import (
	"net/http"

	"github.com/pratikluitel/antipratik/logic"
)

// CORSMiddleware adds permissive CORS headers to all responses.
// OPTIONS preflight requests are handled with 204 No Content.
// Tighten the origin in production by checking r.Header.Get("Origin").
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// JWTAuthMiddleware validates the Bearer token on protected routes.
func JWTAuthMiddleware(authLogic logic.AuthLogic) func(http.Handler) http.Handler {
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
