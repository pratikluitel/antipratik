package logic

import (
	"context"
	"fmt"

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
