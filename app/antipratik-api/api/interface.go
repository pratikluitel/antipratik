// Package api contains the HTTP layer: request parsing, response serialization,
// and static file serving.
package api

import "net/http"

// PostHandler handles HTTP requests for post resources.
type PostHandler interface {
	GetPosts(w http.ResponseWriter, r *http.Request)
	GetPost(w http.ResponseWriter, r *http.Request)
	CreateEssay(w http.ResponseWriter, r *http.Request)
	CreateShort(w http.ResponseWriter, r *http.Request)
	CreateMusic(w http.ResponseWriter, r *http.Request)
	CreatePhoto(w http.ResponseWriter, r *http.Request)
	CreateVideo(w http.ResponseWriter, r *http.Request)
	CreateLinkPost(w http.ResponseWriter, r *http.Request)
	UpdateEssay(w http.ResponseWriter, r *http.Request)
	UpdateShort(w http.ResponseWriter, r *http.Request)
	UpdateMusic(w http.ResponseWriter, r *http.Request)
	UpdatePhoto(w http.ResponseWriter, r *http.Request)
	UpdateVideo(w http.ResponseWriter, r *http.Request)
	UpdateLinkPost(w http.ResponseWriter, r *http.Request)
	DeletePost(w http.ResponseWriter, r *http.Request)
	AddPhotoImage(w http.ResponseWriter, r *http.Request)
	GetPhotoImage(w http.ResponseWriter, r *http.Request)
	UpdatePhotoImage(w http.ResponseWriter, r *http.Request)
	DeletePhotoImage(w http.ResponseWriter, r *http.Request)
}

// LinkHandler handles HTTP requests for external link resources.
type LinkHandler interface {
	GetLinks(w http.ResponseWriter, r *http.Request)
	GetFeaturedLinks(w http.ResponseWriter, r *http.Request)
	CreateLink(w http.ResponseWriter, r *http.Request)
	UpdateLink(w http.ResponseWriter, r *http.Request)
	DeleteLink(w http.ResponseWriter, r *http.Request)
}
