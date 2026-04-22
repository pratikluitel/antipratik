package logic

import (
	"context"

	"github.com/pratikluitel/antipratik/components/auth"
)

// setupLogic implements SetupLogic, providing application bootstrapping
// operations so that main.go never calls the store layer directly.
type setupLogic struct {
	users    auth.UserStore
	settings auth.SettingsStore
}

// NewSetupLogic creates a new setupLogic.
func NewSetupLogic(users auth.UserStore, settings auth.SettingsStore) auth.SetupLogic {
	return &setupLogic{users: users, settings: settings}
}

// UpsertAdminUser creates or updates the admin user with the given password.
func (s *setupLogic) UpsertAdminUser(ctx context.Context, password string) error {
	return s.users.UpsertAdminUser(ctx, password)
}

// GetOrCreateJWTSecret returns the persisted JWT signing secret,
// generating and storing a new one if none exists yet.
func (s *setupLogic) GetOrCreateJWTSecret(ctx context.Context) (string, error) {
	return s.settings.GetOrCreateJWTSecret(ctx)
}
