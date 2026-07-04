package codereview

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const maxStaticToolOutputBytes = 12000

type StaticToolResult struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	ExitCode int    `json:"exit_code"`
	Output   string `json:"output,omitempty"`
	Skipped  bool   `json:"skipped,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

func collectStaticToolResults(ctx context.Context, repoRoot string, facts RepoFacts, opts Options) []StaticToolResult {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("GX_REVIEW_STATIC_TOOLS")), "0") {
		return nil
	}
	if !hasDependencyFile(facts.DependencyFiles, "go.mod") {
		return nil
	}
	pkgs, ok := changedGoPackageArgs(ctx, repoRoot)
	if !ok {
		return nil
	}
	timeout := 90 * time.Second
	if opts.Deep {
		timeout = 180 * time.Second
	}
	tools := []struct {
		progress string
		name     string
		command  []string
	}{
		{progress: "Running go test", name: "go test", command: append([]string{"go", "test"}, pkgs...)},
		{progress: "Running go vet", name: "go vet", command: append([]string{"go", "vet"}, pkgs...)},
	}
	out := make([]StaticToolResult, len(tools))
	if !opts.Deep {
		for i, tool := range tools {
			reviewProgress(opts, tool.progress)
			out[i] = runStaticTool(ctx, repoRoot, timeout, tool.name, tool.command...)
		}
		return out
	}

	var wg sync.WaitGroup
	for i, tool := range tools {
		i, tool := i, tool
		wg.Add(1)
		go func() {
			defer wg.Done()
			reviewProgress(opts, tool.progress)
			out[i] = runStaticTool(ctx, repoRoot, timeout, tool.name, tool.command...)
		}()
	}
	wg.Wait()
	return out
}

func changedGoPackageArgs(ctx context.Context, repoRoot string) ([]string, bool) {
	return goPackageArgsForChangedFiles(repoRoot, reviewChangedFiles(ctx, repoRoot))
}

func goPackageArgsForChangedFiles(repoRoot string, files []string) ([]string, bool) {
	files = normalizedChangedFiles(files)
	if len(files) == 0 {
		return nil, false
	}
	dirs := map[string]struct{}{}
	for _, file := range files {
		if goStaticToolsRequireRepoWide(file) {
			return []string{"./..."}, true
		}
		if !strings.HasSuffix(file, ".go") {
			continue
		}
		dir := filepath.ToSlash(filepath.Dir(file))
		if dir == "" {
			return []string{"./..."}, true
		}
		if dir == "." {
			dirs["."] = struct{}{}
			continue
		}
		if _, err := os.Stat(filepath.Join(repoRoot, filepath.FromSlash(dir))); err != nil {
			return []string{"./..."}, true
		}
		dirs[dir] = struct{}{}
	}
	if len(dirs) == 0 {
		return nil, false
	}
	if len(dirs) > 20 {
		return []string{"./..."}, true
	}
	out := make([]string, 0, len(dirs))
	for dir := range dirs {
		if dir == "." {
			out = append(out, ".")
		} else {
			out = append(out, "./"+dir)
		}
	}
	sort.Strings(out)
	return out, true
}

func goStaticToolsRequireRepoWide(file string) bool {
	file = filepath.ToSlash(strings.TrimSpace(file))
	base := filepath.Base(file)
	switch base {
	case "go.mod", "go.sum", "go.work", "go.work.sum", "Makefile", "Taskfile.yml", "Taskfile.yaml", "Dockerfile", "magefile.go":
		return true
	}
	return strings.HasPrefix(file, ".github/workflows/") ||
		strings.HasPrefix(file, "build/") ||
		strings.HasPrefix(file, "scripts/")
}

func runStaticTool(ctx context.Context, repoRoot string, timeout time.Duration, name string, command ...string) StaticToolResult {
	if len(command) == 0 {
		return StaticToolResult{Name: name, Skipped: true, Reason: "empty command"}
	}
	toolCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	cmd := exec.CommandContext(toolCtx, command[0], command[1:]...)
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), "GOCACHE="+filepath.Join(os.TempDir(), "gx-review-gocache"))
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	err := cmd.Run()
	text := output.String()
	if len(text) > maxStaticToolOutputBytes {
		text = text[:maxStaticToolOutputBytes] + "\n[truncated]\n"
	}
	result := StaticToolResult{
		Name:    name,
		Command: strings.Join(command, " "),
		Output:  strings.TrimSpace(text),
	}
	if err != nil {
		result.ExitCode = 1
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		}
		if toolCtx.Err() == context.DeadlineExceeded {
			result.Skipped = true
			result.Reason = "timed out"
		}
	}
	return result
}

func hasDependencyFile(files []string, base string) bool {
	for _, file := range files {
		if filepath.Base(file) == base {
			return true
		}
	}
	return false
}
