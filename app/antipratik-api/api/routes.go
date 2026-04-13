package api

import (
	"net/http"
	"os"

	"github.com/pratikluitel/antipratik/logic"
)

// RegisterRoutes registers all HTTP routes on mux.
func RegisterRoutes(mux *http.ServeMux, postH PostHandler, linkH LinkHandler, authH *AuthHandlerImpl, authSvc logic.AuthLogic, fileH *FileServingHandler, newsletterH *NewsletterHandlerImpl, openAPIPath, swaggerPath string) {
	// Public file serving routes
	mux.HandleFunc("GET /files/{fileId}", fileH.ServeFile)
	mux.HandleFunc("GET /thumbnails/{thumbnailId}", fileH.ServeThumbnail)

	// Health check
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "OK"})
	})

	// Public read routes
	mux.HandleFunc("GET /api/posts/{slug}", postH.GetPost)
	mux.HandleFunc("GET /api/posts", postH.GetPosts)
	mux.HandleFunc("GET /api/links/featured", linkH.GetFeaturedLinks)
	mux.HandleFunc("GET /api/links", linkH.GetLinks)

	// Auth
	mux.HandleFunc("POST /api/auth/login", authH.Login)

	// Newsletter
	mux.HandleFunc("POST /api/subscribe", newsletterH.Subscribe)

	// OpenAPI spec + Swagger UI
	mux.HandleFunc("GET /api/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(openAPIPath)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/yaml")
		w.Write(data)
	})
	mux.HandleFunc("GET /api/index.html", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(swaggerPath)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(data)
	})

	// Protected write routes
	protect := JWTAuthMiddleware(authSvc)
	mux.Handle("POST /api/posts/essay", protect(http.HandlerFunc(postH.CreateEssay)))
	mux.Handle("POST /api/posts/short", protect(http.HandlerFunc(postH.CreateShort)))
	mux.Handle("POST /api/posts/music", protect(http.HandlerFunc(postH.CreateMusic)))
	mux.Handle("POST /api/posts/photo", protect(http.HandlerFunc(postH.CreatePhoto)))
	mux.Handle("POST /api/posts/video", protect(http.HandlerFunc(postH.CreateVideo)))
	mux.Handle("POST /api/posts/link", protect(http.HandlerFunc(postH.CreateLinkPost)))
	mux.Handle("PUT /api/posts/essay/{id}", protect(http.HandlerFunc(postH.UpdateEssay)))
	mux.Handle("PUT /api/posts/short/{id}", protect(http.HandlerFunc(postH.UpdateShort)))
	mux.Handle("PUT /api/posts/music/{id}", protect(http.HandlerFunc(postH.UpdateMusic)))
	mux.Handle("PUT /api/posts/photo/{id}", protect(http.HandlerFunc(postH.UpdatePhoto)))
	mux.Handle("PUT /api/posts/video/{id}", protect(http.HandlerFunc(postH.UpdateVideo)))
	mux.Handle("PUT /api/posts/link/{id}", protect(http.HandlerFunc(postH.UpdateLinkPost)))
	mux.Handle("DELETE /api/posts/{id}", protect(http.HandlerFunc(postH.DeletePost)))
	mux.Handle("POST /api/links", protect(http.HandlerFunc(linkH.CreateLink)))
	mux.Handle("PUT /api/links/{id}", protect(http.HandlerFunc(linkH.UpdateLink)))
	mux.Handle("DELETE /api/links/{id}", protect(http.HandlerFunc(linkH.DeleteLink)))
}
