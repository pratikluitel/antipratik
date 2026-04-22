package handlers

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"github.com/pratikluitel/antipratik/components/auth"
	authapi "github.com/pratikluitel/antipratik/components/auth/api"
	"github.com/pratikluitel/antipratik/components/broadcaster"
	"github.com/pratikluitel/antipratik/components/files"
	"github.com/pratikluitel/antipratik/components/posts"
)

// RegisterRoutes registers all HTTP routes on mux.
func RegisterRoutes(
	mux *http.ServeMux,
	postH posts.PostHandler,
	linkH posts.LinkHandler,
	authH auth.AuthAPI,
	authSvc auth.AuthLogic,
	fileH files.FilesAPI,
	broadcasterH broadcaster.BroadcasterAPI,
	openAPIPath, swaggerPath string,
) {
	// Public file serving
	mux.HandleFunc("GET /files/{fileId}", fileH.ServeFile)
	mux.HandleFunc("GET /thumbnails/{thumbnailId}", fileH.ServeThumbnail)

	// Health check
	mux.HandleFunc("GET /api/health", HealthHandler)

	// Public read routes
	mux.HandleFunc("GET /api/posts/{slug}", postH.GetPost)
	mux.HandleFunc("GET /api/posts", postH.GetPosts)
	mux.HandleFunc("GET /api/tags", postH.GetTags)
	mux.HandleFunc("GET /api/links/featured", linkH.GetFeaturedLinks)
	mux.HandleFunc("GET /api/links", linkH.GetLinks)

	// Auth
	mux.HandleFunc("POST /api/auth/login", authH.Login)

	// Public broadcaster endpoints — rate limited
	subscribeRL := RateLimitMiddleware(rate.Every(time.Hour/3), 3, time.Hour)
	mux.Handle("POST /api/subscribe", subscribeRL(http.HandlerFunc(broadcasterH.Subscribe)))
	contactRL := RateLimitMiddleware(rate.Every(time.Hour/3), 3, time.Hour)
	mux.Handle("POST /api/contact", contactRL(http.HandlerFunc(broadcasterH.Contact)))

	// Token-based subscriber actions (no auth, no rate limit — tokens are one-time-use)
	mux.HandleFunc("GET /api/confirm", broadcasterH.Confirm)
	mux.HandleFunc("GET /api/unsubscribe", broadcasterH.Unsubscribe)

	// OpenAPI spec + Swagger UI
	mux.HandleFunc("GET /api/openapi.yaml", OpenAPIHandler(openAPIPath))
	mux.HandleFunc("GET /api/index.html", SwaggerHandler(swaggerPath))

	// Protected write routes
	protect := authapi.JWTAuthMiddleware(authSvc)
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

	// Protected broadcaster endpoints
	mux.Handle("POST /api/broadcasts", protect(http.HandlerFunc(broadcasterH.CreateBroadcast)))
	mux.Handle("PUT /api/broadcasts/{id}", protect(http.HandlerFunc(broadcasterH.UpdateBroadcast)))
	mux.Handle("GET /api/broadcasts", protect(http.HandlerFunc(broadcasterH.GetBroadcasts)))
	mux.Handle("POST /api/broadcasts/{id}/dispatch", protect(http.HandlerFunc(broadcasterH.DispatchBroadcast)))
	mux.Handle("GET /api/broadcasts/{id}/sends", protect(http.HandlerFunc(broadcasterH.GetBroadcastSendDetails)))
	mux.Handle("POST /api/subscribers/resend-confirmation", protect(http.HandlerFunc(broadcasterH.ResendConfirmation)))
	mux.Handle("GET /api/subscribers", protect(http.HandlerFunc(broadcasterH.GetSubscribers)))
}
