// Package config loads and exposes runtime configuration from a YAML file.
package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config holds all runtime configuration for the server.
type Config struct {
	Server        ServerConfig  `yaml:"server"`
	DB            DBConfig      `yaml:"db"`
	AdminPassword string        `yaml:"admin_password"`
	Static        StaticConfig  `yaml:"static"`
	Storage       StorageConfig `yaml:"storage"`
	Logging       LoggingConfig `yaml:"logging"`
}

// LoggingConfig controls log verbosity. Level accepts "debug", "info", "warn", or "error".
type LoggingConfig struct {
	Level string `yaml:"level"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// DBConfig holds database connection settings.
type DBConfig struct {
	Path string `yaml:"path"`
}


// StaticConfig holds settings for serving the frontend static build.
type StaticConfig struct {
	Dir string `yaml:"dir"`
}

// StorageConfig holds settings for the pluggable file storage backend.
type StorageConfig struct {
	Backend  string   `yaml:"backend"`   // "local" or "r2"
	LocalDir string   `yaml:"local_dir"` // used when backend=local
	R2       R2Config `yaml:"r2"`
}

// R2Config holds Cloudflare R2 credentials and settings.
// Files are always served via the backend's own /files/ and /thumbnails/ endpoints
// regardless of storage backend — R2 object URLs are never exposed to clients.
type R2Config struct {
	Endpoint        string `yaml:"endpoint"`
	Bucket          string `yaml:"bucket"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
}

// Load reads and parses the YAML config file at the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config %q: %w", path, err)
	}

	if portStr := os.Getenv("ANTIPRATIK_PORT"); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("parse ANTIPRATIK_PORT %q: %w", portStr, err)
		}
		cfg.Server.Port = port
	}
	if v := os.Getenv("ANTIPRATIK_HOST"); v != "" {
		cfg.Server.Host = v
	}
	if v := os.Getenv("ANTIPRATIK_DB_PATH"); v != "" {
		cfg.DB.Path = v
	}
	if v := os.Getenv("ANTIPRATIK_ADMIN_PASSWORD"); v != "" {
		cfg.AdminPassword = v
	}
	if v := os.Getenv("ANTIPRATIK_R2_ENDPOINT"); v != "" {
		cfg.Storage.R2.Endpoint = v
	}
	if v := os.Getenv("ANTIPRATIK_R2_BUCKET"); v != "" {
		cfg.Storage.R2.Bucket = v
	}
	if v := os.Getenv("ANTIPRATIK_R2_ACCESS_KEY_ID"); v != "" {
		cfg.Storage.R2.AccessKeyID = v
	}
	if v := os.Getenv("ANTIPRATIK_R2_SECRET_ACCESS_KEY"); v != "" {
		cfg.Storage.R2.SecretAccessKey = v
	}

	return &cfg, nil
}

// Addr returns the host:port string for http.ListenAndServe.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
