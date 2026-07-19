package forgejo

import (
	"fmt"

	"github.com/mzeahmed/noticoel/internal/modules/event"
)

// toEvent converts a Forgejo release payload into the internal Event
// model. It is the only place that knows both shapes.
func toEvent(p Payload) event.Event {
	return event.Event{
		Source:   "forgejo",
		Category: "release",
		Type:     "release." + p.Action,
		Severity: event.SeverityInfo,
		Title:    fmt.Sprintf("Release %s", p.Action),
		Message:  fmt.Sprintf("%s: release %s", p.Repository.FullName, p.Release.TagName),
		Metadata: map[string]string{
			"repository": p.Repository.FullName,
			"tag":        p.Release.TagName,
			"author":     p.Release.Author.Login,
			"url":        p.Release.HTMLURL,
		},
	}
}
