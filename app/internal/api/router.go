// Package api exposes Noticeal's minimal HTTP surface: a health check and
// a version endpoint.
package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func newRouter(appVersion string) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", handleHealth)
	r.Get("/version", handleVersion(appVersion))

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
