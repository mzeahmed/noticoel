package gitlab

import (
	"fmt"

	"github.com/mzeahmed/noticoel/internal/modules/event"
)

// toEvent converts a GitLab pipeline payload into the internal Event
// model. It is the only place that knows both shapes.
func toEvent(p Payload) event.Event {
	status := p.ObjectAttributes.Status

	return event.Event{
		Source:   "gitlab",
		Category: "ci",
		Type:     "pipeline." + status,
		Severity: severityFor(status),
		Title:    fmt.Sprintf("Pipeline %s", status),
		Message:  fmt.Sprintf("%s: pipeline on %s", p.Project.PathWithNamespace, p.ObjectAttributes.Ref),
		Metadata: map[string]string{
			"project": p.Project.PathWithNamespace,
			"ref":     p.ObjectAttributes.Ref,
			"author":  p.User.Username,
			"url":     fmt.Sprintf("%s/-/pipelines/%d", p.Project.WebURL, p.ObjectAttributes.ID),
		},
	}
}

// severityFor maps a GitLab pipeline status to a Severity.
func severityFor(status string) event.Severity {
	switch status {
	case "success":
		return event.SeverityInfo
	case "failed":
		return event.SeverityError
	case "canceled", "skipped":
		return event.SeverityWarning
	default:
		return event.SeverityWarning
	}
}
