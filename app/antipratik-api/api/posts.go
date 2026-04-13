package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/pratikluitel/antipratik/common/logging"
	"github.com/pratikluitel/antipratik/logic"
	"github.com/pratikluitel/antipratik/models"
)

// PostHandlerImpl implements PostHandler.
type PostHandlerImpl struct {
	logic   logic.PostLogic
	uploads logic.UploadLogic
	log     logging.Logger
}

// NewPostHandler creates a new PostHandlerImpl.
// uploads handles file storage for photo, music, video, and link post types.
func NewPostHandler(l logic.PostLogic, u logic.UploadLogic, log logging.Logger) *PostHandlerImpl {
	return &PostHandlerImpl{logic: l, uploads: u, log: log}
}

// GetPosts handles GET /api/posts
// Query params: type (repeatable), tag (repeatable)
func (h *PostHandlerImpl) GetPosts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := models.FilterState{
		ActiveTypes: q["type"],
		ActiveTags:  q["tag"],
	}

	posts, err := h.logic.GetPosts(r.Context(), filter)
	if err != nil {
		handleLogicError(w, h.log, "GetPosts", err)
		return
	}

	writeJSON(w, http.StatusOK, posts)
}

// GetPost handles GET /api/posts/{slug}
func (h *PostHandlerImpl) GetPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	post, err := h.logic.GetPost(r.Context(), slug)
	if err != nil {
		handleLogicError(w, h.log, "GetPost", err)
		return
	}
	if post == nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	writeJSON(w, http.StatusOK, post)
}

// ── JSON-body write handlers (essay, short) ───────────────────────────────────

func (h *PostHandlerImpl) CreateEssay(w http.ResponseWriter, r *http.Request) {
	var input models.CreateEssayPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	post, err := h.logic.CreateEssay(r.Context(), input)
	if err != nil {
		handleLogicError(w, h.log, "CreateEssay", err)
		return
	}
	writeJSON(w, http.StatusCreated, post)
}

func (h *PostHandlerImpl) CreateShort(w http.ResponseWriter, r *http.Request) {
	var input models.CreateShortPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	post, err := h.logic.CreateShort(r.Context(), input)
	if err != nil {
		handleLogicError(w, h.log, "CreateShort", err)
		return
	}
	writeJSON(w, http.StatusCreated, post)
}

func (h *PostHandlerImpl) UpdateEssay(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input models.UpdateEssayPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	post, err := h.logic.UpdateEssay(r.Context(), id, input)
	if err != nil {
		handleLogicError(w, h.log, "UpdateEssay", err)
		return
	}
	writeJSON(w, http.StatusOK, post)
}

func (h *PostHandlerImpl) UpdateShort(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input models.UpdateShortPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	post, err := h.logic.UpdateShort(r.Context(), id, input)
	if err != nil {
		handleLogicError(w, h.log, "UpdateShort", err)
		return
	}
	writeJSON(w, http.StatusOK, post)
}

// ── Multipart write handlers (music, photo, video, link) ─────────────────────
//
// Form field conventions:
//   Music:  audioFile (binary), albumArtFile (binary, optional),
//           title, duration (int), album (optional), tags[] (repeated)
//   Photo:  images[] (binary, repeated), alt[] (repeated), caption[] (repeated, optional),
//           location (optional), tags[] (repeated)
//   Video:  thumbnailFile (binary, optional), title, videoURL, duration (int),
//           playlist (optional), tags[] (repeated)
//   Link:   thumbnailFile (binary, optional), title, url, domain,
//           description (optional), category (optional), tags[] (repeated)

func (h *PostHandlerImpl) CreateMusic(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	postID := uuid.New().String()

	audioFile, audioHeader, err := r.FormFile("audioFile")
	if err != nil {
		writeError(w, http.StatusBadRequest, "audioFile is required")
		return
	}
	defer audioFile.Close()

	var albumArtInput *logic.FileInput
	if artFile, artHeader, err := r.FormFile("albumArtFile"); err == nil {
		defer artFile.Close()
		fi := logic.FileInput{File: artFile, Header: artHeader}
		albumArtInput = &fi
	}

	audioInput := &logic.FileInput{File: audioFile, Header: audioHeader}
	uploaded, err := h.uploads.UploadMusicFiles(r.Context(), postID, audioInput, albumArtInput)
	if err != nil {
		handleLogicError(w, h.log, "CreateMusic upload", err)
		return
	}

	duration, err := strconv.Atoi(r.FormValue("duration"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "duration must be an integer")
		return
	}

	input := models.CreateMusicPost{
		Title:           r.FormValue("title"),
		AudioURL:        uploaded.AudioURL,
		AlbumArt:        uploaded.AlbumArtURL,
		AlbumArtTinyURL: uploaded.AlbumArtTinyURL,
		Duration:        duration,
		Tags:            r.Form["tags"],
	}
	if album := r.FormValue("album"); album != "" {
		input.Album = &album
	}

	post, err := h.logic.CreateMusic(r.Context(), postID, input)
	if err != nil {
		handleLogicError(w, h.log, "CreateMusic", err)
		return
	}
	writeJSON(w, http.StatusCreated, post)
}

func (h *PostHandlerImpl) UpdateMusic(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	input := models.UpdateMusicPost{Tags: formTags(r)}
	if title := r.FormValue("title"); title != "" {
		input.Title = &title
	}
	if durStr := r.FormValue("duration"); durStr != "" {
		if d, err := strconv.Atoi(durStr); err == nil {
			input.Duration = &d
		} else {
			writeError(w, http.StatusBadRequest, "duration must be an integer")
			return
		}
	}
	if albumStr := r.FormValue("album"); albumStr != "" {
		input.Album = &albumStr
	}

	post, err := h.logic.UpdateMusic(r.Context(), postID, input)
	if err != nil {
		handleLogicError(w, h.log, "UpdateMusic", err)
		return
	}
	writeJSON(w, http.StatusOK, post)
}

func (h *PostHandlerImpl) CreatePhoto(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	fileHeaders := r.MultipartForm.File["images[]"]
	if len(fileHeaders) == 0 {
		writeError(w, http.StatusBadRequest, "at least one image file is required")
		return
	}

	postID := uuid.New().String()

	fileInputs := make([]logic.FileInput, 0, len(fileHeaders))
	for _, fh := range fileHeaders {
		f, err := fh.Open()
		if err != nil {
			writeError(w, http.StatusBadRequest, "could not read uploaded image")
			return
		}
		defer f.Close()
		fileInputs = append(fileInputs, logic.FileInput{File: f, Header: fh})
	}

	uploadResults, err := h.uploads.UploadPhotoFiles(r.Context(), postID, fileInputs)
	if err != nil {
		handleLogicError(w, h.log, "CreatePhoto upload", err)
		return
	}

	alts := r.Form["alt[]"]
	captions := r.Form["caption[]"]
	images := make([]models.PhotoImage, len(uploadResults))
	for i, u := range uploadResults {
		alt := ""
		if i < len(alts) {
			alt = alts[i]
		}
		tiny, small, med, large := u.ThumbnailTinyURL, u.ThumbnailSmallURL, u.ThumbnailMedURL, u.ThumbnailLargeURL
		images[i] = models.PhotoImage{
			URL:              u.OriginalURL,
			Alt:              alt,
			ThumbnailTinyURL:  &tiny,
			ThumbnailSmallURL: &small,
			ThumbnailMedURL:   &med,
			ThumbnailLargeURL: &large,
		}
		if i < len(captions) && captions[i] != "" {
			c := captions[i]
			images[i].Caption = &c
		}
	}

	input := models.CreatePhotoPost{Images: images, Tags: r.Form["tags[]"]}
	if loc := r.FormValue("location"); loc != "" {
		input.Location = &loc
	}

	post, err := h.logic.CreatePhoto(r.Context(), postID, input)
	if err != nil {
		handleLogicError(w, h.log, "CreatePhoto", err)
		return
	}
	writeJSON(w, http.StatusCreated, post)
}

func (h *PostHandlerImpl) UpdatePhoto(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	input := models.UpdatePhotoPost{Tags: formTags(r)}
	if loc := r.FormValue("location"); loc != "" {
		input.Location = &loc
	}

	post, err := h.logic.UpdatePhoto(r.Context(), postID, input)
	if err != nil {
		handleLogicError(w, h.log, "UpdatePhoto", err)
		return
	}
	writeJSON(w, http.StatusOK, post)
}

func (h *PostHandlerImpl) CreateVideo(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	postID := uuid.New().String()

	var thumbnailURL, thumbnailTinyURL string
	if thumbFile, thumbHeader, err := r.FormFile("thumbnailFile"); err == nil {
		defer thumbFile.Close()
		result, err := h.uploads.UploadThumbnail(r.Context(), postID, "thumb",
			logic.FileInput{File: thumbFile, Header: thumbHeader})
		if err != nil {
			handleLogicError(w, h.log, "CreateVideo thumbnail upload", err)
			return
		}
		thumbnailURL = result.URL
		thumbnailTinyURL = result.TinyURL
	}

	duration, err := strconv.Atoi(r.FormValue("duration"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "duration must be an integer")
		return
	}
	input := models.CreateVideoPost{
		Title:            r.FormValue("title"),
		VideoURL:         r.FormValue("videoURL"),
		ThumbnailURL:     thumbnailURL,
		ThumbnailTinyURL: thumbnailTinyURL,
		Duration:         duration,
		Tags:             r.Form["tags"],
	}
	if pl := r.FormValue("playlist"); pl != "" {
		input.Playlist = &pl
	}

	post, err := h.logic.CreateVideo(r.Context(), postID, input)
	if err != nil {
		handleLogicError(w, h.log, "CreateVideo", err)
		return
	}
	writeJSON(w, http.StatusCreated, post)
}

func (h *PostHandlerImpl) UpdateVideo(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	input := models.UpdateVideoPost{Tags: formTags(r)}

	if title := r.FormValue("title"); title != "" {
		input.Title = &title
	}
	if videoURL := r.FormValue("videoURL"); videoURL != "" {
		input.VideoURL = &videoURL
	}
	if durStr := r.FormValue("duration"); durStr != "" {
		if d, err := strconv.Atoi(durStr); err == nil {
			input.Duration = &d
		} else {
			writeError(w, http.StatusBadRequest, "duration must be an integer")
			return
		}
	}
	if pl := r.FormValue("playlist"); pl != "" {
		input.Playlist = &pl
	}

	post, err := h.logic.UpdateVideo(r.Context(), postID, input)
	if err != nil {
		handleLogicError(w, h.log, "UpdateVideo", err)
		return
	}
	writeJSON(w, http.StatusOK, post)
}

func (h *PostHandlerImpl) CreateLinkPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	postID := uuid.New().String()

	var thumbnailURL, thumbnailTinyURL *string
	if thumbFile, thumbHeader, err := r.FormFile("thumbnailFile"); err == nil {
		defer thumbFile.Close()
		result, err := h.uploads.UploadThumbnail(r.Context(), postID, "thumb",
			logic.FileInput{File: thumbFile, Header: thumbHeader})
		if err != nil {
			handleLogicError(w, h.log, "CreateLinkPost thumbnail upload", err)
			return
		}
		thumbnailURL = &result.URL
		thumbnailTinyURL = &result.TinyURL
	}

	input := models.CreateLinkPost{
		Title:            r.FormValue("title"),
		URL:              r.FormValue("url"),
		ThumbnailURL:     thumbnailURL,
		ThumbnailTinyURL: thumbnailTinyURL,
		Tags:             r.Form["tags"],
	}
	if desc := r.FormValue("description"); desc != "" {
		input.Description = &desc
	}
	if cat := r.FormValue("category"); cat != "" {
		input.Category = &cat
	}

	post, err := h.logic.CreateLinkPost(r.Context(), postID, input)
	if err != nil {
		handleLogicError(w, h.log, "CreateLinkPost", err)
		return
	}
	writeJSON(w, http.StatusCreated, post)
}

func (h *PostHandlerImpl) UpdateLinkPost(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	input := models.UpdateLinkPost{Tags: formTags(r)}
	if title := r.FormValue("title"); title != "" {
		input.Title = &title
	}
	if url := r.FormValue("url"); url != "" {
		input.URL = &url
	}
	if desc := r.FormValue("description"); desc != "" {
		input.Description = &desc
	}
	if cat := r.FormValue("category"); cat != "" {
		input.Category = &cat
	}

	post, err := h.logic.UpdateLinkPost(r.Context(), postID, input)
	if err != nil {
		handleLogicError(w, h.log, "UpdateLinkPost", err)
		return
	}
	writeJSON(w, http.StatusOK, post)
}

func (h *PostHandlerImpl) DeletePost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.logic.DeletePost(r.Context(), id); err != nil {
		handleLogicError(w, h.log, "DeletePost", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
