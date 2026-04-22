package auth

import (
	"context"
	"net/http"
	"time"
)

// AuthLogic defines authentication operations.
type AuthLogic interface {
	Login(ctx context.Context, username, password string) (string, error)
	ValidateToken(ctx context.Context, token string) error
}

// SetupLogic defines authentication bootstrapping operations that run at startup.
type SetupLogic interface {
	UpsertAdminUser(ctx context.Context, password string) error
	GetOrCreateJWTSecret(ctx context.Context) (string, error)
}

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

type AuthAPI interface {
	Login(w http.ResponseWriter, r *http.Request)
}
