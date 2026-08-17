package cli

import (
	"io"
	"os"
	"strings"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"

	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/termstyle"
)

// useInteractiveTerminal reports whether both ends of the pipe are a real
// terminal, so a Bubble Tea loader can take over the screen.
func useInteractiveTerminal(in io.Reader, out io.Writer) bool {
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

type reviewLoaderRunFunc func(io.Writer) (codereview.Report, error)

type reviewLoaderModel struct {
	spinner spinner.Model
	phase   string
	phases  chan string
	run     reviewLoaderRunFunc
	report  codereview.Report
	err     error
	done    bool
}

type reviewLoaderPhaseMsg string

type reviewLoaderNoopMsg struct{}

type reviewLoaderResultMsg struct {
	report codereview.Report
	err    error
}

// runReviewWithLoader runs the review under a spinner when both ends are a
// terminal, and plain otherwise. It returns the report; whether the caller
// then prints the linear render or opens the accordion is decided by
// browseReviewInteractively, which needs the same tty answer.
func runReviewWithLoader(in io.Reader, out io.Writer, run reviewLoaderRunFunc) (codereview.Report, error) {
	if !useInteractiveTerminal(in, out) {
		return run(nil)
	}
	phases := make(chan string, 6)
	model := reviewLoaderModel{
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Line),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#6366F1"))),
		),
		phase:  "Starting review",
		phases: phases,
		run:    run,
	}
	program := tea.NewProgram(model, tea.WithInput(in), tea.WithOutput(out))
	finalModel, err := program.Run()
	if err != nil {
		return codereview.Report{}, err
	}
	if result, ok := finalModel.(reviewLoaderModel); ok {
		return result.report, result.err
	}
	return codereview.Report{}, nil
}

func (m reviewLoaderModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, waitReviewLoaderPhase(m.phases), runReviewLoader(m.run, m.phases))
}

func (m reviewLoaderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case reviewLoaderPhaseMsg:
		m.phase = string(msg)
		return m, waitReviewLoaderPhase(m.phases)
	case reviewLoaderResultMsg:
		m.done = true
		m.phase = ""
		m.report = msg.report
		m.err = msg.err
		return m, tea.Quit
	case reviewLoaderNoopMsg:
		return m, nil
	}
	return m, nil
}

func (m reviewLoaderModel) View() tea.View {
	if m.done {
		return tea.NewView("")
	}
	return tea.NewView(m.spinner.View() + " " + muted(m.phase))
}

func waitReviewLoaderPhase(phases <-chan string) tea.Cmd {
	return func() tea.Msg {
		phase, ok := <-phases
		if !ok {
			return reviewLoaderNoopMsg{}
		}
		return reviewLoaderPhaseMsg(phase)
	}
}

func runReviewLoader(run reviewLoaderRunFunc, phases chan<- string) tea.Cmd {
	return func() tea.Msg {
		report, err := run(reviewLoaderProgressWriter{phases: phases})
		close(phases)
		return reviewLoaderResultMsg{report: report, err: err}
	}
}

type reviewLoaderProgressWriter struct {
	phases chan<- string
}

func (w reviewLoaderProgressWriter) Write(p []byte) (int, error) {
	for _, line := range strings.Split(string(p), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		select {
		case w.phases <- line:
		default:
		}
	}
	return len(p), nil
}

// browseReviewInteractively opens the accordion over a finished report when
// both ends are a terminal and the caller has not asked for a machine shape.
// It returns true when it took over the screen — the caller then prints
// nothing else to stdout, because the accordion already showed the report and
// the seam. It returns false when the terminal is not interactive, so the
// caller falls back to RenderReviewText, and the exit code is identical
// either way.
//
// GX_PLAIN_PROMPTS and CI both route here as "not interactive" through
// useInteractiveTerminal, so agents and CI never see a menu.
func browseReviewInteractively(in io.Reader, out io.Writer, report codereview.Report) bool {
	if !useInteractiveTerminal(in, out) {
		return false
	}
	if os.Getenv("CI") != "" {
		return false
	}
	model := newAccordionModel(report, report.Color && termstyle.Enabled())
	program := tea.NewProgram(model, tea.WithInput(in), tea.WithOutput(out))
	if _, err := program.Run(); err != nil {
		// The accordion is presentation. If the TUI cannot run for any
		// reason, fall back to the linear render rather than lose the report.
		return false
	}
	return true
}
