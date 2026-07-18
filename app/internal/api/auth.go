package api

import (
	"net/http"
	"strings"
)

const bearerPrefix = "Bearer "

// authenticate rejects any request whose Authorization header does not
// carry the configured bearer token.
func authenticate(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			if !strings.HasPrefix(header, bearerPrefix) || strings.TrimPrefix(header, bearerPrefix) != token {
				writeError(w, http.StatusUnauthorized, "invalid or missing bearer token")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
