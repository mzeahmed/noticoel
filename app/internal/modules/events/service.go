package events

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/mzeahmed/noticoel/internal/database/sqlc"
)

// Service contains the business logic of the events module. It consumes
// the sqlc-generated Queries directly, with no repository layer in
// between.
type Service struct {
	queries *sqlc.Queries
}

// NewService creates a new events service backed by db.
func NewService(db *sql.DB) *Service {
	return &Service{queries: sqlc.New(db)}
}

// Create persists e and returns it with its generated fields (ID,
// CreatedAt) populated.
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

	return e, nil
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
