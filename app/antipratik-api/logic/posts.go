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

func (s *PostService) CreateEssay(ctx context.Context, input models.CreateEssayPost) (string, error) {
	if err := requireNonEmpty("title", input.Title); err != nil {
		return "", err
	}
	if err := requireNonEmpty("slug", input.Slug); err != nil {
		return "", err
	}
	if err := requireNonEmpty("body", input.Body); err != nil {
		return "", err
	}
	if err := requirePositive("readingTimeMinutes", input.ReadingTimeMinutes); err != nil {
		return "", err
	}

	id := newID()
	if err := s.store.CreatePost(ctx, "essay", id, nowUTC()); err != nil {
		return "", fmt.Errorf("PostService.CreateEssay: %w", err)
	}
	if err := s.store.CreateEssayData(ctx, id, input); err != nil {
		return "", fmt.Errorf("PostService.CreateEssay data: %w", err)
	}
	return id, nil
}

func (s *PostService) CreateShort(ctx context.Context, input models.CreateShortPost) (string, error) {
	if err := requireNonEmpty("body", input.Body); err != nil {
		return "", err
	}

	id := newID()
	if err := s.store.CreatePost(ctx, "short", id, nowUTC()); err != nil {
		return "", fmt.Errorf("PostService.CreateShort: %w", err)
	}
	if err := s.store.CreateShortData(ctx, id, input); err != nil {
		return "", fmt.Errorf("PostService.CreateShort data: %w", err)
	}
	return id, nil
}

func (s *PostService) CreateMusic(ctx context.Context, preID string, input models.CreateMusicPost) (string, error) {
	if err := requireNonEmpty("title", input.Title); err != nil {
		return "", err
	}
	if input.Album != nil && *input.Album != "" {
		if err := requireNonEmpty("albumArt", input.AlbumArt); err != nil {
			return "", err
		}
	}
	if err := requireNonEmpty("audioURL", input.AudioURL); err != nil {
		return "", err
	}
	if err := requirePositive("duration", input.Duration); err != nil {
		return "", err
	}

	id := preID
	if id == "" {
		id = newID()
	}
	if err := s.store.CreatePost(ctx, "music", id, nowUTC()); err != nil {
		return "", fmt.Errorf("PostService.CreateMusic: %w", err)
	}
	if err := s.store.CreateMusicData(ctx, id, input); err != nil {
		return "", fmt.Errorf("PostService.CreateMusic data: %w", err)
	}
	return id, nil
}

func (s *PostService) CreatePhoto(ctx context.Context, preID string, input models.CreatePhotoPost) (string, error) {
	if len(input.Images) == 0 {
		return "", validationErr("images cannot be empty")
	}
	for i, img := range input.Images {
		if err := requireNonEmpty(fmt.Sprintf("images[%d].url", i), img.URL); err != nil {
			return "", err
		}
	}

	id := preID
	if id == "" {
		id = newID()
	}
	if err := s.store.CreatePost(ctx, "photo", id, nowUTC()); err != nil {
		return "", fmt.Errorf("PostService.CreatePhoto: %w", err)
	}
	if err := s.store.CreatePhotoData(ctx, id, input); err != nil {
		return "", fmt.Errorf("PostService.CreatePhoto data: %w", err)
	}
	return id, nil
}

func (s *PostService) CreateVideo(ctx context.Context, preID string, input models.CreateVideoPost) (string, error) {
	if err := requireNonEmpty("title", input.Title); err != nil {
		return "", err
	}
	if err := requireNonEmpty("videoURL", input.VideoURL); err != nil {
		return "", err
	}
	if err := requirePositive("duration", input.Duration); err != nil {
		return "", err
	}

	id := preID
	if id == "" {
		id = newID()
	}
	if err := s.store.CreatePost(ctx, "video", id, nowUTC()); err != nil {
		return "", fmt.Errorf("PostService.CreateVideo: %w", err)
	}
	if err := s.store.CreateVideoData(ctx, id, input); err != nil {
		return "", fmt.Errorf("PostService.CreateVideo data: %w", err)
	}
	return id, nil
}

func (s *PostService) CreateLinkPost(ctx context.Context, preID string, input models.CreateLinkPost) (string, error) {
	if err := requireNonEmpty("title", input.Title); err != nil {
		return "", err
	}
	if err := requireNonEmpty("url", input.URL); err != nil {
		return "", err
	}
	if err := requireNonEmpty("domain", input.Domain); err != nil {
		return "", err
	}

	id := preID
	if id == "" {
		id = newID()
	}
	if err := s.store.CreatePost(ctx, "link", id, nowUTC()); err != nil {
		return "", fmt.Errorf("PostService.CreateLinkPost: %w", err)
	}
	if err := s.store.CreateLinkPostData(ctx, id, input); err != nil {
		return "", fmt.Errorf("PostService.CreateLinkPost data: %w", err)
	}
	return id, nil
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
