package event

import (
	"errors"
	"time"
)

// Event is the internal representation of a notification request.
//
// ID and CreatedAt are populated by Service.Create once the event is
// persisted; a client-supplied value for either is ignored.
type Event struct {
	ID        int64             `json:"id,omitempty"`
	Source    string            `json:"source"`
	Type      string            `json:"type"`
	Status    string            `json:"status"`
	Title     string            `json:"title"`
	Message   string            `json:"message"`
	Data      map[string]string `json:"data,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

// Validate checks that every required field is present.
func (e Event) Validate() error {
	switch {
	case e.Source == "":
		return errors.New("source is required")
	case e.Type == "":
		return errors.New("type is required")
	case e.Status == "":
		return errors.New("status is required")
	case e.Title == "":
		return errors.New("title is required")
	case e.Message == "":
		return errors.New("message is required")
	}

	return nil
}
