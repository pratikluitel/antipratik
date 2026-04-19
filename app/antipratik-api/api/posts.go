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

// GetTags handles GET /api/tags — returns all tag names sorted alphabetically.
func (h *PostHandlerImpl) GetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.logic.GetTags(r.Context())
	if err != nil {
		handleLogicError(w, h.log, "GetTags", err)
		return
	}
	writeJSON(w, http.StatusOK, tags)
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
	var input models.EssayPostInput
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
	var input models.ShortPostInput
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
	defer func() { _ = audioFile.Close() }()

	var albumArtInput *models.FileInput
	if artFile, artHeader, artErr := r.FormFile("albumArtFile"); artErr == nil {
		defer func() { _ = artFile.Close() }()
		albumArtInput = &models.FileInput{File: artFile, Header: artHeader}
	}

	audioInput := &models.FileInput{File: audioFile, Header: audioHeader}
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

	var albumArtTiny, albumArtSmall, albumArtMed, albumArtLarge *string
	if uploaded.AlbumArtTinyURL != "" {
		albumArtTiny = &uploaded.AlbumArtTinyURL
	}
	if uploaded.AlbumArtSmallURL != "" {
		albumArtSmall = &uploaded.AlbumArtSmallURL
	}
	if uploaded.AlbumArtMedURL != "" {
		albumArtMed = &uploaded.AlbumArtMedURL
	}
	if uploaded.AlbumArtLargeURL != "" {
		albumArtLarge = &uploaded.AlbumArtLargeURL
	}
	input := models.MusicPostInput{
		Title:            r.FormValue("title"),
		AudioURL:         uploaded.AudioURL,
		AlbumArt:         uploaded.AlbumArtURL,
		AlbumArtTinyURL:  albumArtTiny,
		AlbumArtSmallURL: albumArtSmall,
		AlbumArtMedURL:   albumArtMed,
		AlbumArtLargeURL: albumArtLarge,
		Duration:         duration,
		Tags:             formTags(r),
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

	fileInputs := make([]models.FileInput, 0, len(fileHeaders))
	for _, fh := range fileHeaders {
		f, err := fh.Open()
		if err != nil {
			writeError(w, http.StatusBadRequest, "could not read uploaded image")
			return
		}
		defer func() { _ = f.Close() }()
		fileInputs = append(fileInputs, models.FileInput{File: f, Header: fh})
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
			URL:               u.OriginalURL,
			Alt:               alt,
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

	input := models.PhotoPostInput{Images: images, Tags: r.Form["tags[]"]}
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

	input := models.PhotoPostInput{Tags: formTags(r)}
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

	var thumbnailURL string
	var thumbnailTinyURL, thumbnailSmallURL, thumbnailMedURL, thumbnailLargeURL *string
	if thumbFile, thumbHeader, err := r.FormFile("thumbnailFile"); err == nil {
		defer func() { _ = thumbFile.Close() }()
		result, uploadErr := h.uploads.UploadThumbnail(r.Context(), postID, "thumb",
			models.FileInput{File: thumbFile, Header: thumbHeader})
		if uploadErr != nil {
			handleLogicError(w, h.log, "CreateVideo thumbnail upload", uploadErr)
			return
		}
		thumbnailURL = result.URL
		thumbnailTinyURL = &result.TinyURL
		thumbnailSmallURL = &result.SmallURL
		thumbnailMedURL = &result.MedURL
		thumbnailLargeURL = &result.LargeURL
	}

	duration, err := strconv.Atoi(r.FormValue("duration"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "duration must be an integer")
		return
	}
	input := models.VideoPostInput{
		Title:             r.FormValue("title"),
		VideoURL:          r.FormValue("videoURL"),
		ThumbnailURL:      thumbnailURL,
		ThumbnailTinyURL:  thumbnailTinyURL,
		ThumbnailSmallURL: thumbnailSmallURL,
		ThumbnailMedURL:   thumbnailMedURL,
		ThumbnailLargeURL: thumbnailLargeURL,
		Duration:          duration,
		Tags:              formTags(r),
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

	var thumbnailURL, thumbnailTinyURL, thumbnailSmallURL, thumbnailMedURL, thumbnailLargeURL *string
	if thumbFile, thumbHeader, err := r.FormFile("thumbnailFile"); err == nil {
		defer func() { _ = thumbFile.Close() }()
		result, uploadErr := h.uploads.UploadThumbnail(r.Context(), postID, "thumb",
			models.FileInput{File: thumbFile, Header: thumbHeader})
		if uploadErr != nil {
			handleLogicError(w, h.log, "CreateLinkPost thumbnail upload", uploadErr)
			return
		}
		thumbnailURL = &result.URL
		thumbnailTinyURL = &result.TinyURL
		thumbnailSmallURL = &result.SmallURL
		thumbnailMedURL = &result.MedURL
		thumbnailLargeURL = &result.LargeURL
	}

	input := models.LinkPostInput{
		Title:             r.FormValue("title"),
		URL:               r.FormValue("url"),
		ThumbnailURL:      thumbnailURL,
		ThumbnailTinyURL:  thumbnailTinyURL,
		ThumbnailSmallURL: thumbnailSmallURL,
		ThumbnailMedURL:   thumbnailMedURL,
		ThumbnailLargeURL: thumbnailLargeURL,
		Tags:             formTags(r),
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

// ── Individual photo image endpoints ─────────────────────────────────────────

func (h *PostHandlerImpl) AddPhotoImage(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	fhs := r.MultipartForm.File["image"]
	if len(fhs) == 0 {
		writeError(w, http.StatusBadRequest, "image file is required")
		return
	}
	f, err := fhs[0].Open()
	if err != nil {
		writeError(w, http.StatusBadRequest, "could not read uploaded image")
		return
	}
	defer func() { _ = f.Close() }()

	uploadResults, err := h.uploads.UploadPhotoFiles(r.Context(), postID, []models.FileInput{{File: f, Header: fhs[0]}})
	if err != nil {
		handleLogicError(w, h.log, "AddPhotoImage upload", err)
		return
	}
	u := uploadResults[0]
	tiny, small, med, large := u.ThumbnailTinyURL, u.ThumbnailSmallURL, u.ThumbnailMedURL, u.ThumbnailLargeURL
	image := models.PhotoImage{
		URL:               u.OriginalURL,
		Alt:               r.FormValue("alt"),
		ThumbnailTinyURL:  &tiny,
		ThumbnailSmallURL: &small,
		ThumbnailMedURL:   &med,
		ThumbnailLargeURL: &large,
	}
	if c := r.FormValue("caption"); c != "" {
		image.Caption = &c
	}

	result, err := h.logic.AddPhotoImage(r.Context(), postID, image)
	if err != nil {
		handleLogicError(w, h.log, "AddPhotoImage", err)
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func (h *PostHandlerImpl) GetPhotoImage(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	imageIDStr := r.PathValue("imageID")

	img, err := h.logic.GetPhotoImage(r.Context(), postID, imageIDStr)
	if err != nil {
		handleLogicError(w, h.log, "GetPhotoImage", err)
		return
	}
	if img == nil {
		writeError(w, http.StatusNotFound, "image not found")
		return
	}
	writeJSON(w, http.StatusOK, img)
}

func (h *PostHandlerImpl) UpdatePhotoImage(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	imageIDStr := r.PathValue("imageID")

	var body struct {
		Caption *string `json:"caption"`
		Alt     *string `json:"alt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	img, err := h.logic.UpdatePhotoImage(r.Context(), postID, imageIDStr, models.UpdatePhotoImage{
		Caption: body.Caption,
		Alt:     body.Alt,
	})
	if err != nil {
		handleLogicError(w, h.log, "UpdatePhotoImage", err)
		return
	}
	if img == nil {
		writeError(w, http.StatusNotFound, "image not found")
		return
	}
	writeJSON(w, http.StatusOK, img)
}

func (h *PostHandlerImpl) DeletePhotoImage(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	imageIDStr := r.PathValue("imageID")

	notFound, err := h.logic.DeletePhotoImage(r.Context(), postID, imageIDStr)
	if err != nil {
		handleLogicError(w, h.log, "DeletePhotoImage", err)
		return
	}
	if notFound {
		writeError(w, http.StatusNotFound, "image not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
