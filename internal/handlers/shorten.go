package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vijote/sdd-backend/internal/shortener"
)

// Shortener is the contract the handler needs from the service layer.
type Shortener interface {
	Shorten(ctx context.Context, longURL string) (*shortener.Link, error)
}

// Shorten handles POST /api/shorten: decodes {"url": "..."}, delegates to the
// service, and returns 201 {"code","long_url"} or 400 {"error"}.
func Shorten(svc Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid json body"})
			return
		}

		link, err := svc.Shorten(r.Context(), req.URL)
		if err != nil {
			if errors.Is(err, shortener.ErrInvalidURL) {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid url"})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal error"})
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":     link.Code,
			"long_url": link.LongURL,
		})
	}
}
