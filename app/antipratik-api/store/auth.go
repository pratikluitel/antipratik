package store

import (
	"context"
	"database/sql"
	"time"
)

// SQLiteUserStore implements UserStore using a SQLite database.
type SQLiteUserStore struct {
	db *sql.DB
}

// NewUserStore creates a new SQLiteUserStore backed by db.
func NewUserStore(db *sql.DB) *SQLiteUserStore {
	return &SQLiteUserStore{db: db}
}

func (s *SQLiteUserStore) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, username, password_hash, current_token, token_expires_at FROM users WHERE username=?`, username)
	return scanUser(row)
}

func (s *SQLiteUserStore) GetUserByToken(ctx context.Context, token string) (*User, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, username, password_hash, current_token, token_expires_at FROM users WHERE current_token=?`, token)
	return scanUser(row)
}

func (s *SQLiteUserStore) UpsertToken(ctx context.Context, username string, token string, expiresAt time.Time) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE users SET current_token=?, token_expires_at=? WHERE username=?`,
		token, expiresAt.UTC().Format(time.RFC3339), username)
	return err
}

func scanUser(row *sql.Row) (*User, error) {
	var u User
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
