package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/satoricorp/gx/internal/authoring"
)

func promptModifySelection(in io.Reader, out io.Writer, candidates []authoring.ChangeInfo) (string, error) {
	if len(candidates) == 0 {
		return "", fmt.Errorf("no mutable changes available to edit")
	}

	fmt.Fprintln(out, section("Select change to edit"))
	for i, candidate := range candidates {
		num := command(fmt.Sprintf("%2d", i+1))
		desc := candidate.Description
		if desc == "" {
			desc = muted("(no description)")
		}
		meta := muted(fmt.Sprintf("[%s %s]", shortID(candidate.ChangeID, 12), shortID(candidate.CommitID, 8)))
		fmt.Fprintf(out, "  %s  %s  %s\n", num, desc, meta)
	}
	fmt.Fprintf(out, "\n%s ", muted(fmt.Sprintf("Enter selection [1-%d, q]:", len(candidates))))

	reader := bufio.NewReader(in)
	raw, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	choice := strings.TrimSpace(raw)
	if choice == "" {
		return candidates[0].ChangeID, nil
	}
	if choice == "q" || choice == "quit" || choice == "exit" {
		return "", context.Canceled
	}
	if index, err := strconv.Atoi(choice); err == nil {
		if index >= 1 && index <= len(candidates) {
			return candidates[index-1].ChangeID, nil
		}
	}
	for _, candidate := range candidates {
		if strings.HasPrefix(candidate.ChangeID, choice) || strings.HasPrefix(candidate.CommitID, choice) {
			return candidate.ChangeID, nil
		}
	}
	if _, err := strconv.Atoi(choice); err == nil {
		return "", fmt.Errorf("selection %s out of range", choice)
	}
	return "", fmt.Errorf("unknown selection %q", choice)
}

func promptContinueGenerate(in io.Reader, out io.Writer) bool {
	fmt.Fprintf(out, "%s ", muted("Should we continue working to resolve these issues? [y/N]"))
	reader := bufio.NewReader(in)
	raw, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "y", "yes":
		return true
	default:
		return false
	}
}
