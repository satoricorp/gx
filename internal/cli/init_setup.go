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

	"github.com/spf13/cobra"
)

const agentsMDSnippetMarker = "Version control: use GX, not `git commit`."

const agentsMDSnippet = `# Agents

Version control: use GX, not ` + "`git commit`" + `.

Default flow:
- Run ` + "`git add`" + ` to stage the files for this revision.
- Run ` + "`gx commit`" + ` (or MCP ` + "`gx_commit`" + `) to record staged work as a GX revision.
- Run ` + "`gx status`" + ` to inspect staged files, local features/revisions, and Remote state.
- Run ` + "`gx push`" + ` to push code changes, sessions, and GX metadata to GitHub/origin and GX Cloud. Plain ` + "`git push`" + ` also works when the GX pre-push hook is installed.

Use GX MCP first: ` + "`gx_sync`" + `, ` + "`gx_status`" + `, then ` + "`gx_commit`" + ` / ` + "`gx_push`" + `.
If MCP is unavailable, use the CLI fallback: ` + "`git add`" + `, ` + "`gx commit`" + `, ` + "`gx status`" + `, then ` + "`git push`" + `.

Use hidden utility commands such as ` + "`gx add`" + `, ` + "`gx switch`" + `, ` + "`gx base`" + `, ` + "`gx edit`" + `, or ` + "`gx modify`" + ` only for explicit surgery or user-directed repair.

Use ` + "`gx generate`" + ` only for bulk organization of large working copies.
`

type initSetupOptions struct {
	RepoRoot string
	Quiet    bool
	In       io.Reader
	Out      io.Writer
	Err      io.Writer
}

func runInitRepoSetup(opts initSetupOptions) error {
	gxPath, err := os.Executable()
	if err != nil {
		return err
	}
	mcpPath, err := resolveMCPBinary(gxPath)
	if err != nil {
		if !opts.Quiet && opts.Err != nil {
			fmt.Fprintln(opts.Err, labelWarningValue("MCP", err.Error()))
		}
	} else {
		registered := registerMCPClients(context.Background(), gxPath, mcpPath)
		if !opts.Quiet && opts.Out != nil {
			if len(registered) == 0 {
				fmt.Fprintln(opts.Out, labelValue("MCP", muted("skipped (no supported agent CLI found)")))
			} else {
				fmt.Fprintln(opts.Out, labelValue("MCP", success("ok")+": "+strings.Join(registered, ", ")))
			}
		}
	}

	if err := offerInitAgentsMD(opts); err != nil && !errors.Is(err, context.Canceled) {
		if !opts.Quiet && opts.Err != nil {
			fmt.Fprintln(opts.Err, labelWarningValue("AGENTS.md", err.Error()))
		}
	}

	// Claude Code Stop/SessionEnd hooks in .claude/settings.json are owned by WP-C.
	return nil
}

func resolveMCPBinary(gxPath string) (string, error) {
	dir := filepath.Dir(gxPath)
	candidate := filepath.Join(dir, "gx-mcp")
	if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
		return candidate, nil
	}
	if mcpPath, err := exec.LookPath("gx-mcp"); err == nil && mcpPath != "" {
		return mcpPath, nil
	}
	home, err := os.UserHomeDir()
	if err == nil {
		candidate = filepath.Join(home, ".local", "bin", "gx-mcp")
		if info, statErr := os.Stat(candidate); statErr == nil && !info.IsDir() {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("gx-mcp not found; install GX CLI package first")
}

func mcpLaunchCommand(gxPath, mcpPath string) []string {
	return []string{"env", "GX_BINARY=" + gxPath, mcpPath}
}

func registerMCPClients(ctx context.Context, gxPath, mcpPath string) []string {
	args := mcpLaunchCommand(gxPath, mcpPath)
	var registered []string
	for _, spec := range []struct {
		name    string
		binary  string
		command []string
	}{
		{
			name:   "Cursor",
			binary: "cursor",
			command: append([]string{"mcp", "add", "gx", "--"}, args...),
		},
		{
			name:   "Claude Code",
			binary: "claude",
			command: append([]string{"mcp", "add", "gx", "--"}, args...),
		},
	} {
		if _, err := exec.LookPath(spec.binary); err != nil {
			continue
		}
		cmd := exec.CommandContext(ctx, spec.binary, spec.command...)
		if err := cmd.Run(); err != nil {
			if output, combinedErr := runCombined(ctx, spec.binary, spec.command...); combinedErr != nil {
				if strings.Contains(strings.ToLower(output), "already") {
					registered = append(registered, spec.name)
				}
				continue
			}
			continue
		}
		registered = append(registered, spec.name)
	}
	if updated, err := mergeCodexMCPServer(gxPath, mcpPath); err == nil && updated {
		registered = append(registered, "Codex")
	}
	return registered
}

func runCombined(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func mergeCodexMCPServer(gxPath, mcpPath string) (bool, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}
	configPath := filepath.Join(home, ".codex", "config.toml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	content := string(data)
	if strings.Contains(content, "[mcp_servers.gx]") {
		return false, nil
	}
	block := strings.Join([]string{
		"",
		"[mcp_servers.gx]",
		`command = "env"`,
		fmt.Sprintf(`args = ["GX_BINARY=%s", %q]`, gxPath, mcpPath),
		"",
	}, "\n")
	if err := os.WriteFile(configPath, append(data, []byte(block)...), 0o644); err != nil {
		return false, err
	}
	return true, nil
}

func offerInitAgentsMD(opts initSetupOptions) error {
	path := filepath.Join(opts.RepoRoot, "AGENTS.md")
	existing, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if agentsMDContainsSnippet(string(existing)) {
		if !opts.Quiet && opts.Out != nil {
			fmt.Fprintln(opts.Out, labelValue("AGENTS.md", success("ok")+": GX instructions already present"))
		}
		return nil
	}

	if opts.Quiet {
		return nil
	}
	in := opts.In
	if in == nil {
		in = os.Stdin
	}
	out := opts.Out
	if out == nil {
		out = os.Stdout
	}
	fmt.Fprintf(out, "%s ", muted("Add GX workflow instructions to AGENTS.md? [Y/n]"))
	reader := bufio.NewReader(in)
	raw, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "y", "yes":
	default:
		fmt.Fprintln(out, labelValue("AGENTS.md", muted("skipped")))
		return nil
	}

	appended, err := appendAgentsMDSnippet(opts.RepoRoot, string(existing))
	if err != nil {
		return err
	}
	if appended && out != nil {
		fmt.Fprintln(out, labelValue("AGENTS.md", success("ok")+": added GX workflow instructions"))
	}
	return nil
}

func agentsMDContainsSnippet(content string) bool {
	return strings.Contains(content, agentsMDSnippetMarker)
}

func appendAgentsMDSnippet(repoRoot, existing string) (bool, error) {
	path := filepath.Join(repoRoot, "AGENTS.md")
	var builder strings.Builder
	trimmed := strings.TrimSpace(existing)
	if trimmed != "" {
		builder.WriteString(trimmed)
		builder.WriteString("\n\n")
	}
	builder.WriteString(strings.TrimSpace(agentsMDSnippet))
	builder.WriteString("\n")
	if err := os.WriteFile(path, []byte(builder.String()), 0o644); err != nil {
		return false, err
	}
	return true, nil
}

func initSetupFromCommand(cmd *cobra.Command, repoRoot string, quiet bool) error {
	return runInitRepoSetup(initSetupOptions{
		RepoRoot: repoRoot,
		Quiet:    quiet,
		In:       initInput(cmd, quiet),
		Out:      initOutput(cmd, quiet),
		Err:      cmd.ErrOrStderr(),
	})
}
