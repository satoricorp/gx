package cli

import (
	"io"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/satoricorp/gx/internal/publication"
)

type publishUploadRunFunc func() (publication.Result, error)

type publishUploadModel struct {
	spinner spinner.Model
	run     publishUploadRunFunc
	result  publication.Result
	err     error
	done    bool
}

type publishUploadResultMsg struct {
	result publication.Result
	err    error
}

func runPublishUploadWithLoader(in io.Reader, out io.Writer, run publishUploadRunFunc) (publication.Result, error) {
	if !useDemuxLoader(in, out) {
		return run()
	}
	model := publishUploadModel{
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Line),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#6366F1"))),
		),
		run: run,
	}
	program := tea.NewProgram(model, tea.WithInput(in), tea.WithOutput(out))
	finalModel, err := program.Run()
	if err != nil {
		return publication.Result{}, err
	}
	if result, ok := finalModel.(publishUploadModel); ok {
		return result.result, result.err
	}
	return publication.Result{}, nil
}

func (m publishUploadModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, runPublishUpload(m.run))
}

func (m publishUploadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case publishUploadResultMsg:
		m.done = true
		m.result = msg.result
		m.err = msg.err
		return m, tea.Quit
	}
	return m, nil
}

func (m publishUploadModel) View() tea.View {
	if m.done {
		return tea.NewView("")
	}
	return tea.NewView(m.spinner.View() + " " + muted("uploading gx session context"))
}

func runPublishUpload(run publishUploadRunFunc) tea.Cmd {
	return func() tea.Msg {
		result, err := run()
		return publishUploadResultMsg{result: result, err: err}
	}
}
