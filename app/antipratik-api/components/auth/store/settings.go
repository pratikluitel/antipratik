package store

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"

	"github.com/pratikluitel/antipratik/components/auth"
)

// sqliteSettingsStore implements SettingsStore using a SQLite database.
type sqliteSettingsStore struct {
	db *sql.DB
}

// NewSettingsStore creates a new sqliteSettingsStore backed by db.
func NewSettingsStore(db *sql.DB) auth.SettingsStore {
	return &sqliteSettingsStore{db: db}
}

// GetOrCreateJWTSecret returns the persisted JWT secret from the settings table,
// generating and storing a new one if none exists yet.
func (s *sqliteSettingsStore) GetOrCreateJWTSecret(ctx context.Context) (string, error) {
	var secret string
	err := s.db.QueryRowContext(ctx, `SELECT value FROM settings WHERE key='jwt_secret'`).Scan(&secret)
	if err == nil {
		return secret, nil
	}
	if err != sql.ErrNoRows {
		return "", fmt.Errorf("read jwt_secret: %w", err)
	}

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate jwt_secret: %w", err)
	}
	secret = hex.EncodeToString(b)

	if _, err := s.db.ExecContext(ctx, `INSERT INTO settings (key, value) VALUES ('jwt_secret', ?)`, secret); err != nil {
		return "", fmt.Errorf("store jwt_secret: %w", err)
	}
	return secret, nil
}
