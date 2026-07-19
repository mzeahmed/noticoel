package notifier

import "time"

type Result struct {
	Notifier string
	Success  bool
	Message  string
	Error    error
	SentAt   time.Time
}
