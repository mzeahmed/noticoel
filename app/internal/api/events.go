package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/mzeahmed/noticeal/internal/event"
)

func handleCreateEvent(log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var e event.Event
		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		if err := e.Validate(); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		log.Info("event received",
			zap.String("source", e.Source),
			zap.String("type", e.Type),
			zap.String("status", e.Status),
		)

		writeJSON(w, http.StatusAccepted, map[string]string{"status": "accepted"})
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
