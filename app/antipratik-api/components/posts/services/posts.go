package services

import (
	"context"

	"github.com/pratikluitel/antipratik/components/posts"
)

type postsService struct {
	logic posts.PostLogic
}

// NewPostsService returns a PostsService backed by the given PostLogic.
func NewPostsService(l posts.PostLogic) posts.PostsService {
	return &postsService{logic: l}
}

func (s *postsService) GetPostsByIDs(ctx context.Context, ids []string) ([]posts.Post, error) {
	return s.logic.GetPostsByIDs(ctx, ids)
}
