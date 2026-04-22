package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pratikluitel/antipratik/components/posts"
)

// sqliteLinkStore implements LinkStore using a SQLite database.
type sqliteLinkStore struct {
	db *sql.DB
}

// NewLinkStore creates a new sqliteLinkStore backed by db.
func NewLinkStore(db *sql.DB) posts.LinkStore {
	return &sqliteLinkStore{db: db}
}

// GetLinks returns all external links.
func (s *sqliteLinkStore) GetLinks(ctx context.Context) ([]posts.ExternalLink, error) {
	const q = `SELECT id, title, url, domain, description, featured, category
	           FROM links ORDER BY category, id`
	return s.queryLinks(ctx, q)
}

// GetFeaturedLinks returns up to 4 featured external links.
func (s *sqliteLinkStore) GetFeaturedLinks(ctx context.Context) ([]posts.ExternalLink, error) {
	const q = `SELECT id, title, url, domain, description, featured, category
	           FROM links WHERE featured = 1 ORDER BY id LIMIT 4`
	return s.queryLinks(ctx, q)
}

// GetLinkByID returns an external link by ID, or an error if not found.
func (s *sqliteLinkStore) GetLinkByID(ctx context.Context, id string) (posts.ExternalLink, error) {
	links, err := s.queryLinks(ctx,
		`SELECT id, title, url, domain, description, featured, category FROM links WHERE id = ?`, id)
	if err != nil {
		return posts.ExternalLink{}, err
	}
	if len(links) == 0 {
		return posts.ExternalLink{}, fmt.Errorf("link %q not found", id)
	}
	return links[0], nil
}

func (s *sqliteLinkStore) CreateLink(ctx context.Context, id string, input posts.CreateExternalLink) error {
	featured := 0
	if input.Featured {
		featured = 1
	}
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO links (id, title, url, domain, description, featured, category) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, input.Title, input.URL, input.Domain, input.Description, featured, input.Category)
	return err
}

func (s *sqliteLinkStore) UpdateLink(ctx context.Context, id string, input posts.CreateExternalLink) error {
	featured := 0
	if input.Featured {
		featured = 1
	}
	_, err := s.db.ExecContext(ctx,
		`UPDATE links SET title=?, url=?, domain=?, description=?, featured=?, category=? WHERE id=?`,
		input.Title, input.URL, input.Domain, input.Description, featured, input.Category, id)
	return err
}

func (s *sqliteLinkStore) DeleteLink(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM links WHERE id=?`, id)
	return err
}

func (s *sqliteLinkStore) queryLinks(ctx context.Context, query string, args ...any) ([]posts.ExternalLink, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("queryLinks: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var result []posts.ExternalLink
	for rows.Next() {
		var l posts.ExternalLink
		var featuredInt int
		if err := rows.Scan(&l.ID, &l.Title, &l.URL, &l.Domain, &l.Description, &featuredInt, &l.Category); err != nil {
			return nil, fmt.Errorf("queryLinks scan: %w", err)
		}
		l.Featured = featuredInt == 1
		result = append(result, l)
	}
	if result == nil {
		result = []posts.ExternalLink{}
	}
	return result, rows.Err()
}
