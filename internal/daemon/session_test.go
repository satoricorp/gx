package daemon

import "testing"

func TestMemorySessionRegistry(t *testing.T) {
	registry := NewMemorySessionRegistry()
	if err := registry.Register("session-1", 1234); err != nil {
		t.Fatal(err)
	}

	sessionID, ok := registry.Resolve(1234)
	if !ok || sessionID != "session-1" {
		t.Fatalf("expected session-1, got %q ok=%v", sessionID, ok)
	}

	if err := registry.Release(1234); err != nil {
		t.Fatal(err)
	}
	if _, ok := registry.Resolve(1234); ok {
		t.Fatal("expected port to be released")
	}
}
