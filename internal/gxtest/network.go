package gxtest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
)

// AllowNetworkEnv opts a test binary out of DenyNetwork, for the `_live_test.go`
// suites that exist precisely to talk to the real thing.
//
// Those tests skip themselves when credentials are absent, which reads as a
// guard but is the opposite of one: it means they arm themselves automatically
// on the machine of anyone who has exported the keys, which is every developer
// working on GX. Requiring this variable makes reaching the network a decision
// someone typed rather than a property of their shell:
//
//	GX_TEST_ALLOW_NETWORK=1 go test ./internal/codereview -run TestLive
const AllowNetworkEnv = "GX_TEST_ALLOW_NETWORK"

// DenyNetwork cuts a test binary off from the internet and returns an accessor
// for whatever still tried to get out. It is meant to be called from TestMain,
// before any test runs, so that it covers tests that forget to opt in — which
// is the whole point, since the tests that leak are exactly the ones nobody
// remembered to sandbox.
//
// It closes two doors, because either one alone leaves a gap:
//
// Clearing credentials is what stops the traffic. Every network path in GX is
// gated on a key being present, so a cleared key makes the code take the same
// branch it takes in CI — which is also the branch these tests mean to be
// testing. This is the door that matters, and it is why the fix is not a list
// of per-feature kill switches: a retriever added tomorrow is off here for the
// same reason it is off in CI, with nobody having to remember it.
//
// Proxying is what makes a miss visible. A path that dials a hardcoded host
// despite having no key would otherwise be invisible — it costs money quietly
// and shows up only as a slow test. Pointing HTTP(S)_PROXY at a local recorder
// turns that into a named, reported failure. Go never proxies loopback, so
// httptest servers the tests stand up themselves are unaffected.
//
// The returned function reports the attempts, so TestMain can fail the package:
//
//	func TestMain(m *testing.M) {
//		egress := gxtest.DenyNetwork()
//		code := m.Run()
//		if attempts := egress(); len(attempts) > 0 {
//			// report and fail
//		}
//		os.Exit(code)
//	}
func DenyNetwork() func() []string {
	if truthyEnv(os.Getenv(AllowNetworkEnv)) {
		return func() []string { return nil }
	}

	var mu sync.Mutex
	seen := map[string]int{}
	tripwire := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CONNECT (a proxied HTTPS tunnel) names its target as host:port with
		// no scheme or path; a plain proxied request carries an absolute URL;
		// a direct one carries only a path, so the Host header names it.
		target := r.RequestURI
		switch {
		case r.Method == http.MethodConnect:
		case strings.Contains(target, "://"):
		default:
			target = r.Host + target
		}
		mu.Lock()
		seen[r.Method+" "+target]++
		mu.Unlock()
		http.Error(w, "network is denied in tests; see gxtest.DenyNetwork", http.StatusForbidden)
	}))

	for _, key := range NetworkCredentialEnv {
		_ = os.Setenv(key, "")
	}
	for _, key := range []string{
		"HTTP_PROXY", "http_proxy",
		"HTTPS_PROXY", "https_proxy",
		// GX's own endpoint overrides, so a path that honours them fails
		// against the tripwire rather than against a real service.
		"GX_OPENAI_BASE_URL", "GX_TPUF_BASE_URL", "GX_API_URL", "GX_CLOUD_URL",
		// GitHub is reached through two independent overrides, and setting only
		// the API one leaves token validation pointed at api.github.com.
		"GX_GITHUB_API_URL", "GX_GITHUB_USER_URL",
	} {
		_ = os.Setenv(key, tripwire.URL)
	}
	// Anything exempted from the proxy would be exempted from the tripwire.
	// Loopback stays exempt regardless, which is what keeps httptest working.
	for _, key := range []string{"NO_PROXY", "no_proxy"} {
		_ = os.Setenv(key, "")
	}

	return func() []string {
		mu.Lock()
		defer mu.Unlock()
		out := make([]string, 0, len(seen))
		for call, count := range seen {
			out = append(out, fmt.Sprintf("%s (%d attempt(s))", call, count))
		}
		sort.Strings(out)
		return out
	}
}

// NetworkCredentialEnv are the keys that arm a real network call. A test binary
// inherits the developer's shell, and on a GX developer's machine these are all
// exported, so leaving any of them set is what turns an offline test into a
// billable one.
var NetworkCredentialEnv = []string{
	"OPENAI_API_KEY",
	"GX_OPENAI_API_KEY",
	"TURBOPUFFER_API_KEY",
	"ANTHROPIC_API_KEY",
	"GX_TPUF_NAMESPACE",
	"GX_API_KEY",
	"GX_UPLOAD_TOKEN",
	"GITHUB_TOKEN",
	"GH_TOKEN",
	"AWS_ACCESS_KEY_ID",
	"AWS_SECRET_ACCESS_KEY",
	"AWS_SESSION_TOKEN",
	"AWS_PROFILE",
}

// FailOnEgress reports any attempt to leave the machine as a test failure.
// TestMain calls it with the accessor DenyNetwork returned.
func FailOnEgress(code int, attempts []string) int {
	if len(attempts) == 0 {
		return code
	}
	fmt.Fprintf(os.Stderr, "\ntests attempted %d outbound call(s); they must run offline:\n", len(attempts))
	for _, attempt := range attempts {
		fmt.Fprintf(os.Stderr, "\t%s\n", attempt)
	}
	fmt.Fprintf(os.Stderr, "if this is a live test, guard it and run it with %s=1\n", AllowNetworkEnv)
	if code == 0 {
		return 1
	}
	return code
}

// RequireNoNetwork skips a live test unless network access was asked for
// explicitly. It is the guard the `_live_test.go` suites need in addition to
// their credential checks: without it, exporting a key is enough to arm them.
func RequireNoNetwork(t *testing.T) {
	t.Helper()
	if !truthyEnv(os.Getenv(AllowNetworkEnv)) {
		t.Skipf("%s is not set; skipping live test that reaches the network", AllowNetworkEnv)
	}
}

func truthyEnv(raw string) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "t", "yes", "y", "on":
		return true
	default:
		return false
	}
}
