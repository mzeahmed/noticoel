package event

import (
	"errors"
	"fmt"
	"time"
)

// Severity represents how urgent an event is. It replaces the earlier
// free-form "status" string with a closed set of values every notifier and
// future routing rule can rely on.
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityError    Severity = "error"
	SeverityCritical Severity = "critical"
)

// Valid reports whether s is one of the defined Severity values.
func (s Severity) Valid() bool {
	switch s {
	case SeverityInfo, SeverityWarning, SeverityError, SeverityCritical:
		return true
	default:
		return false
	}
}

// Event is the internal representation of a notification request. It is
// the common shape every event producer (Forgejo, Yoostart, a monitoring
// system, a cron job...) publishes to Noticoel.
//
// ID and CreatedAt are populated by Service.Create once the event is
// persisted; a client-supplied value for either is ignored.
type Event struct {
	ID int64 `json:"id,omitempty"`

	// Source identifies the application or service that published the
	// event, e.g. "forgejo", "yoostart", "bookingapp".
	Source string `json:"source"`

	// Category groups related event types, e.g. "billing", "ci", "auth",
	// "monitoring". Optional: simple producers with a single kind of
	// event don't need it.
	Category string `json:"category,omitempty"`

	// Type is the specific event within Source/Category, e.g.
	// "subscription.created", "workflow.failed", "user.login".
	Type string `json:"type"`

	// Severity is how urgent the event is. It drives how notifiers
	// present the event (and, later, routing rules).
	Severity Severity `json:"severity"`

	Title   string `json:"title"`
	Message string `json:"message"`

	// Metadata carries arbitrary producer-specific context (repository,
	// tag, url, plan, amount...) that doesn't belong in Title/Message.
	Metadata map[string]string `json:"metadata,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

// Validate checks that every required field is present and well-formed.
func (e Event) Validate() error {
	switch {
	case e.Source == "":
		return errors.New("source is required")
	case e.Type == "":
		return errors.New("type is required")
	case e.Severity == "":
		return errors.New("severity is required")
	case !e.Severity.Valid():
		return fmt.Errorf("severity must be one of: info, warning, error, critical (got %q)", e.Severity)
	case e.Title == "":
		return errors.New("title is required")
	case e.Message == "":
		return errors.New("message is required")
	}

	return nil
}
