package event

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"

	"github.com/mzeahmed/noticoel/internal/database/sqlc"
	"github.com/mzeahmed/noticoel/internal/dispatcher"
	"github.com/mzeahmed/noticoel/internal/notifier"
)

// Service contains the business logic of the events module. It consumes
// the sqlc-generated Queries directly, with no repository layer in
// between.
type Service struct {
	queries    *sqlc.Queries
	dispatcher *dispatcher.Dispatcher
	log        *slog.Logger
}

// NewService creates a new events service backed by db. Every event
// successfully persisted is dispatched to disp's registered notifiers.
func NewService(db *sql.DB, disp *dispatcher.Dispatcher, log *slog.Logger) *Service {
	return &Service{queries: sqlc.New(db), dispatcher: disp, log: log}
}

// Create persists e, returns it with its generated fields (ID, CreatedAt)
// populated, and dispatches it to the registered notifiers.
//
// Notifier failures are logged but do not fail the request: the event is
// already durably stored, and delivery is best-effort.
func (s *Service) Create(ctx context.Context, e Event) (Event, error) {
	data, err := marshalData(e.Data)
	if err != nil {
		return Event{}, err
	}

	row, err := s.queries.CreateEvent(ctx, sqlc.CreateEventParams{
		Source:  e.Source,
		Type:    e.Type,
		Status:  e.Status,
		Title:   e.Title,
		Message: e.Message,
		Data:    data,
	})
	if err != nil {
		return Event{}, err
	}

	e.ID = row.ID
	e.CreatedAt = row.CreatedAt

	s.dispatch(ctx, e)

	return e, nil
}

// dispatch sends e to every registered notifier and logs any failure.
func (s *Service) dispatch(ctx context.Context, e Event) {
	results := s.dispatcher.Dispatch(ctx, notifier.Message{
		Status:  e.Status,
		Title:   e.Title,
		Message: e.Message,
	})

	for _, result := range results {
		if !result.Success {
			s.log.Error("notifier failed",
				"notifier", result.Notifier,
				"event_id", e.ID,
				"message", result.Message,
				"error", result.Error,
			)
		}
	}
}

// List returns a page of events, most recently created first.
func (s *Service) List(ctx context.Context, limit, offset int64) ([]Event, error) {
	rows, err := s.queries.GetEvents(ctx, sqlc.GetEventsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}

	events := make([]Event, len(rows))
	for i, row := range rows {
		data, err := unmarshalData(row.Data)
		if err != nil {
			return nil, err
		}

		events[i] = Event{
			ID:        row.ID,
			Source:    row.Source,
			Type:      row.Type,
			Status:    row.Status,
			Title:     row.Title,
			Message:   row.Message,
			Data:      data,
			CreatedAt: row.CreatedAt,
		}
	}

	return events, nil
}

// marshalData encodes data as JSON for storage in the events.data column,
// leaving it NULL when there is nothing to store.
func marshalData(data map[string]string) (sql.NullString, error) {
	if len(data) == 0 {
		return sql.NullString{}, nil
	}

	b, err := json.Marshal(data)
	if err != nil {
		return sql.NullString{}, err
	}

	return sql.NullString{String: string(b), Valid: true}, nil
}

// unmarshalData decodes the events.data column back into a map, returning
// nil when the column is NULL.
func unmarshalData(data sql.NullString) (map[string]string, error) {
	if !data.Valid {
		return nil, nil
	}

	var m map[string]string
	if err := json.Unmarshal([]byte(data.String), &m); err != nil {
		return nil, err
	}

	return m, nil
}
