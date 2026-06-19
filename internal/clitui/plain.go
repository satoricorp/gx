package clitui

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/satoricorp/gx/internal/vcs"
)

// PickModifyRevisionPlain lists revisions for non-interactive terminals.
func PickModifyRevisionPlain(in io.Reader, out io.Writer, snapshot vcs.StatusSnapshot, revisions []vcs.RevisionSnapshot) (string, error) {
	if len(revisions) == 0 {
		return "", fmt.Errorf("no revisions available to edit")
	}
	fmt.Fprintln(out, RenderModifyInteractive(snapshot, revisions, 0))
	fmt.Fprintf(out, "%s ", hint("Enter revision [change id, q]:"))

	reader := bufio.NewReader(in)
	raw, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	choice := strings.TrimSpace(raw)
	if choice == "" {
		return revisions[0].ChangeID, nil
	}
	if choice == "q" || choice == "quit" || choice == "exit" {
		return "", context.Canceled
	}
	resolved, resolveErr := vcs.ResolveStackRevision(choice, unitsFromRevisions(revisions))
	if resolveErr == nil {
		return resolved, nil
	}
	return "", resolveErr
}

func unitsFromRevisions(revisions []vcs.RevisionSnapshot) []vcs.UnitSummary {
	units := make([]vcs.UnitSummary, 0, len(revisions))
	for _, rev := range revisions {
		units = append(units, vcs.UnitSummary{
			Index:       rev.Index,
			ChangeID:    rev.ChangeID,
			CommitID:    rev.CommitID,
			Description: rev.Description,
			Active:      rev.Active,
			Published:   rev.Published,
		})
	}
	return units
}

func SelectModifyRevision(in io.Reader, out io.Writer, snapshot vcs.StatusSnapshot, revisions []vcs.RevisionSnapshot) (string, error) {
	if UseInteractive(in) {
		return PickModifyRevision(in, snapshot, revisions)
	}
	return PickModifyRevisionPlain(in, out, snapshot, revisions)
}
