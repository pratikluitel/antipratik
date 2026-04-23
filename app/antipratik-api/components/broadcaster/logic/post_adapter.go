package logic

import (
	"context"
	"fmt"

	"github.com/pratikluitel/antipratik/components/posts"
)

// postAdapter adapts posts.PostsService to the PostService interface,
// mapping the posts component's model types to the PostSummary values the broadcaster needs.
type postAdapter struct {
	svc posts.PostsService
}

func newPostAdapter(svc posts.PostsService) PostService {
	return &postAdapter{svc: svc}
}

func (a *postAdapter) GetPostsByIDs(ctx context.Context, ids []string) ([]PostSummary, error) {
	posts, err := a.svc.GetPostsByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("postAdapter.GetPostsByIDs: %w", err)
	}
	out := make([]PostSummary, 0, len(posts))
	for _, p := range posts {
		out = append(out, toPostSummary(p))
	}
	return out, nil
}

func toPostSummary(p posts.Post) PostSummary {
	switch v := p.(type) {
	case posts.EssayPost:
		return PostSummary{
			ID:        v.ID,
			Type:      v.Type,
			Title:     v.Title,
			Slug:      v.Slug,
			Excerpt:   v.Excerpt,
			Body:      markdownToHTML(v.Body),
			CreatedAt: v.CreatedAt,
		}
	case posts.ShortPost:
		return PostSummary{
			ID:        v.ID,
			Type:      v.Type,
			Excerpt:   v.Body,
			CreatedAt: v.CreatedAt,
		}
	case posts.MusicPost:
		var albumArtMed string
		if v.AlbumArtMedURL != nil {
			albumArtMed = *v.AlbumArtMedURL
		}
		return PostSummary{
			ID:                v.ID,
			Type:              v.Type,
			Title:             v.Title,
			AlbumArtMediumURL: albumArtMed,
			CreatedAt:         v.CreatedAt,
		}
	case posts.PhotoPost:
		var imageURL, thumbMedium, thumbLarge string
		if len(v.Images) > 0 {
			imageURL = v.Images[0].URL
			if v.Images[0].ThumbnailMedURL != nil {
				thumbMedium = *v.Images[0].ThumbnailMedURL
			}
			if v.Images[0].ThumbnailLargeURL != nil {
				thumbLarge = *v.Images[0].ThumbnailLargeURL
			}
		}
		return PostSummary{
			ID:                 v.ID,
			Type:               v.Type,
			ImageURL:           imageURL,
			ThumbnailMediumURL: thumbMedium,
			ThumbnailLargeURL:  thumbLarge,
			CreatedAt:          v.CreatedAt,
		}
	case posts.VideoPost:
		var thumbMed string
		if v.ThumbnailMedURL != nil {
			thumbMed = *v.ThumbnailMedURL
		}
		return PostSummary{
			ID:                 v.ID,
			Type:               v.Type,
			Title:              v.Title,
			VideoURL:           v.VideoURL,
			ThumbnailMediumURL: thumbMed,
			CreatedAt:          v.CreatedAt,
		}
	case posts.LinkPost:
		desc := ""
		if v.Description != nil {
			desc = *v.Description
		}
		var thumbMed string
		if v.ThumbnailMedURL != nil {
			thumbMed = *v.ThumbnailMedURL
		}
		return PostSummary{
			ID:                 v.ID,
			Type:               v.Type,
			Title:              v.Title,
			LinkURL:            v.URL,
			Domain:             v.Domain,
			Excerpt:            desc,
			ThumbnailMediumURL: thumbMed,
			CreatedAt:          v.CreatedAt,
		}
	default:
		return PostSummary{}
	}
}
