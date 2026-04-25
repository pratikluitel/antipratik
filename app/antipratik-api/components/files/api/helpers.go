// Package api contains the files component HTTP layer (file serving).
package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

func streamFile(w http.ResponseWriter, r *http.Request, body io.ReadSeekCloser, ct string) {
	defer func() { _ = body.Close() }()
	w.Header().Set("Content-Type", ct)
	http.ServeContent(w, r, "", time.Time{}, body)
}

// streamFileRange writes a range-aware file response.
// If contentRange is non-empty the response is 206 Partial Content; otherwise 200 OK.
func streamFileRange(w http.ResponseWriter, body io.ReadCloser, ct, contentRange string, contentLength int64) {
	defer func() { _ = body.Close() }()
	w.Header().Set("Content-Type", ct)
	w.Header().Set("Accept-Ranges", "bytes")
	if contentLength >= 0 {
		w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
	}
	if contentRange != "" {
		w.Header().Set("Content-Range", contentRange)
		w.WriteHeader(http.StatusPartialContent)
	}
	_, _ = io.Copy(w, body)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
