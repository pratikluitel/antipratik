package logic

import "context"

// AuthLogic defines authentication operations.
type AuthLogic interface {
	Login(ctx context.Context, username, password string) (string, error)
	ValidateToken(ctx context.Context, token string) error
}

// SetupLogic defines application bootstrapping operations that run at startup.
// It ensures main.go never calls the store layer directly.
type SetupLogic interface {
	UpsertAdminUser(ctx context.Context, password string) error
	GetOrCreateJWTSecret(ctx context.Context) (string, error)
}
