package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pratikluitel/antipratik/models"
	"github.com/pratikluitel/antipratik/store"
)

// PostService implements PostLogic.
type PostService struct {
	store store.PostStore
}

// NewPostService creates a new PostService backed by the given store.
func NewPostService(s store.PostStore) *PostService {
	return &PostService{store: s}
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
	if err := requirePositive("readingTimeMinutes", input.ReadingTimeMinutes); err != nil {
		return models.EssayPost{}, err
	}

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

func (s *PostService) CreateMusic(ctx context.Context, preID string, input models.CreateMusicPost) (models.MusicPost, error) {
	if err := requireNonEmpty("title", input.Title); err != nil {
		return models.MusicPost{}, err
	}
	if input.Album != nil && *input.Album != "" {
		if err := requireNonEmpty("albumArt", input.AlbumArt); err != nil {
			return models.MusicPost{}, err
		}
	}
	if err := requireNonEmpty("audioURL", input.AudioURL); err != nil {
		return models.MusicPost{}, err
	}
	if err := requirePositive("duration", input.Duration); err != nil {
		return models.MusicPost{}, err
	}

	id := preID
	if id == "" {
		id = newID()
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
	return models.MusicPost{
		ID: id, Type: "music", CreatedAt: createdAt, Tags: tags,
		Title: input.Title, AlbumArt: input.AlbumArt, AudioURL: input.AudioURL,
		Duration: input.Duration, Album: input.Album,
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
	return models.VideoPost{
		ID: id, Type: "video", CreatedAt: createdAt, Tags: tags,
		Title: input.Title, ThumbnailURL: input.ThumbnailURL,
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
	if err := requireNonEmpty("domain", input.Domain); err != nil {
		return models.LinkPost{}, err
	}

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
		Description: input.Description, ThumbnailURL: input.ThumbnailURL, Category: input.Category,
	}, nil
}

func (s *PostService) UpdateEssay(ctx context.Context, id string, input models.CreateEssayPost) error {
	if err := requireNonEmpty("id", id); err != nil {
		return err
	}
	if err := requireNonEmpty("title", input.Title); err != nil {
		return err
	}
	if err := requireNonEmpty("slug", input.Slug); err != nil {
		return err
	}
	if err := requireNonEmpty("body", input.Body); err != nil {
		return err
	}
	if err := requirePositive("readingTimeMinutes", input.ReadingTimeMinutes); err != nil {
		return err
	}
	return s.store.UpdateEssay(ctx, id, input)
}

func (s *PostService) UpdateShort(ctx context.Context, id string, input models.CreateShortPost) error {
	if err := requireNonEmpty("id", id); err != nil {
		return err
	}
	if err := requireNonEmpty("body", input.Body); err != nil {
		return err
	}
	return s.store.UpdateShort(ctx, id, input)
}

func (s *PostService) UpdateMusic(ctx context.Context, id string, input models.CreateMusicPost) error {
	if err := requireNonEmpty("id", id); err != nil {
		return err
	}
	if err := requireNonEmpty("title", input.Title); err != nil {
		return err
	}
	if input.Album != nil && *input.Album != "" {
		if err := requireNonEmpty("albumArt", input.AlbumArt); err != nil {
			return err
		}
	}
	if err := requireNonEmpty("audioURL", input.AudioURL); err != nil {
		return err
	}
	if err := requirePositive("duration", input.Duration); err != nil {
		return err
	}
	return s.store.UpdateMusic(ctx, id, input)
}

func (s *PostService) UpdatePhoto(ctx context.Context, id string, input models.CreatePhotoPost) error {
	if err := requireNonEmpty("id", id); err != nil {
		return err
	}
	if len(input.Images) == 0 {
		return validationErr("images cannot be empty")
	}
	for i, img := range input.Images {
		if err := requireNonEmpty(fmt.Sprintf("images[%d].url", i), img.URL); err != nil {
			return err
		}
	}
	return s.store.UpdatePhoto(ctx, id, input)
}

func (s *PostService) UpdateVideo(ctx context.Context, id string, input models.CreateVideoPost) error {
	if err := requireNonEmpty("id", id); err != nil {
		return err
	}
	if err := requireNonEmpty("title", input.Title); err != nil {
		return err
	}
	if err := requireNonEmpty("videoURL", input.VideoURL); err != nil {
		return err
	}
	if err := requirePositive("duration", input.Duration); err != nil {
		return err
	}
	return s.store.UpdateVideo(ctx, id, input)
}

func (s *PostService) UpdateLinkPost(ctx context.Context, id string, input models.CreateLinkPost) error {
	if err := requireNonEmpty("id", id); err != nil {
		return err
	}
	if err := requireNonEmpty("title", input.Title); err != nil {
		return err
	}
	if err := requireNonEmpty("url", input.URL); err != nil {
		return err
	}
	if err := requireNonEmpty("domain", input.Domain); err != nil {
		return err
	}
	return s.store.UpdateLinkPost(ctx, id, input)
}

func (s *PostService) DeletePost(ctx context.Context, id string) error {
	if err := requireNonEmpty("id", id); err != nil {
		return err
	}
	return s.store.DeletePost(ctx, id)
}
