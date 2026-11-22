package tool

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeSymlink creates target with content and a symlink at link pointing to it.
func writeSymlink(t *testing.T, target, link, content string) {
	t.Helper()
	if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
		t.Fatalf("seed target: %v", err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks unsupported: %v", err)
	}
}

func TestAppendClaudeImport_refusesSymlink(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	target := filepath.Join(dir, "AGENTS.md")
	link := filepath.Join(dir, "CLAUDE.md")
	writeSymlink(t, target, link, "my project notes\n")

	err := appendClaudeImport(link)
	if err == nil {
		t.Fatal("expected an error for a symlinked CLAUDE.md")
	}
	if !strings.Contains(err.Error(), "symlink") {
		t.Errorf("error should mention the symlink, got: %v", err)
	}

	// The whole point: the target must be byte-for-byte untouched.
	data, readErr := os.ReadFile(target)
	if readErr != nil {
		t.Fatalf("read target: %v", readErr)
	}
	if string(data) != "my project notes\n" {
		t.Errorf("target content = %q, want it unmodified", data)
	}

	// And the link must still be a link, not silently replaced.
	fi, lstatErr := os.Lstat(link)
	if lstatErr != nil {
		t.Fatalf("lstat: %v", lstatErr)
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		t.Error("symlink was replaced by a regular file")
	}
}

func TestAppendClaudeImport_refusesLinkToOwnGeneratedFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	rskDir := filepath.Join(dir, ".rsk")
	if err := os.MkdirAll(rskDir, 0o750); err != nil {
		t.Fatalf("mkdir .rsk: %v", err)
	}

	// The self-import loop: CLAUDE.md → .rsk/CLAUDE.md, which is the file rsk
	// generates. Appending here makes the generated file import itself.
	generated := filepath.Join(rskDir, ClaudeFileName)
	link := filepath.Join(dir, ClaudeFileName)
	writeSymlink(t, generated, link, "@../.claude/skills/x/SKILL.md\n")

	if err := appendClaudeImport(link); err == nil {
		t.Fatal("expected an error when CLAUDE.md points at .rsk/CLAUDE.md")
	}

	data, err := os.ReadFile(generated)
	if err != nil {
		t.Fatalf("read generated: %v", err)
	}
	if strings.Contains(string(data), claudeImportLine) {
		t.Errorf("generated file imports itself: %q", data)
	}
}

func TestRemoveClaudeImport_refusesSymlink(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	target := filepath.Join(dir, "AGENTS.md")
	link := filepath.Join(dir, "CLAUDE.md")
	writeSymlink(t, target, link, "notes\n"+claudeImportLine+"\n")

	if err := removeClaudeImport(link); err == nil {
		t.Fatal("expected an error for a symlinked CLAUDE.md")
	}

	// The rename path would otherwise replace the link with a regular file.
	fi, err := os.Lstat(link)
	if err != nil {
		t.Fatalf("lstat: %v", err)
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		t.Error("symlink was replaced by a regular file")
	}

	data, readErr := os.ReadFile(target)
	if readErr != nil {
		t.Fatalf("read target: %v", readErr)
	}
	if !strings.Contains(string(data), claudeImportLine) {
		t.Errorf("target was modified: %q", data)
	}
}

func TestWritePinnedClaude_refusesSymlink(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	rskDir := filepath.Join(dir, ".rsk")
	if err := os.MkdirAll(rskDir, 0o750); err != nil {
		t.Fatalf("mkdir .rsk: %v", err)
	}

	// .rsk/CLAUDE.md itself symlinked away — the rename would replace the link.
	elsewhere := filepath.Join(dir, "elsewhere.md")
	writeSymlink(t, elsewhere, filepath.Join(rskDir, ClaudeFileName), "old\n")

	if err := writePinnedClaude(rskDir, []string{"go-architect"}); err == nil {
		t.Fatal("expected an error when .rsk/CLAUDE.md is a symlink")
	}

	data, err := os.ReadFile(elsewhere)
	if err != nil {
		t.Fatalf("read elsewhere: %v", err)
	}
	if string(data) != "old\n" {
		t.Errorf("target modified: %q", data)
	}
}

func TestRefuseIfLink_allowsRegularAndMissing(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	regular := filepath.Join(dir, "regular.md")
	if err := os.WriteFile(regular, []byte("hi\n"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := refuseIfLink(regular); err != nil {
		t.Errorf("regular file should be allowed, got: %v", err)
	}

	// A nonexistent path is what `rsk new` hits on a fresh project.
	if err := refuseIfLink(filepath.Join(dir, "missing.md")); err != nil {
		t.Errorf("missing file should be allowed, got: %v", err)
	}
}

func TestRefuseIfLink_refusesDirectory(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	sub := filepath.Join(dir, "subdir")
	if err := os.Mkdir(sub, 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := refuseIfLink(sub); err == nil {
		t.Error("expected an error for a directory")
	}
}

func TestSyncPinned_regularFileEndToEnd(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".rsk"), 0o750); err != nil {
		t.Fatalf("mkdir .rsk: %v", err)
	}
	claudeMD := filepath.Join(dir, ClaudeFileName)
	if err := os.WriteFile(claudeMD, []byte("# notes\n"), 0o644); err != nil {
		t.Fatalf("seed CLAUDE.md: %v", err)
	}
	if err := os.Chmod(claudeMD, 0o644); err != nil {
		t.Fatalf("chmod CLAUDE.md: %v", err)
	}

	ct := &ClaudeTool{}
	if err := ct.SyncPinned(dir, []string{"go-architect"}); err != nil {
		t.Fatalf("SyncPinned: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, ".rsk", ClaudeFileName))
	if err != nil {
		t.Fatalf("read pinned: %v", err)
	}
	if !strings.Contains(string(data), "@SKILLS.md") {
		t.Errorf("pinned file missing the SKILLS.md import: %q", data)
	}

	root, err := os.ReadFile(filepath.Join(dir, ClaudeFileName))
	if err != nil {
		t.Fatalf("read CLAUDE.md: %v", err)
	}
	if !strings.Contains(string(root), claudeImportLine) {
		t.Errorf("CLAUDE.md missing import line: %q", root)
	}

	// Idempotent: a second sync must not append a duplicate.
	if err = ct.SyncPinned(dir, []string{"go-architect"}); err != nil {
		t.Fatalf("second SyncPinned: %v", err)
	}
	root2, err := os.ReadFile(filepath.Join(dir, ClaudeFileName))
	if err != nil {
		t.Fatalf("reread CLAUDE.md: %v", err)
	}
	if strings.Count(string(root2), claudeImportLine) != 1 {
		t.Errorf("import line duplicated: %q", root2)
	}

	// R3: the two-step sync (writePinnedClaude then appendClaudeImport) must
	// not leave CLAUDE.md tighter than it started.
	fi, err := os.Stat(claudeMD)
	if err != nil {
		t.Fatalf("stat CLAUDE.md: %v", err)
	}
	if got := fi.Mode().Perm(); got != 0o644 {
		t.Errorf("CLAUDE.md perm = %o, want 0644 preserved", got)
	}
}

func TestSyncLocalImport_addsAndRemovesOnlyItsOwnLine(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, ClaudeLocalFileName)
	claude := &ClaudeTool{}

	// A personal CLAUDE.local.md is the user's: the import is appended, never replaced.
	if err := os.WriteFile(path, []byte("my own note\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if err := claude.SyncLocal(dir, true); err != nil {
			t.Fatal(err)
		}
	}
	data, _ := os.ReadFile(path)
	if strings.Count(string(data), "@.rsk/SKILLS.local.md") != 1 || !strings.Contains(string(data), "my own note") {
		t.Errorf("CLAUDE.local.md = %q", data)
	}
	if err := claude.SyncLocal(dir, false); err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(path)
	if strings.Contains(string(data), "SKILLS.local.md") || !strings.Contains(string(data), "my own note") {
		t.Errorf("after removal CLAUDE.local.md = %q", data)
	}

	// A file rsk created is removed with its last line.
	os.Remove(path)
	if err := claude.SyncLocal(dir, true); err != nil {
		t.Fatal(err)
	}
	if err := claude.SyncLocal(dir, false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err == nil {
		t.Error("an empty CLAUDE.local.md is left behind")
	}
	// Nothing to remove, nothing to create.
	if err := claude.SyncLocal(dir, false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err == nil {
		t.Error("removing an import created the file")
	}
}

func TestSyncLocalImport_refusesASymlink(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writeSymlink(t, filepath.Join(dir, "elsewhere.md"), filepath.Join(dir, ClaudeLocalFileName), "x\n")
	if err := (&ClaudeTool{}).SyncLocal(dir, true); err == nil {
		t.Fatal("rsk wrote through a symlinked CLAUDE.local.md")
	}
}

func TestOpenCodeSyncLocal_listsTheLocalDigestOnlyWhenThereIsOne(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	oc := &openCodeTool{}
	if err := oc.SyncPinned(dir, []string{"a"}); err != nil {
		t.Fatal(err)
	}
	if err := oc.SyncLocal(dir, true); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(filepath.Join(dir, OpenCodeFileName))
	if !strings.Contains(string(data), ".rsk/SKILLS.local.md") || !strings.Contains(string(data), ".rsk/SKILLS.md") {
		t.Errorf("opencode.json = %s", data)
	}
	// SyncPinned rewrites the managed entries; the local one is re-added by SyncLocal.
	if err := oc.SyncPinned(dir, []string{"a"}); err != nil {
		t.Fatal(err)
	}
	if err := oc.SyncLocal(dir, false); err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(filepath.Join(dir, OpenCodeFileName))
	if strings.Contains(string(data), "SKILLS.local.md") || !strings.Contains(string(data), ".rsk/SKILLS.md") {
		t.Errorf("after removal opencode.json = %s", data)
	}
	if err := oc.RemovePinned(dir); err != nil {
		t.Fatal(err)
	}
}
