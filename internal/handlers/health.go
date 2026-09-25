package handlers

import (
	"context"
	"encoding/json"
	"net/http"
)

// Pinger is the readiness contract: anything that can verify connectivity
// (the real *database.DB, or a fake in tests).
type Pinger interface {
	Ping(ctx context.Context) error
}

// Health responds to liveness probes with a static JSON payload.
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// Ready responds to readiness probes: 200 when the DB ping succeeds, 503
// otherwise. The pod is removed from Service endpoints on 503 but is not
// restarted.
func Ready(pinger Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := pinger.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			if encErr := json.NewEncoder(w).Encode(map[string]string{"error": "database unavailable"}); encErr != nil {
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
