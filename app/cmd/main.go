package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mzeahmed/coelakit/middleware"
	"github.com/mzeahmed/coelakit/server"
	"github.com/mzeahmed/noticoel/internal/config"
	"github.com/mzeahmed/noticoel/internal/database"
	"github.com/mzeahmed/noticoel/internal/dispatcher"
	"github.com/mzeahmed/noticoel/internal/logger"
	"github.com/mzeahmed/noticoel/internal/notifier/telegram"
	"github.com/mzeahmed/noticoel/internal/router"
	"github.com/mzeahmed/noticoel/internal/version"
)

func main() {
	if err := runConfig(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Run loads the configuration, initializes the logger and database, runs
// pending migrations, then starts the HTTP server. It blocks until the
// server stops, and returns an error instead of panicking or exiting.
func runConfig() error {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log := logger.New(cfg.Debug)

	disp := dispatcher.New()

	if cfg.Notifiers.Telegram.Enabled {
		disp.Register(telegram.New(telegram.Config{
			Enabled:  cfg.Notifiers.Telegram.Enabled,
			BotToken: cfg.Notifiers.Telegram.BotToken,
			ChatID:   cfg.Notifiers.Telegram.ChatID,
		}))
	}

	db, err := database.Open(cfg.Database.Path)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() { _ = db.Close() }()

	if err := database.Migrate(db); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	handler := router.New(db, disp, version.Version, cfg.Auth.Token, log)
	handler = middleware.LoggingWith(log)(middleware.RecoveryWith(log)(handler))

	log.Info("starting noticoel",
		"version", version.Version,
		"addr", cfg.Server.Addr(),
	)

	if err := server.Run(server.Config{
		Addr:         cfg.Server.Addr(),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}); err != nil {
		return fmt.Errorf("start server: %w", err)
	}

	return nil
}
