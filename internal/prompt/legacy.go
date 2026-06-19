package prompt

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/satoricorp/gx/internal/termstyle"
)

func requiredValueLegacy(in io.Reader, out io.Writer, label, current string) (string, error) {
	return requiredValueLegacyReader(bufio.NewReader(in), out, label, current)
}

func requiredValueLegacyReader(reader *bufio.Reader, out io.Writer, label, current string) (string, error) {
	for {
		if current != "" {
			fmt.Fprintf(out, "%s [%s]: ", label, current)
		} else {
			fmt.Fprintf(out, "%s: ", label)
		}
		raw, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return "", err
		}
		value := strings.TrimSpace(raw)
		if value == "" {
			value = strings.TrimSpace(current)
		}
		if value != "" {
			return value, nil
		}
		if errors.Is(err, io.EOF) {
			return "", fmt.Errorf("%s is required", strings.ToLower(label))
		}
	}
}

// ModifyCandidate is a change the user can pick when running gx edit.
type ModifyCandidate struct {
	ChangeID    string
	CommitID    string
	Description string
}

func modifyChangeLegacy(in io.Reader, out io.Writer, candidates []ModifyCandidate) (string, error) {
	if len(candidates) == 0 {
		return "", fmt.Errorf("no mutable changes available to edit")
	}
	fmt.Fprintln(out, termstyle.Section("Select change to edit"))
	for i, candidate := range candidates {
		fmt.Fprintf(out, "  %d. %s  [%s %s]\n", i+1, candidate.Description, shortID(candidate.ChangeID, 12), shortID(candidate.CommitID, 8))
	}
	fmt.Fprintf(out, "Enter selection [1-%d, q]: ", len(candidates))

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
	if index, err := strconvAtoi(choice); err == nil {
		if index >= 1 && index <= len(candidates) {
			return candidates[index-1].ChangeID, nil
		}
	}
	for _, candidate := range candidates {
		if strings.HasPrefix(candidate.ChangeID, choice) || strings.HasPrefix(candidate.CommitID, choice) {
			return candidate.ChangeID, nil
		}
	}
	if _, err := strconvAtoi(choice); err == nil {
		return "", fmt.Errorf("selection %s out of range", choice)
	}
	return "", fmt.Errorf("unknown selection %q", choice)
}

// strconvAtoi avoids importing strconv in legacy-only paths used from tests.
func strconvAtoi(s string) (int, error) {
	var n int
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, fmt.Errorf("not a number")
		}
		n = n*10 + int(r-'0')
	}
	return n, nil
}

func shortID(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max]
}

func identityLegacy(in io.Reader, out io.Writer, promptName, promptEmail bool, name, email string, note func(io.Writer)) (string, string, error) {
	if note != nil {
		note(out)
	}
	reader := bufio.NewReader(in)
	var err error
	if promptName {
		name, err = requiredValueLegacyReader(reader, out, "gx name", name)
		if err != nil {
			return "", "", err
		}
	}
	if promptEmail {
		email, err = requiredValueLegacyReader(reader, out, "gx email", email)
		if err != nil {
			return "", "", err
		}
	}
	return name, email, nil
}
