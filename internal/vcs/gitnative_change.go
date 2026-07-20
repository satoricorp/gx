package vcs

import (
	"context"
	"fmt"
	"strings"
)

func (s *Service) gitCurrentChange(ctx context.Context, repoRoot, rev string) (ChangeInfo, error) {
	rev = strings.TrimSpace(rev)
	switch rev {
	case "", "@":
		working, err := s.GitWorkingStatus(ctx, repoRoot)
		if err != nil {
			return ChangeInfo{}, err
		}
		files := append(append([]string{}, working.Staged...), working.Unstaged...)
		files = append(files, working.Untracked...)
		return ChangeInfo{
			ChangeID:    "working",
			Description: "",
			Files:       files,
		}, nil
	case "@-":
		return s.gitChangeAtCommit(ctx, repoRoot, "HEAD")
	default:
		return s.gitChangeBySelector(ctx, repoRoot, rev)
	}
}

func (s *Service) gitChangeAtCommit(ctx context.Context, repoRoot, commitRef string) (ChangeInfo, error) {
	commitOID, err := s.runTrimmed(ctx, repoRoot, "git", "rev-parse", commitRef)
	if err != nil {
		return ChangeInfo{}, fmt.Errorf("resolve commit %s: %w", commitRef, err)
	}
	commitOID = strings.TrimSpace(commitOID)
	message, err := s.runTrimmed(ctx, repoRoot, "git", "log", "-1", "--format=%B", commitOID)
	if err != nil {
		return ChangeInfo{}, fmt.Errorf("read commit message: %w", err)
	}
	revisionIDs := ParseRevisionIDsFromMessage(message)
	changeID := commitOID
	if len(revisionIDs) > 0 {
		changeID = revisionIDs[len(revisionIDs)-1]
	}
	parentOID, _ := s.runTrimmed(ctx, repoRoot, "git", "rev-parse", commitOID+"^")
	parentOID = strings.TrimSpace(parentOID)
	var parentChangeID *string
	if parentOID != "" {
		parentMessage, parentErr := s.runTrimmed(ctx, repoRoot, "git", "log", "-1", "--format=%B", parentOID)
		if parentErr == nil {
			if ids := ParseRevisionIDsFromMessage(parentMessage); len(ids) > 0 {
				parentChangeID = &ids[len(ids)-1]
			}
		}
	}
	filesOut, err := s.runTrimmed(ctx, repoRoot, "git", "diff-tree", "--no-commit-id", "--name-only", "-r", commitOID)
	if err != nil {
		return ChangeInfo{}, fmt.Errorf("list commit files: %w", err)
	}
	subject := strings.TrimSpace(strings.Split(message, "\n")[0])
	return ChangeInfo{
		ChangeID:       changeID,
		CommitID:       commitOID,
		Description:    subject,
		ParentChangeID: parentChangeID,
		Files:          splitLines(filesOut),
	}, nil
}

func (s *Service) gitChangeBySelector(ctx context.Context, repoRoot, selector string) (ChangeInfo, error) {
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return ChangeInfo{}, fmt.Errorf("revision selector is required")
	}
	if len(selector) >= 7 {
		if change, err := s.gitChangeAtCommit(ctx, repoRoot, selector); err == nil {
			return change, nil
		}
	}
	if change, err := s.gitChangeByRevisionID(ctx, repoRoot, selector); err == nil {
		return change, nil
	}
	return s.gitChangeAtCommit(ctx, repoRoot, selector)
}

func (s *Service) gitChangeByRevisionID(ctx context.Context, repoRoot, revisionID string) (ChangeInfo, error) {
	out, err := s.runTrimmed(ctx, repoRoot, "git", "log", "-1", "--format=%H", "--grep="+RevisionTrailerLine(revisionID))
	if err != nil {
		return ChangeInfo{}, err
	}
	for _, line := range splitLines(out) {
		if strings.TrimSpace(line) == "" {
			continue
		}
		return s.gitChangeAtCommit(ctx, repoRoot, strings.TrimSpace(line))
	}
	return ChangeInfo{}, fmt.Errorf("unknown revision %q", revisionID)
}

func (s *Service) gitDiffForRev(ctx context.Context, repoRoot, rev string) (string, error) {
	rev = strings.TrimSpace(rev)
	switch rev {
	case "", "@":
		staged, err := s.runTrimmed(ctx, repoRoot, "git", "diff", "--cached")
		if err != nil {
			return "", err
		}
		unstaged, err := s.runTrimmed(ctx, repoRoot, "git", "diff")
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(staged + unstaged), nil
	case "@-":
		return s.runTrimmed(ctx, repoRoot, "git", "show", "--format=", "--patch", "HEAD")
	default:
		if _, err := s.gitChangeBySelector(ctx, repoRoot, rev); err == nil {
			return s.runTrimmed(ctx, repoRoot, "git", "show", "--format=", "--patch", rev)
		}
		return s.runTrimmed(ctx, repoRoot, "git", "show", "--format=", "--patch", rev)
	}
}
