package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	commonerrors "github.com/pratikluitel/antipratik/common/errors"
	"github.com/pratikluitel/antipratik/components/auth"
)

type authLogic struct {
	users     auth.UserStore
	jwtSecret string
}

// NewAuthLogic creates a new authLogic.
func NewAuthLogic(users auth.UserStore, jwtSecret string) auth.AuthLogic {
	return &authLogic{users: users, jwtSecret: jwtSecret}
}

func (s *authLogic) Login(ctx context.Context, username, password string) (string, error) {
	if err := commonerrors.RequireNonEmpty("username", username); err != nil {
		return "", err
	}
	if err := commonerrors.RequireNonEmpty("password", password); err != nil {
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

func (s *authLogic) ValidateToken(ctx context.Context, token string) error {
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
