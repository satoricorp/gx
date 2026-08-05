package cli

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/satoricorp/lgtm/internal/codereview"
	"github.com/satoricorp/lgtm/internal/github"
	"github.com/satoricorp/lgtm/internal/vcs"
)

const reviewCommentMarker = "<!-- lgtm review summary -->"

func postReviewSummaryComment(ctx context.Context, repo vcs.RepoInfo, report codereview.Report, stderr io.Writer) {
	remoteURL := pointerString(repo.RemoteURL)
	branchName := pointerString(repo.BranchName)
	if strings.TrimSpace(remoteURL) == "" || strings.TrimSpace(branchName) == "" {
		return
	}
	host, owner, repoName, ok := githubRemoteTarget(remoteURL)
	if !ok {
		return
	}
	client, err := github.NewClient(host)
	if err != nil {
		// No GitHub token. Review is read-only and complete without one, so
		// posting the comment is a bonus, not a failure: skip it quietly
		// rather than nag on every review from an unauthenticated checkout.
		return
	}
	pr, err := client.FindPullRequest(ctx, github.CreatePullRequestOptions{
		Host:       host,
		Owner:      owner,
		Repo:       repoName,
		HeadBranch: branchName,
	})
	if err != nil {
		fmt.Fprintln(stderr, labelWarningValue("Warning", fmt.Sprintf("Could not find pull request for %s: %v", branchName, err)))
		return
	}
	if pr == nil || pr.Number == 0 {
		return
	}
	if err := postReviewInlineComments(ctx, client, repo, owner, repoName, pr.Number, report); err != nil {
		fmt.Fprintln(stderr, labelWarningValue("Warning", fmt.Sprintf("Could not post lgtm inline review comment: %v", err)))
	}
	commentBody := reviewCommentMarker + "\n" + codereview.RenderMarkdown(report)
	_, err = client.UpsertIssueComment(ctx, github.IssueCommentOptions{
		Owner:  owner,
		Repo:   repoName,
		Number: pr.Number,
		Body:   commentBody,
		Marker: reviewCommentMarker,
	})
	if err != nil {
		fmt.Fprintln(stderr, labelWarningValue("Warning", fmt.Sprintf("Could not post lgtm review comment: %v", err)))
	}
}

func postReviewInlineComments(ctx context.Context, client *github.Client, repo vcs.RepoInfo, owner string, repoName string, prNumber int, report codereview.Report) error {
	if client == nil || prNumber <= 0 || len(report.Findings) == 0 {
		return nil
	}
	repoRoot := strings.TrimSpace(repo.RootPath)
	if repoRoot == "" {
		repoRoot = strings.TrimSpace(report.RepoRoot)
	}
	commitID := currentHeadCommit(ctx, repoRoot)
	if strings.TrimSpace(commitID) == "" {
		return nil
	}
	seen := map[string]struct{}{}
	var firstErr error
	for _, finding := range report.Findings {
		for _, anchor := range finding.Anchors {
			if !codereview.AnchorMapsToChangedHunk(ctx, repoRoot, anchor) {
				continue
			}
			key := fmt.Sprintf("%s:%d:%s", anchor.File, anchor.Line, finding.ID)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			if err := client.CreatePullRequestReviewComment(ctx, github.PullRequestReviewCommentOptions{
				Owner:    owner,
				Repo:     repoName,
				Number:   prNumber,
				Body:     renderInlineReviewComment(finding),
				CommitID: commitID,
				Path:     anchor.File,
				Line:     anchor.Line,
				Side:     "RIGHT",
			}); err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func renderInlineReviewComment(finding codereview.Finding) string {
	var b strings.Builder
	title := strings.TrimSpace(finding.Title)
	if title == "" {
		title = "lgtm review finding"
	}
	fmt.Fprintf(&b, "**%s**\n\n", title)
	if summary := strings.TrimSpace(finding.Summary); summary != "" {
		fmt.Fprintf(&b, "%s\n\n", summary)
	}
	if recommendation := strings.TrimSpace(finding.Recommendation); recommendation != "" {
		fmt.Fprintf(&b, "**Do next:** %s\n", recommendation)
	}
	return strings.TrimSpace(b.String())
}

func githubRemoteTarget(remoteURL string) (host, owner, repo string, ok bool) {
	remoteURL = strings.TrimSpace(remoteURL)
	if remoteURL == "" {
		return "", "", "", false
	}
	if idx := strings.Index(remoteURL, "github.com/"); idx >= 0 {
		host = "github.com"
		path := strings.TrimSuffix(remoteURL[idx+len("github.com/"):], ".git")
		parts := strings.Split(path, "/")
		if len(parts) >= 2 && parts[0] != "" && parts[1] != "" {
			return host, parts[0], parts[1], true
		}
	}
	if strings.Contains(remoteURL, "://") {
		parsed, err := url.Parse(remoteURL)
		if err != nil {
			return "", "", "", false
		}
		host = parsed.Hostname()
		if port := parsed.Port(); port != "" {
			host = parsed.Host
		}
		path := strings.TrimPrefix(parsed.Path, "/")
		path = strings.TrimSuffix(path, ".git")
		parts := strings.Split(path, "/")
		if len(parts) >= 2 && parts[0] != "" && parts[1] != "" {
			return host, parts[0], parts[1], true
		}
		return "", "", "", false
	}
	if at := strings.LastIndex(remoteURL, "@"); at >= 0 {
		remoteURL = remoteURL[at+1:]
	}
	if colon := strings.Index(remoteURL, ":"); colon >= 0 {
		host = remoteURL[:colon]
		path := strings.TrimSuffix(remoteURL[colon+1:], ".git")
		parts := strings.Split(path, "/")
		if len(parts) >= 2 && parts[0] != "" && parts[1] != "" {
			return host, parts[0], parts[1], true
		}
	}
	return "", "", "", false
}
