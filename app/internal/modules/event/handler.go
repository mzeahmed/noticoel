package event

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/mzeahmed/coelakit/request"
	"github.com/mzeahmed/coelakit/response"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

// Handler handles all HTTP requests related to the events module.
type Handler struct {
	service *Service
	log     *slog.Logger
}

// NewHandler creates a new events handler.
func NewHandler(service *Service, log *slog.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// Create handles POST /api/v1/events.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var e Event
	if !request.Bind(w, r, &e) {
		return
	}

	created, err := h.service.Create(r.Context(), e)
	if err != nil {
		h.log.Error("failed to store event", "error", err)
		response.Error(w, http.StatusInternalServerError, "internal server error")

		return
	}

	h.log.Info("event received",
		"id", created.ID,
		"source", created.Source,
		"type", created.Type,
		"severity", created.Severity,
	)

	response.JSON(w, http.StatusAccepted, created)
}

// List handles GET /api/v1/events.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parsePagination(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	events, err := h.service.List(r.Context(), limit, offset)
	if err != nil {
		h.log.Error("failed to list events", "error", err)
		response.Error(w, http.StatusInternalServerError, "internal server error")

		return
	}

	response.JSON(w, http.StatusOK, events)
}

// parsePagination reads the limit/offset query parameters, defaulting to
// defaultLimit and 0, and capping limit at maxLimit.
func parsePagination(r *http.Request) (limit, offset int64, err error) {
	limit = defaultLimit

	if v := r.URL.Query().Get("limit"); v != "" {
		limit, err = strconv.ParseInt(v, 10, 64)
		if err != nil || limit <= 0 {
			return 0, 0, errors.New("limit must be a positive integer")
		}

		if limit > maxLimit {
			limit = maxLimit
		}
	}

	if v := r.URL.Query().Get("offset"); v != "" {
		offset, err = strconv.ParseInt(v, 10, 64)
		if err != nil || offset < 0 {
			return 0, 0, errors.New("offset must be a non-negative integer")
		}
	}

	return limit, offset, nil
}
