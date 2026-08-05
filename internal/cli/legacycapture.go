package cli

// Legacy ambient-capture cleanup.
//
// lgtm used to install a launchd LaunchAgent (dev.lgtm.capture) that ran a local
// proxy daemon and pointed coding agents at it through launchctl environment
// variables and managed blocks in shell profiles. That capture model is gone —
// transcript capture is hook-driven — but machines that ran the old
// `lgtm ops capture install` still have the agent loaded. This best-effort
// cleanup runs during `lgtm init` and repo auto-init: it stops and removes the
// LaunchAgent, strips the managed shell-profile blocks, clears the lgtm
// launchctl env vars, and reverts ~/.codex/config.toml when it still routes
// Codex through the retired local proxy. It is a silent no-op when nothing
// legacy is present and never fails the caller; problems surface as warnings
// at most.

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	legacyLaunchAgentLabel = "dev.lgtm.capture"
	legacyDefaultProxyAddr = "127.0.0.1:43123"
	legacyShellBlockOpen   = "# >>> lgtm ambient capture >>>"
	legacyShellBlockClose  = "# <<< lgtm ambient capture <<<"
	legacyEnvAnthropicKey  = "ANTHROPIC_BASE_URL"
	legacyEnvOpenAIKey     = "OPENAI_BASE_URL"
	legacyCodexProviderID  = "lgtm-openai"
)

// legacyRunner abstracts launchctl invocations so the cleanup decision logic
// is unit-testable without touching launchd.
type legacyRunner func(ctx context.Context, name string, args ...string) ([]byte, error)

func execLegacyRunner(ctx context.Context, name string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, name, args...).CombinedOutput()
}

type legacyCleanupResult struct {
	CleanedLaunchAgent bool
	CleanedShellBlocks bool
	CleanedEnv         bool
	CleanedCodexConfig bool
	Warnings           []string
}

func (r legacyCleanupResult) cleanedAny() bool {
	return r.CleanedLaunchAgent || r.CleanedShellBlocks || r.CleanedEnv || r.CleanedCodexConfig
}

// legacyLaunchAgentPath returns the plist path the retired
// `lgtm ops capture install` wrote: ~/Library/LaunchAgents/dev.lgtm.capture.plist.
func legacyLaunchAgentPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Library", "LaunchAgents", legacyLaunchAgentLabel+".plist"), nil
}

// legacyLaunchAgentTarget is the launchctl service target for the retired
// agent, e.g. gui/501/dev.lgtm.capture.
func legacyLaunchAgentTarget() string {
	return "gui/" + strconv.Itoa(os.Getuid()) + "/" + legacyLaunchAgentLabel
}

// legacyProxyBaseURLs returns the base URLs the retired proxy advertised.
// LGTM_PROXY_ADDR was the old override knob; it is honored here only so cleanup
// recognizes env values written by a customized install.
func legacyProxyBaseURLs() (anthropicURL, openaiURL string) {
	addr := strings.TrimSpace(os.Getenv("LGTM_PROXY_ADDR"))
	if addr == "" {
		addr = legacyDefaultProxyAddr
	}
	return "http://" + addr, "http://" + addr + "/v1"
}

func legacyShellProfilePaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	return []string{
		filepath.Join(home, ".zshrc"),
		filepath.Join(home, ".bashrc"),
	}
}

// stripLegacyShellBlock removes the first managed ambient-capture block and
// reports whether anything was removed. Text without both markers is returned
// unchanged.
func stripLegacyShellBlock(text string) (string, bool) {
	start := strings.Index(text, legacyShellBlockOpen)
	if start < 0 {
		return text, false
	}
	end := strings.Index(text[start:], legacyShellBlockClose)
	if end < 0 {
		return text, false
	}
	end += start + len(legacyShellBlockClose)
	if end < len(text) && text[end] == '\n' {
		end++
	}
	return strings.TrimRight(text[:start]+text[end:], "\n") + "\n", true
}

func removeLegacyShellBlockFromFile(path string) (bool, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	text := string(data)
	removed := false
	for {
		next, changed := stripLegacyShellBlock(text)
		if !changed {
			break
		}
		text = next
		removed = true
	}
	if !removed {
		return false, nil
	}
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		return false, err
	}
	return true, nil
}

// legacyCodexConfigPath returns the Codex config the retired proxy install
// rewired: ~/.codex/config.toml.
func legacyCodexConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".codex", "config.toml"), nil
}

// legacyCodexProxyHosts returns the host:port values recognized as the
// retired local proxy: the default listen address plus, when set, the old
// LGTM_PROXY_ADDR override a customized install wrote into base_url.
func legacyCodexProxyHosts() []string {
	hosts := []string{legacyDefaultProxyAddr}
	if addr := strings.TrimSpace(os.Getenv("LGTM_PROXY_ADDR")); addr != "" && addr != legacyDefaultProxyAddr {
		hosts = append(hosts, addr)
	}
	return hosts
}

// findLegacyTOMLRootStringValue and findLegacyTOMLTableStringValue are the
// retired service writer's line-based TOML readers (internal/service/codex.go
// before its removal), kept verbatim so every config that writer produced or
// repaired parses the same way here.
func findLegacyTOMLRootStringValue(text, key string) string {
	prefix := key + " "
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") {
			return ""
		}
		if !strings.HasPrefix(trimmed, prefix) && !strings.HasPrefix(trimmed, key+"=") {
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			continue
		}
		return strings.Trim(strings.TrimSpace(parts[1]), `"'`)
	}
	return ""
}

func findLegacyTOMLTableStringValue(text, table, key string) string {
	inTable := false
	prefix := key + " "
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			inTable = trimmed == "["+table+"]"
			continue
		}
		if !inTable {
			continue
		}
		if !strings.HasPrefix(trimmed, prefix) && !strings.HasPrefix(trimmed, key+"=") {
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			continue
		}
		return strings.Trim(strings.TrimSpace(parts[1]), `"'`)
	}
	return ""
}

func legacyCodexHasProviderBlock(text string) bool {
	header := "[model_providers." + legacyCodexProviderID + "]"
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == header {
			return true
		}
	}
	return false
}

// legacyCodexBaseURLPointsAtProxy reports whether a provider base_url targets
// the retired proxy's host:port.
func legacyCodexBaseURLPointsAtProxy(baseURL string) bool {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil || parsed.Host == "" {
		return false
	}
	for _, host := range legacyCodexProxyHosts() {
		if parsed.Host == host {
			return true
		}
	}
	return false
}

// legacyCodexConfigPointsAtRetiredProxy is the gate for the revert: only
// configs that provably still route Codex through the dead proxy are touched.
// With the provider block present that proof is its base_url; without the
// block, a bare `model_provider = "lgtm-openai"` assignment is still the old
// install's leftover (the provider id is lgtm's own) and it breaks Codex by
// selecting a provider that no longer exists.
func legacyCodexConfigPointsAtRetiredProxy(text string) bool {
	if legacyCodexHasProviderBlock(text) {
		baseURL := findLegacyTOMLTableStringValue(text, "model_providers."+legacyCodexProviderID, "base_url")
		return legacyCodexBaseURLPointsAtProxy(baseURL)
	}
	return findLegacyTOMLRootStringValue(text, "model_provider") == legacyCodexProviderID
}

// removeLegacyCodexProviderAssignment drops the root-level
// `model_provider = "lgtm-openai"` line. Assignments selecting any other
// provider are left alone.
func removeLegacyCodexProviderAssignment(lines []string) ([]string, bool) {
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") {
			return lines, false
		}
		if !strings.HasPrefix(trimmed, "model_provider ") && !strings.HasPrefix(trimmed, "model_provider=") {
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			continue
		}
		if strings.Trim(strings.TrimSpace(parts[1]), `"'`) != legacyCodexProviderID {
			continue
		}
		return append(append([]string{}, lines[:i]...), lines[i+1:]...), true
	}
	return lines, false
}

// removeLegacyCodexProviderBlock drops the [model_providers.lgtm-openai] table:
// its header through the line before the next table header, or through EOF.
func removeLegacyCodexProviderBlock(lines []string) ([]string, bool) {
	header := "[model_providers." + legacyCodexProviderID + "]"
	start := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == header {
			start = i
			break
		}
	}
	if start < 0 {
		return lines, false
	}
	end := len(lines)
	for i := start + 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			end = i
			break
		}
	}
	if end == len(lines) {
		// The block runs to EOF. Keep the final empty element that encodes
		// the file's trailing newline, and consume the blank separator the
		// retired writer inserted above the block it appended.
		if end > start && lines[end-1] == "" {
			end--
		}
		for start > 0 && strings.TrimSpace(lines[start-1]) == "" {
			start--
		}
	}
	return append(append([]string{}, lines[:start]...), lines[end:]...), true
}

// revertLegacyCodexConfigText removes the retired proxy wiring — the
// `model_provider = "lgtm-openai"` assignment and the
// [model_providers.lgtm-openai] block — leaving every other line byte-identical.
func revertLegacyCodexConfigText(text string) (string, bool) {
	lines := strings.Split(text, "\n")
	lines, removedAssignment := removeLegacyCodexProviderAssignment(lines)
	lines, removedBlock := removeLegacyCodexProviderBlock(lines)
	if !removedAssignment && !removedBlock {
		return text, false
	}
	return strings.Join(lines, "\n"), true
}

// revertLegacyCodexConfigFile reverts one Codex config when it still routes
// through the retired proxy. It mirrors the retired writer's file handling:
// a timestamped backup beside the file before any modification, then a
// line-based edit. Missing files and configs pointing elsewhere are silent
// no-ops.
func revertLegacyCodexConfigFile(path string) (bool, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	text := string(data)
	if !legacyCodexConfigPointsAtRetiredProxy(text) {
		return false, nil
	}
	reverted, changed := revertLegacyCodexConfigText(text)
	if !changed {
		return false, nil
	}
	backup := fmt.Sprintf("%s.lgtm-backup-%d", path, time.Now().Unix())
	if err := os.WriteFile(backup, data, 0o600); err != nil {
		return false, fmt.Errorf("write codex config backup: %w", err)
	}
	if err := os.WriteFile(path, []byte(reverted), 0o600); err != nil {
		return false, fmt.Errorf("write codex config: %w", err)
	}
	return true, nil
}

// cleanupLegacyAmbientCapture removes what the retired ambient-capture
// service left behind. Every step is best-effort; failures accumulate as
// warnings and never abort the remaining steps.
func cleanupLegacyAmbientCapture(ctx context.Context, run legacyRunner) legacyCleanupResult {
	if run == nil {
		run = execLegacyRunner
	}
	var result legacyCleanupResult
	warnf := func(format string, args ...any) {
		result.Warnings = append(result.Warnings, fmt.Sprintf(format, args...))
	}

	plistPath, err := legacyLaunchAgentPath()
	if err != nil {
		return result
	}
	plistPresent := false
	if _, statErr := os.Stat(plistPath); statErr == nil {
		plistPresent = true
	}
	if plistPresent {
		// Stop the daemon before deleting its definition. Bootout failures
		// (agent not loaded, launchctl missing) are expected and ignored.
		_, _ = run(ctx, "launchctl", "bootout", legacyLaunchAgentTarget())
		if err := os.Remove(plistPath); err != nil && !os.IsNotExist(err) {
			warnf("remove legacy LaunchAgent %s: %v", plistPath, err)
		} else {
			result.CleanedLaunchAgent = true
		}
	}

	for _, path := range legacyShellProfilePaths() {
		removed, err := removeLegacyShellBlockFromFile(path)
		if err != nil {
			warnf("remove lgtm ambient capture block from %s: %v", path, err)
			continue
		}
		if removed {
			result.CleanedShellBlocks = true
		}
	}

	// Touch launchctl env only when other artifacts prove the old install was
	// present, and unset a variable only while it still holds the retired
	// proxy URL, so values the user set themselves survive.
	if plistPresent || result.CleanedShellBlocks {
		anthropicURL, openaiURL := legacyProxyBaseURLs()
		for _, entry := range []struct {
			key  string
			want string
		}{
			{key: legacyEnvAnthropicKey, want: anthropicURL},
			{key: legacyEnvOpenAIKey, want: openaiURL},
		} {
			out, err := run(ctx, "launchctl", "getenv", entry.key)
			if err != nil {
				continue
			}
			if strings.TrimSpace(string(out)) != entry.want {
				continue
			}
			if _, err := run(ctx, "launchctl", "unsetenv", entry.key); err != nil {
				warnf("launchctl unsetenv %s: %v", entry.key, err)
				continue
			}
			result.CleanedEnv = true
		}
	}

	// The old install also rewired ~/.codex/config.toml at the retired proxy;
	// with the daemon gone that routing breaks Codex entirely. This step is
	// gated on the config's own content, so it runs even when the other
	// artifacts are already gone.
	if codexPath, err := legacyCodexConfigPath(); err == nil {
		reverted, err := revertLegacyCodexConfigFile(codexPath)
		if err != nil {
			warnf("revert legacy codex proxy config %s: %v", codexPath, err)
		} else if reverted {
			result.CleanedCodexConfig = true
		}
	}
	return result
}

// cleanupLegacyAmbientCaptureQuiet is the auto-init hook: cleanup with
// warnings only, so hook-driven invocations (e.g. inside git push) stay quiet
// unless something actually went wrong.
func cleanupLegacyAmbientCaptureQuiet(ctx context.Context, errOut io.Writer) {
	result := cleanupLegacyAmbientCapture(ctx, nil)
	for _, warning := range result.Warnings {
		fmt.Fprintln(errOut, labelWarningValue("Legacy capture cleanup", warning))
	}
}

// cleanupLegacyAmbientCaptureFromInit runs during `lgtm init`, reporting a
// one-line note when the retired service was actually removed.
func cleanupLegacyAmbientCaptureFromInit(ctx context.Context, cmd *cobra.Command, quiet bool) {
	result := cleanupLegacyAmbientCapture(ctx, nil)
	for _, warning := range result.Warnings {
		fmt.Fprintln(cmd.ErrOrStderr(), labelWarningValue("Legacy capture", warning))
	}
	if quiet {
		return
	}
	if result.CleanedLaunchAgent || result.CleanedShellBlocks || result.CleanedEnv {
		fmt.Fprintln(cmd.OutOrStdout(), labelValue("Legacy capture", "removed retired background capture service ("+legacyLaunchAgentLabel+")"))
	}
	if result.CleanedCodexConfig {
		fmt.Fprintln(cmd.OutOrStdout(), labelValue("Legacy capture", "reverted ~/.codex/config.toml routing through the retired proxy (backup saved beside it)"))
	}
}
