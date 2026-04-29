package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/pratikluitel/antipratik/common/logging"
	"github.com/pratikluitel/antipratik/common/requests"
	"github.com/pratikluitel/antipratik/components/files"
	"github.com/pratikluitel/antipratik/components/posts"
)

// postHandler implements PostHandler.
type postHandler struct {
	logic   posts.PostLogic
	uploads files.UploaderService
	log     logging.Logger
}

// NewPostHandler creates a new postHandler.
// uploads handles file storage for photo, music, video, and link post types.
func NewPostHandler(l posts.PostLogic, u files.UploaderService, log logging.Logger) posts.PostHandler {
	return &postHandler{logic: l, uploads: u, log: log}
}

// GetPosts handles GET /api/posts
// Query params: type (repeatable), tag (repeatable)
func (h *postHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := posts.FilterState{
		ActiveTypes: q["type"],
		ActiveTags:  q["tag"],
	}

	posts, err := h.logic.GetPosts(r.Context(), filter)
	if err != nil {
		handleLogicError(w, h.log, "GetPosts", err)
		return
	}

	requests.WriteJSON(w, http.StatusOK, posts)
}

// GetTags handles GET /api/tags — returns all tag names sorted alphabetically.
func (h *postHandler) GetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.logic.GetTags(r.Context())
	if err != nil {
		handleLogicError(w, h.log, "GetTags", err)
		return
	}
	requests.WriteJSON(w, http.StatusOK, tags)
}

// GetPost handles GET /api/posts/{slug}
func (h *postHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	post, err := h.logic.GetPost(r.Context(), slug)
	if err != nil {
		handleLogicError(w, h.log, "GetPost", err)
		return
	}
	if post == nil {
		requests.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	requests.WriteJSON(w, http.StatusOK, post)
}

// ── JSON-body write handlers (essay, short) ───────────────────────────────────

func (h *postHandler) CreateEssay(w http.ResponseWriter, r *http.Request) {
	var input posts.EssayPostInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	post, err := h.logic.CreateEssay(r.Context(), input)
	if err != nil {
		handleLogicError(w, h.log, "CreateEssay", err)
		return
	}
	requests.WriteJSON(w, http.StatusCreated, post)
}

func (h *postHandler) CreateShort(w http.ResponseWriter, r *http.Request) {
	var input posts.ShortPostInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	post, err := h.logic.CreateShort(r.Context(), input)
	if err != nil {
		handleLogicError(w, h.log, "CreateShort", err)
		return
	}
	requests.WriteJSON(w, http.StatusCreated, post)
}

func (h *postHandler) UpdateEssay(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input posts.UpdateEssayPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	post, err := h.logic.UpdateEssay(r.Context(), id, input)
	if err != nil {
		handleLogicError(w, h.log, "UpdateEssay", err)
		return
	}
	requests.WriteJSON(w, http.StatusOK, post)
}

func (h *postHandler) UpdateShort(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input posts.UpdateShortPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	post, err := h.logic.UpdateShort(r.Context(), id, input)
	if err != nil {
		handleLogicError(w, h.log, "UpdateShort", err)
		return
	}
	requests.WriteJSON(w, http.StatusOK, post)
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

func (h *postHandler) CreateMusic(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	postID := uuid.New().String()

	audioFile, audioHeader, err := r.FormFile("audioFile")
	if err != nil {
		requests.WriteError(w, http.StatusBadRequest, "audioFile is required")
		return
	}
	defer func() { _ = audioFile.Close() }()

	var albumArtInput *files.FileInput
	if artFile, artHeader, artErr := r.FormFile("albumArtFile"); artErr == nil {
		defer func() { _ = artFile.Close() }()
		albumArtInput = &files.FileInput{File: artFile, Header: artHeader}
	}

	audioInput := &files.FileInput{File: audioFile, Header: audioHeader}
	uploaded, err := h.uploads.UploadMusicFiles(r.Context(), postID, audioInput, albumArtInput)
	if err != nil {
		handleLogicError(w, h.log, "CreateMusic upload", err)
		return
	}

	duration, err := strconv.Atoi(r.FormValue("duration"))
	if err != nil {
		requests.WriteError(w, http.StatusBadRequest, "duration must be an integer")
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
	input := posts.MusicPostInput{
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
	requests.WriteJSON(w, http.StatusCreated, post)
}

func (h *postHandler) UpdateMusic(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	input := posts.UpdateMusicPost{Tags: formTags(r)}
	if title := r.FormValue("title"); title != "" {
		input.Title = &title
	}
	if albumStr := r.FormValue("album"); albumStr != "" {
		input.Album = &albumStr
	}

	if artFile, artHeader, artErr := r.FormFile("albumArtFile"); artErr == nil {
		defer func() { _ = artFile.Close() }()
		uploaded, uploadErr := h.uploads.UploadMusicFiles(r.Context(), postID, nil,
			&files.FileInput{File: artFile, Header: artHeader})
		if uploadErr != nil {
			handleLogicError(w, h.log, "UpdateMusic album art upload", uploadErr)
			return
		}
		input.AlbumArt = &uploaded.AlbumArtURL
		if uploaded.AlbumArtTinyURL != "" {
			v := uploaded.AlbumArtTinyURL
			input.AlbumArtTinyURL = &v
		}
		if uploaded.AlbumArtSmallURL != "" {
			v := uploaded.AlbumArtSmallURL
			input.AlbumArtSmallURL = &v
		}
		if uploaded.AlbumArtMedURL != "" {
			v := uploaded.AlbumArtMedURL
			input.AlbumArtMedURL = &v
		}
		if uploaded.AlbumArtLargeURL != "" {
			v := uploaded.AlbumArtLargeURL
			input.AlbumArtLargeURL = &v
		}
	}

	post, err := h.logic.UpdateMusic(r.Context(), postID, input)
	if err != nil {
		handleLogicError(w, h.log, "UpdateMusic", err)
		return
	}
	requests.WriteJSON(w, http.StatusOK, post)
}

func (h *postHandler) CreatePhoto(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	fileHeaders := r.MultipartForm.File["images[]"]
	if len(fileHeaders) == 0 {
		requests.WriteError(w, http.StatusBadRequest, "at least one image file is required")
		return
	}

	postID := uuid.New().String()

	fileInputs := make([]files.FileInput, 0, len(fileHeaders))
	for _, fh := range fileHeaders {
		f, err := fh.Open()
		if err != nil {
			requests.WriteError(w, http.StatusBadRequest, "could not read uploaded image")
			return
		}
		defer func() { _ = f.Close() }()
		fileInputs = append(fileInputs, files.FileInput{File: f, Header: fh})
	}

	uploadResults, err := h.uploads.UploadPhotoFiles(r.Context(), postID, fileInputs)
	if err != nil {
		handleLogicError(w, h.log, "CreatePhoto upload", err)
		return
	}

	alts := r.Form["alt[]"]
	captions := r.Form["caption[]"]
	images := make([]posts.PhotoImage, len(uploadResults))
	for i, u := range uploadResults {
		alt := ""
		if i < len(alts) {
			alt = alts[i]
		}
		tiny, small, med, large := u.ThumbnailTinyURL, u.ThumbnailSmallURL, u.ThumbnailMedURL, u.ThumbnailLargeURL
		images[i] = posts.PhotoImage{
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

	input := posts.PhotoPostInput{Images: images, Tags: r.Form["tags[]"]}
	if loc := r.FormValue("location"); loc != "" {
		input.Location = &loc
	}

	post, err := h.logic.CreatePhoto(r.Context(), postID, input)
	if err != nil {
		handleLogicError(w, h.log, "CreatePhoto", err)
		return
	}
	requests.WriteJSON(w, http.StatusCreated, post)
}

func (h *postHandler) UpdatePhoto(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	input := posts.PhotoPostInput{Tags: formTags(r)}
	if loc := r.FormValue("location"); loc != "" {
		input.Location = &loc
	}

	post, err := h.logic.UpdatePhoto(r.Context(), postID, input)
	if err != nil {
		handleLogicError(w, h.log, "UpdatePhoto", err)
		return
	}
	requests.WriteJSON(w, http.StatusOK, post)
}

func (h *postHandler) CreateVideo(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	postID := uuid.New().String()

	videoFile, videoHeader, err := r.FormFile("videoFile")
	if err != nil {
		requests.WriteError(w, http.StatusBadRequest, "videoFile is required")
		return
	}
	defer func() { _ = videoFile.Close() }()

	uploaded, err := h.uploads.UploadVideoFile(r.Context(), postID,
		files.FileInput{File: videoFile, Header: videoHeader})
	if err != nil {
		handleLogicError(w, h.log, "CreateVideo video upload", err)
		return
	}

	input := posts.VideoPostInput{
		Title:    r.FormValue("title"),
		VideoURL: uploaded.VideoURL,
		Tags:     formTags(r),
	}
	if desc := r.FormValue("description"); desc != "" {
		input.Description = &desc
	}

	if thumbFile, thumbHeader, thumbErr := r.FormFile("thumbnailFile"); thumbErr == nil {
		defer func() { _ = thumbFile.Close() }()
		result, uploadErr := h.uploads.UploadThumbnail(r.Context(), postID, "thumb",
			files.FileInput{File: thumbFile, Header: thumbHeader})
		if uploadErr != nil {
			handleLogicError(w, h.log, "CreateVideo thumbnail upload", uploadErr)
			return
		}
		input.ThumbnailURL = &result.URL
		input.ThumbnailTinyURL = &result.TinyURL
		input.ThumbnailSmallURL = &result.SmallURL
		input.ThumbnailMedURL = &result.MedURL
		input.ThumbnailLargeURL = &result.LargeURL
	}

	post, err := h.logic.CreateVideo(r.Context(), postID, input)
	if err != nil {
		handleLogicError(w, h.log, "CreateVideo", err)
		return
	}
	requests.WriteJSON(w, http.StatusCreated, map[string]string{"id": post.ID})
}

func (h *postHandler) UpdateVideo(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	// videoFile is silently ignored — editing the uploaded video file is not supported.

	input := posts.UpdateVideoPost{Tags: formTags(r)}

	if title := r.FormValue("title"); title != "" {
		input.Title = &title
	}
	if desc := r.FormValue("description"); desc != "" {
		input.Description = &desc
	}

	if thumbFile, thumbHeader, thumbErr := r.FormFile("thumbnailFile"); thumbErr == nil {
		defer func() { _ = thumbFile.Close() }()
		result, uploadErr := h.uploads.UploadThumbnail(r.Context(), postID, "thumb",
			files.FileInput{File: thumbFile, Header: thumbHeader})
		if uploadErr != nil {
			handleLogicError(w, h.log, "UpdateVideo thumbnail upload", uploadErr)
			return
		}
		input.ThumbnailURL = &result.URL
		input.ThumbnailTinyURL = &result.TinyURL
		input.ThumbnailSmallURL = &result.SmallURL
		input.ThumbnailMedURL = &result.MedURL
		input.ThumbnailLargeURL = &result.LargeURL
	}

	post, err := h.logic.UpdateVideo(r.Context(), postID, input)
	if err != nil {
		handleLogicError(w, h.log, "UpdateVideo", err)
		return
	}
	requests.WriteJSON(w, http.StatusOK, post)
}

func (h *postHandler) CreateLinkPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	postID := uuid.New().String()

	var thumbnailURL, thumbnailTinyURL, thumbnailSmallURL, thumbnailMedURL, thumbnailLargeURL *string
	if thumbFile, thumbHeader, err := r.FormFile("thumbnailFile"); err == nil {
		defer func() { _ = thumbFile.Close() }()
		result, uploadErr := h.uploads.UploadThumbnail(r.Context(), postID, "thumb",
			files.FileInput{File: thumbFile, Header: thumbHeader})
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

	input := posts.LinkPostInput{
		Title:             r.FormValue("title"),
		URL:               r.FormValue("url"),
		ThumbnailURL:      thumbnailURL,
		ThumbnailTinyURL:  thumbnailTinyURL,
		ThumbnailSmallURL: thumbnailSmallURL,
		ThumbnailMedURL:   thumbnailMedURL,
		ThumbnailLargeURL: thumbnailLargeURL,
		Tags:              formTags(r),
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
	requests.WriteJSON(w, http.StatusCreated, post)
}

func (h *postHandler) UpdateLinkPost(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	input := posts.UpdateLinkPost{Tags: formTags(r)}
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

	if thumbFile, thumbHeader, thumbErr := r.FormFile("thumbnailFile"); thumbErr == nil {
		defer func() { _ = thumbFile.Close() }()
		result, uploadErr := h.uploads.UploadThumbnail(r.Context(), postID, "thumb",
			files.FileInput{File: thumbFile, Header: thumbHeader})
		if uploadErr != nil {
			handleLogicError(w, h.log, "UpdateLinkPost thumbnail upload", uploadErr)
			return
		}
		input.ThumbnailURL = &result.URL
		input.ThumbnailTinyURL = &result.TinyURL
		input.ThumbnailSmallURL = &result.SmallURL
		input.ThumbnailMedURL = &result.MedURL
		input.ThumbnailLargeURL = &result.LargeURL
	}

	post, err := h.logic.UpdateLinkPost(r.Context(), postID, input)
	if err != nil {
		handleLogicError(w, h.log, "UpdateLinkPost", err)
		return
	}
	requests.WriteJSON(w, http.StatusOK, post)
}

func (h *postHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.logic.DeletePost(r.Context(), id); err != nil {
		handleLogicError(w, h.log, "DeletePost", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── Individual photo image endpoints ─────────────────────────────────────────

func (h *postHandler) AddPhotoImage(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "request body too large or not multipart/form-data")
		return
	}

	fhs := r.MultipartForm.File["image"]
	if len(fhs) == 0 {
		requests.WriteError(w, http.StatusBadRequest, "image file is required")
		return
	}
	f, err := fhs[0].Open()
	if err != nil {
		requests.WriteError(w, http.StatusBadRequest, "could not read uploaded image")
		return
	}
	defer func() { _ = f.Close() }()

	uploadResults, err := h.uploads.UploadPhotoFiles(r.Context(), postID, []files.FileInput{{File: f, Header: fhs[0]}})
	if err != nil {
		handleLogicError(w, h.log, "AddPhotoImage upload", err)
		return
	}
	u := uploadResults[0]
	tiny, small, med, large := u.ThumbnailTinyURL, u.ThumbnailSmallURL, u.ThumbnailMedURL, u.ThumbnailLargeURL
	image := posts.PhotoImage{
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
	requests.WriteJSON(w, http.StatusCreated, result)
}

func (h *postHandler) GetPhotoImage(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	imageIDStr := r.PathValue("imageID")

	img, err := h.logic.GetPhotoImage(r.Context(), postID, imageIDStr)
	if err != nil {
		handleLogicError(w, h.log, "GetPhotoImage", err)
		return
	}
	if img == nil {
		requests.WriteError(w, http.StatusNotFound, "image not found")
		return
	}
	requests.WriteJSON(w, http.StatusOK, img)
}

func (h *postHandler) UpdatePhotoImage(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	imageIDStr := r.PathValue("imageID")

	var body struct {
		Caption *string `json:"caption"`
		Alt     *string `json:"alt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		requests.WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	img, err := h.logic.UpdatePhotoImage(r.Context(), postID, imageIDStr, posts.UpdatePhotoImage{
		Caption: body.Caption,
		Alt:     body.Alt,
	})
	if err != nil {
		handleLogicError(w, h.log, "UpdatePhotoImage", err)
		return
	}
	if img == nil {
		requests.WriteError(w, http.StatusNotFound, "image not found")
		return
	}
	requests.WriteJSON(w, http.StatusOK, img)
}

func (h *postHandler) DeletePhotoImage(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	imageIDStr := r.PathValue("imageID")

	notFound, err := h.logic.DeletePhotoImage(r.Context(), postID, imageIDStr)
	if err != nil {
		handleLogicError(w, h.log, "DeletePhotoImage", err)
		return
	}
	if notFound {
		requests.WriteError(w, http.StatusNotFound, "image not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
