package prompt

import (
	"fmt"
	"io"
	"strings"

	"charm.land/huh/v2"
)

// Identity collects gx name and email when flags were not provided.
func Identity(in io.Reader, out io.Writer, promptName, promptEmail bool, name, email string, note func(io.Writer)) (string, string, error) {
	if !promptName && !promptEmail {
		return name, email, nil
	}
	if !UseHuh(in) {
		return identityLegacy(in, out, promptName, promptEmail, name, email, note)
	}
	if note != nil {
		note(out)
	}

	nameDefault := name
	emailDefault := email

	var fields []huh.Field
	if promptName {
		fields = append(fields,
			huh.NewInput().
				Title("gx name").
				Value(&name).
				Placeholder(strings.TrimSpace(name)).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" && strings.TrimSpace(nameDefault) == "" {
						return fmt.Errorf("gx name is required")
					}
					return nil
				}),
		)
	}
	if promptEmail {
		fields = append(fields,
			huh.NewInput().
				Title("gx email").
				Value(&email).
				Placeholder(strings.TrimSpace(email)).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" && strings.TrimSpace(emailDefault) == "" {
						return fmt.Errorf("gx email is required")
					}
					if strings.TrimSpace(s) != "" && !strings.Contains(s, "@") {
						return fmt.Errorf("enter a valid email address")
					}
					return nil
				}),
		)
	}

	err := newForm(huh.NewGroup(fields...)).Run()
	if err != nil {
		return "", "", mapCancel(err)
	}
	return finalizeRequired(name, nameDefault), finalizeRequired(email, emailDefault), nil
}
