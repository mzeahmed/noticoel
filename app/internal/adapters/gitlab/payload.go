// Package gitlab converts GitLab's native webhook payloads into
// Noticoel's internal Event model.
package gitlab

import "errors"

// Payload is the subset of GitLab's Pipeline Hook payload Noticoel cares
// about. GitLab sends additional fields on the real payload;
// encoding/json silently ignores whatever isn't declared here.
//
// GitLab reuses one webhook URL for every event type it's configured to
// send, distinguished by ObjectKind (and the "X-Gitlab-Event" header).
// This adapter only handles "pipeline" — enable just the "Pipeline
// events" trigger on the GitLab side for this URL. Handling another kind
// means adding another Payload variant and branching in Handler.Receive.
type Payload struct {
	ObjectKind       string `json:"object_kind"` // "pipeline", "push", "merge_request"...
	ObjectAttributes struct {
		ID     int64  `json:"id"`
		Ref    string `json:"ref"`
		Status string `json:"status"` // "success", "failed", "canceled"...
	} `json:"object_attributes"`
	Project struct {
		PathWithNamespace string `json:"path_with_namespace"`
		WebURL            string `json:"web_url"`
	} `json:"project"`
	User struct {
		Username string `json:"username"`
	} `json:"user"`
}

// Validate checks that every field the mapper depends on is present.
func (p Payload) Validate() error {
	switch {
	case p.ObjectKind != "pipeline":
		return errors.New(`object_kind must be "pipeline"`)
	case p.Project.PathWithNamespace == "":
		return errors.New("project.path_with_namespace is required")
	case p.ObjectAttributes.Status == "":
		return errors.New("object_attributes.status is required")
	}

	return nil
}
