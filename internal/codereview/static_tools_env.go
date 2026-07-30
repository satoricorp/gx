package codereview

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// The static tool stage execs the reviewed repository's own toolchain, and one
// of those tools is `go test`: it compiles and runs the checkout's test
// binaries. That is code from the change under review executing on the
// reviewer's machine, which means the child environment is a disclosure
// channel, not a convenience.
//
// It used to be `append(os.Environ(), ...)`, so the child got everything the
// reviewer's shell had. On a Totality developer machine that is ANTHROPIC_API_KEY,
// OPENAI_API_KEY, TURBOPUFFER_API_KEY, GITHUB_TOKEN, TOTALITY_UPLOAD_TOKEN and the
// AWS_* triple; in CI it is whatever the workflow exported into the step that
// runs `tx review`. A TestMain that prints os.Environ() exfiltrates all of it,
// and the review output would even carry it back — runStaticTool captures the
// child's stdout and stderr into the report.
//
// staticToolChildEnv builds the child environment from an allowlist instead of
// filtering os.Environ() for names that look secret. A denylist has to predict
// every convention a credential might be named under (_KEY, _TOKEN, _SECRET,
// _PASSWORD, _CREDENTIALS, DSN, DATABASE_URL, SENTRY_DSN, npm_config__authToken,
// ...) and it is wrong the first time someone invents a new one; the allowlist
// is wrong only when a toolchain needs something we did not list, which shows
// up as a tool failure in the review output rather than as a silent leak.
//
// What this does NOT protect: HOME is on the list, because the Go toolchain
// derives GOPATH, GOMODCACHE and GOENV from it, and HOME leads to
// ~/.aws/credentials, ~/.netrc, ~/.config/gh/hosts.yml and every other on-disk
// credential. Scrubbing the environment removes the hand-it-over-on-a-plate
// path; it does not make executing a stranger's test suite safe. That larger
// question is a policy decision, not a filter.
//
// Entries are added by measurement, not by guesswork:
// TestRealToolchainsBehaveTheSameUnderTheScrubbedEnvironment runs the real
// binaries twice, once with the reviewer's environment and once through this
// filter, and fails when the scrub is the difference. Dropping PYENV_VERSION
// from the Python group, for instance, makes a pyenv-shimmed `ruff` exit 127
// instead of reporting the lint it found — that is how it got on the list.
//
// The allowlist is matched case-insensitively. Windows environment blocks are
// case-insensitive at the OS level (`Path`, `SystemRoot`, `ProgramFiles(x86)`
// all appear in mixed case), and on Unix admitting an unusual spelling of a
// name that is already allowed costs nothing.
var staticToolEnvAllowed = buildStaticToolEnvAllowlist(
	// Process basics. PATH is how the child finds its own subprocesses — `go`
	// shells out to `git` for module fetches and to `cc` for cgo, cargo shells
	// out to `rustc`. HOME is where every toolchain keeps its caches and
	// per-user config ($HOME/.config/go/env, ~/.npmrc, ~/.cargo/config.toml).
	"PATH", "HOME", "USER", "LOGNAME",
	"TMPDIR", "TEMP", "TMP",
	// Locale and timezone: they change tool output encoding and any test that
	// formats a date. Omitting them does not break tools, it makes their
	// output differ from what the same command prints in the user's shell.
	"LANG", "LANGUAGE", "LC_ALL", "TZ", "TERM",
	// XDG paths relocate the caches the tools above write to. A reviewer who
	// redirects XDG_CACHE_HOME expects tx to honour it like everything else.
	"XDG_CACHE_HOME", "XDG_CONFIG_HOME", "XDG_DATA_HOME", "XDG_STATE_HOME", "XDG_RUNTIME_DIR",
	// Corporate networks put every toolchain's fetches through a proxy and a
	// private CA. Without these, `go test` on a module that is not already in
	// the module cache fails behind any egress proxy. Note the tradeoff: a
	// proxy URL may embed credentials (http://user:pass@proxy.corp), so this
	// is the one group on the list that can carry a secret. It stays because
	// dropping it breaks the tools outright on the networks that use it, and
	// the proxy credential is shared infrastructure rather than the reviewer's
	// personal API keys.
	"HTTP_PROXY", "HTTPS_PROXY", "NO_PROXY", "ALL_PROXY", "FTP_PROXY",
	"SSL_CERT_FILE", "SSL_CERT_DIR", "CURL_CA_BUNDLE", "REQUESTS_CA_BUNDLE",
	// Many linters and test harnesses (eslint, cargo, pytest) change their
	// output when CI is set. tx review is a CI gate; let them know.
	"CI",
	// Windows. Environment variables are how Windows locates the system at
	// all: Go's own runtime needs SystemRoot for crypto/rand, and exec needs
	// ComSpec and PATHEXT to resolve a command name to a .exe or .bat.
	"SystemRoot", "SystemDrive", "ComSpec", "PATHEXT", "windir", "OS",
	"NUMBER_OF_PROCESSORS", "PROCESSOR_ARCHITECTURE", "PROCESSOR_ARCHITEW6432",
	"USERPROFILE", "USERNAME", "HOMEDRIVE", "HOMEPATH", "PUBLIC",
	"APPDATA", "LOCALAPPDATA", "ALLUSERSPROFILE", "ProgramData",
	"ProgramFiles", "ProgramFiles(x86)", "ProgramW6432",
	"CommonProgramFiles", "CommonProgramFiles(x86)", "CommonProgramW6432",
	// macOS toolchain location. cgo and cargo both invoke clang through
	// xcrun, which reads these.
	"DEVELOPER_DIR", "SDKROOT", "MACOSX_DEPLOYMENT_TARGET",
	// Go. Named one at a time on purpose: a "GO*" prefix rule would also match
	// GOOGLE_API_KEY and GOOGLE_APPLICATION_CREDENTIALS, which is exactly the
	// class of variable this list exists to keep out.
	"GOROOT", "GOPATH", "GOBIN", "GOCACHE", "GOMODCACHE", "GOTMPDIR", "GOENV",
	"GOFLAGS", "GOEXPERIMENT", "GOTOOLCHAIN", "GOWORK", "GO111MODULE",
	"GOOS", "GOARCH", "GOHOSTOS", "GOHOSTARCH",
	"GOARM", "GOARM64", "GOAMD64", "GO386", "GOMIPS", "GOMIPS64", "GOPPC64", "GORISCV64", "GOWASM",
	"GODEBUG", "GOMAXPROCS", "GOGC", "GOMEMLIMIT", "GOTRACEBACK",
	"GOPROXY", "GONOPROXY", "GOPRIVATE", "GOSUMDB", "GONOSUMDB", "GONOSUMCHECK",
	"GOINSECURE", "GOVCS", "GOFIPS140",
	// GOAUTH is a policy knob, not a credential: it names *how* the go command
	// authenticates to private module servers. Dropping it would not remove
	// the reviewer's credentials — the default is still netrc — it would only
	// discard a reviewer who set GOAUTH=off, which is the restrictive setting.
	"GOAUTH",
	// C toolchain, for cgo and for anything cargo links against.
	"CC", "CXX", "AR", "FC", "CGO_ENABLED",
	"CGO_CFLAGS", "CGO_CPPFLAGS", "CGO_CXXFLAGS", "CGO_FFLAGS", "CGO_LDFLAGS",
	"PKG_CONFIG", "PKG_CONFIG_PATH", "PKG_CONFIG_LIBDIR",
	// Node. npm_config_* is deliberately absent: npm materializes registry
	// auth into npm_config__authToken and npm_config__auth when it runs a
	// script, so a tx invoked from an npm script would leak the registry token
	// under a prefix rule.
	"NODE_PATH", "NODE_OPTIONS", "NODE_ENV", "NODE_EXTRA_CA_CERTS",
	"NVM_DIR", "NVM_BIN",
	// Python. VIRTUAL_ENV is how ruff and mypy find the interpreter whose site
	// packages they are meant to resolve imports against.
	"VIRTUAL_ENV", "PYTHONPATH", "PYTHONHOME", "PYTHONDONTWRITEBYTECODE",
	"PYTHONUNBUFFERED", "PYENV_ROOT", "PYENV_VERSION",
	"CONDA_PREFIX", "CONDA_DEFAULT_ENV",
	"MYPY_CACHE_DIR", "RUFF_CACHE_DIR",
	// Rust.
	"CARGO_HOME", "RUSTUP_HOME", "RUSTUP_TOOLCHAIN", "RUSTC", "RUSTC_WRAPPER",
	"RUSTFLAGS", "CARGO_TARGET_DIR", "CARGO_BUILD_TARGET", "CARGO_NET_OFFLINE",
	"RUST_BACKTRACE",
)

// staticToolEnvAllowedPrefixes covers the one family whose members cannot be
// enumerated: LC_CTYPE, LC_NUMERIC, LC_TIME and the rest of POSIX locale
// categories, which vary by platform.
var staticToolEnvAllowedPrefixes = []string{"LC_"}

func buildStaticToolEnvAllowlist(names ...string) map[string]struct{} {
	out := make(map[string]struct{}, len(names))
	for _, name := range names {
		out[strings.ToUpper(name)] = struct{}{}
	}
	return out
}

// staticToolChildEnv is the environment every static tool child runs with. It
// is built from scratch rather than filtered in place so that the result is
// exactly the allowlist plus the overrides below, with no duplicate keys for
// the child's getenv to disambiguate.
func staticToolChildEnv(parent []string) []string {
	values := map[string]string{}
	for _, entry := range parent {
		name, value, ok := strings.Cut(entry, "=")
		if !ok || name == "" {
			continue
		}
		upper := strings.ToUpper(name)
		if _, ok := staticToolEnvAllowed[upper]; !ok && !hasStaticToolEnvAllowedPrefix(upper) {
			continue
		}
		values[name] = value
	}
	// GOCACHE is forced rather than inherited: the reviewed repo's build
	// artifacts must not land in the reviewer's own build cache, and a review
	// of someone else's checkout should not evict the entries the reviewer's
	// day job depends on.
	values["GOCACHE"] = filepath.Join(os.TempDir(), "tx-review-gocache")
	out := make([]string, 0, len(values))
	for name, value := range values {
		out = append(out, name+"="+value)
	}
	sort.Strings(out)
	return out
}

func hasStaticToolEnvAllowedPrefix(upperName string) bool {
	for _, prefix := range staticToolEnvAllowedPrefixes {
		if strings.HasPrefix(upperName, prefix) {
			return true
		}
	}
	return false
}
