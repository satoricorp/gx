package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/authoring"
)

func TestPromptModifySelection(t *testing.T) {
	candidates := []authoring.ChangeInfo{
		{ChangeID: "abc123change", CommitID: "commit1111", Description: "first change"},
		{ChangeID: "def456change", CommitID: "commit2222", Description: "second change"},
	}

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr string
	}{
		{name: "default first", input: "\n", want: "abc123change"},
		{name: "numeric selection", input: "2\n", want: "def456change"},
		{name: "prefix selection", input: "def\n", want: "def456change"},
		{name: "quit", input: "q\n", wantErr: context.Canceled.Error()},
		{name: "out of range", input: "9\n", wantErr: "out of range"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			got, err := promptModifySelection(strings.NewReader(tt.input), &out, candidates)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("promptModifySelection() error = %v, want substring %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("promptModifySelection() unexpected error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("promptModifySelection() = %q, want %q", got, tt.want)
			}
		})
	}
}
