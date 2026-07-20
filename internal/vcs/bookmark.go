package vcs

import (
	"context"
	"strings"
)

type gitBranchTarget struct {
	Name     string
	CommitID string
}

const (
	legacyGXStackBookmarkPrefix = "gx/"
	legacyGXDraftBookmarkPrefix = "gx/draft/"
)

func isGXStackBookmark(name string) bool {
	name = strings.TrimSpace(name)
	if legacyGXInternalCheckoutRef(name) {
		return false
	}
	return isConventionalStackBookmark(name) || strings.HasPrefix(name, legacyGXStackBookmarkPrefix)
}

func legacyGXInternalCheckoutRef(name string) bool {
	name = strings.TrimSpace(name)
	return name == "gx/base" || name == "gx/edit" ||
		strings.HasPrefix(name, "gx/base/") ||
		strings.HasPrefix(name, "gx/edit/")
}

func baseCheckoutRef(baseRef string) string {
	baseRef = cleanRefName(baseRef)
	if legacyStackBookmarkName(baseRef) {
		return stackBookmarkName(baseRef, "")
	}
	if baseRef == "" {
		return "main"
	}
	return baseRef
}

func gxAuthoringBaseFromCheckoutRef(name string) (string, bool) {
	return "", false
}

func stackBookmarkName(name, headChangeID string) string {
	name = cleanRefName(name)
	name = strings.TrimPrefix(name, legacyGXDraftBookmarkPrefix)
	name = strings.TrimPrefix(name, legacyGXStackBookmarkPrefix)
	if kind, rest, ok := splitConventionalStackBookmark(name); ok {
		slug := bookmarkSlug(rest)
		if slug == "" {
			slug = fallbackStackSlug(headChangeID)
		}
		return kind + "/" + slug
	}
	slug := bookmarkSlug(name)
	if slug == "" {
		slug = fallbackStackSlug(headChangeID)
	}
	return conventionalStackKindForName(name) + "/" + slug
}

func stackNameFromBookmark(name string) string {
	name = cleanRefName(name)
	name = strings.TrimPrefix(name, legacyGXDraftBookmarkPrefix)
	name = strings.TrimPrefix(name, legacyGXStackBookmarkPrefix)
	if _, rest, ok := splitConventionalStackBookmark(name); ok {
		return rest
	}
	return name
}

func legacyStackBookmarkName(name string) bool {
	name = cleanRefName(name)
	return strings.HasPrefix(name, legacyGXDraftBookmarkPrefix) ||
		strings.HasPrefix(name, legacyGXStackBookmarkPrefix)
}

func cleanRefName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.TrimPrefix(name, "refs/heads/")
	name = strings.TrimPrefix(name, "origin/")
	return name
}

func splitConventionalStackBookmark(name string) (string, string, bool) {
	name = cleanRefName(name)
	kind, rest, ok := strings.Cut(name, "/")
	if !ok || strings.TrimSpace(rest) == "" || !isConventionalStackKind(kind) {
		return "", "", false
	}
	return kind, rest, true
}

func isConventionalStackBookmark(name string) bool {
	_, _, ok := splitConventionalStackBookmark(name)
	return ok
}

func isConventionalStackKind(kind string) bool {
	switch strings.TrimSpace(kind) {
	case "feature", "bug", "docs", "test", "chore":
		return true
	default:
		return false
	}
}

func conventionalStackKindForName(name string) string {
	slug := bookmarkSlug(name)
	switch {
	case strings.HasPrefix(slug, "bug-") ||
		strings.HasPrefix(slug, "fix-") ||
		strings.HasPrefix(slug, "repair-") ||
		strings.Contains(slug, "regression") ||
		strings.Contains(slug, "failure") ||
		strings.Contains(slug, "panic") ||
		strings.Contains(slug, "crash"):
		return "bug"
	case strings.HasPrefix(slug, "docs-") ||
		strings.HasPrefix(slug, "doc-") ||
		strings.HasPrefix(slug, "readme-") ||
		strings.Contains(slug, "documentation"):
		return "docs"
	case strings.HasPrefix(slug, "test-") ||
		strings.HasPrefix(slug, "tests-") ||
		strings.Contains(slug, "-test-") ||
		strings.Contains(slug, "-tests-"):
		return "test"
	case strings.HasPrefix(slug, "chore-") ||
		strings.HasPrefix(slug, "deps-") ||
		strings.HasPrefix(slug, "build-") ||
		strings.HasPrefix(slug, "release-"):
		return "chore"
	default:
		return "feature"
	}
}

func fallbackStackSlug(headChangeID string) string {
	if id := shortID(strings.TrimSpace(headChangeID), 8); id != "" {
		return "change-" + id
	}
	return "draft"
}

func (s *Service) gitBranchTargets(ctx context.Context, repoRoot string) ([]gitBranchTarget, error) {
	out, err := s.runStdoutTrimmed(ctx, repoRoot, "git", "for-each-ref", "--format=%(refname:short)|%(objectname)", "refs/heads")
	if err != nil {
		return nil, err
	}
	targets := make([]gitBranchTarget, 0)
	for _, line := range splitLines(out) {
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}
		name := strings.TrimSpace(parts[0])
		commitID := strings.TrimSpace(parts[1])
		if name == "" || commitID == "" {
			continue
		}
		targets = append(targets, gitBranchTarget{Name: name, CommitID: commitID})
	}
	return targets, nil
}
