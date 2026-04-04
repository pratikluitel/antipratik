package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pratikluitel/antipratik/models"
)

// SQLiteLinkStore implements LinkStore using a SQLite database.
type SQLiteLinkStore struct {
	db *sql.DB
}

// NewLinkStore creates a new SQLiteLinkStore backed by db.
func NewLinkStore(db *sql.DB) *SQLiteLinkStore {
	return &SQLiteLinkStore{db: db}
}

// GetLinks returns all external links.
func (s *SQLiteLinkStore) GetLinks(ctx context.Context) ([]models.ExternalLink, error) {
	const q = `SELECT id, title, url, domain, description, featured, category
	           FROM links ORDER BY category, id`
	return s.queryLinks(ctx, q)
}

// GetFeaturedLinks returns up to 4 featured external links.
func (s *SQLiteLinkStore) GetFeaturedLinks(ctx context.Context) ([]models.ExternalLink, error) {
	const q = `SELECT id, title, url, domain, description, featured, category
	           FROM links WHERE featured = 1 ORDER BY id LIMIT 4`
	return s.queryLinks(ctx, q)
}

func (s *SQLiteLinkStore) queryLinks(ctx context.Context, query string, args ...any) ([]models.ExternalLink, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("queryLinks: %w", err)
	}
	defer rows.Close()

	var result []models.ExternalLink
	for rows.Next() {
		var l models.ExternalLink
		var featuredInt int
		if err := rows.Scan(&l.ID, &l.Title, &l.URL, &l.Domain, &l.Description, &featuredInt, &l.Category); err != nil {
			return nil, fmt.Errorf("queryLinks scan: %w", err)
		}
		l.Featured = featuredInt == 1
		result = append(result, l)
	}
	if result == nil {
		result = []models.ExternalLink{}
	}
	return result, rows.Err()
}
