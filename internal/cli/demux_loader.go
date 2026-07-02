package cli

import (
	"io"
	"os"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"

	"github.com/satoricorp/gx/internal/authoring"
)

type demuxLoaderRunFunc func(io.Writer) (authoring.DemuxPlanPacket, error)
type demuxApplyLoaderRunFunc func() (authoring.ApplyDemuxResult, error)

var demuxBrailleColors = []string{
	"#22D3EE",
	"#38BDF8",
	"#818CF8",
	"#C084FC",
	"#F472B6",
	"#FB7185",
	"#FBBF24",
	"#A3E635",
}

var demuxBrailleSpinner = spinner.Spinner{
	Frames: demuxBrailleFrames(),
	FPS:    70 * time.Millisecond,
}

var demuxLoaderStyle = lipgloss.NewStyle().
	MarginTop(1).
	MarginLeft(2).
	MarginBottom(1)

var demuxDidYouKnowHeaderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#A78BFA")).
	Bold(true)

type demuxLoaderModel struct {
	spinner spinner.Model
	phase   string
	phases  chan string
	run     demuxLoaderRunFunc
	packet  authoring.DemuxPlanPacket
	err     error
	done    bool
}

type demuxApplyLoaderModel struct {
	spinner spinner.Model
	phase   string
	run     demuxApplyLoaderRunFunc
	result  authoring.ApplyDemuxResult
	err     error
	done    bool
}

type demuxLoaderPhaseMsg string

type demuxLoaderNoopMsg struct{}

type demuxLoaderResultMsg struct {
	packet authoring.DemuxPlanPacket
	err    error
}

type demuxApplyLoaderResultMsg struct {
	result authoring.ApplyDemuxResult
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
		spinner: newDemuxBrailleSpinner(),
		phase:   "Grouping changes...",
		phases:  phases,
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

func runDemuxApplyWithLoader(in io.Reader, out io.Writer, phase string, run demuxApplyLoaderRunFunc) (authoring.ApplyDemuxResult, error) {
	model := demuxApplyLoaderModel{
		spinner: newDemuxBrailleSpinner(),
		phase:   phase,
		run:     run,
	}
	program := tea.NewProgram(model, tea.WithInput(in), tea.WithOutput(out))
	finalModel, err := program.Run()
	if err != nil {
		return authoring.ApplyDemuxResult{}, err
	}
	if result, ok := finalModel.(demuxApplyLoaderModel); ok {
		return result.result, result.err
	}
	return authoring.ApplyDemuxResult{}, nil
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
	return tea.NewView(renderDemuxLoaderView(m.spinner.View(), m.phase))
}

func (m demuxApplyLoaderModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, runDemuxApplyLoader(m.run))
}

func (m demuxApplyLoaderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case demuxApplyLoaderResultMsg:
		m.done = true
		m.phase = ""
		m.result = msg.result
		m.err = msg.err
		return m, tea.Quit
	}
	return m, nil
}

func (m demuxApplyLoaderModel) View() tea.View {
	if m.done {
		return tea.NewView("")
	}
	return tea.NewView(renderDemuxLoaderView(m.spinner.View(), m.phase))
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

func runDemuxApplyLoader(run demuxApplyLoaderRunFunc) tea.Cmd {
	return func() tea.Msg {
		result, err := run()
		return demuxApplyLoaderResultMsg{result: result, err: err}
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

func newDemuxBrailleSpinner() spinner.Model {
	return spinner.New(
		spinner.WithSpinner(demuxBrailleSpinner),
	)
}

func demuxBrailleFrames() []string {
	patterns := []string{"⡿", "⣟", "⣯", "⣷", "⣾", "⣽", "⣻", "⢿"}
	frames := make([]string, 0, len(patterns))
	for index, pattern := range patterns {
		color := demuxBrailleColors[index%len(demuxBrailleColors)]
		frames = append(frames, lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(pattern))
	}
	return frames
}

func renderDemuxLoaderView(frame, phase string) string {
	body := strings.Join([]string{
		frame + " " + muted(phase),
		"",
		demuxDidYouKnowHeaderStyle.Render("Did you know?"),
		muted(demuxLoaderFact()),
	}, "\n")
	return demuxLoaderStyle.Render(body)
}

func demuxLoaderFact() string {
	return "GX stores revision and stack metadata with the code changes it creates."
}
