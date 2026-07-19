package gitea

import (
	"fmt"
	"strings"

	"github.com/mzeahmed/noticoel/internal/modules/event"
)

// toEvent converts a Gitea push payload into the internal Event model. It
// is the only place that knows both shapes.
func toEvent(p Payload) event.Event {
	branch := strings.TrimPrefix(p.Ref, "refs/heads/")

	var lastCommit string
	if n := len(p.Commits); n > 0 {
		lastCommit = p.Commits[n-1].ID
	}

	return event.Event{
		Source:   "gitea",
		Category: "vcs",
		Type:     "push",
		Severity: event.SeverityInfo,
		Title:    fmt.Sprintf("Push to %s", branch),
		Message:  fmt.Sprintf("%s: %d commit(s) pushed by %s", p.Repository.FullName, len(p.Commits), p.Pusher.Login),
		Metadata: map[string]string{
			"repository": p.Repository.FullName,
			"branch":     branch,
			"commit":     lastCommit,
			"author":     p.Pusher.Login,
			"url":        p.Repository.HTMLURL,
		},
	}
}
