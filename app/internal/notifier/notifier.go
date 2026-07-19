package notifier

import "context"

// Message is the notifier-facing view of an event: just enough to compose
// a notification, decoupled from the event module's own Event type so
// notifier/dispatcher never depends on it (event.Service depends on
// dispatcher, not the other way around).
type Message struct {
	Severity string
	Title    string
	Message  string
}

type Notifier interface {
	Name() string

	Notify(
		ctx context.Context,
		msg Message,
	) Result
}
