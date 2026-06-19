package capture

import (
	"encoding/json"
	"sort"
	"time"
)

// FieldRecord describes one observed JSON path in session logs.
type FieldRecord struct {
	JSONPath            string `json:"json_path"`
	Occurrences         int    `json:"occurrences"`
	InferredType        string `json:"inferred_type"`
	MapsToSessionEvents string `json:"maps_to_session_events"`
	Notes               string `json:"notes,omitempty"`
}

// ToolInventory aggregates field observations for one agent tool.
type ToolInventory struct {
	SampleFiles []string      `json:"sample_files"`
	Fields      []FieldRecord `json:"fields"`
}

// FieldInventory is the WP-4 input artifact for session_events columns.
type FieldInventory struct {
	SchemaVersion              int                      `json:"schema_version"`
	GeneratedAt                string                   `json:"generated_at"`
	Tools                      map[string]ToolInventory `json:"tools"`
	ProposedSessionEventsCols  []string                 `json:"proposed_session_events_columns"`
	OverflowFields             []string                 `json:"overflow_fields"`
}

// InventoryCollector records JSON field paths while parsing session files.
type InventoryCollector struct {
	tools       map[string]*toolCollector
	sampleFiles map[string]map[string]struct{}
}

type toolCollector struct {
	fields map[string]*FieldRecord
}

// NewInventoryCollector creates an empty field inventory collector.
func NewInventoryCollector() *InventoryCollector {
	return &InventoryCollector{
		tools:       map[string]*toolCollector{},
		sampleFiles: map[string]map[string]struct{}{},
	}
}

// RecordSampleFile notes a parsed session file for a tool.
func (c *InventoryCollector) RecordSampleFile(tool, path string) {
	if c.sampleFiles[tool] == nil {
		c.sampleFiles[tool] = map[string]struct{}{}
	}
	c.sampleFiles[tool][path] = struct{}{}
}

// RecordPath increments occurrence count for a JSON path under a tool.
func (c *InventoryCollector) RecordPath(tool, jsonPath, inferredType, mapsTo, notes string) {
	tc := c.tools[tool]
	if tc == nil {
		tc = &toolCollector{fields: map[string]*FieldRecord{}}
		c.tools[tool] = tc
	}
	rec, ok := tc.fields[jsonPath]
	if !ok {
		rec = &FieldRecord{
			JSONPath:            jsonPath,
			InferredType:        inferredType,
			MapsToSessionEvents: mapsTo,
			Notes:               notes,
		}
		tc.fields[jsonPath] = rec
	}
	rec.Occurrences++
}

// Build returns the finalized field inventory document.
func (c *InventoryCollector) Build() FieldInventory {
	tools := map[string]ToolInventory{}
	for tool, tc := range c.tools {
		fields := make([]FieldRecord, 0, len(tc.fields))
		for _, rec := range tc.fields {
			fields = append(fields, *rec)
		}
		sort.Slice(fields, func(i, j int) bool {
			if fields[i].Occurrences != fields[j].Occurrences {
				return fields[i].Occurrences > fields[j].Occurrences
			}
			return fields[i].JSONPath < fields[j].JSONPath
		})
		samples := make([]string, 0)
		for path := range c.sampleFiles[tool] {
			samples = append(samples, path)
		}
		sort.Strings(samples)
		tools[tool] = ToolInventory{
			SampleFiles: samples,
			Fields:      fields,
		}
	}
	return FieldInventory{
		SchemaVersion: 1,
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
		Tools:         tools,
		ProposedSessionEventsCols: []string{
			"session_id", "tool", "model", "ts", "kind",
			"file_path", "old_text", "new_text", "prompt_context",
			"raw_json",
		},
		OverflowFields: []string{
			"usage", "requestId", "parentUuid", "uuid", "promptId",
			"permissionMode", "gitBranch", "version", "entrypoint",
			"structuredPatch", "toolUseResult", "attachment", "payload.git",
		},
	}
}

// MarshalJSON writes the inventory to JSON bytes.
func (inv FieldInventory) MarshalJSON() ([]byte, error) {
	type alias FieldInventory
	return json.MarshalIndent(alias(inv), "", "  ")
}
