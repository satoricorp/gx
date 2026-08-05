package publication

import (
	"context"

	"github.com/satoricorp/lgtm/internal/reviewbundle"
)

// UpdateGitHubPullRequestBody used to rewrite the GitHub PR body with a rich
// lgtm summary (and previously could seed a thin overview/"Summary" section).
// Summaries are now owned by lgtm Cloud, which appends one rich block below any
// human PR description. The CLI never mutates PR bodies on publish.
func UpdateGitHubPullRequestBody(ctx context.Context, artifact reviewbundle.Artifact) (bool, error) {
	_ = ctx
	_ = artifact
	return false, nil
}
