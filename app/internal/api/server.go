package api

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Server wraps the HTTP server serving Noticeal's API.
type Server struct {
	httpServer *http.Server
}

// NewServer builds an HTTP server listening on addr, serving the router.
func NewServer(addr, appVersion, authToken string, log *zap.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           newRouter(appVersion, authToken, log),
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      10 * time.Second,
		},
	}
}

// Start begins serving HTTP requests. It blocks until the server stops.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
