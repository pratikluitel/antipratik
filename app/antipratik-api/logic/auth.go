package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pratikluitel/antipratik/store"
	"golang.org/x/crypto/bcrypt"
)

// AuthService implements AuthLogic.
type AuthService struct {
	users     store.UserStore
	jwtSecret string
}

// NewAuthService creates a new AuthService.
func NewAuthService(users store.UserStore, jwtSecret string) *AuthService {
	return &AuthService{users: users, jwtSecret: jwtSecret}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	if err := requireNonEmpty("username", username); err != nil {
		return "", err
	}
	if err := requireNonEmpty("password", password); err != nil {
		return "", err
	}

	user, err := s.users.GetUserByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("auth service login: %w", err)
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Username,
		"exp": expiresAt.Unix(),
	})
	tokenStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("auth service sign token: %w", err)
	}

	if err = s.users.UpsertToken(ctx, username, tokenStr, expiresAt); err != nil {
		return "", fmt.Errorf("auth service store token: %w", err)
	}
	return tokenStr, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) error {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !parsed.Valid {
		return errors.New("invalid token")
	}
	user, err := s.users.GetUserByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("auth service validate: %w", err)
	}
	if user == nil {
		return errors.New("token not found")
	}
	return nil
}
