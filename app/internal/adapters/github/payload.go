// Package github converts GitHub's native webhook payloads into
// Noticoel's internal Event model.
package github

import "errors"

// Payload is the subset of GitHub's "workflow_run" webhook payload
// Noticoel cares about. GitHub sends additional fields on the real
// payload; encoding/json silently ignores whatever isn't declared here.
//
// Extending this adapter to other GitHub events (push, release...)
// means adding another Payload variant and branching on the
// "X-GitHub-Event" header in Handler.Receive — it does not affect any
// other package.
type Payload struct {
	Action      string `json:"action"` // "requested", "in_progress", "completed"...
	WorkflowRun struct {
		Name       string `json:"name"`
		HTMLURL    string `json:"html_url"`
		Status     string `json:"status"`
		Conclusion string `json:"conclusion"` // "success", "failure", "cancelled"... empty until completed
		HeadBranch string `json:"head_branch"`
		HeadSHA    string `json:"head_sha"`
	} `json:"workflow_run"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
}

// Validate checks that every field the mapper depends on is present.
func (p Payload) Validate() error {
	switch {
	case p.Action == "":
		return errors.New("action is required")
	case p.Repository.FullName == "":
		return errors.New("repository.full_name is required")
	case p.WorkflowRun.Name == "":
		return errors.New("workflow_run.name is required")
	}

	return nil
}
