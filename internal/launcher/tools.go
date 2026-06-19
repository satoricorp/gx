package launcher

import (
	"fmt"
	"os"
	"os/exec"
)

var wrappedTools = map[string]struct{}{
	"cursor": {},
}

type ToolConfig interface {
	Name() string
	EnvVars(port int) map[string]string
	CommandArgs(port int, args []string) []string
	Supports() bool
}

type toolAdapter struct {
	name     string
	supports bool
}

func (t toolAdapter) Name() string { return t.name }

func (t toolAdapter) Supports() bool { return t.supports }

func (t toolAdapter) EnvVars(port int) map[string]string {
	return map[string]string{
		"ANTHROPIC_BASE_URL": fmt.Sprintf("http://127.0.0.1:%d", port),
		"OPENAI_BASE_URL":    fmt.Sprintf("http://127.0.0.1:%d/v1", port),
	}
}

func (t toolAdapter) CommandArgs(port int, args []string) []string {
	if t.name != "codex" {
		return args
	}

	providerArg := `model_provider="gx-openai"`
	providerNameArg := `model_providers.gx-openai.name="GX OpenAI Proxy"`
	baseURLArg := fmt.Sprintf(`model_providers.gx-openai.base_url="http://127.0.0.1:%d/v1"`, port)
	wireAPIArg := `model_providers.gx-openai.wire_api="responses"`
	authArg := `model_providers.gx-openai.requires_openai_auth=true`
	transportArg := `transport="responses_http"`
	prefix := []string{
		"-c", providerArg,
		"-c", providerNameArg,
		"-c", baseURLArg,
		"-c", wireAPIArg,
		"-c", authArg,
		"-c", transportArg,
	}
	return append(prefix, args[1:]...)
}

func toolFor(name string) ToolConfig {
	switch name {
	case "claude", "codex":
		return toolAdapter{name: name, supports: true}
	case "cursor":
		return toolAdapter{name: name, supports: false}
	default:
		return toolAdapter{name: name, supports: true}
	}
}

func IsWrappedTool(name string) bool {
	_, ok := wrappedTools[name]
	return ok
}

func runCursor(args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
