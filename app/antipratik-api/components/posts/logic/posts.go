package logic

import (
	"context"
	"fmt"
	"strconv"

	commonerrors "github.com/pratikluitel/antipratik/common/errors"
	"github.com/pratikluitel/antipratik/common/logging"
	"github.com/pratikluitel/antipratik/components/files"
	"github.com/pratikluitel/antipratik/components/posts"
)

// postLogic implements PostLogic.
type postLogic struct {
	store posts.PostStore
	files files.StorageService
	log   logging.Logger
}

// NewpostLogic creates a new postLogic backed by the given store and storage service.
func NewPostLogic(s posts.PostStore, files files.StorageService, log logging.Logger) posts.PostLogic {
	return &postLogic{store: s, files: files, log: log}
}

var validTypes = map[posts.PostType]bool{
	posts.PostTypeEssay: true,
	posts.PostTypeShort: true,
	posts.PostTypeMusic: true,
	posts.PostTypePhoto: true,
	posts.PostTypeVideo: true,
	posts.PostTypeLink:  true,
}

// GetPosts validates the filter and delegates to the store.
func (s *postLogic) GetPosts(ctx context.Context, filter posts.FilterState) ([]posts.Post, error) {
	types := make([]string, 0, len(filter.ActiveTypes))
	for _, t := range filter.ActiveTypes {
		if validTypes[t] {
			types = append(types, t)
		}
	}

	posts, err := s.store.GetPosts(ctx, types, filter.ActiveTags)
	if err != nil {
		return nil, fmt.Errorf("postLogic.GetPosts: %w", err)
	}
	return posts, nil
}

// GetTags returns all tag names sorted alphabetically.
func (s *postLogic) GetTags(ctx context.Context) ([]string, error) {
	return s.store.GetAllTags(ctx)
}

// GetPost validates the slug and delegates to the store.
// Returns nil if the post does not exist.
func (s *postLogic) GetPost(ctx context.Context, slug string) (*posts.EssayPost, error) {
	if slug == "" {
		return nil, nil
	}
	post, err := s.store.GetPostBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("postLogic.GetPost: %w", err)
	}
	return post, nil
}

// GetPostsByIDs returns posts for each given ID, preserving order. Not-found IDs are skipped.
func (s *postLogic) GetPostsByIDs(ctx context.Context, ids []string) ([]posts.Post, error) {
	posts, err := s.store.GetPostsByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("postLogic.GetPostsByIDs: %w", err)
	}
	return posts, nil
}

// ── Write methods ─────────────────────────────────────────────────────────────

func (s *postLogic) CreateEssay(ctx context.Context, input posts.EssayPostInput) (posts.EssayPost, error) {
	if err := commonerrors.RequireNonEmpty("title", input.Title); err != nil {
		return posts.EssayPost{}, err
	}
	if err := commonerrors.RequireNonEmpty("slug", input.Slug); err != nil {
		return posts.EssayPost{}, err
	}
	if err := commonerrors.RequireNonEmpty("body", input.Body); err != nil {
		return posts.EssayPost{}, err
	}
	input.ReadingTimeMinutes = computeReadingTime(input.Body)

	id, createdAt := newID(), nowUTC()
	if err := s.store.CreatePost(ctx, posts.PostTypeEssay, id, createdAt); err != nil {
		return posts.EssayPost{}, fmt.Errorf("postLogic.CreateEssay: %w", err)
	}
	if err := s.store.CreateEssayData(ctx, id, input); err != nil {
		return posts.EssayPost{}, fmt.Errorf("postLogic.CreateEssay data: %w", err)
	}
	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}
	return posts.EssayPost{
		ID: id, Type: posts.PostTypeEssay, CreatedAt: createdAt, Tags: tags,
		Title: input.Title, Slug: input.Slug, Excerpt: input.Excerpt,
		Body: input.Body, ReadingTimeMinutes: input.ReadingTimeMinutes,
	}, nil
}

func (s *postLogic) CreateShort(ctx context.Context, input posts.ShortPostInput) (posts.ShortPost, error) {
	if err := commonerrors.RequireNonEmpty("body", input.Body); err != nil {
		return posts.ShortPost{}, err
	}

	id, createdAt := newID(), nowUTC()
	if err := s.store.CreatePost(ctx, posts.PostTypeShort, id, createdAt); err != nil {
		return posts.ShortPost{}, fmt.Errorf("postLogic.CreateShort: %w", err)
	}
	if err := s.store.CreateShortData(ctx, id, input); err != nil {
		return posts.ShortPost{}, fmt.Errorf("postLogic.CreateShort data: %w", err)
	}
	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}
	return posts.ShortPost{ID: id, Type: posts.PostTypeShort, CreatedAt: createdAt, Tags: tags, Body: input.Body}, nil
}

func (s *postLogic) CreateMusic(ctx context.Context, id string, input posts.MusicPostInput) (posts.MusicPost, error) {
	if err := commonerrors.RequireNonEmpty("title", input.Title); err != nil {
		return posts.MusicPost{}, err
	}
	if err := commonerrors.RequireNonEmpty("audioURL", input.AudioURL); err != nil {
		return posts.MusicPost{}, err
	}
	if err := commonerrors.RequirePositive("duration", input.Duration); err != nil {
		return posts.MusicPost{}, err
	}

	createdAt := nowUTC()
	if err := s.store.CreatePost(ctx, posts.PostTypeMusic, id, createdAt); err != nil {
		return posts.MusicPost{}, fmt.Errorf("postLogic.CreateMusic: %w", err)
	}
	if err := s.store.CreateMusicData(ctx, id, input); err != nil {
		return posts.MusicPost{}, fmt.Errorf("postLogic.CreateMusic data: %w", err)
	}

	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}
	return posts.MusicPost{
		ID: id, Type: posts.PostTypeMusic, CreatedAt: createdAt, Tags: tags,
		Title: input.Title, AlbumArt: input.AlbumArt, AlbumArtTinyURL: input.AlbumArtTinyURL,
		AudioURL: input.AudioURL, Duration: input.Duration, Album: input.Album,
	}, nil
}

func (s *postLogic) CreatePhoto(ctx context.Context, preID string, input posts.PhotoPostInput) (posts.PhotoPost, error) {
	if len(input.Images) == 0 {
		return posts.PhotoPost{}, commonerrors.New("images cannot be empty")
	}
	for i, img := range input.Images {
		if err := commonerrors.RequireNonEmpty(fmt.Sprintf("images[%d].url", i), img.URL); err != nil {
			return posts.PhotoPost{}, err
		}
	}

	id, createdAt := preID, nowUTC()
	if id == "" {
		id = newID()
	}
	if err := s.store.CreatePost(ctx, posts.PostTypePhoto, id, createdAt); err != nil {
		return posts.PhotoPost{}, fmt.Errorf("postLogic.CreatePhoto: %w", err)
	}
	if err := s.store.CreatePhotoData(ctx, id, input); err != nil {
		return posts.PhotoPost{}, fmt.Errorf("postLogic.CreatePhoto data: %w", err)
	}

	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}
	return posts.PhotoPost{
		ID: id, Type: posts.PostTypePhoto, CreatedAt: createdAt, Tags: tags,
		Images: input.Images, Location: input.Location,
	}, nil
}

func (s *postLogic) CreateVideo(ctx context.Context, preID string, input posts.VideoPostInput) (posts.VideoPost, error) {
	if err := commonerrors.RequireNonEmpty("title", input.Title); err != nil {
		return posts.VideoPost{}, err
	}
	if err := commonerrors.RequireNonEmpty("videoURL", input.VideoURL); err != nil {
		return posts.VideoPost{}, err
	}
	if err := commonerrors.RequirePositive("duration", input.Duration); err != nil {
		return posts.VideoPost{}, err
	}

	id, createdAt := preID, nowUTC()
	if id == "" {
		id = newID()
	}
	if err := s.store.CreatePost(ctx, posts.PostTypeVideo, id, createdAt); err != nil {
		return posts.VideoPost{}, fmt.Errorf("postLogic.CreateVideo: %w", err)
	}
	if err := s.store.CreateVideoData(ctx, id, input); err != nil {
		return posts.VideoPost{}, fmt.Errorf("postLogic.CreateVideo data: %w", err)
	}

	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}
	return posts.VideoPost{
		ID: id, Type: posts.PostTypeVideo, CreatedAt: createdAt, Tags: tags,
		Title: input.Title, ThumbnailURL: input.ThumbnailURL, ThumbnailTinyURL: input.ThumbnailTinyURL,
		VideoURL: input.VideoURL, Duration: input.Duration, Playlist: input.Playlist,
	}, nil
}

func (s *postLogic) CreateLinkPost(ctx context.Context, preID string, input posts.LinkPostInput) (posts.LinkPost, error) {
	if err := commonerrors.RequireNonEmpty("title", input.Title); err != nil {
		return posts.LinkPost{}, err
	}
	if err := commonerrors.RequireNonEmpty("url", input.URL); err != nil {
		return posts.LinkPost{}, err
	}
	domain, err := extractDomain(input.URL)
	if err != nil {
		return posts.LinkPost{}, err
	}
	input.Domain = domain

	id, createdAt := preID, nowUTC()
	if id == "" {
		id = newID()
	}
	if err := s.store.CreatePost(ctx, posts.PostTypeLink, id, createdAt); err != nil {
		return posts.LinkPost{}, fmt.Errorf("postLogic.CreateLinkPost: %w", err)
	}
	if err := s.store.CreateLinkPostData(ctx, id, input); err != nil {
		return posts.LinkPost{}, fmt.Errorf("postLogic.CreateLinkPost data: %w", err)
	}

	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}
	return posts.LinkPost{
		ID: id, Type: posts.PostTypeLink, CreatedAt: createdAt, Tags: tags,
		Title: input.Title, URL: input.URL, Domain: input.Domain,
		Description: input.Description, ThumbnailURL: input.ThumbnailURL,
		ThumbnailTinyURL: input.ThumbnailTinyURL, Category: input.Category,
	}, nil
}

func (s *postLogic) UpdateEssay(ctx context.Context, id string, input posts.UpdateEssayPost) (posts.EssayPost, error) {
	if err := commonerrors.RequireNonEmpty("id", id); err != nil {
		return posts.EssayPost{}, err
	}
	post, err := s.store.GetPostByID(ctx, id)
	if err != nil {
		return posts.EssayPost{}, fmt.Errorf("postLogic.UpdateEssay: %w", err)
	}
	cur, ok := post.(posts.EssayPost)
	if !ok {
		return posts.EssayPost{}, commonerrors.New("post is not an essay")
	}

	merged := posts.EssayPostInput{Title: cur.Title, Slug: cur.Slug, Excerpt: cur.Excerpt, Body: cur.Body, Tags: cur.Tags}
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

	if err := commonerrors.RequireNonEmpty("title", merged.Title); err != nil {
		return posts.EssayPost{}, err
	}
	if err := commonerrors.RequireNonEmpty("slug", merged.Slug); err != nil {
		return posts.EssayPost{}, err
	}
	if err := commonerrors.RequireNonEmpty("body", merged.Body); err != nil {
		return posts.EssayPost{}, err
	}
	merged.ReadingTimeMinutes = computeReadingTime(merged.Body)

	if err := s.store.UpdateEssay(ctx, id, merged); err != nil {
		return posts.EssayPost{}, fmt.Errorf("postLogic.UpdateEssay: %w", err)
	}
	tags := merged.Tags
	if tags == nil {
		tags = []string{}
	}
	return posts.EssayPost{
		ID: id, Type: posts.PostTypeEssay, CreatedAt: cur.CreatedAt, Tags: tags,
		Title: merged.Title, Slug: merged.Slug, Excerpt: merged.Excerpt,
		Body: merged.Body, ReadingTimeMinutes: merged.ReadingTimeMinutes,
	}, nil
}

func (s *postLogic) UpdateShort(ctx context.Context, id string, input posts.UpdateShortPost) (posts.ShortPost, error) {
	if err := commonerrors.RequireNonEmpty("id", id); err != nil {
		return posts.ShortPost{}, err
	}
	post, err := s.store.GetPostByID(ctx, id)
	if err != nil {
		return posts.ShortPost{}, fmt.Errorf("postLogic.UpdateShort: %w", err)
	}
	cur, ok := post.(posts.ShortPost)
	if !ok {
		return posts.ShortPost{}, commonerrors.New("post is not a short post")
	}

	merged := posts.ShortPostInput{Body: cur.Body, Tags: cur.Tags}
	if input.Body != nil {
		merged.Body = *input.Body
	}
	if input.Tags != nil {
		merged.Tags = input.Tags
	}

	if err := commonerrors.RequireNonEmpty("body", merged.Body); err != nil {
		return posts.ShortPost{}, err
	}

	if err := s.store.UpdateShort(ctx, id, merged); err != nil {
		return posts.ShortPost{}, fmt.Errorf("postLogic.UpdateShort: %w", err)
	}
	tags := merged.Tags
	if tags == nil {
		tags = []string{}
	}
	return posts.ShortPost{ID: id, Type: posts.PostTypeShort, CreatedAt: cur.CreatedAt, Tags: tags, Body: merged.Body}, nil
}

func (s *postLogic) UpdateMusic(ctx context.Context, id string, input posts.UpdateMusicPost) (posts.MusicPost, error) {
	if err := commonerrors.RequireNonEmpty("id", id); err != nil {
		return posts.MusicPost{}, err
	}

	post, err := s.store.GetPostByID(ctx, id)
	if err != nil {
		return posts.MusicPost{}, fmt.Errorf("postLogic.UpdateMusic: %w", err)
	}
	cur, ok := post.(posts.MusicPost)
	if !ok {
		return posts.MusicPost{}, commonerrors.New("post is not a music post")
	}

	merged := posts.MusicPostInput{
		Title: cur.Title, AudioURL: cur.AudioURL, AlbumArt: cur.AlbumArt,
		AlbumArtTinyURL: cur.AlbumArtTinyURL, AlbumArtSmallURL: cur.AlbumArtSmallURL,
		AlbumArtMedURL: cur.AlbumArtMedURL, AlbumArtLargeURL: cur.AlbumArtLargeURL,
		Duration: cur.Duration, Album: cur.Album, Tags: cur.Tags,
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
	if input.Album != nil {
		merged.Album = input.Album
	}
	if input.Tags != nil {
		merged.Tags = input.Tags
	}

	if err := commonerrors.RequireNonEmpty("title", merged.Title); err != nil {
		return posts.MusicPost{}, err
	}
	if err := commonerrors.RequireNonEmpty("audioURL", merged.AudioURL); err != nil {
		return posts.MusicPost{}, err
	}
	if err := commonerrors.RequirePositive("duration", merged.Duration); err != nil {
		return posts.MusicPost{}, err
	}

	var oldArtKeys []string
	if input.AlbumArt != nil && cur.AlbumArt != "" && cur.AlbumArt != *input.AlbumArt {
		oldArtKeys = albumArtFileKeys(cur)
	}

	if err := s.store.UpdateMusic(ctx, id, merged); err != nil {
		return posts.MusicPost{}, fmt.Errorf("postLogic.UpdateMusic: %w", err)
	}

	if len(oldArtKeys) > 0 {
		log := s.log
		files := s.files
		go func() {
			for _, key := range oldArtKeys {
				if err := files.Delete(context.Background(), key); err != nil {
					log.Error("UpdateMusic: failed to delete old album art", "key", key, "err", err)
				}
			}
		}()
	}

	tags := merged.Tags
	if tags == nil {
		tags = []string{}
	}
	return posts.MusicPost{
		ID: id, Type: posts.PostTypeMusic, CreatedAt: cur.CreatedAt, Tags: tags,
		Title: merged.Title, AlbumArt: merged.AlbumArt, AlbumArtTinyURL: merged.AlbumArtTinyURL,
		AlbumArtSmallURL: merged.AlbumArtSmallURL, AlbumArtMedURL: merged.AlbumArtMedURL,
		AlbumArtLargeURL: merged.AlbumArtLargeURL,
		AudioURL:         merged.AudioURL, Duration: merged.Duration, Album: merged.Album,
	}, nil
}

func (s *postLogic) UpdatePhoto(ctx context.Context, id string, input posts.PhotoPostInput) (posts.PhotoPost, error) {
	if err := commonerrors.RequireNonEmpty("id", id); err != nil {
		return posts.PhotoPost{}, err
	}

	post, err := s.store.GetPostByID(ctx, id)
	if err != nil {
		return posts.PhotoPost{}, fmt.Errorf("postLogic.UpdatePhoto: %w", err)
	}
	cur, ok := post.(posts.PhotoPost)
	if !ok {
		return posts.PhotoPost{}, commonerrors.New("post is not a photo post")
	}

	merged := posts.PhotoPostInput{Images: cur.Images, Location: cur.Location, Tags: cur.Tags}
	if input.Location != nil {
		merged.Location = input.Location
	}
	if input.Tags != nil {
		merged.Tags = input.Tags
	}

	if err := s.store.UpdatePhoto(ctx, id, merged); err != nil {
		return posts.PhotoPost{}, fmt.Errorf("postLogic.UpdatePhoto: %w", err)
	}

	tags := merged.Tags
	if tags == nil {
		tags = []string{}
	}
	return posts.PhotoPost{
		ID: id, Type: posts.PostTypePhoto, CreatedAt: cur.CreatedAt, Tags: tags,
		Images: merged.Images, Location: merged.Location,
	}, nil
}

func (s *postLogic) UpdateVideo(ctx context.Context, id string, input posts.UpdateVideoPost) (posts.VideoPost, error) {
	if err := commonerrors.RequireNonEmpty("id", id); err != nil {
		return posts.VideoPost{}, err
	}
	post, err := s.store.GetPostByID(ctx, id)
	if err != nil {
		return posts.VideoPost{}, fmt.Errorf("postLogic.UpdateVideo: %w", err)
	}
	cur, ok := post.(posts.VideoPost)
	if !ok {
		return posts.VideoPost{}, commonerrors.New("post is not a video post")
	}

	merged := posts.VideoPostInput{
		Title: cur.Title, ThumbnailURL: cur.ThumbnailURL, ThumbnailTinyURL: cur.ThumbnailTinyURL,
		ThumbnailSmallURL: cur.ThumbnailSmallURL, ThumbnailMedURL: cur.ThumbnailMedURL,
		ThumbnailLargeURL: cur.ThumbnailLargeURL,
		VideoURL:          cur.VideoURL, Duration: cur.Duration, Playlist: cur.Playlist, Tags: cur.Tags,
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

	if err := commonerrors.RequireNonEmpty("title", merged.Title); err != nil {
		return posts.VideoPost{}, err
	}
	if err := commonerrors.RequireNonEmpty("videoURL", merged.VideoURL); err != nil {
		return posts.VideoPost{}, err
	}
	if err := commonerrors.RequirePositive("duration", merged.Duration); err != nil {
		return posts.VideoPost{}, err
	}

	var oldThumbKeys []string
	if input.ThumbnailURL != nil && cur.ThumbnailURL != "" && cur.ThumbnailURL != *input.ThumbnailURL {
		oldThumbKeys = thumbnailFileKeys(cur.ThumbnailURL, cur.ThumbnailTinyURL, cur.ThumbnailSmallURL, cur.ThumbnailMedURL, cur.ThumbnailLargeURL)
	}

	if err := s.store.UpdateVideo(ctx, id, merged); err != nil {
		return posts.VideoPost{}, fmt.Errorf("postLogic.UpdateVideo: %w", err)
	}

	if len(oldThumbKeys) > 0 {
		log := s.log
		files := s.files
		go func() {
			for _, key := range oldThumbKeys {
				if err := files.Delete(context.Background(), key); err != nil {
					log.Error("UpdateVideo: failed to delete old thumbnail", "key", key, "err", err)
				}
			}
		}()
	}

	tags := merged.Tags
	if tags == nil {
		tags = []string{}
	}
	return posts.VideoPost{
		ID: id, Type: posts.PostTypeVideo, CreatedAt: cur.CreatedAt, Tags: tags,
		Title: merged.Title, ThumbnailURL: merged.ThumbnailURL, ThumbnailTinyURL: merged.ThumbnailTinyURL,
		ThumbnailSmallURL: merged.ThumbnailSmallURL, ThumbnailMedURL: merged.ThumbnailMedURL,
		ThumbnailLargeURL: merged.ThumbnailLargeURL,
		VideoURL:          merged.VideoURL, Duration: merged.Duration, Playlist: merged.Playlist,
	}, nil
}

func (s *postLogic) UpdateLinkPost(ctx context.Context, id string, input posts.UpdateLinkPost) (posts.LinkPost, error) {
	if err := commonerrors.RequireNonEmpty("id", id); err != nil {
		return posts.LinkPost{}, err
	}
	post, err := s.store.GetPostByID(ctx, id)
	if err != nil {
		return posts.LinkPost{}, fmt.Errorf("postLogic.UpdateLinkPost: %w", err)
	}
	cur, ok := post.(posts.LinkPost)
	if !ok {
		return posts.LinkPost{}, commonerrors.New("post is not a link post")
	}

	merged := posts.LinkPostInput{
		Title: cur.Title, URL: cur.URL, Domain: cur.Domain,
		Description: cur.Description, ThumbnailURL: cur.ThumbnailURL,
		ThumbnailTinyURL: cur.ThumbnailTinyURL, ThumbnailSmallURL: cur.ThumbnailSmallURL,
		ThumbnailMedURL: cur.ThumbnailMedURL, ThumbnailLargeURL: cur.ThumbnailLargeURL,
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

	if err = commonerrors.RequireNonEmpty("title", merged.Title); err != nil {
		return posts.LinkPost{}, err
	}
	if err = commonerrors.RequireNonEmpty("url", merged.URL); err != nil {
		return posts.LinkPost{}, err
	}
	domain, err := extractDomain(merged.URL)
	if err != nil {
		return posts.LinkPost{}, err
	}
	merged.Domain = domain

	var oldThumbKeys []string
	if input.ThumbnailURL != nil && cur.ThumbnailURL != nil && *cur.ThumbnailURL != "" && *cur.ThumbnailURL != *input.ThumbnailURL {
		oldThumbKeys = thumbnailFileKeys(*cur.ThumbnailURL, cur.ThumbnailTinyURL, cur.ThumbnailSmallURL, cur.ThumbnailMedURL, cur.ThumbnailLargeURL)
	}

	if err = s.store.UpdateLinkPost(ctx, id, merged); err != nil {
		return posts.LinkPost{}, fmt.Errorf("postLogic.UpdateLinkPost: %w", err)
	}

	if len(oldThumbKeys) > 0 {
		log := s.log
		files := s.files
		go func() {
			for _, key := range oldThumbKeys {
				if err := files.Delete(context.Background(), key); err != nil {
					log.Error("UpdateLinkPost: failed to delete old thumbnail", "key", key, "err", err)
				}
			}
		}()
	}
	tags := merged.Tags
	if tags == nil {
		tags = []string{}
	}

	return posts.LinkPost{
		ID: id, Type: posts.PostTypeLink, CreatedAt: cur.CreatedAt, Tags: tags,
		Title: merged.Title, URL: merged.URL, Domain: merged.Domain,
		Description: merged.Description, ThumbnailURL: merged.ThumbnailURL,
		ThumbnailTinyURL: merged.ThumbnailTinyURL, ThumbnailSmallURL: merged.ThumbnailSmallURL,
		ThumbnailMedURL: merged.ThumbnailMedURL, ThumbnailLargeURL: merged.ThumbnailLargeURL,
		Category: merged.Category,
	}, nil
}

func (s *postLogic) AddPhotoImage(ctx context.Context, postID string, image posts.PhotoImage) (*posts.PhotoImage, error) {
	if err := commonerrors.RequireNonEmpty("postID", postID); err != nil {
		return nil, err
	}
	if err := commonerrors.RequireNonEmpty("image url", image.URL); err != nil {
		return nil, err
	}
	post, err := s.store.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if _, ok := post.(posts.PhotoPost); !ok {
		return nil, commonerrors.New("post is not a photo post")
	}
	return s.store.AddPhotoImage(ctx, postID, image)
}

func (s *postLogic) GetPhotoImage(ctx context.Context, postID string, imageIDStr string) (*posts.PhotoImage, error) {
	if err := commonerrors.RequireNonEmpty("postID", postID); err != nil {
		return nil, err
	}
	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil {
		return nil, commonerrors.New("imageID must be an integer")
	}
	// Returns nil, nil if not found — API layer maps this to 404.
	return s.store.GetPhotoImage(ctx, postID, imageID)
}

func (s *postLogic) UpdatePhotoImage(ctx context.Context, postID string, imageIDStr string, input posts.UpdatePhotoImage) (*posts.PhotoImage, error) {
	if err := commonerrors.RequireNonEmpty("postID", postID); err != nil {
		return nil, err
	}
	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil {
		return nil, commonerrors.New("imageID must be an integer")
	}
	post, err := s.store.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if _, ok := post.(posts.PhotoPost); !ok {
		return nil, commonerrors.New("post is not a photo post")
	}
	// Returns nil, nil if image not found — API layer maps this to 404.
	return s.store.UpdatePhotoImage(ctx, postID, imageID, input)
}

func (s *postLogic) DeletePhotoImage(ctx context.Context, postID string, imageIDStr string) (notFound bool, err error) {
	if err = commonerrors.RequireNonEmpty("postID", postID); err != nil {
		return false, err
	}
	imageID, convErr := strconv.Atoi(imageIDStr)
	if convErr != nil {
		return false, commonerrors.New("imageID must be an integer")
	}
	// Fetch the image first to collect file keys for cleanup.
	img, err := s.store.GetPhotoImage(ctx, postID, imageID)
	if err != nil {
		return false, err
	}
	if img == nil {
		return true, nil // not found
	}
	// Ensure the post has more than one image.
	post, err := s.store.GetPostByID(ctx, postID)
	if err != nil {
		return false, err
	}
	photoPost, ok := post.(posts.PhotoPost)
	if !ok {
		return false, commonerrors.New("post is not a photo post")
	}
	if len(photoPost.Images) <= 1 {
		return false, commonerrors.New("cannot delete the only image in a photo post")
	}
	if err = s.store.DeletePhotoImage(ctx, postID, imageID); err != nil {
		return false, err
	}
	// Clean up files in the background (same pattern as DeletePost).
	keys := fileKeysForImage(img)
	if len(keys) > 0 {
		log := s.log
		files := s.files
		go func() {
			for _, key := range keys {
				if err := files.Delete(context.Background(), key); err != nil {
					log.Error("DeletePhotoImage: failed to delete file", "key", key, "err", err)
				}
			}
		}()
	}
	return false, nil
}

func (s *postLogic) DeletePost(ctx context.Context, id string) error {
	if err := commonerrors.RequireNonEmpty("id", id); err != nil {
		return err
	}
	post, err := s.store.GetPostByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete the database record first so the post is immediately gone for readers.
	// File cleanup runs in the background — errors are logged but do not affect the
	// response, since orphaned files are recoverable but a partial delete is not.
	if err := s.store.DeletePost(ctx, id); err != nil {
		return err
	}

	keys := fileKeysForPost(post)
	if len(keys) > 0 {
		log := s.log
		files := s.files
		go func() {
			for _, key := range keys {
				if err := files.Delete(context.Background(), key); err != nil {
					log.Error("DeletePost: failed to delete file", "key", key, "err", err)
				}
			}
		}()
	}

	return nil
}
