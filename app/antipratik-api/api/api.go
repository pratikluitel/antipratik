// Package api contains the HTTP layer: request parsing, response serialization,
// and static file serving.
package api

import "net/http"

// PostHandler handles HTTP requests for post resources.
type PostHandler interface {
	GetPosts(w http.ResponseWriter, r *http.Request)
	GetPost(w http.ResponseWriter, r *http.Request)
}

// LinkHandler handles HTTP requests for external link resources.
type LinkHandler interface {
	GetLinks(w http.ResponseWriter, r *http.Request)
	GetFeaturedLinks(w http.ResponseWriter, r *http.Request)
}
