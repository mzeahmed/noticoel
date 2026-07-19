// Package config loads noticoel's configuration.
//
// Infrastructure configuration (server, database, which notifiers are
// enabled, and their credentials) is deployment-specific, so it is read
// entirely from NOTICOEL_* environment variables, following the
// Twelve-Factor App convention — a Docker deployment needs nothing more
// than a docker-compose.yml with an `environment:` block, no config file
// to mount.
//
// The YAML file passed to Load is optional. It exists for structured
// business configuration that doesn't map cleanly onto environment
// variables (routing rules, notification templates... none of which
// exist yet). If it's absent, Load simply falls back to environment
// variables and defaults. If an older-style file (with `server:`,
// `database:`, `notifiers:` sections) is still present, those values are
// used as a fallback wherever the corresponding environment variable
// isn't set, so upgrading doesn't break an existing deployment.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config is the strongly typed application configuration.
type Config struct {
	Debug     bool
	Server    ServerConfig
	Auth      AuthConfig
	Database  DatabaseConfig
	Notifiers NotifiersConfig
}

// ServerConfig configures the HTTP server.
type ServerConfig struct {
	Host string
	Port int
}

// Addr returns the host:port pair the HTTP server should listen on.
func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// AuthConfig configures the API's bearer token authentication.
type AuthConfig struct {
	// Token authenticates incoming requests. There is no default: Load
	// fails if NOTICOEL_AUTH_TOKEN isn't set.
	Token string
}

// DatabaseConfig configures the SQLite database.
type DatabaseConfig struct {
	Path string
}

// NotifiersConfig configures the notification channels events are
// dispatched to.
type NotifiersConfig struct {
	Telegram TelegramConfig
	Ntfy     NtfyConfig
	Webhook  WebhookConfig
	Discord  DiscordConfig
	Email    EmailConfig
}

// TelegramConfig configures the Telegram notifier. BotToken and ChatID
// have no default: Load fails if they're missing while Telegram is
// enabled.
type TelegramConfig struct {
	Enabled  bool
	BotToken string
	ChatID   string
}

// NtfyConfig configures the ntfy notifier.
type NtfyConfig struct {
	Enabled bool
	Server  string
	Topic   string
}

// WebhookConfig configures the generic webhook notifier.
type WebhookConfig struct {
	Enabled bool
	URL     string
}

// DiscordConfig configures the Discord notifier.
type DiscordConfig struct {
	Enabled bool
	Webhook string
}

// EmailConfig configures the email (SMTP) notifier.
type EmailConfig struct {
	Enabled  bool
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// legacy mirrors the pre-2.0 config.yaml shape, where infrastructure
// settings lived in YAML instead of the environment. It exists only so
// that shape keeps working as a fallback; new deployments don't need a
// YAML file at all.
type legacy struct {
	Debug  bool `yaml:"debug"`
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
	Notifiers struct {
		Telegram struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"telegram"`
		Ntfy struct {
			Enabled bool   `yaml:"enabled"`
			Server  string `yaml:"server"`
			Topic   string `yaml:"topic"`
		} `yaml:"ntfy"`
		Webhook struct {
			Enabled bool   `yaml:"enabled"`
			URL     string `yaml:"url"`
		} `yaml:"webhook"`
		Discord struct {
			Enabled bool   `yaml:"enabled"`
			Webhook string `yaml:"webhook"`
		} `yaml:"discord"`
		Email struct {
			Enabled  bool   `yaml:"enabled"`
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			From     string `yaml:"from"`
		} `yaml:"email"`
	} `yaml:"notifiers"`
}

// loadLegacy reads the YAML file at path if it exists. A missing file is
// not an error — it just means there's nothing to fall back to. found
// reports whether a file was actually read, so Load can tell "no file"
// (meaning "no override") apart from a file whose fields happen to match
// their own zero value.
func loadLegacy(path string) (l legacy, found bool, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return legacy{}, false, nil
		}

		return legacy{}, false, fmt.Errorf("read config: %w", err)
	}

	if err := yaml.Unmarshal(data, &l); err != nil {
		return legacy{}, false, fmt.Errorf("decode config: %w", err)
	}

	return l, true, nil
}

// Load builds the Config from environment variables, falling back to the
// optional legacy YAML file at path for any field whose environment
// variable isn't set, and to a hardcoded default after that.
func Load(path string) (*Config, error) {
	l, found, err := loadLegacy(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Debug: envBool("NOTICOEL_DEBUG", fallbackBool(found, l.Debug, false)),
		Server: ServerConfig{
			Host: envString("NOTICOEL_SERVER_HOST", fallbackString(found, l.Server.Host, "0.0.0.0")),
			Port: envInt("NOTICOEL_SERVER_PORT", fallbackInt(found, l.Server.Port, 8080)),
		},
		Database: DatabaseConfig{
			Path: envString("NOTICOEL_DATABASE_PATH", fallbackString(found, l.Database.Path, "./data/noticoel.db")),
		},
		Notifiers: NotifiersConfig{
			Telegram: TelegramConfig{
				Enabled: envBool("NOTICOEL_TELEGRAM_ENABLED", fallbackBool(found, l.Notifiers.Telegram.Enabled, true)),
			},
			Ntfy: NtfyConfig{
				Enabled: envBool("NOTICOEL_NTFY_ENABLED", fallbackBool(found, l.Notifiers.Ntfy.Enabled, false)),
				Server:  envString("NOTICOEL_NTFY_SERVER", fallbackString(found, l.Notifiers.Ntfy.Server, "https://ntfy.sh")),
				Topic:   envString("NOTICOEL_NTFY_TOPIC", fallbackString(found, l.Notifiers.Ntfy.Topic, "")),
			},
			Webhook: WebhookConfig{
				Enabled: envBool("NOTICOEL_WEBHOOK_ENABLED", fallbackBool(found, l.Notifiers.Webhook.Enabled, false)),
				URL:     envString("NOTICOEL_WEBHOOK_URL", fallbackString(found, l.Notifiers.Webhook.URL, "")),
			},
			Discord: DiscordConfig{
				Enabled: envBool("NOTICOEL_DISCORD_ENABLED", fallbackBool(found, l.Notifiers.Discord.Enabled, false)),
				Webhook: envString("NOTICOEL_DISCORD_WEBHOOK", fallbackString(found, l.Notifiers.Discord.Webhook, "")),
			},
			Email: EmailConfig{
				Enabled:  envBool("NOTICOEL_EMAIL_ENABLED", fallbackBool(found, l.Notifiers.Email.Enabled, false)),
				Host:     envString("NOTICOEL_EMAIL_HOST", fallbackString(found, l.Notifiers.Email.Host, "")),
				Port:     envInt("NOTICOEL_EMAIL_PORT", fallbackInt(found, l.Notifiers.Email.Port, 587)),
				Username: envString("NOTICOEL_EMAIL_USERNAME", fallbackString(found, l.Notifiers.Email.Username, "")),
				Password: envString("NOTICOEL_EMAIL_PASSWORD", fallbackString(found, l.Notifiers.Email.Password, "")),
				From:     envString("NOTICOEL_EMAIL_FROM", fallbackString(found, l.Notifiers.Email.From, "")),
			},
		},
	}

	cfg.Auth.Token = os.Getenv("NOTICOEL_AUTH_TOKEN")
	if cfg.Auth.Token == "" {
		return nil, errors.New("NOTICOEL_AUTH_TOKEN environment variable is required")
	}

	if cfg.Notifiers.Telegram.Enabled {
		cfg.Notifiers.Telegram.BotToken = os.Getenv("NOTICOEL_TELEGRAM_BOT_TOKEN")
		cfg.Notifiers.Telegram.ChatID = os.Getenv("NOTICOEL_TELEGRAM_CHAT_ID")

		if cfg.Notifiers.Telegram.BotToken == "" || cfg.Notifiers.Telegram.ChatID == "" {
			return nil, errors.New("NOTICOEL_TELEGRAM_BOT_TOKEN and NOTICOEL_TELEGRAM_CHAT_ID environment variables are required when Telegram is enabled")
		}
	}

	return cfg, nil
}

// fallbackString returns the legacy YAML value when a file was found,
// otherwise def.
func fallbackString(found bool, yamlValue, def string) string {
	if found {
		return yamlValue
	}

	return def
}

// fallbackBool returns the legacy YAML value when a file was found,
// otherwise def.
func fallbackBool(found bool, yamlValue, def bool) bool {
	if found {
		return yamlValue
	}

	return def
}

// fallbackInt returns the legacy YAML value when a file was found,
// otherwise def.
func fallbackInt(found bool, yamlValue, def int) int {
	if found {
		return yamlValue
	}

	return def
}

// envString returns the environment variable's value, or def if unset.
func envString(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return def
}

// envBool returns the environment variable parsed as a bool, or def if
// unset or unparseable.
func envBool(key string, def bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}

	return b
}

// envInt returns the environment variable parsed as an int, or def if
// unset or unparseable.
func envInt(key string, def int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}

	return n
}
