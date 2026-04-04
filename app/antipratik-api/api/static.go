package api

import (
	"net/http"
	"os"
	"path/filepath"
)

// NewSPAHandler returns an http.Handler that serves static files from dir.
// If the requested path does not exist as a file or directory, it falls back
// to serving dir/index.html (Next.js static export SPA routing).
func NewSPAHandler(dir string) http.Handler {
	fs := http.Dir(dir)
	fileServer := http.FileServer(fs)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(dir, filepath.Clean("/"+r.URL.Path))
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(dir, "index.html"))
			return
		}
		fileServer.ServeHTTP(w, r)
	})
}
