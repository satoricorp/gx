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
	summary, err := GitHubPullRequestBodyFromArtifact(ctx, artifact)
	if err != nil {
		return false, fmt.Errorf("build GitHub PR summary for %s: %w", ref.URL, err)
	}
	client, err := githubapi.NewClient(ref.Host)
	if err != nil {
		return false, fmt.Errorf("prepare GitHub PR summary update for %s: %w", ref.URL, err)
	}
	existing, err := client.GetPullRequest(ctx, ref.Owner, ref.Repo, ref.Number)
	if err != nil {
		return false, fmt.Errorf("load GitHub PR before summary update for %s: %w", ref.URL, err)
	}
	if existing == nil {
		return false, nil
	}
	authorNotes, gxOwned := splitGeneratedPRBody(existing.Body)
	desired := appendAuthorNotes(summary, authorNotes)
	if strings.TrimSpace(existing.Body) == strings.TrimSpace(desired) {
		return false, nil
	}
	if !gxOwned && strings.TrimSpace(authorNotes) == "" {
		authorNotes = strings.TrimSpace(existing.Body)
		desired = appendAuthorNotes(summary, authorNotes)
	}
	if _, err := client.UpdatePullRequest(ctx, githubapi.UpdatePullRequestOptions{
		Owner:  ref.Owner,
		Repo:   ref.Repo,
		Number: ref.Number,
		Body:   desired,
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
