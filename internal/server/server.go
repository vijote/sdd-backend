package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/vijote/sdd-backend/internal/config"
	"github.com/vijote/sdd-backend/internal/handlers"
)

// Server wraps an http.Server configured from application config.
type Server struct {
	cfg *config.Config
}

// New creates a Server bound to the given configuration.
func New(cfg *config.Config) *Server {
	return &Server{cfg: cfg}
}

// NewRouter builds the chi router with all application routes.
func NewRouter(cfg *config.Config) http.Handler {
	r := chi.NewRouter()
	r.Get("/healthz", handlers.Health)
	return r
}

// Start runs the HTTP server and blocks until it is shut down.
func (s *Server) Start() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      NewRouter(s.cfg),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return srv.ListenAndServe()
}
