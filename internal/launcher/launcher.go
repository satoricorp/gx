package launcher

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/satoricorp/gx/internal/inference"
	"github.com/satoricorp/gx/internal/vcs"
	"github.com/satoricorp/gx/internal/version"
)

func Run(ctx context.Context, args []string) error {
	tool := toolFor(args[0])
	if tool.Name() == "cursor" && !tool.Supports() {
		return runCursor(args)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	repoResult, err := vcs.NewService().EnsureRepoAtPath(ctx, cwd)
	if err != nil {
		return err
	}
	if repoResult.Initialized {
		fmt.Fprintf(os.Stderr, "Initialized gx repo at %s\n", repoResult.Repo.RootPath)
	}
	manager := NewDaemonManager()
	controlURL, err := manager.Ensure(ctx)
	if err != nil {
		return err
	}
	session, err := manager.CreateSession(ctx, controlURL, strings.Join(args, " "), cwd, version.Current())
	if err != nil {
		return err
	}

	childArgs := tool.CommandArgs(session.Port, args)
	cmd := exec.CommandContext(ctx, args[0], childArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	envVars := tool.EnvVars(session.Port)
	envVars["GX_SESSION_ID"] = session.SessionID
	for key, value := range inference.EnvVars() {
		envVars[key] = value
	}
	cmd.Env = append(baseEnvironmentWithoutInferenceKeys(), flattenEnv(envVars)...)

	if err := cmd.Start(); err != nil {
		_ = manager.EndSession(context.Background(), controlURL, session.SessionID, 127)
		return fmt.Errorf("start child process: %w", err)
	}

	go relaySignals(cmd.Process)

	waitErr := cmd.Wait()
	exitCode := exitCode(waitErr)
	endErr := manager.EndSession(context.Background(), controlURL, session.SessionID, exitCode)
	if waitErr != nil {
		if endErr != nil {
			return fmt.Errorf("%v; %w", waitErr, endErr)
		}
		return waitErr
	}
	return endErr
}

func baseEnvironmentWithoutInferenceKeys() []string {
	env := os.Environ()
	out := make([]string, 0, len(env))
	for _, entry := range env {
		if strings.HasPrefix(entry, "ANTHROPIC_API_KEY=") || strings.HasPrefix(entry, "OPENAI_API_KEY=") {
			continue
		}
		out = append(out, entry)
	}
	return out
}

func flattenEnv(values map[string]string) []string {
	env := make([]string, 0, len(values))
	for key, value := range values {
		env = append(env, key+"="+value)
	}
	return env
}

func relaySignals(process *os.Process) {
	signals := make(chan os.Signal, 4)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	defer signal.Stop(signals)
	for sig := range signals {
		if process == nil {
			return
		}
		_ = process.Signal(sig)
	}
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return 1
	}
	status, ok := exitErr.Sys().(syscall.WaitStatus)
	if !ok {
		return 1
	}
	if status.Signaled() {
		return 128 + int(status.Signal())
	}
	return status.ExitStatus()
}
