package auth

import "time"

// User represents an authenticated user.
type User struct {
	CurrentToken   string
	TokenExpiresAt time.Time
	ID             string
	Username       string
	PasswordHash   string
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}
