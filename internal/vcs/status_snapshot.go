package vcs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RevisionSnapshot is one ordered revision in a bookmark stack.
type RevisionSnapshot struct {
	Index       int
	ChangeID    string
	CommitID    string
	Description string
	ShortID     string
	Active      bool
	Published   bool
	Working     bool
	SyncNote    string
}

// BookmarkSnapshot summarizes one gx bookmark/stack.
type BookmarkSnapshot struct {
	Stack         StackInfo
	Current       bool
	ChangeCount   int
	ApprovedCount int
	FileCount     int
	Units         []UnitSummary
	Revisions     []RevisionSnapshot
}

// StatusSnapshot is the structured input for stack rendering in internal/clitui.
type StatusSnapshot struct {
	RepoLabel string
	RepoRoot  string
	BaseRef   string
	Bookmarks []BookmarkSnapshot
}

func hasRemoteSync(body StackInfo) bool {
	if body.RemoteRef != nil && strings.TrimSpace(*body.RemoteRef) != "" {
		return true
	}
	return body.Status == "published"
}

func fmtChanges(n int) string {
	if n == 1 {
		return "1 change"
	}
	return fmt.Sprintf("%d changes", n)
}

func fmtApproved(approved, total int) string {
	if total == 0 {
		return "0/0 approved"
	}
	return fmt.Sprintf("%d/%d approved", approved, total)
}

func repoLabel(repo RepoInfo) string {
	if repo.RemoteURL != nil {
		if full := githubRepoFullName(*repo.RemoteURL); full != "" {
			return full
		}
	}
	root := strings.TrimSpace(repo.RootPath)
	if root == "" {
		return "."
	}
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		if rel, relErr := filepath.Rel(home, root); relErr == nil && !strings.HasPrefix(rel, "..") {
			return "~/" + rel
		}
	}
	return root
}

func githubRepoFullName(remoteURL string) string {
	remoteURL = strings.TrimSpace(remoteURL)
	if idx := strings.Index(remoteURL, "github.com"); idx >= 0 {
		suffix := remoteURL[idx+len("github.com"):]
		suffix = strings.TrimPrefix(suffix, ":")
		suffix = strings.TrimPrefix(suffix, "/")
		suffix = strings.TrimSuffix(suffix, ".git")
		parts := strings.Split(suffix, "/")
		if len(parts) >= 2 && parts[0] != "" && parts[1] != "" {
			return parts[0] + "/" + parts[1]
		}
	}
	return ""
}
