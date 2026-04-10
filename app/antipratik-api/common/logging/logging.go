// Package logging provides a structured logger backed by log/slog.
// The Logger interface is passed through the factory pattern so each layer
// can emit levelled, structured log lines without importing slog directly.
package logging

import (
	"log/slog"
	"os"
	"strings"
)

// Logger is the structured logging interface threaded through the factory pattern.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type slogLogger struct {
	l *slog.Logger
}

func (s *slogLogger) Debug(msg string, args ...any) { s.l.Debug(msg, args...) }
func (s *slogLogger) Info(msg string, args ...any)  { s.l.Info(msg, args...) }
func (s *slogLogger) Warn(msg string, args ...any)  { s.l.Warn(msg, args...) }
func (s *slogLogger) Error(msg string, args ...any) { s.l.Error(msg, args...) }

// New creates a Logger at the given level string ("debug", "info", "warn", "error").
// Defaults to INFO for any unrecognised value.
// Output is text format (suitable for journald/systemd on Hetzner).
func New(level string) Logger {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn", "warning":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	return &slogLogger{l: slog.New(h)}
}
