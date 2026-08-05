package cloud

import "strings"

// CloudBaseURL returns the lgtm-cloud API origin.
func CloudBaseURL() string {
	return strings.TrimRight(strings.TrimSpace(CloudURL()), "/")
}

func CloudURLWithPath(path string) string {
	return cloudURLWithPath(CloudURL(), path)
}

func cloudURLWithPath(base, path string) string {
	base = strings.TrimRight(strings.TrimSpace(base), "/")
	if base == "" {
		return ""
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return base + path
}
