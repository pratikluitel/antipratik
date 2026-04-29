package logic

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	commonerrors "github.com/pratikluitel/antipratik/common/errors"
	"github.com/pratikluitel/antipratik/components/auth"
)

type authLogic struct {
	userStore auth.UserStore
	jwtSecret string
}

// NewAuthLogic creates a new authLogic.
func NewAuthLogic(userStore auth.UserStore, jwtSecret string) auth.AuthLogic {
	return &authLogic{userStore: userStore, jwtSecret: jwtSecret}
}

func (s *authLogic) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.validateUser(ctx, username)
	if err != nil {
		return "", err
	}

	if err = s.validatePassword(password, user.PasswordHash); err != nil {
		return "", err
	}

	token, expiresAt, err := s.createToken(user.Username)
	if err != nil {
		return "", err
	}

	if err = s.userStore.UpsertToken(ctx, user.Username, token, expiresAt); err != nil {
		return "", fmt.Errorf("failed to store token: %w", err)
	}
	return token, nil
}

func (s *authLogic) validateUser(ctx context.Context, username string) (*auth.User, error) {
	if err := commonerrors.RequireNonEmpty("username", username); err != nil {
		return nil, err
	}

	user, err := s.userStore.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *authLogic) validatePassword(password, hashedPassword string) error {
	if err := commonerrors.RequireNonEmpty("password", password); err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return errors.New("invalid credentials")
	}
	return nil
}
