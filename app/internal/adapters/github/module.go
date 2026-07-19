package github

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/mzeahmed/noticoel/internal/dispatcher"
	"github.com/mzeahmed/noticoel/internal/modules/event"
)

// Module wires together the GitHub adapter's dependencies and exposes
// its HTTP route.
type Module struct {
	handler *Handler
}

// New builds a GitHub adapter Module. Events converted from GitHub's
// native payload are persisted and dispatched through disp's registered
// notifiers, exactly like events submitted through the generic API.
func New(db *sql.DB, disp *dispatcher.Dispatcher, log *slog.Logger) *Module {
	service := event.NewService(db, disp, log)

	return &Module{handler: NewHandler(service, log)}
}

// RegisterRoutes registers the GitHub adapter's route on the given
// mux.
//
// authenticate guards the route, requiring a valid bearer token; the
// caller (see router.New) is expected to pass auth.Authenticate(token).
func (m *Module) RegisterRoutes(mux *http.ServeMux, authenticate func(http.Handler) http.Handler) {
	mux.Handle("POST /api/v1/adapters/github", authenticate(http.HandlerFunc(m.handler.Receive)))
}
