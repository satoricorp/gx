package git

import "strings"

// RevParseCommit resolves a git ref to a full commit SHA.
func RevParseCommit(repoRoot, ref string) (string, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		ref = "HEAD"
	}
	out, err := runGit(repoRoot, "rev-parse", ref)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}
