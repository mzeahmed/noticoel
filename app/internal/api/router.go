// Package api exposes Noticeal's HTTP surface: a health check, a version
// endpoint, and the event ingestion endpoint.
package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func newRouter(appVersion, authToken string, log *slog.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", handleHealth)
	r.Get("/version", handleVersion(appVersion))

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(authenticate(authToken))
		r.Post("/events", handleCreateEvent(log))
	})

	return r
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func handleVersion(appVersion string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(appVersion))
	}
}
