package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *authLogic) createToken(username string) (string, time.Time, error) {
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": expiresAt.Unix(),
	})

	tokenStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("auth service sign token: %w", err)
	}
	return tokenStr, expiresAt, nil
}

func (s *authLogic) ValidateToken(ctx context.Context, token string) error {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !parsed.Valid {
		return fmt.Errorf("invalid token")
	}
	user, err := s.userStore.GetUserByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("auth service validate: %w", err)
	}
	if user == nil {
		return fmt.Errorf("token not found")
	}
	return nil
}
