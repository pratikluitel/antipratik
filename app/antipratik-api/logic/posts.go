package logic

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pratikluitel/antipratik/common/logging"
	"github.com/pratikluitel/antipratik/models"
	"github.com/pratikluitel/antipratik/store"
)

// PostService implements PostLogic.
type PostService struct {
	store store.PostStore
	files store.FileStore
	log   logging.Logger
}

// NewPostService creates a new PostService backed by the given store and file store.
func NewPostService(s store.PostStore, files store.FileStore, log logging.Logger) *PostService {
	return &PostService{store: s, files: files, log: log}
}

var validTypes = map[string]bool{
	"essay": true, "short": true, "music": true,
	"photo": true, "video": true, "link": true,
}

// GetPosts validates the filter and delegates to the store.
func (s *PostService) GetPosts(ctx context.Context, filter models.FilterState) ([]models.Post, error) {
	types := make([]string, 0, len(filter.ActiveTypes))
	for _, t := range filter.ActiveTypes {
		if validTypes[t] {
			types = append(types, t)
		}
	}

	posts, err := s.store.GetPosts(ctx, types, filter.ActiveTags)
	if err != nil {
		return nil, fmt.Errorf("PostService.GetPosts: %w", err)
	}
	return posts, nil
}

// GetPost validates the slug and delegates to the store.
// Returns nil if the post does not exist.
func (s *PostService) GetPost(ctx context.Context, slug string) (*models.EssayPost, error) {
	if slug == "" {
		return nil, nil
	}
	post, err := s.store.GetPostBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("PostService.GetPost: %w", err)
	}
	return post, nil
}

// ── Write methods ─────────────────────────────────────────────────────────────

func newID() string  { return uuid.New().String() }
func nowUTC() string { return time.Now().UTC().Format(time.RFC3339) }

// extractDomain parses rawURL and returns the hostname with www. stripped.
// Returns a ValidationError if the URL is malformed or missing scheme/host.
func extractDomain(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", validationErr("url must be a valid absolute URL (e.g. https://example.com/path)")
	}
	host := strings.TrimPrefix(u.Hostname(), "www.")
	return host, nil
}

// computeReadingTime returns ceil(wordCount / 200), minimum 1.
func computeReadingTime(body string) int {
	words := len(strings.Fields(body))
	if words == 0 {
		return 1
	}
	return (words + 199) / 200
}

func (s *PostService) CreateEssay(ctx context.Context, input models.CreateEssayPost) (models.EssayPost, error) {
	if err := requireNonEmpty("title", input.Title); err != nil {
		return models.EssayPost{}, err
	}
	if err := requireNonEmpty("slug", input.Slug); err != nil {
		return models.EssayPost{}, err
	}
	if err := requireNonEmpty("body", input.Body); err != nil {
		return models.EssayPost{}, err
	}
	input.ReadingTimeMinutes = computeReadingTime(input.Body)

	id, createdAt := newID(), nowUTC()
	if err := s.store.CreatePost(ctx, "essay", id, createdAt); err != nil {
		return models.EssayPost{}, fmt.Errorf("PostService.CreateEssay: %w", err)
	}
	if err := s.store.CreateEssayData(ctx, id, input); err != nil {
		return models.EssayPost{}, fmt.Errorf("PostService.CreateEssay data: %w", err)
	}
	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}
	return models.EssayPost{
		ID: id, Type: "essay", CreatedAt: createdAt, Tags: tags,
		Title: input.Title, Slug: input.Slug, Excerpt: input.Excerpt,
		Body: input.Body, ReadingTimeMinutes: input.ReadingTimeMinutes,
	}, nil
}

func (s *PostService) CreateShort(ctx context.Context, input models.CreateShortPost) (models.ShortPost, error) {
	if err := requireNonEmpty("body", input.Body); err != nil {
		return models.ShortPost{}, err
	}

	id, createdAt := newID(), nowUTC()
	if err := s.store.CreatePost(ctx, "short", id, createdAt); err != nil {
		return models.ShortPost{}, fmt.Errorf("PostService.CreateShort: %w", err)
	}
	if err := s.store.CreateShortData(ctx, id, input); err != nil {
		return models.ShortPost{}, fmt.Errorf("PostService.CreateShort data: %w", err)
	}
	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}
	return models.ShortPost{ID: id, Type: "short", CreatedAt: createdAt, Tags: tags, Body: input.Body}, nil
}

func (s *PostService) CreateMusic(ctx context.Context, id string, input models.CreateMusicPost) (models.MusicPost, error) {
	if err := requireNonEmpty("title", input.Title); err != nil {
		return models.MusicPost{}, err
	}
	if err := requireNonEmpty("audioURL", input.AudioURL); err != nil {
		return models.MusicPost{}, err
	}
	if err := requirePositive("duration", input.Duration); err != nil {
		return models.MusicPost{}, err
	}

	createdAt := nowUTC()
	if err := s.store.CreatePost(ctx, "music", id, createdAt); err != nil {
		return models.MusicPost{}, fmt.Errorf("PostService.CreateMusic: %w", err)
	}
	if err := s.store.CreateMusicData(ctx, id, input); err != nil {
		return models.MusicPost{}, fmt.Errorf("PostService.CreateMusic data: %w", err)
	}

	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}
	var albumArtTiny *string
	if input.AlbumArtTinyURL != "" {
		albumArtTiny = &input.AlbumArtTinyURL
	}
	return models.MusicPost{
		ID: id, Type: "music", CreatedAt: createdAt, Tags: tags,
		Title: input.Title, AlbumArt: input.AlbumArt, AlbumArtTinyURL: albumArtTiny,
		AudioURL: input.AudioURL, Duration: input.Duration, Album: input.Album,
	}, nil
}

func (s *PostService) CreatePhoto(ctx context.Context, preID string, input models.CreatePhotoPost) (models.PhotoPost, error) {
	if len(input.Images) == 0 {
		return models.PhotoPost{}, validationErr("images cannot be empty")
	}
	for i, img := range input.Images {
		if err := requireNonEmpty(fmt.Sprintf("images[%d].url", i), img.URL); err != nil {
			return models.PhotoPost{}, err
		}
	}

	id := preID
	if id == "" {
		id = newID()
	}
	createdAt := nowUTC()
	if err := s.store.CreatePost(ctx, "photo", id, createdAt); err != nil {
		return models.PhotoPost{}, fmt.Errorf("PostService.CreatePhoto: %w", err)
	}
	if err := s.store.CreatePhotoData(ctx, id, input); err != nil {
		return models.PhotoPost{}, fmt.Errorf("PostService.CreatePhoto data: %w", err)
	}
	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}
	return models.PhotoPost{
		ID: id, Type: "photo", CreatedAt: createdAt, Tags: tags,
		Images: input.Images, Location: input.Location,
	}, nil
}

func (s *PostService) CreateVideo(ctx context.Context, preID string, input models.CreateVideoPost) (models.VideoPost, error) {
	if err := requireNonEmpty("title", input.Title); err != nil {
		return models.VideoPost{}, err
	}
	if err := requireNonEmpty("videoURL", input.VideoURL); err != nil {
		return models.VideoPost{}, err
	}
	if err := requirePositive("duration", input.Duration); err != nil {
		return models.VideoPost{}, err
	}

	id := preID
	if id == "" {
		id = newID()
	}
	createdAt := nowUTC()
	if err := s.store.CreatePost(ctx, "video", id, createdAt); err != nil {
		return models.VideoPost{}, fmt.Errorf("PostService.CreateVideo: %w", err)
	}
	if err := s.store.CreateVideoData(ctx, id, input); err != nil {
		return models.VideoPost{}, fmt.Errorf("PostService.CreateVideo data: %w", err)
	}
	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}
	var thumbTiny *string
	if input.ThumbnailTinyURL != "" {
		thumbTiny = &input.ThumbnailTinyURL
	}
	return models.VideoPost{
		ID: id, Type: "video", CreatedAt: createdAt, Tags: tags,
		Title: input.Title, ThumbnailURL: input.ThumbnailURL, ThumbnailTinyURL: thumbTiny,
		VideoURL: input.VideoURL, Duration: input.Duration, Playlist: input.Playlist,
	}, nil
}

func (s *PostService) CreateLinkPost(ctx context.Context, preID string, input models.CreateLinkPost) (models.LinkPost, error) {
	if err := requireNonEmpty("title", input.Title); err != nil {
		return models.LinkPost{}, err
	}
	if err := requireNonEmpty("url", input.URL); err != nil {
		return models.LinkPost{}, err
	}
	domain, err := extractDomain(input.URL)
	if err != nil {
		return models.LinkPost{}, err
	}
	input.Domain = domain

	id := preID
	if id == "" {
		id = newID()
	}
	createdAt := nowUTC()
	if err := s.store.CreatePost(ctx, "link", id, createdAt); err != nil {
		return models.LinkPost{}, fmt.Errorf("PostService.CreateLinkPost: %w", err)
	}
	if err := s.store.CreateLinkPostData(ctx, id, input); err != nil {
		return models.LinkPost{}, fmt.Errorf("PostService.CreateLinkPost data: %w", err)
	}
	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}
	return models.LinkPost{
		ID: id, Type: "link", CreatedAt: createdAt, Tags: tags,
		Title: input.Title, URL: input.URL, Domain: input.Domain,
		Description: input.Description, ThumbnailURL: input.ThumbnailURL,
		ThumbnailTinyURL: input.ThumbnailTinyURL, Category: input.Category,
	}, nil
}

func (s *PostService) UpdateEssay(ctx context.Context, id string, input models.UpdateEssayPost) (models.EssayPost, error) {
	if err := requireNonEmpty("id", id); err != nil {
		return models.EssayPost{}, err
	}
	post, err := s.store.GetPostByID(ctx, id)
	if err != nil {
		return models.EssayPost{}, fmt.Errorf("PostService.UpdateEssay: %w", err)
	}
	cur, ok := post.(models.EssayPost)
	if !ok {
		return models.EssayPost{}, validationErr("post is not an essay")
	}

	merged := models.CreateEssayPost{Title: cur.Title, Slug: cur.Slug, Excerpt: cur.Excerpt, Body: cur.Body, Tags: cur.Tags}
	if input.Title != nil {
		merged.Title = *input.Title
	}
	if input.Slug != nil {
		merged.Slug = *input.Slug
	}
	if input.Excerpt != nil {
		merged.Excerpt = *input.Excerpt
	}
	if input.Body != nil {
		merged.Body = *input.Body
	}
	if input.Tags != nil {
		merged.Tags = input.Tags
	}

	if err := requireNonEmpty("title", merged.Title); err != nil {
		return models.EssayPost{}, err
	}
	if err := requireNonEmpty("slug", merged.Slug); err != nil {
		return models.EssayPost{}, err
	}
	if err := requireNonEmpty("body", merged.Body); err != nil {
		return models.EssayPost{}, err
	}
	merged.ReadingTimeMinutes = computeReadingTime(merged.Body)

	if err := s.store.UpdateEssay(ctx, id, merged); err != nil {
		return models.EssayPost{}, fmt.Errorf("PostService.UpdateEssay: %w", err)
	}
	tags := merged.Tags
	if tags == nil {
		tags = []string{}
	}
	return models.EssayPost{
		ID: id, Type: "essay", CreatedAt: cur.CreatedAt, Tags: tags,
		Title: merged.Title, Slug: merged.Slug, Excerpt: merged.Excerpt,
		Body: merged.Body, ReadingTimeMinutes: merged.ReadingTimeMinutes,
	}, nil
}

func (s *PostService) UpdateShort(ctx context.Context, id string, input models.UpdateShortPost) (models.ShortPost, error) {
	if err := requireNonEmpty("id", id); err != nil {
		return models.ShortPost{}, err
	}
	post, err := s.store.GetPostByID(ctx, id)
	if err != nil {
		return models.ShortPost{}, fmt.Errorf("PostService.UpdateShort: %w", err)
	}
	cur, ok := post.(models.ShortPost)
	if !ok {
		return models.ShortPost{}, validationErr("post is not a short post")
	}

	merged := models.CreateShortPost{Body: cur.Body, Tags: cur.Tags}
	if input.Body != nil {
		merged.Body = *input.Body
	}
	if input.Tags != nil {
		merged.Tags = input.Tags
	}

	if err := requireNonEmpty("body", merged.Body); err != nil {
		return models.ShortPost{}, err
	}

	if err := s.store.UpdateShort(ctx, id, merged); err != nil {
		return models.ShortPost{}, fmt.Errorf("PostService.UpdateShort: %w", err)
	}
	tags := merged.Tags
	if tags == nil {
		tags = []string{}
	}
	return models.ShortPost{ID: id, Type: "short", CreatedAt: cur.CreatedAt, Tags: tags, Body: merged.Body}, nil
}

func (s *PostService) UpdateMusic(ctx context.Context, id string, input models.UpdateMusicPost) (models.MusicPost, error) {
	if err := requireNonEmpty("id", id); err != nil {
		return models.MusicPost{}, err
	}
	post, err := s.store.GetPostByID(ctx, id)
	if err != nil {
		return models.MusicPost{}, fmt.Errorf("PostService.UpdateMusic: %w", err)
	}
	cur, ok := post.(models.MusicPost)
	if !ok {
		return models.MusicPost{}, validationErr("post is not a music post")
	}

	curAlbumArtTiny := ""
	if cur.AlbumArtTinyURL != nil {
		curAlbumArtTiny = *cur.AlbumArtTinyURL
	}
	merged := models.CreateMusicPost{
		Title: cur.Title, AudioURL: cur.AudioURL, AlbumArt: cur.AlbumArt,
		AlbumArtTinyURL: curAlbumArtTiny, Duration: cur.Duration, Album: cur.Album, Tags: cur.Tags,
	}
	if input.Title != nil {
		merged.Title = *input.Title
	}
	if input.AudioURL != nil {
		merged.AudioURL = *input.AudioURL
	}
	if input.AlbumArt != nil {
		merged.AlbumArt = *input.AlbumArt
	}
	if input.Duration != nil {
		merged.Duration = *input.Duration
	}
	if input.Album != nil {
		merged.Album = input.Album
	}
	if input.Tags != nil {
		merged.Tags = input.Tags
	}

	if err := requireNonEmpty("title", merged.Title); err != nil {
		return models.MusicPost{}, err
	}
	if err := requireNonEmpty("audioURL", merged.AudioURL); err != nil {
		return models.MusicPost{}, err
	}
	if err := requirePositive("duration", merged.Duration); err != nil {
		return models.MusicPost{}, err
	}

	if err := s.store.UpdateMusic(ctx, id, merged); err != nil {
		return models.MusicPost{}, fmt.Errorf("PostService.UpdateMusic: %w", err)
	}
	tags := merged.Tags
	if tags == nil {
		tags = []string{}
	}
	var mergedAlbumArtTiny *string
	if merged.AlbumArtTinyURL != "" {
		mergedAlbumArtTiny = &merged.AlbumArtTinyURL
	}
	return models.MusicPost{
		ID: id, Type: "music", CreatedAt: cur.CreatedAt, Tags: tags,
		Title: merged.Title, AlbumArt: merged.AlbumArt, AlbumArtTinyURL: mergedAlbumArtTiny,
		AudioURL: merged.AudioURL, Duration: merged.Duration, Album: merged.Album,
	}, nil
}

func (s *PostService) UpdatePhoto(ctx context.Context, id string, input models.UpdatePhotoPost) (models.PhotoPost, error) {
	if err := requireNonEmpty("id", id); err != nil {
		return models.PhotoPost{}, err
	}
	post, err := s.store.GetPostByID(ctx, id)
	if err != nil {
		return models.PhotoPost{}, fmt.Errorf("PostService.UpdatePhoto: %w", err)
	}
	cur, ok := post.(models.PhotoPost)
	if !ok {
		return models.PhotoPost{}, validationErr("post is not a photo post")
	}

	merged := models.CreatePhotoPost{Images: cur.Images, Location: cur.Location, Tags: cur.Tags}
	if input.Location != nil {
		merged.Location = input.Location
	}
	if input.Tags != nil {
		merged.Tags = input.Tags
	}

	if err := s.store.UpdatePhoto(ctx, id, merged); err != nil {
		return models.PhotoPost{}, fmt.Errorf("PostService.UpdatePhoto: %w", err)
	}
	tags := merged.Tags
	if tags == nil {
		tags = []string{}
	}
	return models.PhotoPost{
		ID: id, Type: "photo", CreatedAt: cur.CreatedAt, Tags: tags,
		Images: merged.Images, Location: merged.Location,
	}, nil
}

func (s *PostService) UpdateVideo(ctx context.Context, id string, input models.UpdateVideoPost) (models.VideoPost, error) {
	if err := requireNonEmpty("id", id); err != nil {
		return models.VideoPost{}, err
	}
	post, err := s.store.GetPostByID(ctx, id)
	if err != nil {
		return models.VideoPost{}, fmt.Errorf("PostService.UpdateVideo: %w", err)
	}
	cur, ok := post.(models.VideoPost)
	if !ok {
		return models.VideoPost{}, validationErr("post is not a video post")
	}

	curThumbTiny := ""
	if cur.ThumbnailTinyURL != nil {
		curThumbTiny = *cur.ThumbnailTinyURL
	}
	merged := models.CreateVideoPost{
		Title: cur.Title, ThumbnailURL: cur.ThumbnailURL, ThumbnailTinyURL: curThumbTiny,
		VideoURL: cur.VideoURL, Duration: cur.Duration, Playlist: cur.Playlist, Tags: cur.Tags,
	}
	if input.Title != nil {
		merged.Title = *input.Title
	}
	if input.VideoURL != nil {
		merged.VideoURL = *input.VideoURL
	}
	if input.Duration != nil {
		merged.Duration = *input.Duration
	}
	if input.Playlist != nil {
		merged.Playlist = input.Playlist
	}
	if input.Tags != nil {
		merged.Tags = input.Tags
	}

	if err := requireNonEmpty("title", merged.Title); err != nil {
		return models.VideoPost{}, err
	}
	if err := requireNonEmpty("videoURL", merged.VideoURL); err != nil {
		return models.VideoPost{}, err
	}
	if err := requirePositive("duration", merged.Duration); err != nil {
		return models.VideoPost{}, err
	}

	if err := s.store.UpdateVideo(ctx, id, merged); err != nil {
		return models.VideoPost{}, fmt.Errorf("PostService.UpdateVideo: %w", err)
	}
	tags := merged.Tags
	if tags == nil {
		tags = []string{}
	}
	var mergedThumbTiny *string
	if merged.ThumbnailTinyURL != "" {
		mergedThumbTiny = &merged.ThumbnailTinyURL
	}
	return models.VideoPost{
		ID: id, Type: "video", CreatedAt: cur.CreatedAt, Tags: tags,
		Title: merged.Title, ThumbnailURL: merged.ThumbnailURL, ThumbnailTinyURL: mergedThumbTiny,
		VideoURL: merged.VideoURL, Duration: merged.Duration, Playlist: merged.Playlist,
	}, nil
}

func (s *PostService) UpdateLinkPost(ctx context.Context, id string, input models.UpdateLinkPost) (models.LinkPost, error) {
	if err := requireNonEmpty("id", id); err != nil {
		return models.LinkPost{}, err
	}
	post, err := s.store.GetPostByID(ctx, id)
	if err != nil {
		return models.LinkPost{}, fmt.Errorf("PostService.UpdateLinkPost: %w", err)
	}
	cur, ok := post.(models.LinkPost)
	if !ok {
		return models.LinkPost{}, validationErr("post is not a link post")
	}

	merged := models.CreateLinkPost{
		Title: cur.Title, URL: cur.URL, Domain: cur.Domain,
		Description: cur.Description, ThumbnailURL: cur.ThumbnailURL,
		ThumbnailTinyURL: cur.ThumbnailTinyURL,
		Category: cur.Category, Tags: cur.Tags,
	}
	if input.Title != nil {
		merged.Title = *input.Title
	}
	if input.URL != nil {
		merged.URL = *input.URL
	}
	if input.Description != nil {
		merged.Description = input.Description
	}
	if input.Category != nil {
		merged.Category = input.Category
	}
	if input.Tags != nil {
		merged.Tags = input.Tags
	}

	if err := requireNonEmpty("title", merged.Title); err != nil {
		return models.LinkPost{}, err
	}
	if err := requireNonEmpty("url", merged.URL); err != nil {
		return models.LinkPost{}, err
	}
	domain, err := extractDomain(merged.URL)
	if err != nil {
		return models.LinkPost{}, err
	}
	merged.Domain = domain

	if err := s.store.UpdateLinkPost(ctx, id, merged); err != nil {
		return models.LinkPost{}, fmt.Errorf("PostService.UpdateLinkPost: %w", err)
	}
	tags := merged.Tags
	if tags == nil {
		tags = []string{}
	}
	return models.LinkPost{
		ID: id, Type: "link", CreatedAt: cur.CreatedAt, Tags: tags,
		Title: merged.Title, URL: merged.URL, Domain: merged.Domain,
		Description: merged.Description, ThumbnailURL: merged.ThumbnailURL,
		ThumbnailTinyURL: merged.ThumbnailTinyURL, Category: merged.Category,
	}, nil
}

func (s *PostService) DeletePost(ctx context.Context, id string) error {
	if err := requireNonEmpty("id", id); err != nil {
		return err
	}
	post, err := s.store.GetPostByID(ctx, id)
	if err != nil {
		return err
	}
	for _, key := range fileKeysForPost(post) {
		if err := s.files.Delete(ctx, key); err != nil {
			s.log.Error("DeletePost: failed to delete file", "key", key, "err", err)
		}
	}
	return s.store.DeletePost(ctx, id)
}

func fileKeysForPost(post models.Post) []string {
	var keys []string
	switch p := post.(type) {
	case models.MusicPost:
		if p.AudioURL != "" {
			keys = append(keys, urlToStorageKey(p.AudioURL))
		}
		if p.AlbumArt != "" {
			keys = append(keys, urlToStorageKey(p.AlbumArt))
		}
	case models.PhotoPost:
		for _, img := range p.Images {
			keys = append(keys, urlToStorageKey(img.URL))
			if img.ThumbnailSmallURL != nil {
				keys = append(keys, urlToStorageKey(*img.ThumbnailSmallURL))
			}
			if img.ThumbnailMedURL != nil {
				keys = append(keys, urlToStorageKey(*img.ThumbnailMedURL))
			}
			if img.ThumbnailLargeURL != nil {
				keys = append(keys, urlToStorageKey(*img.ThumbnailLargeURL))
			}
		}
	case models.VideoPost:
		if p.ThumbnailURL != "" {
			keys = append(keys, urlToStorageKey(p.ThumbnailURL))
		}
	case models.LinkPost:
		if p.ThumbnailURL != nil && *p.ThumbnailURL != "" {
			keys = append(keys, urlToStorageKey(*p.ThumbnailURL))
		}
	}
	return keys
}

// urlToStorageKey converts a serving URL (/files/<id> or /thumbnails/<id>)
// to a storage key (photos/<id>, music/<id>, or thumbnails/<id>).
func urlToStorageKey(u string) string {
	if after, ok := strings.CutPrefix(u, "/thumbnails/"); ok {
		return "thumbnails/" + after
	}
	if after, ok := strings.CutPrefix(u, "/files/"); ok {
		switch strings.ToLower(filepath.Ext(after)) {
		case ".mp3", ".wav", ".ogg", ".m4a":
			return "music/" + after
		default:
			return "photos/" + after
		}
	}
	return u
}
