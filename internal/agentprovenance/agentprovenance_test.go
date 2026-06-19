package agentprovenance

import "testing"

func TestResolveInfersToolAndPreservesModel(t *testing.T) {
	source := "ambient"
	process := "codex"
	got := Resolve(Source{
		SessionID:   "session-one",
		Command:     "codex exec",
		Source:      &source,
		ProcessName: &process,
		Provider:    " openai ",
		ModelID:     " gpt-5.1-code ",
		CreatedAt:   12,
	})
	if got.SessionID != "session-one" || got.AgentTool != "codex" || got.Provider != "openai" || got.ModelID != "gpt-5.1-code" || got.CreatedAt != 12 {
		t.Fatalf("Resolve() = %#v", got)
	}
	if got.Source == nil || *got.Source != source || got.ProcessName == nil || *got.ProcessName != process {
		t.Fatalf("Resolve() source/process = %#v/%#v", got.Source, got.ProcessName)
	}
}

func TestInferToolPriority(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		source      string
		processName string
		want        string
	}{
		{name: "cursor source", source: "cursor", command: "codex", want: "cursor"},
		{name: "claude command", command: "claude code", want: "claude"},
		{name: "codex process", processName: "codex", want: "codex"},
		{name: "custom source", source: "aura", want: "aura"},
		{name: "unknown", want: "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InferTool(tt.command, tt.source, tt.processName); got != tt.want {
				t.Fatalf("InferTool() = %q, want %q", got, tt.want)
			}
		})
	}
}
