// Package config loads noticoel's configuration from a YAML file into a
// strongly typed Config.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the strongly typed application configuration.
type Config struct {
	Debug    bool           `yaml:"debug"`
	Server   ServerConfig   `yaml:"server"`
	Auth     AuthConfig     `yaml:"auth"`
	Database DatabaseConfig `yaml:"database"`
}

// ServerConfig configures the HTTP server.
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// Addr returns the host:port pair the HTTP server should listen on.
func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// AuthConfig configures the API's bearer token authentication.
type AuthConfig struct {
	// Token authenticates incoming requests. It is read from the
	// AUTH_TOKEN environment variable rather than the YAML file, so the
	// secret never needs to be committed to version control.
	Token string `yaml:"-"`
}

// DatabaseConfig configures the SQLite database.
type DatabaseConfig struct {
	Path string `yaml:"path"`
}

// Load reads and decodes the YAML configuration file at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	cfg.Auth.Token = os.Getenv("AUTH_TOKEN")
	if cfg.Auth.Token == "" {
		return nil, fmt.Errorf("AUTH_TOKEN environment variable is required")
	}

	return &cfg, nil
}
