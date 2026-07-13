package publication

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/satoricorp/gx/internal/reviewbundle"
)

type githubPullRequestRef struct {
	Host   string
	Owner  string
	Repo   string
	Number int
	URL    string
}

// UpdateGitHubPullRequestBody used to rewrite the GitHub PR body with a rich
// GX summary (and previously could seed a thin overview/"Summary" section).
// Summaries are now owned by GX Cloud, which appends one rich block below any
// human PR description. The CLI never mutates PR bodies on publish.
func UpdateGitHubPullRequestBody(ctx context.Context, artifact reviewbundle.Artifact) (bool, error) {
	_ = ctx
	_ = artifact
	return false, nil
}

func githubPullRequestRefFromArtifact(artifact reviewbundle.Artifact) (githubPullRequestRef, bool) {
	prURL := pullRequestURL(artifact)
	if prURL == "" {
		return githubPullRequestRef{}, false
	}
	parsed, err := url.Parse(prURL)
	if err != nil || parsed.Host == "" {
		return githubPullRequestRef{}, false
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) < 4 || parts[0] == "" || parts[1] == "" || parts[2] != "pull" {
		return githubPullRequestRef{}, false
	}
	number, err := strconv.Atoi(parts[3])
	if err != nil || number <= 0 {
		return githubPullRequestRef{}, false
	}
	return githubPullRequestRef{
		Host:   parsed.Host,
		Owner:  parts[0],
		Repo:   strings.TrimSuffix(parts[1], ".git"),
		Number: number,
		URL:    strings.TrimSpace(prURL),
	}, true
}
