// Package forgejo converts Forgejo's native webhook payloads into
// Noticoel's internal Event model.
package forgejo

import "errors"

// Payload is the subset of Forgejo's "release" webhook payload Noticoel
// cares about. Forgejo (Gitea-compatible) sends additional fields on the
// real payload; encoding/json silently ignores whatever isn't declared
// here, so only fields we actually use need to be listed.
//
// Extending this adapter to other Forgejo webhook events (push,
// workflow_run...) means adding another Payload variant and branching on
// the "X-Forgejo-Event" header in Handler.Receive — it does not affect any
// other package.
type Payload struct {
	Action  string `json:"action"` // "published", "updated", "deleted"...
	Release struct {
		TagName string `json:"tag_name"`
		Name    string `json:"name"`
		HTMLURL string `json:"html_url"`
		Author  struct {
			Login string `json:"login"`
		} `json:"author"`
	} `json:"release"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

// Validate checks that every field the mapper depends on is present.
func (p Payload) Validate() error {
	switch {
	case p.Action == "":
		return errors.New("action is required")
	case p.Repository.FullName == "":
		return errors.New("repository.full_name is required")
	case p.Release.TagName == "":
		return errors.New("release.tag_name is required")
	}

	return nil
}
