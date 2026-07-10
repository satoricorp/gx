package vcs

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseRevisionIDsFromMessage(t *testing.T) {
	id := "klpxswmqonvr"
	message := "add feature\n\n" + RevisionTrailerLine(id)
	got := ParseRevisionIDsFromMessage(message)
	if len(got) != 1 || got[0] != id {
		t.Fatalf("ParseRevisionIDsFromMessage() = %v, want [%s]", got, id)
	}
	dup := message + "\n\n" + RevisionTrailerLine(id)
	got = ParseRevisionIDsFromMessage(dup)
	if len(got) != 1 {
		t.Fatalf("ParseRevisionIDsFromMessage() duplicated = %v, want one id", got)
	}
}

func TestRevisionIDsInGitRange(t *testing.T) {
	repo := t.TempDir()
	initTrailerTestRepo(t, repo)

	idA := "revaaaaaaaaaa"
	idB := "revbbbbbbbbbb"
	commitWithTrailer(t, repo, "first", idA)
	commitWithTrailer(t, repo, "second", idB)

	base, err := exec.Command("git", "-C", repo, "rev-parse", "HEAD~2").Output()
	if err != nil {
		t.Fatal(err)
	}
	head, err := exec.Command("git", "-C", repo, "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatal(err)
	}
	refRange := strings.TrimSpace(string(base)) + ".." + strings.TrimSpace(string(head))

	got, err := RevisionIDsInGitRange(context.Background(), repo, refRange)
	if err != nil {
		t.Fatalf("RevisionIDsInGitRange() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("RevisionIDsInGitRange() = %v, want two ids", got)
	}
	want := map[string]bool{idA: true, idB: true}
	for _, id := range got {
		if !want[id] {
			t.Fatalf("RevisionIDsInGitRange() = %v, want ids %s and %s", got, idA, idB)
		}
	}

	single, err := RevisionIDsInGitRange(context.Background(), repo, strings.TrimSpace(string(head))+"~1.."+strings.TrimSpace(string(head)))
	if err != nil {
		t.Fatalf("RevisionIDsInGitRange(single) error = %v", err)
	}
	if len(single) != 1 || single[0] != idB {
		t.Fatalf("RevisionIDsInGitRange(single) = %v, want [%s]", single, idB)
	}
}

func initTrailerTestRepo(t *testing.T, dir string) {
	t.Helper()
	for _, args := range [][]string{
		{"git", "init", "-b", "main"},
		{"git", "config", "user.email", "test@example.com"},
		{"git", "config", "user.name", "Test"},
	} {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%v: %s", err, out)
		}
	}
	readme := filepath.Join(dir, "README.md")
	if err := os.WriteFile(readme, []byte("# test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("git", "add", "README.md")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git add: %v: %s", err, out)
	}
	cmd = exec.Command("git", "commit", "-m", "initial")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git commit: %v: %s", err, out)
	}
}

func commitWithTrailer(t *testing.T, repo, subject, revisionID string) {
	t.Helper()
	file := filepath.Join(repo, revisionID+".txt")
	if err := os.WriteFile(file, []byte(revisionID+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("git", "add", filepath.Base(file))
	cmd.Dir = repo
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git add: %v: %s", err, out)
	}
	message := subject + "\n\n" + RevisionTrailerLine(revisionID)
	cmd = exec.Command("git", "commit", "-m", message)
	cmd.Dir = repo
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git commit: %v: %s", err, out)
	}
}
