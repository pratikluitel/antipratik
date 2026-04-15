package logic

import (
	"context"
	"fmt"
	"strconv"

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

var validTypes = map[models.PostType]bool{
	models.PostTypeEssay: true,
	models.PostTypeShort: true,
	models.PostTypeMusic: true,
	models.PostTypePhoto: true,
	models.PostTypeVideo: true,
	models.PostTypeLink:  true,
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

func (s *PostService) CreateEssay(ctx context.Context, input models.EssayPostInput) (models.EssayPost, error) {
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
	if err := s.store.CreatePost(ctx, models.PostTypeEssay, id, createdAt); err != nil {
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
		ID: id, Type: models.PostTypeEssay, CreatedAt: createdAt, Tags: tags,
		Title: input.Title, Slug: input.Slug, Excerpt: input.Excerpt,
		Body: input.Body, ReadingTimeMinutes: input.ReadingTimeMinutes,
	}, nil
}

func (s *PostService) CreateShort(ctx context.Context, input models.ShortPostInput) (models.ShortPost, error) {
	if err := requireNonEmpty("body", input.Body); err != nil {
		return models.ShortPost{}, err
	}

	id, createdAt := newID(), nowUTC()
	if err := s.store.CreatePost(ctx, models.PostTypeShort, id, createdAt); err != nil {
		return models.ShortPost{}, fmt.Errorf("PostService.CreateShort: %w", err)
	}
	if err := s.store.CreateShortData(ctx, id, input); err != nil {
		return models.ShortPost{}, fmt.Errorf("PostService.CreateShort data: %w", err)
	}
	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}
	return models.ShortPost{ID: id, Type: models.PostTypeShort, CreatedAt: createdAt, Tags: tags, Body: input.Body}, nil
}

func (s *PostService) CreateMusic(ctx context.Context, id string, input models.MusicPostInput) (models.MusicPost, error) {
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
	if err := s.store.CreatePost(ctx, models.PostTypeMusic, id, createdAt); err != nil {
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
		ID: id, Type: models.PostTypeMusic, CreatedAt: createdAt, Tags: tags,
		Title: input.Title, AlbumArt: input.AlbumArt, AlbumArtTinyURL: input.AlbumArtTinyURL,
		AudioURL: input.AudioURL, Duration: input.Duration, Album: input.Album,
	}, nil
}

func (s *PostService) CreatePhoto(ctx context.Context, preID string, input models.PhotoPostInput) (models.PhotoPost, error) {
	if len(input.Images) == 0 {
		return models.PhotoPost{}, validationErr("images cannot be empty")
	}
	for i, img := range input.Images {
		if err := requireNonEmpty(fmt.Sprintf("images[%d].url", i), img.URL); err != nil {
			return models.PhotoPost{}, err
		}
	}

	id, createdAt := preID, nowUTC()
	if id == "" {
		id = newID()
	}
	if err := s.store.CreatePost(ctx, models.PostTypePhoto, id, createdAt); err != nil {
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
		ID: id, Type: models.PostTypePhoto, CreatedAt: createdAt, Tags: tags,
		Images: input.Images, Location: input.Location,
	}, nil
}

func (s *PostService) CreateVideo(ctx context.Context, preID string, input models.VideoPostInput) (models.VideoPost, error) {
	if err := requireNonEmpty("title", input.Title); err != nil {
		return models.VideoPost{}, err
	}
	if err := requireNonEmpty("videoURL", input.VideoURL); err != nil {
		return models.VideoPost{}, err
	}
	if err := requirePositive("duration", input.Duration); err != nil {
		return models.VideoPost{}, err
	}

	id, createdAt := preID, nowUTC()
	if id == "" {
		id = newID()
	}
	if err := s.store.CreatePost(ctx, models.PostTypeVideo, id, createdAt); err != nil {
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
		ID: id, Type: models.PostTypeVideo, CreatedAt: createdAt, Tags: tags,
		Title: input.Title, ThumbnailURL: input.ThumbnailURL, ThumbnailTinyURL: input.ThumbnailTinyURL,
		VideoURL: input.VideoURL, Duration: input.Duration, Playlist: input.Playlist,
	}, nil
}

func (s *PostService) CreateLinkPost(ctx context.Context, preID string, input models.LinkPostInput) (models.LinkPost, error) {
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

	id, createdAt := preID, nowUTC()
	if id == "" {
		id = newID()
	}
	if err := s.store.CreatePost(ctx, models.PostTypeLink, id, createdAt); err != nil {
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
		ID: id, Type: models.PostTypeLink, CreatedAt: createdAt, Tags: tags,
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

	merged := models.EssayPostInput{Title: cur.Title, Slug: cur.Slug, Excerpt: cur.Excerpt, Body: cur.Body, Tags: cur.Tags}
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
		ID: id, Type: models.PostTypeEssay, CreatedAt: cur.CreatedAt, Tags: tags,
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

	merged := models.ShortPostInput{Body: cur.Body, Tags: cur.Tags}
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
	return models.ShortPost{ID: id, Type: models.PostTypeShort, CreatedAt: cur.CreatedAt, Tags: tags, Body: merged.Body}, nil
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

	merged := models.MusicPostInput{
		Title: cur.Title, AudioURL: cur.AudioURL, AlbumArt: cur.AlbumArt,
		AlbumArtTinyURL: cur.AlbumArtTinyURL, Duration: cur.Duration, Album: cur.Album, Tags: cur.Tags,
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
	return models.MusicPost{
		ID: id, Type: models.PostTypeMusic, CreatedAt: cur.CreatedAt, Tags: tags,
		Title: merged.Title, AlbumArt: merged.AlbumArt, AlbumArtTinyURL: merged.AlbumArtTinyURL,
		AudioURL: merged.AudioURL, Duration: merged.Duration, Album: merged.Album,
	}, nil
}

func (s *PostService) UpdatePhoto(ctx context.Context, id string, input models.PhotoPostInput) (models.PhotoPost, error) {
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

	merged := models.PhotoPostInput{Images: cur.Images, Location: cur.Location, Tags: cur.Tags}
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
		ID: id, Type: models.PostTypePhoto, CreatedAt: cur.CreatedAt, Tags: tags,
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

	merged := models.VideoPostInput{
		Title: cur.Title, ThumbnailURL: cur.ThumbnailURL, ThumbnailTinyURL: cur.ThumbnailTinyURL,
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
	return models.VideoPost{
		ID: id, Type: models.PostTypeVideo, CreatedAt: cur.CreatedAt, Tags: tags,
		Title: merged.Title, ThumbnailURL: merged.ThumbnailURL, ThumbnailTinyURL: merged.ThumbnailTinyURL,
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

	merged := models.LinkPostInput{
		Title: cur.Title, URL: cur.URL, Domain: cur.Domain,
		Description: cur.Description, ThumbnailURL: cur.ThumbnailURL,
		ThumbnailTinyURL: cur.ThumbnailTinyURL,
		Category:         cur.Category, Tags: cur.Tags,
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
		ID: id, Type: models.PostTypeLink, CreatedAt: cur.CreatedAt, Tags: tags,
		Title: merged.Title, URL: merged.URL, Domain: merged.Domain,
		Description: merged.Description, ThumbnailURL: merged.ThumbnailURL,
		ThumbnailTinyURL: merged.ThumbnailTinyURL, Category: merged.Category,
	}, nil
}

func (s *PostService) AddPhotoImage(ctx context.Context, postID string, image models.PhotoImage) (*models.PhotoImage, error) {
	if err := requireNonEmpty("postID", postID); err != nil {
		return nil, err
	}
	if err := requireNonEmpty("image url", image.URL); err != nil {
		return nil, err
	}
	if err := requireNonEmpty("image alt", image.Alt); err != nil {
		return nil, err
	}
	post, err := s.store.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if _, ok := post.(models.PhotoPost); !ok {
		return nil, validationErr("post is not a photo post")
	}
	return s.store.AddPhotoImage(ctx, postID, image)
}

func (s *PostService) GetPhotoImage(ctx context.Context, postID string, imageIDStr string) (*models.PhotoImage, error) {
	if err := requireNonEmpty("postID", postID); err != nil {
		return nil, err
	}
	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil {
		return nil, validationErr("imageID must be an integer")
	}
	// Returns nil, nil if not found — API layer maps this to 404.
	return s.store.GetPhotoImage(ctx, postID, imageID)
}

func (s *PostService) UpdatePhotoImage(ctx context.Context, postID string, imageIDStr string, input models.UpdatePhotoImage) (*models.PhotoImage, error) {
	if err := requireNonEmpty("postID", postID); err != nil {
		return nil, err
	}
	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil {
		return nil, validationErr("imageID must be an integer")
	}
	post, err := s.store.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if _, ok := post.(models.PhotoPost); !ok {
		return nil, validationErr("post is not a photo post")
	}
	// Returns nil, nil if image not found — API layer maps this to 404.
	return s.store.UpdatePhotoImage(ctx, postID, imageID, input)
}

func (s *PostService) DeletePhotoImage(ctx context.Context, postID string, imageIDStr string) (notFound bool, err error) {
	if err := requireNonEmpty("postID", postID); err != nil {
		return false, err
	}
	imageID, convErr := strconv.Atoi(imageIDStr)
	if convErr != nil {
		return false, validationErr("imageID must be an integer")
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
	photoPost, ok := post.(models.PhotoPost)
	if !ok {
		return false, validationErr("post is not a photo post")
	}
	if len(photoPost.Images) <= 1 {
		return false, validationErr("cannot delete the only image in a photo post")
	}
	if err := s.store.DeletePhotoImage(ctx, postID, imageID); err != nil {
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

func (s *PostService) DeletePost(ctx context.Context, id string) error {
	if err := requireNonEmpty("id", id); err != nil {
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
