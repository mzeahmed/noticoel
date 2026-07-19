package forgejo

import (
	"log/slog"
	"net/http"

	"github.com/mzeahmed/coelakit/request"
	"github.com/mzeahmed/coelakit/response"
	"github.com/mzeahmed/noticoel/internal/modules/event"
)

// Handler handles Forgejo webhook requests.
type Handler struct {
	service *event.Service
	log     *slog.Logger
}

// NewHandler creates a new Forgejo webhook handler.
func NewHandler(service *event.Service, log *slog.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// Receive handles POST /api/v1/adapters/forgejo. It decodes Forgejo's
// native payload, converts it into an Event, then hands it to the same
// event pipeline used by the generic events API.
func (h *Handler) Receive(w http.ResponseWriter, r *http.Request) {
	var p Payload
	if !request.Bind(w, r, &p) {
		return
	}

	created, err := h.service.Create(r.Context(), toEvent(p))
	if err != nil {
		h.log.Error("failed to store event", "adapter", "forgejo", "error", err)
		response.Error(w, http.StatusInternalServerError, "internal server error")

		return
	}

	h.log.Info("event received",
		"adapter", "forgejo",
		"id", created.ID,
		"type", created.Type,
	)

	response.JSON(w, http.StatusAccepted, created)
}
