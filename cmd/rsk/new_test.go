package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureGitignoreEntry_createsFileAndKeepsOneEntryPerLine(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	if err := ensureGitignoreEntry(dir, ".rsk/rsk.lock"); err != nil {
		t.Fatalf("ensureGitignoreEntry: %v", err)
	}
	got := readFile(t, filepath.Join(dir, ".gitignore"))
	if got != ".rsk/rsk.lock\n" {
		t.Fatalf("gitignore: got %q", got)
	}

	// Idempotent: a second run adds nothing.
	if err := ensureGitignoreEntry(dir, ".rsk/rsk.lock"); err != nil {
		t.Fatalf("ensureGitignoreEntry (again): %v", err)
	}
	if again := readFile(t, filepath.Join(dir, ".gitignore")); again != got {
		t.Fatalf("gitignore changed on second run: %q", again)
	}
}

func TestEnsureGitignoreEntry_preservesExistingContent(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")
	if err := os.WriteFile(path, []byte("node_modules/\n.claude/"), 0o600); err != nil {
		t.Fatalf("seed: %v", err)
	}

	if err := ensureGitignoreEntry(dir, ".rsk/rsk.lock"); err != nil {
		t.Fatalf("ensureGitignoreEntry: %v", err)
	}
	got := readFile(t, path)
	want := "node_modules/\n.claude/\n.rsk/rsk.lock\n"
	if got != want {
		t.Fatalf("gitignore: got %q want %q", got, want)
	}
}

func TestEnsureGitignoreEntry_skipsExistingEntryAmongOthers(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")
	seed := "node_modules/\n  .rsk/rsk.lock  \n.env\n"
	if err := os.WriteFile(path, []byte(seed), 0o600); err != nil {
		t.Fatalf("seed: %v", err)
	}

	if err := ensureGitignoreEntry(dir, ".rsk/rsk.lock"); err != nil {
		t.Fatalf("ensureGitignoreEntry: %v", err)
	}
	if got := readFile(t, path); got != seed {
		t.Fatalf("gitignore changed: %q", got)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
