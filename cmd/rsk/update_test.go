package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/digest"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

// writeLinkedSkill seeds the linked skill a pinned entry points at: a SKILL.md
// and the DIGEST.md it ships.
func writeLinkedSkill(t *testing.T, dir, name, version, body string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	skillMD := "---\nname: " + name + "\nversion: " + version + "\ndescription: does " + name + "\n---\nbody\n"
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(skillMD), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, digest.FileName), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestRefreshPinnedDigest covers the bug where `rsk update` re-linked a skill
// but left .rsk/SKILLS.md holding the previous digest until the next `rsk pin`:
// a changed DIGEST.md must land in the generated file, and the call must be a
// silent no-op outside a project.
func TestRefreshPinnedDigest(t *testing.T) {
	projDir := t.TempDir()
	t.Chdir(projDir)

	var buf, errBuf bytes.Buffer
	refreshPinnedDigest(&buf, &errBuf)
	if got := buf.String() + errBuf.String(); got != "" {
		t.Errorf("outside a project: expected no output, got %q", got)
	}

	if err := manifest.WriteMod(filepath.Join(projDir, ".rsk"), manifest.Mod{Pinned: []string{"tdd"}}); err != nil {
		t.Fatalf("write rsk.mod: %v", err)
	}
	skillDir := filepath.Join(projDir, ".claude", "skills", "tdd")
	index := filepath.Join(projDir, ".rsk", digest.IndexFileName)

	writeLinkedSkill(t, skillDir, "tdd", "1.0.0", "Rules\n- old rule\n")
	refreshPinnedDigest(&buf, &errBuf)
	if got := buf.String() + errBuf.String(); got != "" {
		t.Errorf("expected no warning, got %q", got)
	}
	assertIndexHas(t, index, "## tdd v1.0.0", "- old rule")

	// A pull rewrites DIGEST.md in place; the next refresh must carry the new
	// rule into SKILLS.md without a second `rsk pin`.
	writeLinkedSkill(t, skillDir, "tdd", "1.0.1", "Rules\n- new rule\n")
	refreshPinnedDigest(&buf, &errBuf)
	assertIndexHas(t, index, "## tdd v1.0.1", "- new rule")

	data, err := os.ReadFile(index)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(data, []byte("- old rule")) {
		t.Errorf("SKILLS.md still has the previous digest:\n%s", data)
	}
}

// TestNeedsUpdate covers the drift k3vx7p reported: `rsk update` decided
// "up to date" from the symlink alone, so a lock left behind by a restore or
// an aborted update never repaired itself and the digest re-render never ran.
// TestProjectLockedVersions checks that the lock is read from the project and
// ignored for a global run, where no single project owns the state.
func TestProjectLockedVersions(t *testing.T) {
	projDir := t.TempDir()
	t.Chdir(projDir)

	if got := projectLockedVersions(false); len(got) != 0 {
		t.Fatalf("outside a project: got %v, want empty", got)
	}

	rskDir := filepath.Join(projDir, ".rsk")
	if err := manifest.WriteMod(rskDir, manifest.Mod{}); err != nil {
		t.Fatalf("write rsk.mod: %v", err)
	}
	lock := manifest.Lock{Skills: []manifest.LockEntry{
		{Name: "go-architect", Version: "1.20.2", Source: skill.SourceRegistry},
	}}
	if err := manifest.WriteLock(rskDir, lock); err != nil {
		t.Fatalf("write rsk.lock: %v", err)
	}

	if got := projectLockedVersions(false); got["go-architect"] != "1.20.2" {
		t.Errorf("project scope: got %v, want go-architect 1.20.2", got)
	}
	if got := projectLockedVersions(true); len(got) != 0 {
		t.Errorf("global scope: got %v, want empty", got)
	}
}

// assertIndexHas fails when SKILLS.md lacks any of wants.
func assertIndexHas(t *testing.T, path string, wants ...string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	for _, want := range wants {
		if !bytes.Contains(data, []byte(want)) {
			t.Errorf("SKILLS.md lacks %q:\n%s", want, data)
		}
	}
}
