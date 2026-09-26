package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/vijote/sdd-backend/internal/shortener"
)

// Redirecter is the contract the handler needs from the service layer.
type Redirecter interface {
	Resolve(ctx context.Context, code string) (*shortener.Link, error)
}

// Redirect handles GET /{code}: resolves the code via the service and issues
// a 301 to the stored long URL, or 404 {"error":"not found"} for unknown or
// malformed codes.
func Redirect(svc Redirecter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")

		link, err := svc.Resolve(r.Context(), code)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
			return
		}

		http.Redirect(w, r, link.LongURL, http.StatusMovedPermanently)
	}
}
