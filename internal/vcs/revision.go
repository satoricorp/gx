package vcs

import (
	"fmt"
	"strconv"
	"strings"
)

// ResolveStackRevision maps lgtm stack labels like r2 to a JJ change ID.
func ResolveStackRevision(rev string, units []UnitSummary) (string, error) {
	rev = strings.TrimSpace(rev)
	if rev == "" {
		return "", fmt.Errorf("revision is required")
	}
	if strings.HasPrefix(strings.ToLower(rev), "r") {
		indexText := strings.TrimPrefix(strings.ToLower(rev), "r")
		index, err := strconv.Atoi(indexText)
		if err != nil || index <= 0 {
			return "", fmt.Errorf("unknown revision %q", rev)
		}
		for _, unit := range units {
			if unit.Index == index {
				return unit.ChangeID, nil
			}
		}
		return "", fmt.Errorf("revision %q is not in the current bookmark stack", rev)
	}
	for _, unit := range units {
		if unit.ChangeID == rev || strings.HasPrefix(unit.ChangeID, rev) {
			return unit.ChangeID, nil
		}
		if unit.CommitID == rev || strings.HasPrefix(unit.CommitID, rev) {
			return unit.ChangeID, nil
		}
	}
	return rev, nil
}
