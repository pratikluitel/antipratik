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
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
	Static StaticConfig `yaml:"static"`
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

	if host := os.Getenv("ANTIPRATIK_HOST"); host != "" {
		cfg.Server.Host = host
	}

	return &cfg, nil
}

// Addr returns the host:port string for http.ListenAndServe.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
