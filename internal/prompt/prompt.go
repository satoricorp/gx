package prompt

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"charm.land/huh/v2"
	"github.com/charmbracelet/x/term"
)

// UseHuh reports whether gx should use huh TUI prompts instead of plain stdin lines.
func UseHuh(in io.Reader) bool {
	if os.Getenv("GX_PLAIN_PROMPTS") != "" {
		return false
	}
	f, ok := in.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(f.Fd())
}

func accessible() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("GX_ACCESSIBLE"))) {
	case "1", "true", "yes":
		return true
	default:
		return strings.TrimSpace(os.Getenv("ACCESSIBLE")) != ""
	}
}

func newForm(groups ...*huh.Group) *huh.Form {
	f := huh.NewForm(groups...)
	f = f.WithTheme(huh.ThemeFunc(huh.ThemeCharm))
	if accessible() {
		f = f.WithAccessible(true)
	}
	return f
}

func mapCancel(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, huh.ErrUserAborted) {
		return context.Canceled
	}
	return err
}

func finalizeRequired(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value != "" {
		return value
	}
	return strings.TrimSpace(fallback)
}
