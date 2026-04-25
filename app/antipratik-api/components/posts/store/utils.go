package store

import (
	"context"
	"fmt"

	"github.com/pratikluitel/antipratik/components/posts"
)

// ── Per-type fetchers ─────────────────────────────────────────────────────────

type essayData struct {
	title, slug, excerpt, body string
	readingTimeMinutes         int
}

func (s *sqlitePostStore) fetchEssayData(ctx context.Context, ids []string) (map[string]essayData, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	q := "SELECT post_id, title, slug, excerpt, body, reading_time_minutes FROM essay_posts WHERE post_id IN (" + placeholders(len(ids)) + ")"
	rows, err := s.db.QueryContext(ctx, q, stringsToAny(ids)...)
	if err != nil {
		return nil, fmt.Errorf("fetchEssayData: %w", err)
	}
	defer func() { _ = rows.Close() }()

	m := make(map[string]essayData)
	for rows.Next() {
		var id string
		var d essayData
		if err := rows.Scan(&id, &d.title, &d.slug, &d.excerpt, &d.body, &d.readingTimeMinutes); err != nil {
			return nil, fmt.Errorf("fetchEssayData scan: %w", err)
		}
		m[id] = d
	}
	return m, rows.Err()
}

type shortData struct{ body string }

func (s *sqlitePostStore) fetchShortData(ctx context.Context, ids []string) (map[string]shortData, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	q := "SELECT post_id, body FROM short_posts WHERE post_id IN (" + placeholders(len(ids)) + ")"
	rows, err := s.db.QueryContext(ctx, q, stringsToAny(ids)...)
	if err != nil {
		return nil, fmt.Errorf("fetchShortData: %w", err)
	}
	defer func() { _ = rows.Close() }()

	m := make(map[string]shortData)
	for rows.Next() {
		var id string
		var d shortData
		if err := rows.Scan(&id, &d.body); err != nil {
			return nil, fmt.Errorf("fetchShortData scan: %w", err)
		}
		m[id] = d
	}
	return m, rows.Err()
}

type musicData struct {
	albumArtTiny  *string
	albumArtSmall *string
	albumArtMed   *string
	albumArtLarge *string
	album         *string
	title         string
	albumArt      string
	audioURL      string
	duration      int
}

func (s *sqlitePostStore) fetchMusicData(ctx context.Context, ids []string) (map[string]musicData, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	q := "SELECT post_id, title, album_art, album_art_tiny, album_art_small, album_art_medium, album_art_large, audio_url, duration, album FROM music_posts WHERE post_id IN (" + placeholders(len(ids)) + ")"
	rows, err := s.db.QueryContext(ctx, q, stringsToAny(ids)...)
	if err != nil {
		return nil, fmt.Errorf("fetchMusicData: %w", err)
	}
	defer func() { _ = rows.Close() }()

	m := make(map[string]musicData)
	for rows.Next() {
		var id string
		var d musicData
		if err := rows.Scan(&id, &d.title, &d.albumArt, &d.albumArtTiny, &d.albumArtSmall, &d.albumArtMed, &d.albumArtLarge, &d.audioURL, &d.duration, &d.album); err != nil {
			return nil, fmt.Errorf("fetchMusicData scan: %w", err)
		}
		m[id] = d
	}
	return m, rows.Err()
}

type photoMeta struct{ location *string }

func (s *sqlitePostStore) fetchPhotoMeta(ctx context.Context, ids []string) (map[string]photoMeta, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	q := "SELECT post_id, location FROM photo_posts WHERE post_id IN (" + placeholders(len(ids)) + ")"
	rows, err := s.db.QueryContext(ctx, q, stringsToAny(ids)...)
	if err != nil {
		return nil, fmt.Errorf("fetchPhotoMeta: %w", err)
	}
	defer func() { _ = rows.Close() }()

	m := make(map[string]photoMeta)
	for rows.Next() {
		var id string
		var d photoMeta
		if err := rows.Scan(&id, &d.location); err != nil {
			return nil, fmt.Errorf("fetchPhotoMeta scan: %w", err)
		}
		m[id] = d
	}
	return m, rows.Err()
}

func (s *sqlitePostStore) fetchPhotoImages(ctx context.Context, ids []string) (map[string][]posts.PhotoImage, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	q := "SELECT id, post_id, url, alt, caption, thumbnail_tiny_url, thumbnail_small_url, thumbnail_medium_url, thumbnail_large_url FROM photo_images WHERE post_id IN (" + placeholders(len(ids)) + ") ORDER BY post_id, sort_order"
	rows, err := s.db.QueryContext(ctx, q, stringsToAny(ids)...)
	if err != nil {
		return nil, fmt.Errorf("fetchPhotoImages: %w", err)
	}
	defer func() { _ = rows.Close() }()

	m := make(map[string][]posts.PhotoImage)
	for rows.Next() {
		var postID string
		var img posts.PhotoImage
		if err := rows.Scan(&img.ID, &postID, &img.URL, &img.Alt, &img.Caption, &img.ThumbnailTinyURL, &img.ThumbnailSmallURL, &img.ThumbnailMedURL, &img.ThumbnailLargeURL); err != nil {
			return nil, fmt.Errorf("fetchPhotoImages scan: %w", err)
		}
		m[postID] = append(m[postID], img)
	}
	return m, rows.Err()
}

type videoData struct {
	description       *string
	thumbnailURL      *string
	thumbnailTinyURL  *string
	thumbnailSmallURL *string
	thumbnailMedURL   *string
	thumbnailLargeURL *string
	title             string
	videoURL          string
}

func (s *sqlitePostStore) fetchVideoData(ctx context.Context, ids []string) (map[string]videoData, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	q := "SELECT post_id, title, description, thumbnail_url, thumbnail_tiny_url, thumbnail_small_url, thumbnail_medium_url, thumbnail_large_url, video_url FROM video_posts WHERE post_id IN (" + placeholders(len(ids)) + ")"
	rows, err := s.db.QueryContext(ctx, q, stringsToAny(ids)...)
	if err != nil {
		return nil, fmt.Errorf("fetchVideoData: %w", err)
	}
	defer func() { _ = rows.Close() }()

	m := make(map[string]videoData)
	for rows.Next() {
		var id string
		var d videoData
		if err := rows.Scan(&id, &d.title, &d.description, &d.thumbnailURL, &d.thumbnailTinyURL, &d.thumbnailSmallURL, &d.thumbnailMedURL, &d.thumbnailLargeURL, &d.videoURL); err != nil {
			return nil, fmt.Errorf("fetchVideoData scan: %w", err)
		}
		m[id] = d
	}
	return m, rows.Err()
}

type linkPostData struct {
	description       *string
	thumbnailURL      *string
	thumbnailTinyURL  *string
	thumbnailSmallURL *string
	thumbnailMedURL   *string
	thumbnailLargeURL *string
	category          *string
	title             string
	url               string
	domain            string
}

func (s *sqlitePostStore) fetchLinkPostData(ctx context.Context, ids []string) (map[string]linkPostData, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	q := "SELECT post_id, title, url, domain, description, thumbnail_url, thumbnail_tiny_url, thumbnail_small_url, thumbnail_medium_url, thumbnail_large_url, category FROM link_posts WHERE post_id IN (" + placeholders(len(ids)) + ")"
	rows, err := s.db.QueryContext(ctx, q, stringsToAny(ids)...)
	if err != nil {
		return nil, fmt.Errorf("fetchLinkPostData: %w", err)
	}
	defer func() { _ = rows.Close() }()

	m := make(map[string]linkPostData)
	for rows.Next() {
		var id string
		var d linkPostData
		if err := rows.Scan(&id, &d.title, &d.url, &d.domain, &d.description, &d.thumbnailURL, &d.thumbnailTinyURL, &d.thumbnailSmallURL, &d.thumbnailMedURL, &d.thumbnailLargeURL, &d.category); err != nil {
			return nil, fmt.Errorf("fetchLinkPostData scan: %w", err)
		}
		m[id] = d
	}
	return m, rows.Err()
}

// ── Assembly ──────────────────────────────────────────────────────────────────

func (s *sqlitePostStore) assembleAll(
	ctx context.Context,
	baseRows []baseRow,
	byType map[string][]string,
	tagsMap map[string][]string,
) ([]posts.Post, error) {
	essayMap, err := s.fetchEssayData(ctx, byType["essay"])
	if err != nil {
		return nil, err
	}
	shortMap, err := s.fetchShortData(ctx, byType["short"])
	if err != nil {
		return nil, err
	}
	musicMap, err := s.fetchMusicData(ctx, byType["music"])
	if err != nil {
		return nil, err
	}
	photoMeta, err := s.fetchPhotoMeta(ctx, byType["photo"])
	if err != nil {
		return nil, err
	}
	photoImages, err := s.fetchPhotoImages(ctx, byType["photo"])
	if err != nil {
		return nil, err
	}
	videoMap, err := s.fetchVideoData(ctx, byType["video"])
	if err != nil {
		return nil, err
	}
	linkPostMap, err := s.fetchLinkPostData(ctx, byType["link"])
	if err != nil {
		return nil, err
	}

	result := make([]posts.Post, 0, len(baseRows))
	for _, r := range baseRows {
		tags := coalesceStringSlice(tagsMap[r.ID])
		var post posts.Post

		switch r.Type {
		case "essay":
			d := essayMap[r.ID]
			post = posts.EssayPost{
				ID: r.ID, Type: r.Type, CreatedAt: r.CreatedAt, Tags: tags,
				Title: d.title, Slug: d.slug, Excerpt: d.excerpt, Body: d.body,
				ReadingTimeMinutes: d.readingTimeMinutes,
			}
		case "short":
			d := shortMap[r.ID]
			post = posts.ShortPost{
				ID: r.ID, Type: r.Type, CreatedAt: r.CreatedAt, Tags: tags,
				Body: d.body,
			}
		case "music":
			d := musicMap[r.ID]
			post = posts.MusicPost{
				ID: r.ID, Type: r.Type, CreatedAt: r.CreatedAt, Tags: tags,
				Title: d.title, AlbumArt: d.albumArt, AlbumArtTinyURL: d.albumArtTiny,
				AlbumArtSmallURL: d.albumArtSmall, AlbumArtMedURL: d.albumArtMed,
				AlbumArtLargeURL: d.albumArtLarge,
				AudioURL:         d.audioURL, Duration: d.duration, Album: d.album,
			}
		case "photo":
			meta := photoMeta[r.ID]
			imgs := photoImages[r.ID]
			if imgs == nil {
				imgs = []posts.PhotoImage{}
			}
			post = posts.PhotoPost{
				ID: r.ID, Type: r.Type, CreatedAt: r.CreatedAt, Tags: tags,
				Images: imgs, Location: meta.location,
			}
		case "video":
			d := videoMap[r.ID]
			post = posts.VideoPost{
				ID: r.ID, Type: r.Type, CreatedAt: r.CreatedAt, Tags: tags,
				Title: d.title, Description: d.description,
				ThumbnailURL: d.thumbnailURL, ThumbnailTinyURL: d.thumbnailTinyURL,
				ThumbnailSmallURL: d.thumbnailSmallURL, ThumbnailMedURL: d.thumbnailMedURL,
				ThumbnailLargeURL: d.thumbnailLargeURL,
				VideoURL:          d.videoURL,
			}
		case "link":
			d := linkPostMap[r.ID]
			post = posts.LinkPost{
				ID: r.ID, Type: r.Type, CreatedAt: r.CreatedAt, Tags: tags,
				Title: d.title, URL: d.url, Domain: d.domain,
				Description: d.description, ThumbnailURL: d.thumbnailURL,
				ThumbnailTinyURL: d.thumbnailTinyURL, ThumbnailSmallURL: d.thumbnailSmallURL,
				ThumbnailMedURL: d.thumbnailMedURL, ThumbnailLargeURL: d.thumbnailLargeURL,
				Category: d.category,
			}
		default:
			continue
		}
		result = append(result, post)
	}
	return result, nil
}
