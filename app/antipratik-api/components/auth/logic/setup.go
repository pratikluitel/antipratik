package logic

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/pratikluitel/antipratik/components/auth"
)

// setupLogic implements SetupLogic, providing application bootstrapping
// operations so that main.go never calls the store layer directly.
type setupLogic struct {
	users         auth.UserStore
	settingsStore auth.SettingsStore
}

// NewSetupLogic creates a new setupLogic.
func NewSetupLogic(users auth.UserStore, settingsStore auth.SettingsStore) auth.SetupLogic {
	return &setupLogic{users: users, settingsStore: settingsStore}
}

// UpsertAdminUser creates or updates the admin user with the given password.
func (s *setupLogic) UpsertAdminUser(ctx context.Context, password string) error {
	return s.users.UpsertAdminUser(ctx, password)
}

// GetOrCreateJWTSecret returns the persisted JWT signing secret,
// generating and storing a new one if none exists yet.
func (s *setupLogic) GetOrCreateJWTSecret(ctx context.Context) (string, error) {
	fallbackSecret, err := generateSecureSecret(32)
	if err != nil {
		return "", fmt.Errorf("generate fallback JWT secret: %w", err)
	}

	return s.settingsStore.GetOrCreateJWTSecret(ctx, fallbackSecret)
}

func generateSecureSecret(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate secure bytes: %w", err)
	}
	return hex.EncodeToString(b), nil
}
