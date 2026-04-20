package logic

import (
	"context"

	"github.com/pratikluitel/antipratik/components/auth/store"
)

// SetupService implements SetupLogic, providing application bootstrapping
// operations so that main.go never calls the store layer directly.
type SetupService struct {
	users    store.UserStore
	settings store.SettingsStore
}

// NewSetupService creates a new SetupService.
func NewSetupService(users store.UserStore, settings store.SettingsStore) *SetupService {
	return &SetupService{users: users, settings: settings}
}

// UpsertAdminUser creates or updates the admin user with the given password.
func (s *SetupService) UpsertAdminUser(ctx context.Context, password string) error {
	return s.users.UpsertAdminUser(ctx, password)
}

// GetOrCreateJWTSecret returns the persisted JWT signing secret,
// generating and storing a new one if none exists yet.
func (s *SetupService) GetOrCreateJWTSecret(ctx context.Context) (string, error) {
	return s.settings.GetOrCreateJWTSecret(ctx)
}
