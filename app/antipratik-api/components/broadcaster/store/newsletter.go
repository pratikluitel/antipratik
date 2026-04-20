package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// SQLiteNewsletterStore is the SQLite implementation of NewsletterStore.
type SQLiteNewsletterStore struct {
	db *sql.DB
}

// NewNewsletterStore creates a new SQLiteNewsletterStore.
func NewNewsletterStore(db *sql.DB) *SQLiteNewsletterStore {
	return &SQLiteNewsletterStore{db: db}
}

// Subscribe inserts the email (lowercase-trimmed) into newsletter_subscribers.
// Returns ErrDuplicate if the email is already subscribed.
func (s *SQLiteNewsletterStore) Subscribe(ctx context.Context, email string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	res, err := s.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO newsletter_subscribers (email) VALUES (?)`, email)
	if err != nil {
		return fmt.Errorf("NewsletterStore.Subscribe: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("NewsletterStore.Subscribe rows affected: %w", err)
	}
	if rows == 0 {
		return ErrDuplicate
	}
	return nil
}
