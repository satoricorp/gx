package providers

func asInt(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	default:
		return 0
	}
}

func optInt(v any) *int {
	switch n := v.(type) {
	case float64:
		i := int(n)
		return &i
	case int:
		return &n
	case int64:
		i := int(n)
		return &i
	default:
		return nil
	}
}

func firstInt(values ...any) *int {
	for _, value := range values {
		if out := optInt(value); out != nil {
			return out
		}
	}
	return nil
}

func asString(values ...any) string {
	for _, v := range values {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func fieldForDelta(kind any) string {
	if asString(kind) == "thinking_delta" {
		return "thinking"
	}
	return "text"
}

func appendString(target map[string]any, field, fragment string) {
	if fragment == "" {
		return
	}
	current, _ := target[field].(string)
	target[field] = current + fragment
}

func cloneMap(src map[string]any) map[string]any {
	out := make(map[string]any, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}

func mergeUsage(dst *Usage, src Usage) {
	if src.InputTokens != nil {
		dst.InputTokens = src.InputTokens
	}
	if src.OutputTokens != nil {
		dst.OutputTokens = src.OutputTokens
	}
	if src.CacheReadTokens != nil {
		dst.CacheReadTokens = src.CacheReadTokens
	}
	if src.CacheWriteTokens != nil {
		dst.CacheWriteTokens = src.CacheWriteTokens
	}
}

func derefString(v *string) any {
	if v == nil {
		return nil
	}
	return *v
}
