package service

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/satoricorp/gx/internal/storage"
)

const (
	Label          = "dev.gx.capture"
	DefaultProxy   = "127.0.0.1:43123"
	DefaultControl = "127.0.0.1:43124"
	EnvAnthropic   = "ANTHROPIC_BASE_URL"
	EnvOpenAI      = "OPENAI_BASE_URL"
	shellBlockOpen = "# >>> gx ambient capture >>>"
	shellBlockEnd  = "# <<< gx ambient capture <<<"
)

type Runner interface {
	Run(ctx context.Context, name string, args ...string) ([]byte, error)
}

type ExecRunner struct{}

func (ExecRunner) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.Bytes(), err
}

type Manager struct {
	Runner Runner
}

func NewManager() Manager {
	return Manager{Runner: ExecRunner{}}
}

func ProxyAddress() string {
	if value := strings.TrimSpace(os.Getenv("GX_PROXY_ADDR")); value != "" {
		return value
	}
	return DefaultProxy
}

func ControlAddress() string {
	if value := strings.TrimSpace(os.Getenv("GX_DAEMON_ADDR")); value != "" {
		return value
	}
	return DefaultControl
}

func ControlURL() string {
	return "http://" + ControlAddress()
}

func AnthropicBaseURL() string {
	return "http://" + ProxyAddress()
}

func OpenAIBaseURL() string {
	return "http://" + ProxyAddress() + "/v1"
}

func LaunchAgentPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Library", "LaunchAgents", Label+".plist"), nil
}

func (m Manager) Install(ctx context.Context, gxPath string) error {
	if strings.TrimSpace(gxPath) == "" {
		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("resolve gx executable: %w", err)
		}
		gxPath = exe
	}
	agentPath, err := LaunchAgentPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(agentPath), 0o755); err != nil {
		return fmt.Errorf("create launch agents dir: %w", err)
	}
	if err := os.MkdirAll(logDir(), 0o755); err != nil {
		return fmt.Errorf("create log dir: %w", err)
	}
	if err := os.WriteFile(agentPath, []byte(LaunchAgentPlist(gxPath)), 0o644); err != nil {
		return fmt.Errorf("write launch agent: %w", err)
	}
	if err := m.ApplyEnv(ctx); err != nil {
		return err
	}
	if err := InstallShellBlocks(); err != nil {
		return err
	}
	target := guiTarget()
	_, _ = m.Runner.Run(ctx, "launchctl", "bootout", target+"/"+Label)
	return m.loadAndStart(ctx, target)
}

func (m Manager) loadAndStart(ctx context.Context, gui string) error {
	agentPath, err := LaunchAgentPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(agentPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("launch agent not installed at %s: run `gx ops capture install`", agentPath)
		}
		return fmt.Errorf("stat launch agent: %w", err)
	}
	_, _ = m.Runner.Run(ctx, "launchctl", "bootout", gui+"/"+Label)
	if out, err := m.Runner.Run(ctx, "launchctl", "bootstrap", gui, agentPath); err != nil {
		return fmt.Errorf("launchctl bootstrap: %w: %s", err, strings.TrimSpace(string(out)))
	}
	if out, err := m.Runner.Run(ctx, "launchctl", "enable", gui+"/"+Label); err != nil {
		return fmt.Errorf("launchctl enable: %w: %s", err, strings.TrimSpace(string(out)))
	}
	if out, err := m.Runner.Run(ctx, "launchctl", "kickstart", "-k", gui+"/"+Label); err != nil {
		return fmt.Errorf("launchctl kickstart: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func launchctlServiceMissing(out []byte, err error) bool {
	msg := strings.ToLower(strings.TrimSpace(string(out)))
	if err != nil {
		msg += " " + strings.ToLower(err.Error())
	}
	return strings.Contains(msg, "could not find service") ||
		strings.Contains(msg, "no such process") ||
		strings.Contains(msg, "not found")
}

func (m Manager) Uninstall(ctx context.Context) error {
	target := guiTarget()
	_, _ = m.Runner.Run(ctx, "launchctl", "bootout", target+"/"+Label)
	if err := m.ClearEnv(ctx); err != nil {
		return err
	}
	if err := RemoveShellBlocks(); err != nil {
		return err
	}
	agentPath, err := LaunchAgentPath()
	if err != nil {
		return err
	}
	if err := os.Remove(agentPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove launch agent: %w", err)
	}
	return nil
}

func (m Manager) Start(ctx context.Context) error {
	agentPath, err := LaunchAgentPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(agentPath); err != nil {
		if os.IsNotExist(err) {
			return m.Install(ctx, "")
		}
		return fmt.Errorf("stat launch agent: %w", err)
	}
	if err := m.ApplyEnv(ctx); err != nil {
		return err
	}
	gui := guiTarget()
	target := gui + "/" + Label
	if out, err := m.Runner.Run(ctx, "launchctl", "kickstart", "-k", target); err == nil {
		return nil
	} else if !launchctlServiceMissing(out, err) {
		return fmt.Errorf("launchctl kickstart: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return m.loadAndStart(ctx, gui)
}

func (m Manager) Stop(ctx context.Context) error {
	target := guiTarget() + "/" + Label
	if out, err := m.Runner.Run(ctx, "launchctl", "bootout", target); err != nil {
		if launchctlServiceMissing(out, err) {
			return nil
		}
		return fmt.Errorf("launchctl bootout: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func (m Manager) Status(ctx context.Context) (string, error) {
	out, err := m.Runner.Run(ctx, "launchctl", "print", guiTarget()+"/"+Label)
	if err != nil {
		return strings.TrimSpace(string(out)), err
	}
	return strings.TrimSpace(string(out)), nil
}

func (m Manager) ApplyEnv(ctx context.Context) error {
	for key, value := range map[string]string{
		EnvAnthropic: AnthropicBaseURL(),
		EnvOpenAI:    OpenAIBaseURL(),
	} {
		if out, err := m.Runner.Run(ctx, "launchctl", "setenv", key, value); err != nil {
			return fmt.Errorf("launchctl setenv %s: %w: %s", key, err, strings.TrimSpace(string(out)))
		}
	}
	return nil
}

func (m Manager) ClearEnv(ctx context.Context) error {
	for _, key := range []string{EnvAnthropic, EnvOpenAI} {
		if out, err := m.Runner.Run(ctx, "launchctl", "unsetenv", key); err != nil {
			return fmt.Errorf("launchctl unsetenv %s: %w: %s", key, err, strings.TrimSpace(string(out)))
		}
	}
	return nil
}

func (m Manager) EnvStatus(ctx context.Context) (map[string]string, error) {
	values := map[string]string{}
	for _, key := range []string{EnvAnthropic, EnvOpenAI} {
		out, err := m.Runner.Run(ctx, "launchctl", "getenv", key)
		if err != nil {
			return values, fmt.Errorf("launchctl getenv %s: %w: %s", key, err, strings.TrimSpace(string(out)))
		}
		values[key] = strings.TrimSpace(string(out))
	}
	return values, nil
}

func LaunchAgentPlist(gxPath string) string {
	home, _ := os.UserHomeDir()
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>%s</string>
  <key>ProgramArguments</key>
  <array>
    <string>%s</string>
    <string>__gx-daemon</string>
    <string>--ambient</string>
  </array>
  <key>RunAtLoad</key>
  <true/>
  <key>KeepAlive</key>
  <true/>
  <key>StandardOutPath</key>
  <string>%s</string>
  <key>StandardErrorPath</key>
  <string>%s</string>
  <key>WorkingDirectory</key>
  <string>%s</string>
</dict>
</plist>
`, Label, gxPath, filepath.Join(logDir(), "daemon.out.log"), filepath.Join(logDir(), "daemon.err.log"), home)
}

func InstallShellBlocks() error {
	for _, path := range shellProfilePaths() {
		if err := upsertBlock(path, ShellBlock()); err != nil {
			return err
		}
	}
	return nil
}

func RemoveShellBlocks() error {
	for _, path := range shellProfilePaths() {
		if err := removeBlock(path); err != nil {
			return err
		}
	}
	return nil
}

func ShellBlock() string {
	return strings.Join([]string{
		shellBlockOpen,
		"export " + EnvAnthropic + "=" + shellQuote(AnthropicBaseURL()),
		"export " + EnvOpenAI + "=" + shellQuote(OpenAIBaseURL()),
		shellBlockEnd,
		"",
	}, "\n")
}

func shellProfilePaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	return []string{
		filepath.Join(home, ".zshrc"),
		filepath.Join(home, ".bashrc"),
	}
}

func upsertBlock(path, block string) error {
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read %s: %w", path, err)
	}
	text := string(data)
	text = stripBlock(text)
	if strings.TrimSpace(text) != "" && !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	text += block
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func removeBlock(path string) error {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	text := stripBlock(string(data))
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func stripBlock(text string) string {
	start := strings.Index(text, shellBlockOpen)
	if start < 0 {
		return text
	}
	end := strings.Index(text[start:], shellBlockEnd)
	if end < 0 {
		return text
	}
	end += start + len(shellBlockEnd)
	if end < len(text) && text[end] == '\n' {
		end++
	}
	return strings.TrimRight(text[:start]+text[end:], "\n") + "\n"
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func guiTarget() string {
	return "gui/" + strconv.Itoa(os.Getuid())
}

func logDir() string {
	dir, err := storage.DefaultDir()
	if err != nil {
		return filepath.Join(os.TempDir(), "gx", "logs")
	}
	return filepath.Join(dir, "logs")
}
