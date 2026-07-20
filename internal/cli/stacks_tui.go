package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/term"

	"github.com/satoricorp/gx/internal/authoring"
)

type stacksMode int

const (
	stacksModeStacks stacksMode = iota
	stacksModeRevisions
)

type stacksAction struct {
	Kind   string
	Target string
}

type stacksModel struct {
	stack       authoring.StackSummary
	unrecorded  *authoring.ChangeInfo
	hiddenEmpty int
	mode        stacksMode
	height      int
	stackCursor int
	revCursor   int
	action      stacksAction
}

func useStatusInteractive(in io.Reader, out io.Writer) bool {
	if os.Getenv("GX_PLAIN_PROMPTS") != "" {
		return false
	}
	input, ok := in.(*os.File)
	if !ok || !term.IsTerminal(input.Fd()) {
		return false
	}
	output, ok := out.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(output.Fd())
}

func runStacksInteractive(in io.Reader, out io.Writer, stack authoring.StackSummary, unrecorded *authoring.ChangeInfo, hiddenEmpty int) (stacksAction, error) {
	model := newStacksModelWithHidden(stack, unrecorded, hiddenEmpty)
	program := tea.NewProgram(model, tea.WithInput(in), tea.WithOutput(out))
	final, err := program.Run()
	if err != nil {
		return stacksAction{}, err
	}
	if model, ok := final.(stacksModel); ok {
		return model.action, nil
	}
	return stacksAction{}, nil
}

func newStacksModel(stack authoring.StackSummary, unrecorded *authoring.ChangeInfo) stacksModel {
	return newStacksModelWithHidden(stack, unrecorded, 0)
}

func newStacksModelWithHidden(stack authoring.StackSummary, unrecorded *authoring.ChangeInfo, hiddenEmpty int) stacksModel {
	stack = stackSummaryWithDisplayFallback(stack)
	model := stacksModel{stack: stack, unrecorded: unrecorded, hiddenEmpty: hiddenEmpty, stackCursor: currentStackIndex(stack), mode: stacksModeStacks}
	if model.stackCursor < 0 {
		model.stackCursor = 0
	}
	model.revCursor = latestRevisionDisplayIndex(stack.Revisions, unrecorded)
	return model
}

func (m stacksModel) Init() tea.Cmd {
	return nil
}

func (m stacksModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			m.move(-1)
		case "down", "j":
			m.move(1)
		case "enter":
			if m.mode == stacksModeStacks {
				if len(m.selectedStackRevisions()) > 0 {
					m.mode = stacksModeRevisions
					m.revCursor = latestRevisionDisplayIndex(m.selectedStackRevisions(), m.selectedStackUnrecorded())
					return m, nil
				}
			}
		case "esc":
			if m.mode == stacksModeRevisions {
				m.mode = stacksModeStacks
				return m, nil
			}
		case "ctrl+c", "q":
			m.action = stacksAction{Kind: "quit"}
			return m, tea.Quit
		}
	}
	m.clamp()
	return m, nil
}

func (m *stacksModel) move(delta int) {
	switch m.mode {
	case stacksModeRevisions:
		m.revCursor += delta
	default:
		m.stackCursor += delta
	}
	m.clamp()
}

func (m *stacksModel) clamp() {
	stacks := orderedStacks(m.stack)
	if len(stacks) == 0 {
		m.stackCursor = 0
	} else {
		if m.stackCursor < 0 {
			m.stackCursor = 0
		}
		if m.stackCursor >= len(stacks) {
			m.stackCursor = len(stacks) - 1
		}
	}
	revisions := m.visibleRevisions()
	if len(revisions) == 0 {
		m.revCursor = 0
		return
	}
	if m.revCursor < 0 {
		m.revCursor = 0
	}
	if m.revCursor >= len(revisions) {
		m.revCursor = len(revisions) - 1
	}
}

func (m stacksModel) selectedStack() *authoring.StackInfo {
	stacks := orderedStacks(m.stack)
	if m.stackCursor < 0 || m.stackCursor >= len(stacks) {
		return nil
	}
	return &stacks[m.stackCursor]
}

func (m stacksModel) selectedStackSelector() string {
	stack := m.selectedStack()
	if stack == nil {
		return ""
	}
	return firstNonEmptyString(stack.Alias, stack.BookmarkName, stack.Name)
}

func (m stacksModel) selectedStackBookmark() string {
	stack := m.selectedStack()
	if stack == nil {
		return ""
	}
	return stack.BookmarkName
}

func (m stacksModel) selectedStackIsCurrent() bool {
	stack := m.selectedStack()
	return stack != nil && m.stack.Stack != nil && stack.BookmarkName == m.stack.Stack.BookmarkName
}

func (m stacksModel) selectedStackRevisions() []authoring.RevisionSummary {
	stack := m.selectedStack()
	if stack == nil {
		return nil
	}
	return stackEntryRevisions(m.stack, *stack)
}

func (m stacksModel) selectedStackUnrecorded() *authoring.ChangeInfo {
	if !m.selectedStackIsCurrent() {
		return nil
	}
	return m.unrecorded
}

func (m stacksModel) visibleRevisions() []authoring.RevisionSummary {
	switch m.mode {
	case stacksModeRevisions:
		return m.selectedStackRevisions()
	default:
		return m.stack.Revisions
	}
}

func (m stacksModel) selectedRevision() string {
	revisions := m.visibleRevisions()
	if m.revCursor < 0 || m.revCursor >= len(revisions) {
		return ""
	}
	return firstNonEmptyString(revisions[m.revCursor].CommitID, revisions[m.revCursor].ChangeID)
}

func (m stacksModel) View() tea.View {
	var buf bytes.Buffer
	switch m.mode {
	case stacksModeRevisions:
		stack := m.selectedStack()
		if stack == nil {
			fmt.Fprint(&buf, renderStacksSummary(m.stack, m.unrecorded, m.stackCursor, false, m.revCursor))
		} else {
			lines := renderStackRevisionLines(m.stack, *stack, m.selectedStackUnrecorded(), m.revCursor)
			for _, line := range lines {
				fmt.Fprintln(&buf, line)
			}
		}
	default:
		fmt.Fprint(&buf, renderStacksSummary(m.stack, m.unrecorded, m.stackCursor, true, m.revCursor))
	}
	fmt.Fprint(&buf, muted("\n↑/↓ navigate · enter open revisions · q quit"))
	return tea.NewView(buf.String())
}
