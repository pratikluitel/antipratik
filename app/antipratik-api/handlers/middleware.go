package handlers

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimitMiddleware returns a middleware that enforces per-IP rate limiting.
func RateLimitMiddleware(r rate.Limit, burst int, cleanupInterval time.Duration) func(http.Handler) http.Handler {
	var (
		mu      sync.Mutex
		clients = make(map[string]*ipLimiter)
	)

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
