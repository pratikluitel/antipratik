package api

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/pratikluitel/antipratik/logic"
)

// ipLimiter holds a rate limiter and the last time it was seen,
// so idle entries can be evicted periodically.
type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimitMiddleware returns a middleware that enforces per-IP rate limiting
// on the wrapped handler. Each IP is allowed burst requests initially, then
// refills at r tokens/second. IPs that haven't been seen for cleanupInterval
// are evicted to prevent unbounded memory growth.
func RateLimitMiddleware(r rate.Limit, burst int, cleanupInterval time.Duration) func(http.Handler) http.Handler {
	var (
		mu      sync.Mutex
		clients = make(map[string]*ipLimiter)
	)

	// Background goroutine evicts stale entries.
	go func() {
		for {
			time.Sleep(cleanupInterval)
			mu.Lock()
			for ip, cl := range clients {
				if time.Since(cl.lastSeen) > cleanupInterval {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	getLimiter := func(ip string) *rate.Limiter {
		mu.Lock()
		defer mu.Unlock()
		cl, ok := clients[ip]
		if !ok {
			cl = &ipLimiter{limiter: rate.NewLimiter(r, burst)}
			clients[ip] = cl
		}
		cl.lastSeen = time.Now()
		return cl.limiter
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ip, _, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				ip = req.RemoteAddr
			}
			if !getLimiter(ip).Allow() {
				writeError(w, http.StatusTooManyRequests, "too many requests")
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}

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
