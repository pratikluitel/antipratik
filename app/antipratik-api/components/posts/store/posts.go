package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pratikluitel/antipratik/components/posts"
)

// sqlitePostStore implements PostStore using a SQLite database.
type sqlitePostStore struct {
	db *sql.DB
}

// NewPostStore creates a new sqlitePostStore backed by db.
func NewPostStore(db *sql.DB) posts.PostStore {
	return &sqlitePostStore{db: db}
}

// ── Write methods ─────────────────────────────────────────────────────────────

func (s *sqlitePostStore) CreatePost(ctx context.Context, postType string, id string, createdAt string) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO posts (id, type, created_at) VALUES (?, ?, ?)`, id, postType, createdAt)
	return err
}

func (s *sqlitePostStore) insertTags(ctx context.Context, tx *sql.Tx, id string, tags []string) error {
	for _, tag := range tags {
		// Upsert tag into the normalized tags table.
		if _, err := tx.ExecContext(ctx, `INSERT OR IGNORE INTO tags (name) VALUES (?)`, tag); err != nil {
			return err
		}
		var tagID int64
		if err := tx.QueryRowContext(ctx, `SELECT id FROM tags WHERE name = ?`, tag).Scan(&tagID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `INSERT OR IGNORE INTO post_tags (post_id, tag_id) VALUES (?, ?)`, id, tagID); err != nil {
			return err
		}
	}
	return nil
}

func (s *sqlitePostStore) CreateEssayData(ctx context.Context, id string, input posts.EssayPostInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	_, err = tx.ExecContext(ctx, `INSERT INTO essay_posts (post_id, title, slug, excerpt, body, reading_time_minutes) VALUES (?, ?, ?, ?, ?, ?)`,
		id, input.Title, input.Slug, input.Excerpt, input.Body, input.ReadingTimeMinutes)
	if err != nil {
		return err
	}
	if err := s.insertTags(ctx, tx, id, input.Tags); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *sqlitePostStore) CreateShortData(ctx context.Context, id string, input posts.ShortPostInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	_, err = tx.ExecContext(ctx, `INSERT INTO short_posts (post_id, body) VALUES (?, ?)`, id, input.Body)
	if err != nil {
		return err
	}
	if err := s.insertTags(ctx, tx, id, input.Tags); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *sqlitePostStore) CreateMusicData(ctx context.Context, id string, input posts.MusicPostInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	_, err = tx.ExecContext(ctx, `INSERT INTO music_posts (post_id, title, album_art, album_art_tiny, album_art_small, album_art_medium, album_art_large, audio_url, duration, album) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, input.Title, input.AlbumArt, input.AlbumArtTinyURL, input.AlbumArtSmallURL, input.AlbumArtMedURL, input.AlbumArtLargeURL, input.AudioURL, input.Duration, input.Album)
	if err != nil {
		return err
	}
	if err := s.insertTags(ctx, tx, id, input.Tags); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *sqlitePostStore) CreatePhotoData(ctx context.Context, id string, input posts.PhotoPostInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	_, err = tx.ExecContext(ctx, `INSERT INTO photo_posts (post_id, location) VALUES (?, ?)`, id, input.Location)
	if err != nil {
		return err
	}
	for i, img := range input.Images {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO photo_images (post_id, url, alt, caption, sort_order, thumbnail_tiny_url, thumbnail_small_url, thumbnail_medium_url, thumbnail_large_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			id, img.URL, img.Alt, img.Caption, i, img.ThumbnailTinyURL, img.ThumbnailSmallURL, img.ThumbnailMedURL, img.ThumbnailLargeURL)
		if err != nil {
			return err
		}
	}
	if err := s.insertTags(ctx, tx, id, input.Tags); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *sqlitePostStore) CreateVideoData(ctx context.Context, id string, input posts.VideoPostInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	_, err = tx.ExecContext(ctx, `INSERT INTO video_posts (post_id, title, thumbnail_url, thumbnail_tiny_url, thumbnail_small_url, thumbnail_medium_url, thumbnail_large_url, video_url, duration, playlist) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, input.Title, input.ThumbnailURL, input.ThumbnailTinyURL, input.ThumbnailSmallURL, input.ThumbnailMedURL, input.ThumbnailLargeURL, input.VideoURL, input.Duration, input.Playlist)
	if err != nil {
		return err
	}
	if err := s.insertTags(ctx, tx, id, input.Tags); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *sqlitePostStore) CreateLinkPostData(ctx context.Context, id string, input posts.LinkPostInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	_, err = tx.ExecContext(ctx, `INSERT INTO link_posts (post_id, title, url, domain, description, thumbnail_url, thumbnail_tiny_url, thumbnail_small_url, thumbnail_medium_url, thumbnail_large_url, category) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, input.Title, input.URL, input.Domain, input.Description, input.ThumbnailURL, input.ThumbnailTinyURL, input.ThumbnailSmallURL, input.ThumbnailMedURL, input.ThumbnailLargeURL, input.Category)
	if err != nil {
		return err
	}
	if err := s.insertTags(ctx, tx, id, input.Tags); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *sqlitePostStore) UpdateEssay(ctx context.Context, id string, input posts.EssayPostInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	_, err = tx.ExecContext(ctx, `UPDATE essay_posts SET title=?, slug=?, excerpt=?, body=?, reading_time_minutes=? WHERE post_id=?`,
		input.Title, input.Slug, input.Excerpt, input.Body, input.ReadingTimeMinutes, id)
	if err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM post_tags WHERE post_id=?`, id); err != nil {
		return err
	}
	if err := s.insertTags(ctx, tx, id, input.Tags); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *sqlitePostStore) UpdateShort(ctx context.Context, id string, input posts.ShortPostInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err = tx.ExecContext(ctx, `UPDATE short_posts SET body=? WHERE post_id=?`, input.Body, id); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM post_tags WHERE post_id=?`, id); err != nil {
		return err
	}
	if err := s.insertTags(ctx, tx, id, input.Tags); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *sqlitePostStore) UpdateMusic(ctx context.Context, id string, input posts.MusicPostInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err = tx.ExecContext(ctx, `UPDATE music_posts SET title=?, album_art=?, album_art_tiny=?, album_art_small=?, album_art_medium=?, album_art_large=?, audio_url=?, duration=?, album=? WHERE post_id=?`,
		input.Title, input.AlbumArt, input.AlbumArtTinyURL, input.AlbumArtSmallURL, input.AlbumArtMedURL, input.AlbumArtLargeURL, input.AudioURL, input.Duration, input.Album, id); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM post_tags WHERE post_id=?`, id); err != nil {
		return err
	}
	if err := s.insertTags(ctx, tx, id, input.Tags); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *sqlitePostStore) UpdatePhoto(ctx context.Context, id string, input posts.PhotoPostInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err = tx.ExecContext(ctx, `UPDATE photo_posts SET location=? WHERE post_id=?`, input.Location, id); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM photo_images WHERE post_id=?`, id); err != nil {
		return err
	}
	for i, img := range input.Images {
		if _, err = tx.ExecContext(ctx,
			`INSERT INTO photo_images (post_id, url, alt, caption, sort_order, thumbnail_tiny_url, thumbnail_small_url, thumbnail_medium_url, thumbnail_large_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			id, img.URL, img.Alt, img.Caption, i, img.ThumbnailTinyURL, img.ThumbnailSmallURL, img.ThumbnailMedURL, img.ThumbnailLargeURL); err != nil {
			return err
		}
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM post_tags WHERE post_id=?`, id); err != nil {
		return err
	}
	if err := s.insertTags(ctx, tx, id, input.Tags); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *sqlitePostStore) UpdateVideo(ctx context.Context, id string, input posts.VideoPostInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err = tx.ExecContext(ctx, `UPDATE video_posts SET title=?, thumbnail_url=?, thumbnail_tiny_url=?, thumbnail_small_url=?, thumbnail_medium_url=?, thumbnail_large_url=?, video_url=?, duration=?, playlist=? WHERE post_id=?`,
		input.Title, input.ThumbnailURL, input.ThumbnailTinyURL, input.ThumbnailSmallURL, input.ThumbnailMedURL, input.ThumbnailLargeURL, input.VideoURL, input.Duration, input.Playlist, id); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM post_tags WHERE post_id=?`, id); err != nil {
		return err
	}
	if err := s.insertTags(ctx, tx, id, input.Tags); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *sqlitePostStore) UpdateLinkPost(ctx context.Context, id string, input posts.LinkPostInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err = tx.ExecContext(ctx, `UPDATE link_posts SET title=?, url=?, domain=?, description=?, thumbnail_url=?, thumbnail_tiny_url=?, thumbnail_small_url=?, thumbnail_medium_url=?, thumbnail_large_url=?, category=? WHERE post_id=?`,
		input.Title, input.URL, input.Domain, input.Description, input.ThumbnailURL, input.ThumbnailTinyURL, input.ThumbnailSmallURL, input.ThumbnailMedURL, input.ThumbnailLargeURL, input.Category, id); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM post_tags WHERE post_id=?`, id); err != nil {
		return err
	}
	if err := s.insertTags(ctx, tx, id, input.Tags); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *sqlitePostStore) DeletePost(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM posts WHERE id=?`, id)
	return err
}

func (s *sqlitePostStore) AddPhotoImage(ctx context.Context, postID string, image posts.PhotoImage) (*posts.PhotoImage, error) {
	// Determine the next sort_order for this post.
	var maxOrder int
	row := s.db.QueryRowContext(ctx, `SELECT COALESCE(MAX(sort_order), -1) FROM photo_images WHERE post_id=?`, postID)
	if err := row.Scan(&maxOrder); err != nil {
		return nil, fmt.Errorf("AddPhotoImage max sort_order: %w", err)
	}
	result, err := s.db.ExecContext(ctx,
		`INSERT INTO photo_images (post_id, url, alt, caption, sort_order, thumbnail_tiny_url, thumbnail_small_url, thumbnail_medium_url, thumbnail_large_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		postID, image.URL, image.Alt, image.Caption, maxOrder+1, image.ThumbnailTinyURL, image.ThumbnailSmallURL, image.ThumbnailMedURL, image.ThumbnailLargeURL,
	)
	if err != nil {
		return nil, fmt.Errorf("AddPhotoImage insert: %w", err)
	}
	newID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("AddPhotoImage last insert id: %w", err)
	}
	image.ID = int(newID)
	return &image, nil
}

func (s *sqlitePostStore) GetPhotoImage(ctx context.Context, postID string, imageID int) (*posts.PhotoImage, error) {
	var img posts.PhotoImage
	err := s.db.QueryRowContext(ctx,
		`SELECT id, url, alt, caption, thumbnail_tiny_url, thumbnail_small_url, thumbnail_medium_url, thumbnail_large_url FROM photo_images WHERE id=? AND post_id=?`,
		imageID, postID,
	).Scan(&img.ID, &img.URL, &img.Alt, &img.Caption, &img.ThumbnailTinyURL, &img.ThumbnailSmallURL, &img.ThumbnailMedURL, &img.ThumbnailLargeURL)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetPhotoImage: %w", err)
	}
	return &img, nil
}

func (s *sqlitePostStore) UpdatePhotoImage(ctx context.Context, postID string, imageID int, input posts.UpdatePhotoImage) (*posts.PhotoImage, error) {
	// Build a dynamic SET clause for only the non-nil fields.
	setClauses := []string{}
	args := []any{}
	if input.Caption != nil {
		setClauses = append(setClauses, "caption=?")
		args = append(args, *input.Caption)
	}
	if input.Alt != nil {
		setClauses = append(setClauses, "alt=?")
		args = append(args, *input.Alt)
	}
	if len(setClauses) > 0 {
		q := "UPDATE photo_images SET " + strings.Join(setClauses, ", ") + " WHERE id=? AND post_id=?"
		args = append(args, imageID, postID)
		if _, err := s.db.ExecContext(ctx, q, args...); err != nil {
			return nil, fmt.Errorf("UpdatePhotoImage: %w", err)
		}
	}
	return s.GetPhotoImage(ctx, postID, imageID)
}

func (s *sqlitePostStore) DeletePhotoImage(ctx context.Context, postID string, imageID int) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM photo_images WHERE id=? AND post_id=?`, imageID, postID)
	return err
}

// ── Public methods ────────────────────────────────────────────────────────────

// GetPosts returns all posts matching the optional type and tag filters.
func (s *sqlitePostStore) GetPosts(ctx context.Context, types, tags []string) ([]posts.Post, error) {
	baseRows, err := s.queryBaseRows(ctx, types, tags)
	if err != nil {
		return nil, err
	}
	if len(baseRows) == 0 {
		return []posts.Post{}, nil
	}

	ids := extractIDs(baseRows)
	byType := groupByType(baseRows)

	tagsMap, err := s.fetchTagsMap(ctx, ids)
	if err != nil {
		return nil, err
	}

	return s.assembleAll(ctx, baseRows, byType, tagsMap)
}

// GetPostByID returns any post type by ID, or an error if not found.
func (s *sqlitePostStore) GetPostByID(ctx context.Context, id string) (posts.Post, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, type, created_at FROM posts WHERE id = ?`, id)
	var r baseRow
	if err := row.Scan(&r.ID, &r.Type, &r.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post %q not found", id)
		}
		return nil, fmt.Errorf("GetPostByID: %w", err)
	}
	baseRows := []baseRow{r}
	byType := groupByType(baseRows)
	tagsMap, err := s.fetchTagsMap(ctx, []string{r.ID})
	if err != nil {
		return nil, err
	}
	posts, err := s.assembleAll(ctx, baseRows, byType, tagsMap)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return nil, fmt.Errorf("post %q not found", id)
	}
	return posts[0], nil
}

// GetPostsByIDs returns posts for each given ID, in the order of ids.
// IDs not found in the database are silently skipped.
func (s *sqlitePostStore) GetPostsByIDs(ctx context.Context, ids []string) ([]posts.Post, error) {
	var out []posts.Post
	for _, id := range ids {
		p, err := s.GetPostByID(ctx, id)
		if err != nil {
			// Not found — skip silently.
			continue
		}
		out = append(out, p)
	}
	return out, nil
}

// GetPostBySlug returns the essay with the given slug, or nil if not found.
func (s *sqlitePostStore) GetPostBySlug(ctx context.Context, slug string) (*posts.EssayPost, error) {
	const q = `
		SELECT p.id, p.created_at,
		       e.title, e.slug, e.excerpt, e.body, e.reading_time_minutes
		FROM posts p
		JOIN essay_posts e ON p.id = e.post_id
		WHERE e.slug = ?`

	row := s.db.QueryRowContext(ctx, q, slug)
	var (
		id, createdAt, title, slugVal, excerpt, body string
		readingTime                                  int
	)
	if err := row.Scan(&id, &createdAt, &title, &slugVal, &excerpt, &body, &readingTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetPostBySlug: %w", err)
	}

	tags, err := s.fetchTagsMap(ctx, []string{id})
	if err != nil {
		return nil, err
	}

	return &posts.EssayPost{
		ID:                 id,
		Type:               "essay",
		CreatedAt:          createdAt,
		Tags:               coalesceStringSlice(tags[id]),
		Title:              title,
		Slug:               slugVal,
		Excerpt:            excerpt,
		Body:               body,
		ReadingTimeMinutes: readingTime,
	}, nil
}

// ── Internal query helpers ────────────────────────────────────────────────────

type baseRow struct {
	ID        string
	Type      string
	CreatedAt string
}

// queryBaseRows fetches (id, type, created_at) rows matching filters.
func (s *sqlitePostStore) queryBaseRows(ctx context.Context, types, tags []string) ([]baseRow, error) {
	var sb strings.Builder
	var args []any

	if len(tags) > 0 {
		sb.WriteString("SELECT DISTINCT p.id, p.type, p.created_at FROM posts p JOIN post_tags pt ON p.id = pt.post_id JOIN tags t ON t.id = pt.tag_id")
	} else {
		sb.WriteString("SELECT id, type, created_at FROM posts")
	}

	var conditions []string
	if len(types) > 0 {
		if len(tags) > 0 {
			conditions = append(conditions, "p.type IN ("+placeholders(len(types))+")")
		} else {
			conditions = append(conditions, "type IN ("+placeholders(len(types))+")")
		}
		for _, t := range types {
			args = append(args, t)
		}
	}
	if len(tags) > 0 {
		conditions = append(conditions, "t.name IN ("+placeholders(len(tags))+")")
		for _, t := range tags {
			args = append(args, t)
		}
	}

	if len(conditions) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(conditions, " AND "))
	}
	sb.WriteString(" ORDER BY created_at DESC")

	rows, err := s.db.QueryContext(ctx, sb.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("queryBaseRows: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var result []baseRow
	for rows.Next() {
		var r baseRow
		if err := rows.Scan(&r.ID, &r.Type, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("queryBaseRows scan: %w", err)
		}
		result = append(result, r)
	}
	return result, rows.Err()
}

// fetchTagsMap returns a map of post_id → []tag for the given post IDs.
func (s *sqlitePostStore) fetchTagsMap(ctx context.Context, ids []string) (map[string][]string, error) {
	if len(ids) == 0 {
		return map[string][]string{}, nil
	}
	q := "SELECT pt.post_id, t.name FROM post_tags pt JOIN tags t ON t.id = pt.tag_id WHERE pt.post_id IN (" + placeholders(len(ids)) + ") ORDER BY pt.post_id"
	args := stringsToAny(ids)
	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("fetchTagsMap: %w", err)
	}
	defer func() { _ = rows.Close() }()

	result := make(map[string][]string)
	for rows.Next() {
		var postID, tag string
		if err := rows.Scan(&postID, &tag); err != nil {
			return nil, fmt.Errorf("fetchTagsMap scan: %w", err)
		}
		result[postID] = append(result[postID], tag)
	}
	return result, rows.Err()
}

// GetAllTags returns all tag names sorted alphabetically.
func (s *sqlitePostStore) GetAllTags(ctx context.Context) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT name FROM tags ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("GetAllTags: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var result []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("GetAllTags scan: %w", err)
		}
		result = append(result, name)
	}
	if result == nil {
		result = []string{}
	}
	return result, rows.Err()
}
