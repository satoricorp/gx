package daemon

import (
	"net/http"
	"net/url"
	"testing"
)

func TestIsChatGPTAuth(t *testing.T) {
	req := httptestRequestWithHeader("Chatgpt-Account-Id", "acct-123")
	if !isChatGPTAuth(req) {
		t.Fatal("expected ChatGPT auth when account id header is present")
	}
	if isChatGPTAuth(httptestRequestWithHeader("Authorization", "Bearer sk-test")) {
		t.Fatal("expected platform auth without account id header")
	}
}

func TestResolveUpstreamOpenAIPlatform(t *testing.T) {
	req := httptestRequestWithHeader("Authorization", "Bearer sk-test")
	target, ok := resolveUpstream(req, "openai")
	if !ok {
		t.Fatal("expected upstream target")
	}
	if target.host != "api.openai.com" || target.path != "/v1/responses" {
		t.Fatalf("unexpected target: %+v", target)
	}
}

func TestResolveUpstreamChatGPT(t *testing.T) {
	req := httptestRequestWithPath("/v1/responses")
	req.Header.Set("Authorization", "Bearer chatgpt-token")
	req.Header.Set("Chatgpt-Account-Id", "acct-123")

	target, ok := resolveUpstream(req, "openai")
	if !ok {
		t.Fatal("expected upstream target")
	}
	if target.host != chatGPTHost || target.path != "/backend-api/codex/responses" {
		t.Fatalf("unexpected target: %+v", target)
	}
}

func TestChatGPTUpstreamPathMappings(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"/v1/responses", "/backend-api/codex/responses"},
		{"/responses", "/backend-api/codex/responses"},
		{"/v1/models", "/backend-api/codex/models"},
		{"/v1/chat/completions", "/backend-api/codex/chat/completions"},
	}
	for _, tc := range tests {
		if got := chatGPTUpstreamPath(tc.in); got != tc.want {
			t.Fatalf("chatGPTUpstreamPath(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestUpstreamTargetURLPreservesQuery(t *testing.T) {
	target := upstreamTarget{host: chatGPTHost, path: "/backend-api/codex/models"}
	got := target.URL("https", "client_version=0.135.0")
	want := "https://chatgpt.com/backend-api/codex/models?client_version=0.135.0"
	if got != want {
		t.Fatalf("URL() = %q, want %q", got, want)
	}
}

func httptestRequestWithHeader(key, value string) *http.Request {
	req := httptestRequestWithPath("/v1/responses")
	req.Header.Set(key, value)
	return req
}

func httptestRequestWithPath(path string) *http.Request {
	return &http.Request{
		Header: make(http.Header),
		URL:    mustParseURL("http://127.0.0.1:43123" + path),
	}
}

func mustParseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}
