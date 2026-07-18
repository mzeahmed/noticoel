// Package logger builds the application's single zap.Logger instance,
// constructed once at startup and passed down to whatever needs it.
package logger

import "go.uber.org/zap"

// New builds a production-configured zap logger.
func New() (*zap.Logger, error) {
	return zap.NewProduction()
}
