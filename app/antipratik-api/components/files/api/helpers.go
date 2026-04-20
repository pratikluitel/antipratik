// Package api contains the files component HTTP layer (file serving).
package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func streamFile(w http.ResponseWriter, r *http.Request, body io.ReadSeekCloser, ct string) {
	defer func() { _ = body.Close() }()
	w.Header().Set("Content-Type", ct)
	http.ServeContent(w, r, "", time.Time{}, body)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
