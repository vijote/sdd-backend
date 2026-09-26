package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/vijote/sdd-backend/internal/config"
	"github.com/vijote/sdd-backend/internal/database"
	"github.com/vijote/sdd-backend/internal/handlers"
	"github.com/vijote/sdd-backend/internal/shortener"
)

// Server wraps an http.Server configured from application config.
type Server struct {
	cfg *config.Config
	db  *database.DB
}

// New creates a Server bound to the given configuration and database handle.
func New(cfg *config.Config, db *database.DB) *Server {
	return &Server{cfg: cfg, db: db}
}

// NewRouter builds the chi router with all application routes.
func NewRouter(cfg *config.Config, db *database.DB) http.Handler {
	r := chi.NewRouter()
	r.Get("/healthz", handlers.Health)
	if db != nil {
		r.Get("/readyz", handlers.Ready(db))
	}
	r.Post("/shorten", handlers.Shorten(shortener.NewService(shortener.NewRepository(db))))
	r.Get("/{code}", handlers.Redirect(shortener.NewService(shortener.NewRepository(db))))
	r.Head("/{code}", handlers.Redirect(shortener.NewService(shortener.NewRepository(db))))
	return r
}

// Start runs the HTTP server and blocks until it is shut down.
func (s *Server) Start() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      NewRouter(s.cfg, s.db),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return srv.ListenAndServe()
}
