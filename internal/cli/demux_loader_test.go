package cli

import (
	"testing"
)

func TestDemuxLoaderProgressWriterSendsNonEmptyLines(t *testing.T) {
	phases := make(chan string, 2)
	writer := demuxLoaderProgressWriter{phases: phases}

	n, err := writer.Write([]byte("Planning compose proposal...\n\nFixing compose proposal demux-1...\n"))
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if n == 0 {
		t.Fatal("Write() wrote 0 bytes")
	}

	for _, want := range []string{"Planning compose proposal...", "Fixing compose proposal demux-1..."} {
		select {
		case got := <-phases:
			if got != want {
				t.Fatalf("phase = %q, want %q", got, want)
			}
		default:
			t.Fatalf("missing phase %q", want)
		}
	}
}
