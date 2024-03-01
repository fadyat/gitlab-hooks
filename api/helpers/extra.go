package helpers

import (
	"errors"
	"fmt"
	"github.com/fadyat/hooks/api"
	"github.com/fadyat/hooks/api/entities"
	"strings"
)

// GetBranchNameFromRef returns the branch name from the ref
//
// Example: refs/heads/feature-123 -> feature-123
func GetBranchNameFromRef(ref string) string {
	return strings.TrimPrefix(ref, "refs/heads/")
}

// isMergeCommit checks if the commit message is a merge commit
//
// Ignoring merge commits for a push events, because they are not related to the task
func isMergeCommit(message string) bool {
	return strings.HasPrefix(message, "Merge branch")
}

// isCustomMergeCommit checks if the commit message is an imitated merge commit
//
// Used in /api/handlers/gitlab/last_commit_mergebased.go
func isCustomMergeCommit(message string) bool {
	return strings.Contains(message, "is merged into")
}

// ConfigureMessage configures the message for the comment
//
// If the message is a merge commit, returns an error, because merge commits
// are handled by the different webhook.
// Done to avoid duplicate comments on push and merge hooks.
func ConfigureMessage(msg entities.Message) (s string, err error) {
	// todo: add the author of the commit to the end
	//  of the message
	// todo: can provide feature flags to add the author
	//  of the commit to the end of the message

	if isCustomMergeCommit(msg.Text) {
		return fmt.Sprintf("%s\n\n%s", msg.URL, msg.Text), nil
	}

	// todo: can be removed if only the push hook is used
	//  or when custom fields are updated only, without creating comments
	if isMergeCommit(msg.Text) {
		return "", errors.New(api.MergeCommitUnsupported)
	}

	return fmt.Sprintf("%s\n\n%s", msg.URL, removeTaskMentions(msg.Text)), nil
}
