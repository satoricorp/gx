package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/x/term"
	"github.com/satoricorp/gx/internal/termstyle"
	"github.com/spf13/cobra"
)

const demoWorkspace = "/tmp/gx-demo"

func newDemoCommand() *cobra.Command {
	return &cobra.Command{
		Use:    "demo",
		Short:  "Walk through how GX works",
		Hidden: true,
		Long: strings.Join([]string{
			"Walk through how GX works.",
			"",
			"The demo is interactive and self-contained. It explains the core GX workflow,",
			"asks you to type a few commands, then shows the outputs you would see while",
			"working with stacks, revisions, compose proposals, publishing, and MCP setup.",
		}, "\n"),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDemo(cmd.InOrStdin(), cmd.OutOrStdout())
		},
	}
}

type demoSession struct {
	in         *bufio.Reader
	out        io.Writer
	revisionID string
	reviewURL  string
	mcpDir     string
	fullScreen bool
}

func runDemo(in io.Reader, out io.Writer) error {
	session := demoSession{
		in:         bufio.NewReader(in),
		out:        out,
		revisionID: "demo-r1",
		reviewURL:  "https://gx.run/reviews/demo-stack",
		mcpDir:     demoMCPDir(),
		fullScreen: demoUseFullScreen(out),
	}
	return session.run()
}

func (d *demoSession) run() error {
	if err := setupDemoWorkspace(); err != nil {
		return err
	}
	d.enterFullScreen()
	defer d.exitFullScreen()
	if err := d.welcome(); err != nil {
		return err
	}
	if err := d.addPrompt(); err != nil {
		return err
	}
	if err := d.addOutputAndStacksPrompt(); err != nil {
		return err
	}
	if err := d.stacksOutput(); err != nil {
		return err
	}
	if err := d.composePrompt(); err != nil {
		return err
	}
	if err := d.composeOutputAndAcceptPrompt(); err != nil {
		return err
	}
	if err := d.publishPrompt(); err != nil {
		return err
	}
	if err := d.reviewURLScreen(); err != nil {
		return err
	}
	if err := d.mcpScreen(); err != nil {
		return err
	}
	return d.finish()
}

func (d *demoSession) welcome() error {
	d.beginScreen()
	d.println(demoMint("Welcome to gx!"))
	d.println(demoWhite("gx creates stacks and revisions autonomously."))
	d.println(demoWhite("If you're familiar with git, you can think of stacks like branches, and revisions like commits."))
	d.println(demoWhite("The difference is that gx can easily edit, combine and mutate stacks and revisions."))
	d.println("")
	return d.wait()
}

func (d *demoSession) addPrompt() error {
	d.beginScreen()
	d.println(demoWhite("gx is meant to be used through our MCP server, however, you can use it through the CLI."))
	d.println(demoWhite("We will setup the MCP at the end, but first, let's get comfortable with some commands."))
	d.println(demoWhite(fmt.Sprintf("We've setup a code change for you inside %s for this demo. To save it, type:", demoWorkspace)))
	d.println("")
	d.println(demoMint(`git add hello.txt && gx commit -m "initial change for demo"`))
	d.println("")
	_, err := d.promptCommand(`git add hello.txt && gx commit -m "initial change for demo"`, []string{`git add hello.txt && gx commit -m "initial change for demo"`})
	return err
}

func (d *demoSession) addOutputAndStacksPrompt() error {
	d.beginScreen()
	d.println(demoCommandOutput(`$ git add hello.txt && gx commit -m "initial change for demo"`))
	d.println(demoCommandOutput(""))
	d.println(demoCommandOutput("Revision recorded"))
	d.println(demoCommandOutput("Message       initial change for demo"))
	d.println(demoCommandOutput("Revision      " + d.revisionID))
	d.println(demoCommandOutput("Stack         demo-stack"))
	d.println(demoCommandOutput("Files         hello.txt"))
	d.println("")
	d.println(demoMint("You saved our initial change! ") + demoWhite("If you ever want to edit this revision, you can type ") + demoWhite("`gx edit <revision-id>`") + demoMint("."))
	d.println(demoWhite("To see our change, type:"))
	d.println("")
	d.println(demoMint("gx stacks"))
	d.println("")
	_, err := d.promptCommand("gx stacks", []string{"gx stacks"})
	return err
}

func (d *demoSession) stacksOutput() error {
	d.beginScreen()
	d.println(demoCommandOutput("$ gx stacks"))
	d.println(demoCommandOutput("● demo-stack"))
	d.println(demoCommandOutput("    demo-stack  main · draft · ↑1"))
	d.println(demoCommandOutput("    * demo-r1  initial change for demo"))
	d.println(demoCommandOutput("        hello.txt"))
	d.println(demoCommandOutput(""))
	d.println(demoCommandOutput("j/k up/down · d diff · esc stacks · q quit · ● selected · ↑ cloud · ↓ local"))
	d.println("")
	d.println(demoWhite("This command is interactive, and you can view diffs and edit revisions here."))
	d.println("")
	return d.wait()
}

func (d *demoSession) composePrompt() error {
	d.beginScreen()
	d.println(demoWhite("The ") + demoWhite("`gx commit`") + demoMint(" is very similar to other commands you're already used to, like ") + demoWhite("`git commit`") + demoMint("."))
	d.println(demoWhite("`gx commit`") + demoMint(" should be used sparingly, since we want to use gx to save ourselves time."))
	d.println("")
	d.println(demoWhite("Let's automate this process."))
	d.println(demoWhite("We have new changes in our /tmp directory, but instead of running ") + demoWhite("`gx commit`") + demoWhite(", let's try a new command. Type:"))
	d.println("")
	d.println(demoMint("gx generate"))
	d.println("")
	_, err := d.promptCommand("gx generate", []string{"gx generate"})
	return err
}

func (d *demoSession) composeOutputAndAcceptPrompt() error {
	d.beginScreen()
	d.println(demoCommandOutput("$ gx generate"))
	d.println(demoCommandOutput("Generated demo-stack"))
	d.println(demoCommandOutput("● demo-stack"))
	d.println(demoCommandOutput("    r2  add friendly greeting"))
	d.println(demoCommandOutput("        hello.txt"))
	d.println(demoCommandOutput("    r3  document demo workspace"))
	d.println(demoCommandOutput("        README.md"))
	d.println(demoCommandOutput(""))
	d.println(demoCommandOutput("Next: gx status, then gx push"))
	d.println("")
	d.println(demoWhite("`gx generate`") + demoMint(" reviews all of your code changes and creates local stacks and revisions for you."))
	d.println(demoWhite("No more manual commits for this demo."))
	d.println("")
	d.println(demoWhite("You may notice that the revisions are small. This is to help you during the code review process."))
	d.println(demoWhite("Smaller code changes make it easier to understand what you're merging and easier to suggest changes."))
	d.println("")
	d.println(demoWhite("You can inspect generated work with ") + demoWhite("`gx status`") + demoMint(" before pushing."))
	d.println("")
	d.println(demoMint("Type:"))
	d.println("")
	d.println(demoMint("gx status"))
	d.println("")
	_, err := d.promptCommand("gx status", []string{"gx status"})
	return err
}

func (d *demoSession) publishPrompt() error {
	d.beginScreen()
	d.println(demoMint("That was much easier! The benefit of gx is organizing your work to help yourself and your team review code faster."))
	d.println("")
	d.println(demoMint("Our changes are still local (and now viewable in ") + demoWhite("`gx status`") + demoMint("). They aren't on the server yet, so let's push our changes. Type:"))
	d.println("")
	d.println(demoMint("gx push"))
	d.println("")
	if _, err := d.promptCommand("gx push", []string{"gx push"}); err != nil {
		return err
	}
	d.println("")
	d.println(demoWhite("This is just like `git push`, but it sends code, sessions, and GX metadata for review."))
	d.println(demoWhite("When GitHub is configured, gx also creates or updates the matching PR."))
	d.println("")
	return d.wait()
}

func (d *demoSession) reviewURLScreen() error {
	d.beginScreen()
	d.println(demoCommandOutput("$ gx push"))
	d.println(demoCommandOutput("Pushing 1 stack with 3 revisions"))
	d.println(demoCommandOutput("Synced       gx session context"))
	d.println(demoCommandOutput("Review       " + d.reviewURL))
	d.println("")
	d.println(demoMint("You can go review your new change request at the above URL."))
	d.println(demoMint("Open that link, and then come back for one last tip."))
	d.println("")
	return d.wait()
}

func (d *demoSession) mcpScreen() error {
	d.beginScreen()
	d.println(demoWhite("Command line tools are very helpful for both humans and agents,"))
	d.println(demoWhite("but to take full advantage of gx, you'll want to install the MCP server."))
	d.println("")
	d.println(demoWhite("The MCP server exposes gx_sync, gx_generate, gx_status, gx_push,"))
	d.println(demoWhite("gx_review, and gx_set_base so your coding agent can save with GX first."))
	d.println("")
	for _, item := range d.mcpInstructions() {
		d.println(demoWhite("> " + item.Name))
		for _, line := range item.Lines {
			d.println("  " + demoMint(line))
		}
	}
	d.println("")
	d.println(demoWhite(`Ask your agent to "save work", "save using gx", or "save with gx".`))
	d.println(demoWhite("It should run gx_generate, inspect gx_status, and gx_push ready stacks."))
	d.println("")
	return d.wait()
}

func (d *demoSession) finish() error {
	d.beginScreen()
	d.println(demoMint("Congrats for finishing the gx demo!"))
	d.println("")
	d.println(demoWhite("You should have the basic commands to begin using gx through the CLI and MCP."))
	d.println(demoWhite("If you have questions, please say hi: hi@satori.sh"))
	d.println("")
	d.println(demoWhite("Thank you for using gx, and we look forward to building great things together."))
	d.println(demoWhite("- Joe (Founder)"))
	return nil
}

func (d *demoSession) wait() error {
	fmt.Fprintf(d.out, "%s ", demoGold("[Enter] to continue"))
	raw, err := d.in.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	if isDemoCancel(raw) {
		return fmt.Errorf("demo canceled")
	}
	d.println("")
	return nil
}

func (d *demoSession) promptCommand(defaultValue string, accepted []string) (string, error) {
	fmt.Fprintf(d.out, "%s ", demoWhite(">"))
	raw, err := d.in.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	if isDemoCancel(raw) {
		return "", fmt.Errorf("demo canceled")
	}
	value := strings.TrimSpace(raw)
	if value == "" {
		value = defaultValue
	}
	if len(accepted) > 0 && !demoAcceptedCommand(value, accepted) {
		fmt.Fprintf(d.out, "%s\n", demoGold(fmt.Sprintf("Using `%s` for the walkthrough.", defaultValue)))
		value = defaultValue
	}
	return value, nil
}

func (d *demoSession) println(value string) {
	fmt.Fprintln(d.out, value)
}

func (d *demoSession) enterFullScreen() {
	if !d.fullScreen {
		return
	}
	fmt.Fprint(d.out, "\x1b[?1049h\x1b[H")
}

func (d *demoSession) exitFullScreen() {
	if !d.fullScreen {
		return
	}
	fmt.Fprint(d.out, "\x1b[?1049l")
}

func (d *demoSession) beginScreen() {
	if d.fullScreen {
		fmt.Fprint(d.out, "\x1b[2J\x1b[H")
	}
	d.println("")
}

type demoMCPInstruction struct {
	Name  string
	Lines []string
}

func (d *demoSession) mcpInstructions() []demoMCPInstruction {
	start := fmt.Sprintf("npm --prefix %s run start", d.mcpDir)
	return []demoMCPInstruction{
		{Name: "Cursor", Lines: []string{fmt.Sprintf("cursor mcp add gx -- %s", start)}},
		{Name: "Codex", Lines: []string{
			`Add to ~/.codex/config.toml:`,
			`[mcp_servers.gx]`,
			`command = "npm"`,
			fmt.Sprintf(`args = ["--prefix", "%s", "run", "start"]`, d.mcpDir),
		}},
		{Name: "Claude Code", Lines: []string{fmt.Sprintf("claude mcp add gx -- %s", start)}},
		{Name: "Opencode", Lines: []string{fmt.Sprintf("opencode mcp add gx -- %s", start)}},
		{Name: "Antigravity", Lines: []string{
			`Add a stdio MCP server named "gx":`,
			`command: npm`,
			fmt.Sprintf(`args: --prefix %s run start`, d.mcpDir),
		}},
	}
}

func setupDemoWorkspace() error {
	if err := os.MkdirAll(demoWorkspace, 0o755); err != nil {
		return fmt.Errorf("create demo workspace: %w", err)
	}
	files := map[string]string{
		"hello.txt": "hello gx\n",
		"README.md": strings.Join([]string{
			"# GX demo",
			"",
			"This temporary workspace is used by `gx demo`.",
			"",
		}, "\n"),
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(demoWorkspace, name), []byte(content), 0o644); err != nil {
			return fmt.Errorf("write demo file %s: %w", name, err)
		}
	}
	return nil
}

func demoMCPDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "mcp"
	}
	return filepath.Join(cwd, "mcp")
}

func demoAcceptedCommand(value string, accepted []string) bool {
	normalized := normalizeDemoCommand(value)
	for _, candidate := range accepted {
		if normalized == normalizeDemoCommand(candidate) {
			return true
		}
	}
	return false
}

func normalizeDemoCommand(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, `"`, `'`)
	value = strings.Join(strings.Fields(value), " ")
	return strings.ReplaceAll(value, "shift + a", "shift+a")
}

func isDemoCancel(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "q", "quit", "exit":
		return true
	default:
		return false
	}
}

func demoMint(text string) string {
	return accent(text)
}

func demoGold(text string) string {
	return demoANSI("38;2;245;194;107", text)
}

func demoWhite(text string) string {
	return demoANSI("38;2;250;250;250", text)
}

func demoItalic(text string) string {
	if text == "" || !termstyle.Enabled() {
		return text
	}
	return "\x1b[3m" + text + "\x1b[0m"
}

func demoCommandOutput(text string) string {
	return demoANSI("38;2;212;212;216", text)
}

func demoANSI(code, text string) string {
	if text == "" || !termstyle.Enabled() {
		return text
	}
	return "\x1b[" + code + "m" + text + "\x1b[0m"
}

func demoUseFullScreen(out io.Writer) bool {
	file, ok := out.(*os.File)
	if !ok {
		return false
	}
	if strings.EqualFold(os.Getenv("TERM"), "dumb") {
		return false
	}
	return term.IsTerminal(file.Fd())
}
