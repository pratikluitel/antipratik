package auth

import "time"

// User represents an authenticated user.
type User struct {
	CurrentToken   *string
	TokenExpiresAt *time.Time
	ID             string
	Username       string
	PasswordHash   string
}
