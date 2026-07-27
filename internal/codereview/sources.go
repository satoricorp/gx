package codereview

import (
	"fmt"
	"strings"
)

type ResolvedSource struct {
	ID        string `json:"id,omitempty"`
	Kind      string `json:"kind,omitempty"`
	Publisher string `json:"publisher,omitempty"`
	Opaque    bool   `json:"opaque,omitempty"`
	Ref       string `json:"ref,omitempty"`
	File      string `json:"file,omitempty"`
	StartLine int    `json:"start_line,omitempty"`
	EndLine   int    `json:"end_line,omitempty"`
	URL       string `json:"url,omitempty"`
}

func ResolvedSourceLabel(src ResolvedSource) string {
	if src.Opaque {
		if publisher := strings.TrimSpace(src.Publisher); publisher != "" {
			return publisher
		}
		if kind := strings.TrimSpace(src.Kind); kind != "" {
			return kind
		}
		return "review resource"
	}
	if file := strings.TrimSpace(src.File); file != "" {
		if src.StartLine > 0 {
			return fmt.Sprintf("%s:%d", file, src.StartLine)
		}
		return file
	}
	if ref := strings.TrimSpace(src.Ref); ref != "" {
		return ref
	}
	if publisher := strings.TrimSpace(src.Publisher); publisher != "" {
		return publisher
	}
	if kind := strings.TrimSpace(src.Kind); kind != "" {
		return kind
	}
	return ""
}

func resolveSourceLabels(brief ReviewBrief, labels []string) []ResolvedSource {
	refByID := map[string]SourceRef{}
	for _, ref := range brief.SourceRefs {
		id := strings.TrimSpace(ref.ID)
		if id != "" {
			refByID[id] = ref
		}
	}
	catalogByID := map[string]SourceBrief{}
	for _, entry := range brief.SourceCatalog {
		id := strings.TrimSpace(entry.ID)
		if id != "" {
			catalogByID[id] = entry
		}
	}
	staticByID := map[string]Source{}
	for _, source := range sourceCatalog {
		staticByID[source.ID] = source
	}

	seen := map[string]struct{}{}
	var out []ResolvedSource
	for _, label := range labels {
		label = strings.TrimSpace(label)
		if label == "" {
			continue
		}
		if _, ok := seen[label]; ok {
			continue
		}
		if ref, ok := refByID[label]; ok {
			if resolved := resolvedSourceFromRef(ref); resolved.ID != "" || resolved.Kind != "" {
				seen[label] = struct{}{}
				out = append(out, resolved)
			}
			continue
		}
		if entry, ok := catalogByID[label]; ok {
			resolved := resolvedSourceFromCatalog(entry, staticByID[entry.ID])
			if resolved.ID != "" {
				seen[label] = struct{}{}
				out = append(out, resolved)
			}
			continue
		}
		if source, ok := staticByID[label]; ok {
			seen[label] = struct{}{}
			out = append(out, resolvedSourceFromCatalog(SourceBrief{
				ID:        source.ID,
				Publisher: source.Publisher,
				Scopes:    source.Scopes,
			}, source))
		}
	}
	return out
}

func resolvedSourceFromRef(ref SourceRef) ResolvedSource {
	kind := strings.TrimSpace(ref.Kind)
	opaque := sourceRefOpaque(kind)
	publisher := strings.TrimSpace(ref.Publisher)
	if publisher == "" {
		publisher = publisherForSourceRef(ref)
	}
	title := strings.TrimSpace(ref.Title)
	if title == "" && ref.File != "" && ref.StartLine > 0 {
		title = fmt.Sprintf("%s:%d", ref.File, ref.StartLine)
	}
	if title == "" {
		title = strings.TrimSpace(ref.ID)
	}
	return ResolvedSource{
		ID:        strings.TrimSpace(ref.ID),
		Kind:      kind,
		Publisher: publisher,
		Opaque:    opaque,
		Ref:       title,
		File:      strings.TrimSpace(ref.File),
		StartLine: ref.StartLine,
		EndLine:   ref.EndLine,
		URL:       strings.TrimSpace(ref.URL),
	}
}

func resolvedSourceFromCatalog(entry SourceBrief, source Source) ResolvedSource {
	publisher := strings.TrimSpace(entry.Publisher)
	if publisher == "" {
		publisher = strings.TrimSpace(source.Publisher)
	}
	if publisher == "" {
		publisher = publisherFromURLHost(source.URL)
	}
	return ResolvedSource{
		ID:        strings.TrimSpace(entry.ID),
		Kind:      "catalog",
		Publisher: publisher,
		Opaque:    true,
		Ref:       strings.TrimSpace(source.Title),
		URL:       strings.TrimSpace(source.URL),
	}
}

func sourceRefOpaque(kind string) bool {
	switch strings.TrimSpace(kind) {
	case "resource", "reference", "catalog":
		return true
	case "local", "code", "session", "policy":
		return false
	default:
		return strings.Contains(kind, "resource") || strings.Contains(kind, "reference")
	}
}

func publisherForSourceRef(ref SourceRef) string {
	if publisher := strings.TrimSpace(ref.Publisher); publisher != "" {
		return publisher
	}
	switch strings.TrimSpace(ref.Kind) {
	case "local", "policy":
		return "local"
	case "session":
		return "session"
	case "code", "indexed_code":
		return "indexed"
	default:
		if publisher := publisherFromURLHost(ref.URL); publisher != "" {
			return publisher
		}
		if strings.EqualFold(strings.TrimSpace(ref.Source), "local") {
			return "local"
		}
		if strings.HasPrefix(strings.TrimSpace(ref.Source), "turbopuffer:") {
			return "indexed"
		}
		return "unknown"
	}
}
