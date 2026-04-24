package logic

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/pratikluitel/antipratik/components/posts"
)

func newID() string  { return uuid.New().String() }
func nowUTC() string { return time.Now().UTC().Format(time.RFC3339) }

// fileKeysForPost returns all storage keys for files attached to a post.
func fileKeysForPost(post posts.Post) []string {
	var keys []string
	switch p := post.(type) {
	case posts.MusicPost:
		if p.AudioURL != "" {
			keys = append(keys, urlToStorageKey(p.AudioURL))
		}
		if p.AlbumArt != "" {
			keys = append(keys, urlToStorageKey(p.AlbumArt))
		}
	case posts.PhotoPost:
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
	case posts.VideoPost:
		if p.ThumbnailURL != nil && *p.ThumbnailURL != "" {
			keys = append(keys, urlToStorageKey(*p.ThumbnailURL))
		}
	case posts.LinkPost:
		if p.ThumbnailURL != nil && *p.ThumbnailURL != "" {
			keys = append(keys, urlToStorageKey(*p.ThumbnailURL))
		}
	}
	return keys
}

// fileKeysForImage returns all storage keys for a single PhotoImage.
func fileKeysForImage(img *posts.PhotoImage) []string {
	if img == nil {
		return nil
	}
	var keys []string
	if img.URL != "" {
		keys = append(keys, urlToStorageKey(img.URL))
	}
	if img.ThumbnailTinyURL != nil {
		keys = append(keys, urlToStorageKey(*img.ThumbnailTinyURL))
	}
	if img.ThumbnailSmallURL != nil {
		keys = append(keys, urlToStorageKey(*img.ThumbnailSmallURL))
	}
	if img.ThumbnailMedURL != nil {
		keys = append(keys, urlToStorageKey(*img.ThumbnailMedURL))
	}
	if img.ThumbnailLargeURL != nil {
		keys = append(keys, urlToStorageKey(*img.ThumbnailLargeURL))
	}
	return keys
}

// albumArtFileKeys returns all storage keys for a MusicPost's album art (original + 4 sizes).
func albumArtFileKeys(p posts.MusicPost) []string {
	var keys []string
	if p.AlbumArt != "" {
		keys = append(keys, urlToStorageKey(p.AlbumArt))
	}
	for _, u := range []*string{p.AlbumArtTinyURL, p.AlbumArtSmallURL, p.AlbumArtMedURL, p.AlbumArtLargeURL} {
		if u != nil && *u != "" {
			keys = append(keys, urlToStorageKey(*u))
		}
	}
	return keys
}

// thumbnailFileKeys returns all storage keys for a thumbnail set (original + 4 sizes).
func thumbnailFileKeys(url string, tiny, small, med, large *string) []string {
	var keys []string
	if url != "" {
		keys = append(keys, urlToStorageKey(url))
	}
	for _, u := range []*string{tiny, small, med, large} {
		if u != nil && *u != "" {
			keys = append(keys, urlToStorageKey(*u))
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
