// create_user is a standalone CLI tool to create or refresh a user in the antipratik SQLite database.
//
// Usage:
//
//	create_user --db ./data/antipratik.db --username admin --password secret
//	create_user --db ./data/antipratik.db --username admin --refresh
package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	dbPath := flag.String("db", "", "path to SQLite database file (required)")
	username := flag.String("username", "", "username (required)")
	password := flag.String("password", "", "password (required unless --refresh)")
	refresh := flag.Bool("refresh", false, "only refresh token for existing user, no password change")
	flag.Parse()

	if *dbPath == "" || *username == "" {
		fmt.Fprintln(os.Stderr, "usage: create_user --db <path> --username <user> [--password <pass>] [--refresh]")
		os.Exit(1)
	}
	if !*refresh && *password == "" {
		fmt.Fprintln(os.Stderr, "error: --password is required unless --refresh is set")
		os.Exit(1)
	}

	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?_foreign_keys=on", *dbPath))
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	secret, err := getOrCreateJWTSecret(db)
	if err != nil {
		log.Fatalf("jwt secret: %v", err)
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	token, err := generateToken(*username, secret, expiresAt)
	if err != nil {
		log.Fatalf("generate token: %v", err)
	}

	if *refresh {
		res, err := db.Exec(
			`UPDATE users SET current_token=?, token_expires_at=? WHERE username=?`,
			token, expiresAt.UTC().Format(time.RFC3339), *username,
		)
		if err != nil {
			log.Fatalf("update token: %v", err)
		}
		rows, _ := res.RowsAffected()
		if rows == 0 {
			log.Fatalf("user %q not found", *username)
		}
		fmt.Printf("Token refreshed for user %q\nToken: %s\n", *username, token)
		return
	}

	// Check if user already exists
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE username=?`, *username).Scan(&count); err != nil {
		log.Fatalf("check user: %v", err)
	}
	if count > 0 {
		log.Fatalf("user %q already exists; use --refresh to update the token", *username)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}

	id := uuid.New().String()
	_, err = db.Exec(
		`INSERT INTO users (id, username, password_hash, current_token, token_expires_at) VALUES (?, ?, ?, ?, ?)`,
		id, *username, string(hash), token, expiresAt.UTC().Format(time.RFC3339),
	)
	if err != nil {
		log.Fatalf("insert user: %v", err)
	}

	fmt.Printf("User %q created (id=%s)\nToken: %s\n", *username, id, token)
}

func getOrCreateJWTSecret(db *sql.DB) (string, error) {
	var secret string
	err := db.QueryRow(`SELECT value FROM settings WHERE key='jwt_secret'`).Scan(&secret)
	if err == nil {
		return secret, nil
	}
	if err != sql.ErrNoRows {
		return "", fmt.Errorf("read jwt_secret: %w", err)
	}

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate jwt_secret: %w", err)
	}
	secret = hex.EncodeToString(b)

	if _, err := db.Exec(`INSERT INTO settings (key, value) VALUES ('jwt_secret', ?)`, secret); err != nil {
		return "", fmt.Errorf("store jwt_secret: %w", err)
	}
	return secret, nil
}

func generateToken(username, secret string, expiresAt time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": expiresAt.Unix(),
	})
	return token.SignedString([]byte(secret))
}
