package daemon

import (
	"net/http"
	"strings"
)

const chatGPTHost = "chatgpt.com"

type upstreamTarget struct {
	host string
	path string
}

func (t upstreamTarget) URL(scheme, rawQuery string) string {
	url := scheme + "://" + t.host + t.path
	if rawQuery != "" {
		url += "?" + rawQuery
	}
	return url
}

func isChatGPTAuth(r *http.Request) bool {
	return strings.TrimSpace(r.Header.Get("Chatgpt-Account-Id")) != ""
}

func resolveUpstream(r *http.Request, provider string) (upstreamTarget, bool) {
	switch provider {
	case "anthropic":
		return upstreamTarget{host: "api.anthropic.com", path: r.URL.Path}, true
	case "openai":
		if isChatGPTAuth(r) {
			return upstreamTarget{host: chatGPTHost, path: chatGPTUpstreamPath(r.URL.Path)}, true
		}
		return upstreamTarget{host: "api.openai.com", path: upstreamPath(provider, r.URL.Path)}, true
	default:
		return upstreamTarget{}, false
	}
}

func chatGPTUpstreamPath(path string) string {
	normalized := upstreamPath("openai", path)
	switch {
	case strings.HasPrefix(normalized, "/v1/responses"):
		return "/backend-api/codex/responses" + strings.TrimPrefix(normalized, "/v1/responses")
	case strings.HasPrefix(normalized, "/v1/models"):
		return "/backend-api/codex/models" + strings.TrimPrefix(normalized, "/v1/models")
	case strings.HasPrefix(normalized, "/v1/chat/completions"):
		return "/backend-api/codex/chat/completions" + strings.TrimPrefix(normalized, "/v1/chat/completions")
	case strings.HasPrefix(normalized, "/v1/"):
		return "/backend-api/codex/" + strings.TrimPrefix(normalized, "/v1/")
	default:
		return normalized
	}
}
