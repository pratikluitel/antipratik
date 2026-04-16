package main

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"github.com/pratikluitel/antipratik/api"
	"github.com/pratikluitel/antipratik/logic"
)

// RegisterRoutes registers all HTTP routes on mux.
// Middleware (CORS, JWT auth, rate limiting) is applied per-route here;
// it lives in api/middleware.go and is imported rather than defined alongside routes.
func RegisterRoutes(
	mux *http.ServeMux,
	postH api.PostHandler,
	linkH api.LinkHandler,
	authH *api.AuthHandlerImpl,
	authSvc logic.AuthLogic,
	fileH *api.FileServingHandler,
	newsletterH *api.NewsletterHandlerImpl,
	openAPIPath, swaggerPath string,
) {
	// Public file serving
	mux.HandleFunc("GET /files/{fileId}", fileH.ServeFile)
	mux.HandleFunc("GET /thumbnails/{thumbnailId}", fileH.ServeThumbnail)

	// Health check
	mux.HandleFunc("GET /api/health", api.HealthHandler)

	// Public read routes
	mux.HandleFunc("GET /api/posts/{slug}", postH.GetPost)
	mux.HandleFunc("GET /api/posts", postH.GetPosts)
	mux.HandleFunc("GET /api/tags", postH.GetTags)
	mux.HandleFunc("GET /api/links/featured", linkH.GetFeaturedLinks)
	mux.HandleFunc("GET /api/links", linkH.GetLinks)

	// Auth
	mux.HandleFunc("POST /api/auth/login", authH.Login)

	// Newsletter — 3 requests/hour per IP, burst of 3
	subscribeRL := api.RateLimitMiddleware(rate.Every(time.Hour/3), 3, time.Hour)
	mux.Handle("POST /api/subscribe", subscribeRL(http.HandlerFunc(newsletterH.Subscribe)))

	// OpenAPI spec + Swagger UI
	mux.HandleFunc("GET /api/openapi.yaml", api.OpenAPIHandler(openAPIPath))
	mux.HandleFunc("GET /api/index.html", api.SwaggerHandler(swaggerPath))

	// Protected write routes
	protect := api.JWTAuthMiddleware(authSvc)
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
	mux.Handle("POST /api/posts/{id}/images", protect(http.HandlerFunc(postH.AddPhotoImage)))
	mux.Handle("PUT /api/posts/{id}/images/{imageID}", protect(http.HandlerFunc(postH.UpdatePhotoImage)))
	mux.Handle("DELETE /api/posts/{id}/images/{imageID}", protect(http.HandlerFunc(postH.DeletePhotoImage)))
	mux.HandleFunc("GET /api/posts/{id}/images/{imageID}", postH.GetPhotoImage)
	mux.Handle("POST /api/links", protect(http.HandlerFunc(linkH.CreateLink)))
	mux.Handle("PUT /api/links/{id}", protect(http.HandlerFunc(linkH.UpdateLink)))
	mux.Handle("DELETE /api/links/{id}", protect(http.HandlerFunc(linkH.DeleteLink)))
}
