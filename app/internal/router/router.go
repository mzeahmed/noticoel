// Package router assembles the application's HTTP handler by wiring up the
// routes exposed by each module.
package router

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/mzeahmed/noticoel/internal/dispatcher"
	"github.com/mzeahmed/noticoel/internal/modules/auth"
	"github.com/mzeahmed/noticoel/internal/modules/event"
	"github.com/mzeahmed/noticoel/internal/modules/health"
)

// New builds and returns the application's top-level http.Handler, with all
// module routes registered on a fresh http.ServeMux.
func New(db *sql.DB, disp *dispatcher.Dispatcher, appVersion, authToken string, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	health.New(appVersion).RegisterRoutes(mux)
	event.New(db, disp, log).RegisterRoutes(mux, auth.Authenticate(authToken))

	return mux
}
