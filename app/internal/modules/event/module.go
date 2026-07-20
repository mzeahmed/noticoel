// Package events receives notification events over HTTP.
package event

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/mzeahmed/noticoel/internal/dispatcher"
)

// Module wires together the events module's dependencies and exposes its
// HTTP routes.
type Module struct {
	handler *Handler
}

// New builds an events Module with its service and handler dependencies
// initialized. Events created through this module are dispatched to disp's
// registered notifiers.
func New(db *sql.DB, disp *dispatcher.Dispatcher, log *slog.Logger) *Module {
	service := NewService(db, disp, log)
	handler := NewHandler(service, log)

	return &Module{
		handler: handler,
	}
}

// RegisterRoutes registers the events module's routes on the given mux.
//
// authenticate guards the route, requiring a valid bearer token; the
// caller (see router.New) is expected to pass auth.Authenticate(token).
func (m *Module) RegisterRoutes(mux *http.ServeMux, authenticate func(http.Handler) http.Handler) {
	mux.Handle("POST /api/v1/events/create", authenticate(http.HandlerFunc(m.handler.Create)))
	mux.Handle("GET /api/v1/events/list", authenticate(http.HandlerFunc(m.handler.List)))
}
