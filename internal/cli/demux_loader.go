package cli

import (
	"io"
	"os"
	"strings"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"

	"github.com/satoricorp/gx/internal/authoring"
)

type demuxLoaderRunFunc func(io.Writer) (authoring.DemuxPlanPacket, error)

type demuxLoaderModel struct {
	spinner spinner.Model
	phase   string
	phases  chan string
	run     demuxLoaderRunFunc
	packet  authoring.DemuxPlanPacket
	err     error
	done    bool
}

type demuxLoaderPhaseMsg string

type demuxLoaderNoopMsg struct{}

type demuxLoaderResultMsg struct {
	packet authoring.DemuxPlanPacket
	err    error
}

func useDemuxLoader(in io.Reader, out io.Writer) bool {
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

func runDemuxWithLoader(in io.Reader, out io.Writer, run demuxLoaderRunFunc) (authoring.DemuxPlanPacket, error) {
	phases := make(chan string, 4)
	model := demuxLoaderModel{
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Line),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#6366F1"))),
		),
		phase:  "Starting compose...",
		phases: phases,
		run: func(progress io.Writer) (authoring.DemuxPlanPacket, error) {
			return run(progress)
		},
	}
	program := tea.NewProgram(model, tea.WithInput(in), tea.WithOutput(out))
	finalModel, err := program.Run()
	if err != nil {
		return authoring.DemuxPlanPacket{}, err
	}
	if result, ok := finalModel.(demuxLoaderModel); ok {
		return result.packet, result.err
	}
	return authoring.DemuxPlanPacket{}, nil
}

func (m demuxLoaderModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, waitDemuxLoaderPhase(m.phases), runDemuxLoader(m.run, m.phases))
}

func (m demuxLoaderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case demuxLoaderPhaseMsg:
		m.phase = string(msg)
		return m, waitDemuxLoaderPhase(m.phases)
	case demuxLoaderResultMsg:
		m.done = true
		m.phase = ""
		m.packet = msg.packet
		m.err = msg.err
		return m, tea.Quit
	case demuxLoaderNoopMsg:
		return m, nil
	}
	return m, nil
}

func (m demuxLoaderModel) View() tea.View {
	if m.done {
		return tea.NewView("")
	}
	return tea.NewView(m.spinner.View() + " " + muted(m.phase))
}

func waitDemuxLoaderPhase(phases <-chan string) tea.Cmd {
	return func() tea.Msg {
		phase, ok := <-phases
		if !ok {
			return demuxLoaderNoopMsg{}
		}
		return demuxLoaderPhaseMsg(phase)
	}
}

func runDemuxLoader(run demuxLoaderRunFunc, phases chan<- string) tea.Cmd {
	return func() tea.Msg {
		packet, err := run(demuxLoaderProgressWriter{phases: phases})
		close(phases)
		return demuxLoaderResultMsg{packet: packet, err: err}
	}
}

type demuxLoaderProgressWriter struct {
	phases chan<- string
}

func (w demuxLoaderProgressWriter) Write(p []byte) (int, error) {
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
