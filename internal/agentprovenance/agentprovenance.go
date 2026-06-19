package agentprovenance

import "strings"

type Source struct {
	SessionID   string
	Command     string
	Source      *string
	ProcessName *string
	Provider    string
	ModelID     string
	CreatedAt   int64
}

type Record struct {
	SessionID   string
	AgentTool   string
	Provider    string
	ModelID     string
	Source      *string
	ProcessName *string
	CreatedAt   int64
}

func Resolve(source Source) Record {
	return Record{
		SessionID:   source.SessionID,
		AgentTool:   InferTool(source.Command, stringPtrValue(source.Source), stringPtrValue(source.ProcessName)),
		Provider:    strings.TrimSpace(source.Provider),
		ModelID:     strings.TrimSpace(source.ModelID),
		Source:      source.Source,
		ProcessName: source.ProcessName,
		CreatedAt:   source.CreatedAt,
	}
}

func InferTool(command, source, processName string) string {
	source = strings.TrimSpace(strings.ToLower(source))
	command = strings.ToLower(command)
	processName = strings.ToLower(processName)
	switch {
	case source == "cursor" || strings.Contains(command, "cursor") || strings.Contains(processName, "cursor"):
		return "cursor"
	case strings.Contains(command, "claude") || strings.Contains(processName, "claude"):
		return "claude"
	case strings.Contains(command, "codex") || strings.Contains(processName, "codex"):
		return "codex"
	case source != "":
		return source
	default:
		return "unknown"
	}
}

func stringPtrValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
