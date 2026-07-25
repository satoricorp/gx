package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"

	"github.com/satoricorp/gx/internal/cloud"
)

// demoSkipExec is set by tests to exercise copy without running real commands.
var demoSkipExec bool

func newDemoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "demo",
		Short: "Interactive first-run walkthrough",
		Long: strings.Join([]string{
			"Walk through the core gx loop in a scratch repository.",
			"",
			"Your real working copy is never touched. Press Enter to advance",
			"each step; type q to quit.",
		}, "\n"),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDemo(cmd.Context(), cmd.InOrStdin(), cmd.OutOrStdout())
		},
	}
}

type demoSession struct {
	ctx     context.Context
	in      *bufio.Reader
	out     io.Writer
	scratch string
	gxBin   string
}

func runDemo(ctx context.Context, in io.Reader, out io.Writer) error {
	scratch, err := os.MkdirTemp("", "gx-demo-*")
	if err != nil {
		return fmt.Errorf("create demo scratch repo: %w", err)
	}
	defer os.RemoveAll(scratch)

	session := demoSession{
		ctx:     ctx,
		in:      bufio.NewReader(in),
		out:     out,
		scratch: scratch,
		gxBin:   demoGxBinary(),
	}
	return session.run()
}

func (d *demoSession) run() error {
	if err := d.seedScratch(); err != nil {
		return err
	}
	d.printWelcome()
	if err := d.waitContinue(); err != nil {
		return err
	}
	if err := d.stepInit(); err != nil {
		return err
	}
	if err := d.stepAuth(); err != nil {
		return err
	}
	if err := d.stepStage(); err != nil {
		return err
	}
	if err := d.stepCommit(); err != nil {
		return err
	}
	if err := d.stepPush(); err != nil {
		return err
	}
	d.printDone()
	return nil
}

func (d *demoSession) seedScratch() error {
	if demoSkipExec {
		return nil
	}
	if err := runInDir(d.scratch, "git", "init", "-q"); err != nil {
		return fmt.Errorf("git init scratch repo: %w", err)
	}
	if err := runInDir(d.scratch, "git", "config", "user.name", "gx demo"); err != nil {
		return err
	}
	if err := runInDir(d.scratch, "git", "config", "user.email", "demo@gx.local"); err != nil {
		return err
	}
	readme := filepath.Join(d.scratch, "README.md")
	if err := os.WriteFile(readme, []byte("# gx demo\n"), 0o644); err != nil {
		return fmt.Errorf("seed demo readme: %w", err)
	}
	if err := runInDir(d.scratch, "git", "add", "README.md"); err != nil {
		return err
	}
	if err := runInDir(d.scratch, "git", "commit", "-q", "-m", "initial commit"); err != nil {
		return err
	}
	path := filepath.Join(d.scratch, "rate_limiter.go")
	content := "package main\n\nfunc rateLimit() {}\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("seed demo file: %w", err)
	}
	return nil
}

func (d *demoSession) printWelcome() {
	d.println("")
	d.println("  " + accent("gx") + " — version control that captures how code was made")
	d.println("")
	d.println("  " + muted("5 steps · 2 minutes · runs in a scratch repo, your work is untouched"))
	d.println("")
}

func (d *demoSession) stepInit() error {
	d.printStepHeader("1/5", "init")
	d.println("  Set up gx in a repository.")
	d.println("")
	d.println("    $ " + command("gx init"))
	d.println("")
	d.println("  Installs hooks and starts ambient capture. From here on,")
	d.println("  every AI session that touches this repo is recorded —")
	d.println("  models, tokens, and context ride along with your commits.")
	d.println("")
	d.println("  " + muted("↵ run it"))
	d.println("")
	if err := d.waitContinue(); err != nil {
		return err
	}
	// Accept identity/setup defaults via a dedicated stdin so the walkthrough
	// stdin is not consumed by init prompts. Screen copy still shows `gx init`.
	return d.runAndShowDisplay("gx init", d.gxBin, "init")
}

func (d *demoSession) stepAuth() error {
	d.printStepHeader("2/5", "sign in")
	d.println("  Connect to gx cloud.")
	d.println("")
	d.println("    $ " + command("gx auth login"))
	d.println("")
	d.println("  GitHub device login. Captures stay on your machine")
	d.println("  until you choose to share them.")
	d.println("")

	if demoAlreadySignedIn() {
		d.println("  " + success("already signed in ✓"))
		d.println("")
		d.println("  " + muted("↵ continue"))
		d.println("")
		return d.waitContinue()
	}

	d.println("  " + muted("↵ run it"))
	d.println("")
	if err := d.waitContinue(); err != nil {
		return err
	}
	if demoSkipExec || !demoIsTTY(d.out) {
		d.println("  " + muted("Skipped in this session — run gx auth login when you are ready."))
		d.println("")
		return nil
	}
	if !d.confirm("Run gx auth login now?") {
		d.println("  " + muted("Skipped — run gx auth login when you are ready."))
		d.println("")
		return nil
	}
	return d.runAndShowDisplay("gx auth login", d.gxBin, "auth", "login")
}

func (d *demoSession) stepStage() error {
	d.printStepHeader("3/5", "stage")
	d.println("  Choose what goes in — plain git.")
	d.println("")
	d.println("    $ " + command("git add -p"))
	d.println("")
	d.println("  gx doesn't replace staging. Pick files and hunks")
	d.println("  exactly like you always have.")
	d.println("")
	d.println("  " + muted("↵ run it"))
	d.println("")
	if err := d.waitContinue(); err != nil {
		return err
	}
	// Non-interactive stand-in for git add -p: stage the seeded file.
	if err := d.runAndShow("git", "add", "rate_limiter.go"); err != nil {
		return err
	}
	d.println("  " + muted("(demo staged the seeded file; use git add -p in your own repo)"))
	d.println("")
	return nil
}

func (d *demoSession) stepCommit() error {
	d.printStepHeader("4/5", "commit")
	d.println("  Commit with git — nothing new to learn.")
	d.println("")
	d.println("    $ " + command(`git commit -m "add rate limiter"`))
	d.println("")
	d.println("  The hooks gx installed record the revision as you commit:")
	d.println("  the sessions, models, and tokens behind the change —")
	d.println("  provenance your reviewer can actually use.")
	d.println("")
	d.println("  " + muted("↵ run it"))
	d.println("")
	if err := d.waitContinue(); err != nil {
		return err
	}
	return d.runAndShowDisplay(`git commit -m "add rate limiter"`, "git", "commit", "-m", "add rate limiter")
}

func (d *demoSession) stepPush() error {
	d.printStepHeader("5/5", "push")
	d.println("  Ship it.")
	d.println("")
	d.println("    $ " + command("git push"))
	d.println("")
	d.println("  Your branch goes up like always. gx attaches the context,")
	d.println("  so reviewers see the why — not just the diff.")
	d.println("")
	d.println("  " + muted("↵ run it"))
	d.println("")
	if err := d.waitContinue(); err != nil {
		return err
	}
	d.println("  $ git push")
	d.println("  " + muted("To github.com:you/your-repo.git"))
	d.println("  " + muted(" * [new branch]      feature/rate-limiter -> feature/rate-limiter"))
	d.println("  " + success("gx attached session context for review"))
	d.println("")
	d.println("  " + muted("(simulated — this scratch repo has no remote)"))
	d.println("")
	return nil
}

func (d *demoSession) printDone() {
	d.println("  " + section("── done ───────────────────────────────────────────"))
	d.println("")
	d.println("  That's the loop: stage and commit with git, gx records the rest.")
	d.println("")
	d.println("    " + command("git status") + "   your working tree at a glance")
	d.println("    " + command("gx review") + "    context-aware code review")
	d.println("    " + command("gx doctor") + "    check and fix your setup")
	d.println("")
}

func (d *demoSession) printStepHeader(step, title string) {
	d.println("")
	d.println("  " + section(fmt.Sprintf("── %s · %s ─────────────────────────────────────────", step, title)))
	d.println("")
}

func (d *demoSession) waitContinue() error {
	raw, err := d.in.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	if isDemoCancel(raw) {
		return fmt.Errorf("demo canceled")
	}
	return nil
}

func (d *demoSession) confirm(prompt string) bool {
	fmt.Fprintf(d.out, "  %s [y/N] ", muted(prompt))
	raw, err := d.in.ReadString('\n')
	if err != nil {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "y", "yes":
		return true
	default:
		return false
	}
}

func (d *demoSession) runAndShow(name string, args ...string) error {
	display := strings.Join(append([]string{filepath.Base(name)}, args...), " ")
	return d.runAndShowDisplay(display, name, args...)
}

func (d *demoSession) runAndShowDisplay(display, name string, args ...string) error {
	d.println("  $ " + display)
	if demoSkipExec {
		d.println("  " + muted("(skipped)"))
		d.println("")
		return nil
	}
	cmd := exec.CommandContext(d.ctx, name, args...)
	cmd.Dir = d.scratch
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	// Keep demo stdin for Enter/q; feed defaults to child prompts.
	if name == d.gxBin && len(args) > 0 && args[0] == "init" {
		cmd.Stdin = strings.NewReader("\n\n\n")
	}
	output, err := cmd.CombinedOutput()
	text := strings.TrimRight(string(output), "\n")
	if text != "" {
		for _, line := range strings.Split(text, "\n") {
			d.println("  " + line)
		}
	}
	d.println("")
	if err != nil {
		return fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}
	return nil
}

func (d *demoSession) println(value string) {
	fmt.Fprintln(d.out, value)
}

func demoAlreadySignedIn() bool {
	creds, err := cloud.LoadCloudCredentials()
	if err != nil || creds == nil {
		return false
	}
	return strings.TrimSpace(creds.Login) != "" || strings.TrimSpace(creds.UserID) != ""
}

func demoGxBinary() string {
	if bin := strings.TrimSpace(os.Getenv("GX_DEMO_BIN")); bin != "" {
		return bin
	}
	if exe, err := os.Executable(); err == nil {
		base := filepath.Base(exe)
		if base == "gx" || strings.HasPrefix(base, "gx-") {
			return exe
		}
	}
	if path, err := exec.LookPath("gx"); err == nil {
		return path
	}
	return "gx"
}

func demoIsTTY(out io.Writer) bool {
	file, ok := out.(*os.File)
	if !ok {
		return false
	}
	if strings.EqualFold(os.Getenv("TERM"), "dumb") {
		return false
	}
	return term.IsTerminal(file.Fd())
}

func runInDir(dir, name string, args ...string) error {
	if demoSkipExec && name != "git" {
		return nil
	}
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s %s: %w\n%s", name, strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return nil
}

func isDemoCancel(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "q", "quit", "exit":
		return true
	default:
		return false
	}
}
