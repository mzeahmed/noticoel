// Package app wires Noticeal's startup sequence together: configuration,
// logging, database, migrations, and the HTTP server.
package app

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/mzeahmed/noticeal/internal/api"
	"github.com/mzeahmed/noticeal/internal/config"
	"github.com/mzeahmed/noticeal/internal/database"
	"github.com/mzeahmed/noticeal/internal/logger"
	"github.com/mzeahmed/noticeal/internal/version"
)

// Run loads the configuration, initializes the logger and database, runs
// pending migrations, then starts the HTTP server. It blocks until the
// server stops, and returns an error instead of panicking or exiting.
func Run() error {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log, err := logger.New()
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	defer func() { _ = log.Sync() }()

	db, err := database.Open(cfg.Database.Path)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() { _ = db.Close() }()

	if err := database.Migrate(db); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	server := api.NewServer(cfg.Server.Addr(), version.Version, cfg.Auth.Token, log)

	log.Info("starting noticeal",
		zap.String("version", version.Version),
		zap.String("addr", cfg.Server.Addr()),
	)

	if err := server.Start(); err != nil {
		return fmt.Errorf("start server: %w", err)
	}

	return nil
}
