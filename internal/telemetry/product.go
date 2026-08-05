package telemetry

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/satoricorp/lgtm/internal/cloud"
	"github.com/satoricorp/lgtm/internal/storage"
	"github.com/satoricorp/lgtm/internal/version"
)

// stateWritesKey marks a context as belonging to a command that must not
// create lgtm state on the machine.
type stateWritesKeyType struct{}

var stateWritesKey stateWritesKeyType

// WithoutStateWrites marks ctx as read-only for the machine. Telemetry still
// reports under it; it just never writes anything down to do so. `lgtm review`
// runs as a CI gate and on checkouts the reviewer does not own, and minting a
// machine ID to label an event would leave $LGTM_HOME behind on a machine that
// never ran lgtm — buying an analytics dimension with the promise the command
// makes.
func WithoutStateWrites(ctx context.Context) context.Context {
	return context.WithValue(ctx, stateWritesKey, true)
}

func stateWritesAllowed(ctx context.Context) bool {
	if ctx == nil {
		return true
	}
	blocked, _ := ctx.Value(stateWritesKey).(bool)
	return !blocked
}

// EmitProductEvent sends low-cardinality product analytics when PostHog is configured.
func EmitProductEvent(ctx context.Context, event string, properties map[string]any) {
	if !Configured() {
		return
	}
	NewFromEnv().EmitEvent(ctx, event, properties)
}

func EmitInstallOnce(ctx context.Context) {
	if !Configured() {
		return
	}
	// The install event is remembered by writing a sentinel under $LGTM_HOME, so
	// there is no way to send it without creating lgtm state. A command that
	// promised not to skips it rather than sending the event and forgetting.
	if !stateWritesAllowed(ctx) {
		return
	}
	source := Entrypoint()
	path, err := installSentinelPath(source)
	if err != nil {
		return
	}
	if _, err := os.Stat(path); err == nil {
		return
	}
	EmitProductEvent(ctx, EventCLIInstall, map[string]any{
		"source": source,
		"os":     runtime.GOOS,
		"arch":   runtime.GOARCH,
	})
	_ = os.MkdirAll(filepath.Dir(path), 0o700)
	_ = os.WriteFile(path, []byte("1\n"), 0o600)
}

func ProductProperties(ctx context.Context, properties map[string]any) map[string]any {
	out := map[string]any{
		"lgtm_version": version.Current(),
		"entrypoint": Entrypoint(),
	}
	if properties != nil {
		for key, value := range properties {
			out[key] = value
		}
	}

	if creds, err := cloud.LoadCloudCredentials(); err == nil && creds != nil {
		if id := strings.TrimSpace(creds.UserID); id != "" {
			out["user_id"] = id
			if _, ok := out["distinct_id"]; !ok {
				out["distinct_id"] = id
			}
		}
		if login := strings.TrimSpace(creds.Login); login != "" {
			out["login"] = login
		}
		if machineID := strings.TrimSpace(creds.MachineID); machineID != "" {
			out["machine_id"] = machineID
			if _, ok := out["distinct_id"]; !ok {
				out["distinct_id"] = machineID
			}
		}
	} else if _, ok := out["distinct_id"]; !ok {
		// Signed out: fall back to the machine ID. Minting one writes
		// machine_id.json and creates $LGTM_HOME, so a read-only command reads
		// the existing ID and otherwise reports anonymously. The first command
		// that legitimately writes lgtm state mints it for everyone after.
		machineID, err := "", error(nil)
		if stateWritesAllowed(ctx) {
			machineID, err = cloud.DefaultMachineID()
		} else {
			machineID, err = cloud.ExistingMachineID()
		}
		if err == nil && strings.TrimSpace(machineID) != "" {
			out["machine_id"] = strings.TrimSpace(machineID)
			out["distinct_id"] = strings.TrimSpace(machineID)
		}
	}
	if _, ok := out["distinct_id"]; !ok {
		out["distinct_id"] = "anonymous"
	}
	return out
}

func installSentinelPath(source string) (string, error) {
	dir, err := storage.DefaultDir()
	if err != nil {
		return "", err
	}
	name := strings.NewReplacer("/", "_", "\\", "_", ":", "_", " ", "_").Replace(source + "-" + version.Current())
	if name == "" {
		name = "unknown"
	}
	return filepath.Join(dir, "telemetry", "install-"+name), nil
}

func Entrypoint() string {
	if strings.TrimSpace(os.Getenv("LGTM_MCP")) != "" {
		return "mcp"
	}
	return "cli"
}
