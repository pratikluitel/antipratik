package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pratikluitel/antipratik/models"
)

// ── Write methods ─────────────────────────────────────────────────────────────

func (s *SQLitePostStore) CreatePost(ctx context.Context, postType string, id string, createdAt string) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO posts (id, type, created_at) VALUES (?, ?, ?)`, id, postType, createdAt)
	return err
}

func (s *SQLitePostStore) insertTags(ctx context.Context, tx *sql.Tx, id string, tags []string) error {
	for _, tag := range tags {
		_, err := tx.ExecContext(ctx, `INSERT OR IGNORE INTO post_tags (post_id, tag) VALUES (?, ?)`, id, tag)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SQLitePostStore) CreateEssayData(ctx context.Context, id string, input models.EssayPostInput) error {
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

func (s *SQLitePostStore) CreateShortData(ctx context.Context, id string, input models.ShortPostInput) error {
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

func (s *SQLitePostStore) CreateMusicData(ctx context.Context, id string, input models.MusicPostInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	_, err = tx.ExecContext(ctx, `INSERT INTO music_posts (post_id, title, album_art, album_art_tiny, audio_url, duration, album) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, input.Title, input.AlbumArt, input.AlbumArtTinyURL, input.AudioURL, input.Duration, input.Album)
	if err != nil {
		return err
	}
	if err := s.insertTags(ctx, tx, id, input.Tags); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *SQLitePostStore) CreatePhotoData(ctx context.Context, id string, input models.PhotoPostInput) error {
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

func (s *SQLitePostStore) CreateVideoData(ctx context.Context, id string, input models.VideoPostInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	_, err = tx.ExecContext(ctx, `INSERT INTO video_posts (post_id, title, thumbnail_url, thumbnail_tiny_url, video_url, duration, playlist) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, input.Title, input.ThumbnailURL, input.ThumbnailTinyURL, input.VideoURL, input.Duration, input.Playlist)
	if err != nil {
		return err
	}
	if err := s.insertTags(ctx, tx, id, input.Tags); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *SQLitePostStore) CreateLinkPostData(ctx context.Context, id string, input models.LinkPostInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	_, err = tx.ExecContext(ctx, `INSERT INTO link_posts (post_id, title, url, domain, description, thumbnail_url, thumbnail_tiny_url, category) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id, input.Title, input.URL, input.Domain, input.Description, input.ThumbnailURL, input.ThumbnailTinyURL, input.Category)
	if err != nil {
		return err
	}
	if err := s.insertTags(ctx, tx, id, input.Tags); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *SQLitePostStore) UpdateEssay(ctx context.Context, id string, input models.EssayPostInput) error {
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

func (s *SQLitePostStore) UpdateShort(ctx context.Context, id string, input models.ShortPostInput) error {
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

func (s *SQLitePostStore) UpdateMusic(ctx context.Context, id string, input models.MusicPostInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err = tx.ExecContext(ctx, `UPDATE music_posts SET title=?, album_art=?, album_art_tiny=?, audio_url=?, duration=?, album=? WHERE post_id=?`,
		input.Title, input.AlbumArt, input.AlbumArtTinyURL, input.AudioURL, input.Duration, input.Album, id); err != nil {
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

func (s *SQLitePostStore) UpdatePhoto(ctx context.Context, id string, input models.PhotoPostInput) error {
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

func (s *SQLitePostStore) UpdateVideo(ctx context.Context, id string, input models.VideoPostInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err = tx.ExecContext(ctx, `UPDATE video_posts SET title=?, thumbnail_url=?, thumbnail_tiny_url=?, video_url=?, duration=?, playlist=? WHERE post_id=?`,
		input.Title, input.ThumbnailURL, input.ThumbnailTinyURL, input.VideoURL, input.Duration, input.Playlist, id); err != nil {
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

func (s *SQLitePostStore) UpdateLinkPost(ctx context.Context, id string, input models.LinkPostInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err = tx.ExecContext(ctx, `UPDATE link_posts SET title=?, url=?, domain=?, description=?, thumbnail_url=?, thumbnail_tiny_url=?, category=? WHERE post_id=?`,
		input.Title, input.URL, input.Domain, input.Description, input.ThumbnailURL, input.ThumbnailTinyURL, input.Category, id); err != nil {
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

func (s *SQLitePostStore) DeletePost(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM posts WHERE id=?`, id)
	return err
}

func (s *SQLitePostStore) AddPhotoImage(ctx context.Context, postID string, image models.PhotoImage) (*models.PhotoImage, error) {
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

func (s *SQLitePostStore) GetPhotoImage(ctx context.Context, postID string, imageID int) (*models.PhotoImage, error) {
	var img models.PhotoImage
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

func (s *SQLitePostStore) UpdatePhotoImage(ctx context.Context, postID string, imageID int, input models.UpdatePhotoImage) (*models.PhotoImage, error) {
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

func (s *SQLitePostStore) DeletePhotoImage(ctx context.Context, postID string, imageID int) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM photo_images WHERE id=? AND post_id=?`, imageID, postID)
	return err
}

// SQLitePostStore implements PostStore using a SQLite database.
type SQLitePostStore struct {
	db *sql.DB
}

// NewPostStore creates a new SQLitePostStore backed by db.
func NewPostStore(db *sql.DB) *SQLitePostStore {
	return &SQLitePostStore{db: db}
}

// ── Public methods ────────────────────────────────────────────────────────────

// GetPosts returns all posts matching the optional type and tag filters.
func (s *SQLitePostStore) GetPosts(ctx context.Context, types, tags []string) ([]models.Post, error) {
	baseRows, err := s.queryBaseRows(ctx, types, tags)
	if err != nil {
		return nil, err
	}
	if len(baseRows) == 0 {
		return []models.Post{}, nil
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
func (s *SQLitePostStore) GetPostByID(ctx context.Context, id string) (models.Post, error) {
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

// GetPostBySlug returns the essay with the given slug, or nil if not found.
func (s *SQLitePostStore) GetPostBySlug(ctx context.Context, slug string) (*models.EssayPost, error) {
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

	return &models.EssayPost{
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
func (s *SQLitePostStore) queryBaseRows(ctx context.Context, types, tags []string) ([]baseRow, error) {
	var sb strings.Builder
	var args []any

	if len(tags) > 0 {
		sb.WriteString("SELECT DISTINCT p.id, p.type, p.created_at FROM posts p JOIN post_tags pt ON p.id = pt.post_id")
	} else {
		sb.WriteString("SELECT id, p_alias.type, p_alias.created_at FROM posts p_alias")
		// rewrite to avoid aliasing issue — use plain
		sb.Reset()
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
		conditions = append(conditions, "pt.tag IN ("+placeholders(len(tags))+")")
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
func (s *SQLitePostStore) fetchTagsMap(ctx context.Context, ids []string) (map[string][]string, error) {
	if len(ids) == 0 {
		return map[string][]string{}, nil
	}
	q := "SELECT post_id, tag FROM post_tags WHERE post_id IN (" + placeholders(len(ids)) + ") ORDER BY post_id"
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

// ── Per-type fetchers ─────────────────────────────────────────────────────────

type essayData struct {
	title, slug, excerpt, body string
	readingTimeMinutes         int
}

func (s *SQLitePostStore) fetchEssayData(ctx context.Context, ids []string) (map[string]essayData, error) {
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

func (s *SQLitePostStore) fetchShortData(ctx context.Context, ids []string) (map[string]shortData, error) {
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
	albumArtTiny *string
	album        *string
	title        string
	albumArt     string
	audioURL     string
	duration     int
}

func (s *SQLitePostStore) fetchMusicData(ctx context.Context, ids []string) (map[string]musicData, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	q := "SELECT post_id, title, album_art, album_art_tiny, audio_url, duration, album FROM music_posts WHERE post_id IN (" + placeholders(len(ids)) + ")"
	rows, err := s.db.QueryContext(ctx, q, stringsToAny(ids)...)
	if err != nil {
		return nil, fmt.Errorf("fetchMusicData: %w", err)
	}
	defer func() { _ = rows.Close() }()

	m := make(map[string]musicData)
	for rows.Next() {
		var id string
		var d musicData
		if err := rows.Scan(&id, &d.title, &d.albumArt, &d.albumArtTiny, &d.audioURL, &d.duration, &d.album); err != nil {
			return nil, fmt.Errorf("fetchMusicData scan: %w", err)
		}
		m[id] = d
	}
	return m, rows.Err()
}

type photoMeta struct{ location *string }

func (s *SQLitePostStore) fetchPhotoMeta(ctx context.Context, ids []string) (map[string]photoMeta, error) {
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

func (s *SQLitePostStore) fetchPhotoImages(ctx context.Context, ids []string) (map[string][]models.PhotoImage, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	q := "SELECT id, post_id, url, alt, caption, thumbnail_tiny_url, thumbnail_small_url, thumbnail_medium_url, thumbnail_large_url FROM photo_images WHERE post_id IN (" + placeholders(len(ids)) + ") ORDER BY post_id, sort_order"
	rows, err := s.db.QueryContext(ctx, q, stringsToAny(ids)...)
	if err != nil {
		return nil, fmt.Errorf("fetchPhotoImages: %w", err)
	}
	defer func() { _ = rows.Close() }()

	m := make(map[string][]models.PhotoImage)
	for rows.Next() {
		var postID string
		var img models.PhotoImage
		if err := rows.Scan(&img.ID, &postID, &img.URL, &img.Alt, &img.Caption, &img.ThumbnailTinyURL, &img.ThumbnailSmallURL, &img.ThumbnailMedURL, &img.ThumbnailLargeURL); err != nil {
			return nil, fmt.Errorf("fetchPhotoImages scan: %w", err)
		}
		m[postID] = append(m[postID], img)
	}
	return m, rows.Err()
}

type videoData struct {
	thumbnailTinyURL *string
	playlist         *string
	title            string
	thumbnailURL     string
	videoURL         string
	duration         int
}

func (s *SQLitePostStore) fetchVideoData(ctx context.Context, ids []string) (map[string]videoData, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	q := "SELECT post_id, title, thumbnail_url, thumbnail_tiny_url, video_url, duration, playlist FROM video_posts WHERE post_id IN (" + placeholders(len(ids)) + ")"
	rows, err := s.db.QueryContext(ctx, q, stringsToAny(ids)...)
	if err != nil {
		return nil, fmt.Errorf("fetchVideoData: %w", err)
	}
	defer func() { _ = rows.Close() }()

	m := make(map[string]videoData)
	for rows.Next() {
		var id string
		var d videoData
		if err := rows.Scan(&id, &d.title, &d.thumbnailURL, &d.thumbnailTinyURL, &d.videoURL, &d.duration, &d.playlist); err != nil {
			return nil, fmt.Errorf("fetchVideoData scan: %w", err)
		}
		m[id] = d
	}
	return m, rows.Err()
}

type linkPostData struct {
	description      *string
	thumbnailURL     *string
	thumbnailTinyURL *string
	category         *string
	title            string
	url              string
	domain           string
}

func (s *SQLitePostStore) fetchLinkPostData(ctx context.Context, ids []string) (map[string]linkPostData, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	q := "SELECT post_id, title, url, domain, description, thumbnail_url, thumbnail_tiny_url, category FROM link_posts WHERE post_id IN (" + placeholders(len(ids)) + ")"
	rows, err := s.db.QueryContext(ctx, q, stringsToAny(ids)...)
	if err != nil {
		return nil, fmt.Errorf("fetchLinkPostData: %w", err)
	}
	defer func() { _ = rows.Close() }()

	m := make(map[string]linkPostData)
	for rows.Next() {
		var id string
		var d linkPostData
		if err := rows.Scan(&id, &d.title, &d.url, &d.domain, &d.description, &d.thumbnailURL, &d.thumbnailTinyURL, &d.category); err != nil {
			return nil, fmt.Errorf("fetchLinkPostData scan: %w", err)
		}
		m[id] = d
	}
	return m, rows.Err()
}

// ── Assembly ──────────────────────────────────────────────────────────────────

func (s *SQLitePostStore) assembleAll(
	ctx context.Context,
	baseRows []baseRow,
	byType map[string][]string,
	tagsMap map[string][]string,
) ([]models.Post, error) {
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

	result := make([]models.Post, 0, len(baseRows))
	for _, r := range baseRows {
		tags := coalesceStringSlice(tagsMap[r.ID])
		var post models.Post

		switch r.Type {
		case "essay":
			d := essayMap[r.ID]
			post = models.EssayPost{
				ID: r.ID, Type: r.Type, CreatedAt: r.CreatedAt, Tags: tags,
				Title: d.title, Slug: d.slug, Excerpt: d.excerpt, Body: d.body,
				ReadingTimeMinutes: d.readingTimeMinutes,
			}
		case "short":
			d := shortMap[r.ID]
			post = models.ShortPost{
				ID: r.ID, Type: r.Type, CreatedAt: r.CreatedAt, Tags: tags,
				Body: d.body,
			}
		case "music":
			d := musicMap[r.ID]
			post = models.MusicPost{
				ID: r.ID, Type: r.Type, CreatedAt: r.CreatedAt, Tags: tags,
				Title: d.title, AlbumArt: d.albumArt, AlbumArtTinyURL: d.albumArtTiny,
				AudioURL: d.audioURL, Duration: d.duration, Album: d.album,
			}
		case "photo":
			meta := photoMeta[r.ID]
			imgs := photoImages[r.ID]
			if imgs == nil {
				imgs = []models.PhotoImage{}
			}
			post = models.PhotoPost{
				ID: r.ID, Type: r.Type, CreatedAt: r.CreatedAt, Tags: tags,
				Images: imgs, Location: meta.location,
			}
		case "video":
			d := videoMap[r.ID]
			post = models.VideoPost{
				ID: r.ID, Type: r.Type, CreatedAt: r.CreatedAt, Tags: tags,
				Title: d.title, ThumbnailURL: d.thumbnailURL, ThumbnailTinyURL: d.thumbnailTinyURL,
				VideoURL: d.videoURL, Duration: d.duration, Playlist: d.playlist,
			}
		case "link":
			d := linkPostMap[r.ID]
			post = models.LinkPost{
				ID: r.ID, Type: r.Type, CreatedAt: r.CreatedAt, Tags: tags,
				Title: d.title, URL: d.url, Domain: d.domain,
				Description: d.description, ThumbnailURL: d.thumbnailURL,
				ThumbnailTinyURL: d.thumbnailTinyURL, Category: d.category,
			}
		default:
			continue
		}
		result = append(result, post)
	}
	return result, nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func placeholders(n int) string {
	if n == 0 {
		return ""
	}
	return strings.Repeat("?,", n-1) + "?"
}

func stringsToAny(ss []string) []any {
	out := make([]any, len(ss))
	for i, s := range ss {
		out[i] = s
	}
	return out
}

func extractIDs(rows []baseRow) []string {
	ids := make([]string, len(rows))
	for i, r := range rows {
		ids[i] = r.ID
	}
	return ids
}

func groupByType(rows []baseRow) map[string][]string {
	m := make(map[string][]string)
	for _, r := range rows {
		m[r.Type] = append(m[r.Type], r.ID)
	}
	return m
}

func coalesceStringSlice(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
