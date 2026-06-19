package clitui

import (
	"context"
	"fmt"
	"io"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/term"
	"github.com/satoricorp/gx/internal/vcs"
)

func useInteractive(in any) bool {
	if os.Getenv("GX_PLAIN_PROMPTS") != "" {
		return false
	}
	f, ok := in.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(f.Fd())
}

type pickerKind int

const (
	pickerStatus pickerKind = iota
	pickerModify
)

type pickerModel struct {
	kind      pickerKind
	snapshot  vcs.StatusSnapshot
	revisions []vcs.RevisionSnapshot
	cursor    int
	confirmed bool
	canceled  bool
	selection string
}

// PickModifyRevision runs the interactive gx edit picker.
func PickModifyRevision(in io.Reader, snapshot vcs.StatusSnapshot, revisions []vcs.RevisionSnapshot) (string, error) {
	return runPicker(pickerModify, snapshot, revisions, in)
}

// ExploreStatus runs the interactive gx status explorer.
func ExploreStatus(in io.Reader, snapshot vcs.StatusSnapshot) error {
	_, err := runPicker(pickerStatus, snapshot, nil, in)
	return err
}

func runPicker(kind pickerKind, snapshot vcs.StatusSnapshot, revisions []vcs.RevisionSnapshot, in io.Reader) (string, error) {
	if !useInteractive(in) {
		return "", fmt.Errorf("interactive picker requires a terminal")
	}

	model := pickerModel{
		kind:      kind,
		snapshot:  snapshot,
		revisions: revisions,
	}
	program := tea.NewProgram(model)
	final, runErr := program.Run()
	if runErr != nil {
		return "", runErr
	}
	m, ok := final.(pickerModel)
	if !ok {
		return "", fmt.Errorf("unexpected picker result")
	}
	if m.canceled || !m.confirmed {
		return "", context.Canceled
	}
	return m.selection, nil
}

func (m pickerModel) Init() tea.Cmd {
	return nil
}

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < m.maxCursor() {
				m.cursor++
			}
		case "enter":
			if m.kind == pickerStatus {
				m.confirmed = true
				return m, tea.Quit
			}
			m.confirmed = true
			m.selection = m.selectedValue()
			return m, tea.Quit
		case "ctrl+c", "q", "esc":
			m.canceled = true
			return m, tea.Quit
		}
	}
	if m.cursor > m.maxCursor() {
		m.cursor = m.maxCursor()
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	return m, nil
}

func (m pickerModel) maxCursor() int {
	switch m.kind {
	case pickerModify:
		if len(m.revisions) == 0 {
			return 0
		}
		return len(m.revisions) - 1
	default:
		if len(m.snapshot.Bookmarks) == 0 {
			return 0
		}
		return len(m.snapshot.Bookmarks) - 1
	}
}

func (m pickerModel) selectedValue() string {
	switch m.kind {
	case pickerModify:
		if m.cursor >= 0 && m.cursor < len(m.revisions) {
			return m.revisions[m.cursor].ChangeID
		}
	}
	return ""
}

func (m pickerModel) View() tea.View {
	var content string
	switch m.kind {
	case pickerStatus:
		content = RenderStatusInteractive(m.snapshot, m.cursor)
	case pickerModify:
		content = RenderModifyInteractive(m.snapshot, m.revisions, m.cursor)
	}
	return tea.NewView(content)
}
