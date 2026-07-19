// Package gitea converts Gitea's native webhook payloads into Noticoel's
// internal Event model.
package gitea

import "errors"

// Payload is the subset of Gitea's "push" webhook payload Noticoel cares
// about. Gitea sends additional fields on the real payload; encoding/json
// silently ignores whatever isn't declared here.
//
// Extending this adapter to other Gitea events (release,
// workflow_run...) means adding another Payload variant and branching on
// the "X-Gitea-Event" header in Handler.Receive — it does not affect any
// other package.
type Payload struct {
	Ref        string `json:"ref"` // "refs/heads/main"
	Repository struct {
		FullName string `json:"full_name"`
		HTMLURL  string `json:"html_url"`
	} `json:"repository"`
	Pusher struct {
		Login string `json:"login"`
	} `json:"pusher"`
	Commits []struct {
		ID      string `json:"id"`
		Message string `json:"message"`
	} `json:"commits"`
}

// Validate checks that every field the mapper depends on is present.
func (p Payload) Validate() error {
	switch {
	case p.Ref == "":
		return errors.New("ref is required")
	case p.Repository.FullName == "":
		return errors.New("repository.full_name is required")
	}

	return nil
}
