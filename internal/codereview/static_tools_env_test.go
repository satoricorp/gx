package codereview

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// reviewerCredentialEnv is the set this whole file exists to keep out of the
// reviewed checkout's processes. These are the real names a Totality developer has
// exported (see internal/cloud/credentials.go and internal/telemetry) plus the
// two cloud conventions a CI runner exports into the step that runs tx review.
var reviewerCredentialEnv = map[string]string{
	"ANTHROPIC_API_KEY":              "sk-ant-leaked",
	"OPENAI_API_KEY":                 "sk-openai-leaked",
	"TURBOPUFFER_API_KEY":            "tpuf-leaked",
	"GITHUB_TOKEN":                   "ghp-leaked",
	"GH_TOKEN":                       "gho-leaked",
	"TOTALITY_UPLOAD_TOKEN":          "totality-leaked",
	"TOTALITY_API_URL":               "https://internal.example.invalid",
	"AWS_ACCESS_KEY_ID":              "AKIALEAKED",
	"AWS_SECRET_ACCESS_KEY":          "aws-secret-leaked",
	"AWS_SESSION_TOKEN":              "aws-session-leaked",
	"GOOGLE_API_KEY":                 "google-leaked",
	"GOOGLE_APPLICATION_CREDENTIALS": "/home/reviewer/.config/gcloud/creds.json",
	"NPM_TOKEN":                      "npm-leaked",
	"npm_config__authToken":          "npm-registry-leaked",
	"DATABASE_URL":                   "postgres://user:pw@db.internal/prod",
	"SSH_AUTH_SOCK":                  "/private/tmp/ssh-agent.sock",
	"SENTRY_DSN":                     "https://key@sentry.example.invalid/1",
}

// The static tool stage runs `go test` from the checkout under review, so its
// environment is readable by code the reviewer did not write. This pins the
// contract: none of the reviewer's credentials cross that boundary.
func TestStaticToolChildEnvDropsTheReviewersCredentials(t *testing.T) {
	parent := []string{"PATH=/usr/bin", "HOME=/home/reviewer"}
	for name, value := range reviewerCredentialEnv {
		parent = append(parent, name+"="+value)
	}

	child := staticToolChildEnv(parent)
	joined := strings.Join(child, "\n")
	for name, value := range reviewerCredentialEnv {
		if strings.Contains(joined, name+"=") {
			t.Fatalf("child environment carries %s; the reviewed checkout must not see it:\n%s", name, joined)
		}
		if strings.Contains(joined, value) {
			t.Fatalf("child environment carries the value of %s; the reviewed checkout must not see it:\n%s", name, joined)
		}
	}
	// GOOGLE_API_KEY above is the reason the Go entries are listed one by one
	// instead of matched with a "GO" prefix: a prefix rule would have admitted
	// both Google variables while looking like it only allowed the toolchain.
	if strings.Contains(joined, "GOOGLE") {
		t.Fatalf("a GO-prefixed rule leaked a Google credential:\n%s", joined)
	}
}

func TestStaticToolChildEnvKeepsTheVariablesTheToolchainsNeed(t *testing.T) {
	parent := []string{
		"PATH=/usr/bin:/bin",
		"HOME=/home/reviewer",
		"LANG=en_US.UTF-8",
		"LC_CTYPE=en_US.UTF-8",
		"GOPATH=/home/reviewer/go",
		"GOMODCACHE=/home/reviewer/go/pkg/mod",
		"GOFLAGS=-mod=mod",
		"GOPROXY=https://proxy.example",
		"GOTOOLCHAIN=local",
		"CGO_ENABLED=1",
		"CC=clang",
		"HTTPS_PROXY=http://proxy.corp:3128",
		"SSL_CERT_FILE=/etc/ssl/corp.pem",
		"CARGO_HOME=/home/reviewer/.cargo",
		"RUSTUP_HOME=/home/reviewer/.rustup",
		"NODE_PATH=/home/reviewer/lib/node",
		"VIRTUAL_ENV=/repo/.venv",
		"PYENV_ROOT=/home/reviewer/.pyenv",
		"SystemRoot=C:\\Windows",
		"ProgramFiles(x86)=C:\\Program Files (x86)",
		"TMPDIR=/scratch",
		"CI=true",
	}

	child := staticToolChildEnv(parent)
	got := map[string]string{}
	for _, entry := range child {
		name, value, _ := strings.Cut(entry, "=")
		got[name] = value
	}
	for _, entry := range parent {
		name, value, _ := strings.Cut(entry, "=")
		// GOCACHE is the one deliberate override; nothing else in this list is.
		if got[name] != value {
			t.Fatalf("child environment lost %s=%q (got %q); the toolchains need it", name, value, got[name])
		}
	}
	if got["GOCACHE"] != filepath.Join(os.TempDir(), "tx-review-gocache") {
		t.Fatalf("GOCACHE = %q, want the review-owned build cache", got["GOCACHE"])
	}
	// Duplicate keys make the child's getenv pick a winner we did not choose,
	// so the environment is rebuilt from a map rather than appended to.
	seen := map[string]int{}
	for _, entry := range child {
		name, _, _ := strings.Cut(entry, "=")
		seen[name]++
		if seen[name] > 1 {
			t.Fatalf("child environment sets %s twice:\n%s", name, strings.Join(child, "\n"))
		}
	}
}

// leaked, the fixture's test fails and reports the variable by name.
func TestGoStaticToolsRunWithoutHandingTheCheckoutTheReviewersCredentials(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skipf("no go on PATH: %v (this test is the empirical check that the scrub does not break the toolchain)", err)
	}
	root := t.TempDir()
	writeFile(t, root, "go.mod", "module example.com/fixture\n\ngo 1.21\n")
	writeFile(t, root, "app/app.go", "package app\n\n// Run exists so go vet has a package to type-check.\nfunc Run() string { return \"ok\" }\n")
	writeFile(t, root, "app/app_test.go", staticToolFixtureTest)

	for name, value := range reviewerCredentialEnv {
		t.Setenv(name, value)
	}
	t.Setenv("TOTALITY_REVIEW_STATIC_TOOLS", "1")
	// GOTOOLCHAIN=local keeps the fixture from trying to download a toolchain,
	// and doubles as the positive control: it is allowlisted and set here to a
	// known value, so the fixture reading it back proves the child got a
	// populated environment rather than an empty one that would make every
	// credential assertion above it vacuous.
	t.Setenv("GOTOOLCHAIN", "local")

	results := collectStaticToolResults(context.Background(), root, RepoFacts{DependencyFiles: []string{"go.mod"}}, Options{}, []string{"app/app.go"})
	if len(results) != 2 {
		t.Fatalf("results = %#v, want go test and go vet", results)
	}
	for _, result := range results {
		if result.Skipped {
			t.Fatalf("%s was skipped (%s); the scrubbed environment has to keep the toolchain usable", result.Name, result.Reason)
		}
		if result.ExitCode != 0 {
			t.Fatalf("%s exited %d after the environment scrub:\n%s", result.Command, result.ExitCode, result.Output)
		}
	}
	if results[0].Command != "go test ./app" || results[1].Command != "go vet ./app" {
		t.Fatalf("commands = %q / %q, want the scoped go test and go vet", results[0].Command, results[1].Command)
	}
	joined := results[0].Output + "\n" + results[1].Output
	for name, value := range reviewerCredentialEnv {
		if value != "" && strings.Contains(joined, value) {
			t.Fatalf("%s value surfaced in captured static tool output:\n%s", name, joined)
		}
	}
}

// The other half of a leak is the return path: runStaticTool captures the
// child's stdout and stderr into StaticToolResult.Output, which becomes part of
// the AI reviewer's prompt and of a published report. A checkout that printed
// the reviewer's environment would be exfiltrating it through tx itself.
//
// The fixture here fails on purpose, because that is the only way the output
// travels: `go test` without -v discards a passing test's stdout, so the
// captured output for a green package is the single "ok <pkg>" line. Measured,
// not assumed — the passing fixture above produces exactly that.
func TestFailingCheckoutTestsCarryNoReviewerCredentialsBackIntoTheReport(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skipf("no go on PATH: %v", err)
	}
	root := t.TempDir()
	writeFile(t, root, "go.mod", "module example.com/leaky\n\ngo 1.21\n")
	writeFile(t, root, "app/app.go", "package app\n\nfunc Run() string { return \"ok\" }\n")
	writeFile(t, root, "app/app_test.go", staticToolLeakFixtureTest)

	for name, value := range reviewerCredentialEnv {
		t.Setenv(name, value)
	}
	t.Setenv("TOTALITY_REVIEW_STATIC_TOOLS", "1")
	t.Setenv("GOTOOLCHAIN", "local")

	results := collectStaticToolResults(context.Background(), root, RepoFacts{DependencyFiles: []string{"go.mod"}}, Options{}, []string{"app/app.go"})
	if len(results) == 0 || results[0].Command != "go test ./app" {
		t.Fatalf("results = %#v, want go test", results)
	}
	output := results[0].Output
	// Guard first: without the dump in hand, the searches below would pass
	// against an empty string and prove nothing.
	if !strings.Contains(output, "CHILD-ENV PATH=") {
		t.Fatalf("the fixture's environment dump did not reach the captured output:\n%s", output)
	}
	for name, value := range reviewerCredentialEnv {
		if strings.Contains(output, "CHILD-ENV "+name+"=") {
			t.Fatalf("the reviewed checkout read %s and tx carried it back into the report:\n%s", name, output)
		}
		if value != "" && strings.Contains(output, value) {
			t.Fatalf("the value of %s surfaced in the captured report output:\n%s", name, output)
		}
	}
}

// staticToolFixtureTest is the checkout's own test: the same position a hostile
// pull request would be in, reading the environment it was handed. Written as a
// literal rather than a testdata file so the credential names it checks sit
// next to reviewerCredentialEnv above.
const staticToolFixtureTest = `package app

import (
	"os"
	"testing"
)

func TestWhatTheReviewedCheckoutCanSee(t *testing.T) {
	for _, name := range []string{
		"ANTHROPIC_API_KEY", "OPENAI_API_KEY", "TURBOPUFFER_API_KEY",
		"GITHUB_TOKEN", "GH_TOKEN", "TOTALITY_UPLOAD_TOKEN", "TOTALITY_API_URL",
		"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_SESSION_TOKEN",
		"GOOGLE_API_KEY", "GOOGLE_APPLICATION_CREDENTIALS",
		"NPM_TOKEN", "npm_config__authToken", "DATABASE_URL",
		"SSH_AUTH_SOCK", "SENTRY_DSN",
	} {
		if _, ok := os.LookupEnv(name); ok {
			t.Errorf("the reviewed checkout can read %s", name)
		}
	}
	// Positive control: without this, a child that got an empty environment
	// would pass every assertion above while proving nothing.
	if got := os.Getenv("GOTOOLCHAIN"); got != "local" {
		t.Errorf("GOTOOLCHAIN = %q, want local: the child environment is not populated", got)
	}
	if os.Getenv("PATH") == "" {
		t.Error("PATH is empty; the child cannot exec anything")
	}
	if os.Getenv("HOME") == "" {
		t.Error("HOME is empty; go test cannot resolve GOPATH or GOENV")
	}
}
`

// staticToolLeakFixtureTest is the shape a hostile pull request would take: an
// ordinary-looking test that prints everything it was handed and then fails, so
// that tx captures the dump and carries it into the review report.
const staticToolLeakFixtureTest = `package app

import (
	"fmt"
	"os"
	"sort"
	"testing"
)

func TestExfiltrateTheReviewersEnvironment(t *testing.T) {
	env := os.Environ()
	sort.Strings(env)
	for _, entry := range env {
		fmt.Println("CHILD-ENV", entry)
	}
	t.Fatal("failing on purpose so the dump above reaches the review report")
}
`

// The non-Go runners exec a binary the checkout can ship (node_modules/.bin),
// so the scrub has to cover them too. A shell stand-in is used deliberately:
// it dumps the whole environment, which is a stronger assertion than any real
// linter's exit code, and it runs on every machine.
func TestNonGoStaticToolsRunWithoutHandingTheCheckoutTheReviewersCredentials(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "eslint.config.js", "export default []\n")
	writeFile(t, root, "src/a.js", "export const a = 1\n")
	writeExecutable(t, filepath.Join(root, "node_modules", ".bin", "eslint"), "#!/bin/sh\nenv\n")

	for name, value := range reviewerCredentialEnv {
		t.Setenv(name, value)
	}
	t.Setenv("TOTALITY_REVIEW_STATIC_TOOLS", "1")

	results := collectStaticToolResults(context.Background(), root, RepoFacts{}, Options{}, []string{"src/a.js"})
	if len(results) != 1 || results[0].Command != "eslint ./src/a.js" {
		t.Fatalf("results = %#v, want the repo-local eslint", results)
	}
	for name := range reviewerCredentialEnv {
		if strings.Contains(results[0].Output, name+"=") {
			t.Fatalf("the repo-local eslint was handed %s:\n%s", name, results[0].Output)
		}
	}
	// It got a working environment, not an empty one.
	for _, want := range []string{"PATH=", "HOME="} {
		if !strings.Contains(results[0].Output, want) {
			t.Fatalf("the repo-local eslint environment is missing %s:\n%s", want, results[0].Output)
		}
	}
}

// The non-Go toolchains most likely to break under a scrubbed environment are
// the ones invoked through a shim that resolves itself from the environment:
// cargo is a rustup shim, ruff on a developer machine is usually a pyenv, uv
// or .venv shim. Reasoning about which variables those shims read is how you
// get an allowlist that looks right and is not, so each case here runs the
// real binary twice — once with the environment the reviewer has, once through
// the real static-tool path — and fails only if the scrub is the difference.
// A tool that is missing, or already broken in the reviewer's own environment
// (a pyenv global that does not have ruff installed, say), can say nothing
// about the scrub and is skipped with the reason it could not be used.
func TestRealToolchainsBehaveTheSameUnderTheScrubbedEnvironment(t *testing.T) {
	cases := []struct {
		tool    string
		files   map[string]string
		changed []string
		command string
	}{
		{
			tool: "cargo",
			files: map[string]string{
				"Cargo.toml": "[package]\nname = \"fixture\"\nversion = \"0.1.0\"\nedition = \"2021\"\n",
				"src/lib.rs": "pub fn run() -> u32 { 1 }\n",
			},
			changed: []string{"src/lib.rs"},
			command: "cargo check",
		},
		{
			tool: "ruff",
			// An unused import: ruff reports F401 with no configuration at all,
			// so the baseline exit code is a non-trivial 1 rather than 0.
			files:   map[string]string{"app/main.py": "import os\n\nx = 1\n"},
			changed: []string{"app/main.py"},
			command: "ruff check ./app/main.py",
		},
	}

	for _, test := range cases {
		t.Run(test.tool, func(t *testing.T) {
			bin, err := exec.LookPath(test.tool)
			if err != nil {
				t.Skipf("no %s on PATH: %v", test.tool, err)
			}
			root := t.TempDir()
			for rel, content := range test.files {
				writeFile(t, root, rel, content)
			}
			baseline := runWithInheritedEnv(t, root, bin, strings.Fields(test.command)[1:])

			for name, value := range reviewerCredentialEnv {
				t.Setenv(name, value)
			}
			t.Setenv("TOTALITY_REVIEW_STATIC_TOOLS", "1")
			results := collectStaticToolResults(context.Background(), root, RepoFacts{}, Options{}, test.changed)
			if len(results) != 1 || results[0].Command != test.command {
				t.Fatalf("results = %#v, want %q", results, test.command)
			}
			if results[0].Skipped {
				t.Fatalf("%s was skipped (%s) under the scrubbed environment", test.tool, results[0].Reason)
			}
			if results[0].ExitCode != baseline.exitCode {
				t.Fatalf("%s exited %d under the scrubbed environment but %d with the reviewer's own:\nscrubbed:\n%s\ninherited:\n%s",
					test.tool, results[0].ExitCode, baseline.exitCode, results[0].Output, baseline.output)
			}
			for name := range reviewerCredentialEnv {
				if strings.Contains(results[0].Output, name+"=") {
					t.Fatalf("%s output carries %s:\n%s", test.tool, name, results[0].Output)
				}
			}
		})
	}
}

type inheritedEnvRun struct {
	exitCode int
	output   string
}

// runWithInheritedEnv is the baseline half of the differential above: the same
// command the static tool stage would run, with the environment tx itself has.
// It runs before the credential fixtures are set so the baseline is the
// reviewer's real environment and not a doctored one.
func runWithInheritedEnv(t *testing.T, root, bin string, args []string) inheritedEnvRun {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Dir = root
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	run := inheritedEnvRun{output: string(out)}
	if err != nil {
		run.exitCode = 1
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			run.exitCode = exitErr.ExitCode()
		} else {
			t.Skipf("%s is not usable in this environment even before the scrub: %v\n%s", bin, err, out)
		}
	}
	// 126 and 127 are "found the shim, could not run the tool" — a pyenv whose
	// global version has no ruff installed reports exactly that. Comparing two
	// runs that both failed to start anything would pass while proving nothing.
	if run.exitCode == 126 || run.exitCode == 127 {
		t.Skipf("%s does not resolve to a working tool in this environment (exit %d), so it cannot test the scrub:\n%s", bin, run.exitCode, out)
	}
	return run
}
