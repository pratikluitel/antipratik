package store

import (
	"context"
	"time"
)

// UserStore handles user authentication database operations.
type UserStore interface {
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByToken(ctx context.Context, token string) (*User, error)
	UpsertToken(ctx context.Context, username string, token string, expiresAt time.Time) error
	UpsertAdminUser(ctx context.Context, password string) error
}

// SettingsStore handles application settings persistence.
type SettingsStore interface {
	GetOrCreateJWTSecret(ctx context.Context) (string, error)
}
