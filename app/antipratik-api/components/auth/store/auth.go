package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/pratikluitel/antipratik/components/auth"
)

// sqliteUserStore implements UserStore using a SQLite database.
type sqliteUserStore struct {
	db *sql.DB
}

// NewUserStore creates a new sqliteUserStore backed by db.
func NewUserStore(db *sql.DB) auth.UserStore {
	return &sqliteUserStore{db: db}
}

func (s *sqliteUserStore) GetUserByUsername(ctx context.Context, username string) (*auth.User, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, username, password_hash, current_token, token_expires_at FROM users WHERE username=?`, username)
	return scanUser(row)
}

func (s *sqliteUserStore) GetUserByToken(ctx context.Context, token string) (*auth.User, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, username, password_hash, current_token, token_expires_at FROM users WHERE current_token=?`, token)
	return scanUser(row)
}

func (s *sqliteUserStore) UpsertToken(ctx context.Context, username string, token string, expiresAt time.Time) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE users SET current_token=?, token_expires_at=? WHERE username=?`,
		token, expiresAt.UTC().Format(time.RFC3339), username)
	return err
}

// UpsertAdminUser ensures an admin user exists with the given password.
// Creates the user if absent; updates the password hash if it has changed.
func (s *sqliteUserStore) UpsertAdminUser(ctx context.Context, password string) error {
	var id, hash string
	err := s.db.QueryRowContext(ctx, `SELECT id, password_hash FROM users WHERE username = ?`, "admin").Scan(&id, &hash)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("upsert admin: %w", err)
	}

	if err == sql.ErrNoRows {
		var newHash []byte
		newHash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash admin password: %w", err)
		}
		_, err = s.db.ExecContext(ctx, `INSERT INTO users (id, username, password_hash) VALUES (?, ?, ?)`,
			uuid.New().String(), "admin", string(newHash))
		return err
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil {
		return nil // password unchanged
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `UPDATE users SET password_hash = ? WHERE username = ?`, string(newHash), "admin")
	return err
}

func scanUser(row *sql.Row) (*auth.User, error) {
	var u auth.User
	var currentToken sql.NullString
	var tokenExpiresAt sql.NullString
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &currentToken, &tokenExpiresAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if currentToken.Valid {
		u.CurrentToken = &currentToken.String
	}
	if tokenExpiresAt.Valid {
		t, err := time.Parse(time.RFC3339, tokenExpiresAt.String)
		if err == nil {
			u.TokenExpiresAt = &t
		}
	}
	return &u, nil
}
