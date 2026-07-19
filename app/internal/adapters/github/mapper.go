package github

import (
	"fmt"

	"github.com/mzeahmed/noticoel/internal/modules/event"
)

// toEvent converts a GitHub workflow_run payload into the internal Event
// model. It is the only place that knows both shapes.
func toEvent(p Payload) event.Event {
	// The run may still be in progress, in which case Conclusion is empty
	// and Action ("requested", "in_progress"...) is the best summary.
	outcome := p.WorkflowRun.Conclusion
	if outcome == "" {
		outcome = p.Action
	}

	return event.Event{
		Source:   "github",
		Category: "ci",
		Type:     "workflow_run." + outcome,
		Severity: severityFor(outcome),
		Title:    fmt.Sprintf("Workflow run %s: %s", outcome, p.WorkflowRun.Name),
		Message:  fmt.Sprintf("%s: %s on %s", p.Repository.FullName, p.WorkflowRun.Name, p.WorkflowRun.HeadBranch),
		Metadata: map[string]string{
			"repository": p.Repository.FullName,
			"branch":     p.WorkflowRun.HeadBranch,
			"workflow":   p.WorkflowRun.Name,
			"commit":     p.WorkflowRun.HeadSHA,
			"author":     p.Sender.Login,
			"url":        p.WorkflowRun.HTMLURL,
		},
	}
}

// severityFor maps a GitHub workflow_run conclusion (or, if not yet
// concluded, its action) to a Severity.
func severityFor(outcome string) event.Severity {
	switch outcome {
	case "success":
		return event.SeverityInfo
	case "failure", "timed_out":
		return event.SeverityError
	case "cancelled":
		return event.SeverityWarning
	default:
		return event.SeverityWarning
	}
}
