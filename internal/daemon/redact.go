package daemon

import (
	"encoding/json"
	"net/http"
	"strings"
)

var redactedHeaders = map[string]struct{}{
	"authorization":       {},
	"x-api-key":           {},
	"anthropic-api-key":   {},
	"openai-api-key":      {},
	"openai-organization": {},
	"cookie":              {},
	"set-cookie":          {},
}

func redactHeaders(headers http.Header) string {
	safe := make(map[string]any, len(headers))
	for key, values := range headers {
		lower := strings.ToLower(key)
		if _, ok := redactedHeaders[lower]; ok {
			safe[key] = "[redacted]"
			continue
		}
		if len(values) == 1 {
			safe[key] = values[0]
			continue
		}
		copyValues := make([]string, len(values))
		copy(copyValues, values)
		safe[key] = copyValues
	}
	data, _ := json.Marshal(safe)
	return string(data)
}
