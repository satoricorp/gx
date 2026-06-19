package prompt

import (
	"fmt"
	"io"

	"charm.land/huh/v2"
)

// ModifyChange asks the user to pick a JJ change to modify.
func ModifyChange(in io.Reader, out io.Writer, candidates []ModifyCandidate) (string, error) {
	if len(candidates) == 0 {
		return "", fmt.Errorf("no mutable changes available to edit")
	}
	if !UseHuh(in) {
		return modifyChangeLegacy(in, out, candidates)
	}

	selected := candidates[0].ChangeID
	options := make([]huh.Option[string], len(candidates))
	for i, candidate := range candidates {
		label := fmt.Sprintf("%d. %s  [%s %s]", i+1, candidate.Description, shortID(candidate.ChangeID, 12), shortID(candidate.CommitID, 8))
		options[i] = huh.NewOption(label, candidate.ChangeID)
	}

	height := len(candidates)
	if height > 12 {
		height = 12
	}

	field := huh.NewSelect[string]().
		Title("Select change to edit").
		Description("Type to filter • ↑↓ to browse • enter to select • ctrl+c to cancel").
		Options(options...).
		Value(&selected).
		Height(height)

	err := newForm(huh.NewGroup(field)).Run()
	if err != nil {
		return "", mapCancel(err)
	}
	return selected, nil
}
