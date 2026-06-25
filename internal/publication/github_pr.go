package publication

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	githubapi "github.com/satoricorp/gx/internal/github"
	"github.com/satoricorp/gx/internal/reviewbundle"
)

type githubPullRequestRef struct {
	Host   string
	Owner  string
	Repo   string
	Number int
	URL    string
}

func UpdateGitHubPullRequestBody(ctx context.Context, artifact reviewbundle.Artifact) (bool, error) {
	ref, ok := githubPullRequestRefFromArtifact(artifact)
	if !ok {
		return false, nil
	}
	body := GitHubPullRequestBodyFromArtifact(artifact)
	client, err := githubapi.NewClient(ref.Host)
	if err != nil {
		return false, fmt.Errorf("prepare GitHub PR summary update for %s: %w", ref.URL, err)
	}
	existing, err := client.GetPullRequest(ctx, githubapi.GetPullRequestOptions{
		Owner:  ref.Owner,
		Repo:   ref.Repo,
		Number: ref.Number,
	})
	if err != nil {
		return false, fmt.Errorf("load GitHub PR before summary update for %s: %w", ref.URL, err)
	}
	if existing != nil && !shouldUpdateGitHubPullRequestBody(existing.Body, body) {
		return false, nil
	}
	if _, err := client.UpdatePullRequest(ctx, githubapi.UpdatePullRequestOptions{
		Owner:  ref.Owner,
		Repo:   ref.Repo,
		Number: ref.Number,
		Body:   body,
	}); err != nil {
		return false, fmt.Errorf("update GitHub PR summary for %s: %w", ref.URL, err)
	}
	return true, nil
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

func shouldUpdateGitHubPullRequestBody(existing, desired string) bool {
	existing = strings.TrimSpace(existing)
	desired = strings.TrimSpace(desired)
	if desired == "" || existing == desired {
		return false
	}
	return existing == "" || strings.HasPrefix(existing, "Published by GX.")
}
