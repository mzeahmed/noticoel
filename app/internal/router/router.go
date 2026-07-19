// Package router assembles the application's HTTP handler by wiring up the
// routes exposed by each module.
package router

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/mzeahmed/noticoel/internal/adapters/forgejo"
	"github.com/mzeahmed/noticoel/internal/adapters/gitea"
	"github.com/mzeahmed/noticoel/internal/adapters/github"
	"github.com/mzeahmed/noticoel/internal/adapters/gitlab"
	"github.com/mzeahmed/noticoel/internal/dispatcher"
	"github.com/mzeahmed/noticoel/internal/modules/auth"
	"github.com/mzeahmed/noticoel/internal/modules/event"
	"github.com/mzeahmed/noticoel/internal/modules/health"
)

// New builds and returns the application's top-level http.Handler, with all
// module and adapter routes registered on a fresh http.ServeMux.
//
// Native Event producers (any application that already speaks Noticoel's
// Event model) use event's routes directly. Third-party systems with their
// own payload format go through their adapter's route instead, which
// converts that payload into an Event before it reaches the same event
// pipeline.
func New(db *sql.DB, disp *dispatcher.Dispatcher, appVersion, authToken string, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	authenticate := auth.Authenticate(authToken)

	health.New(appVersion).RegisterRoutes(mux)
	event.New(db, disp, log).RegisterRoutes(mux, authenticate)

	forgejo.New(db, disp, log).RegisterRoutes(mux, authenticate)
	github.New(db, disp, log).RegisterRoutes(mux, authenticate)
	gitlab.New(db, disp, log).RegisterRoutes(mux, authenticate)
	gitea.New(db, disp, log).RegisterRoutes(mux, authenticate)

	return mux
}
