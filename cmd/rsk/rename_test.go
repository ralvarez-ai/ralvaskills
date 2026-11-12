package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
)

const (
	modBody = `version = 1
tools = ["claude-code"]
pinned = ["hexagonal-arch"]

[skills]
  hexagonal-arch = "*"
  go-architect = "*"
`
	lockBody = `version = 1

[[skills]]
  name = "hexagonal-arch"
  version = "1.2.1"
  source = "registry"
`
)

func writeProject(t *testing.T) string {
	t.Helper()

	root := t.TempDir()
	t.Chdir(root)
	rskDir := filepath.Join(root, ".rsk")
	if err := os.MkdirAll(rskDir, fsperm.Dir); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(rskDir, "rsk.mod"), []byte(modBody), fsperm.File); err != nil {
		t.Fatalf("write rsk.mod: %v", err)
	}
	if err := os.WriteFile(filepath.Join(rskDir, "rsk.lock"), []byte(lockBody), fsperm.File); err != nil {
		t.Fatalf("write rsk.lock: %v", err)
	}
	return rskDir
}

//nolint:paralleltest // t.Chdir forbids parallel tests
func TestMigrateProjectRenamesDryRunWritesNothing(t *testing.T) {
	rskDir := writeProject(t)
	aliases := map[string]string{"hexagonal-arch": "hexagonal-architecture"}

	renamed, err := migrateProjectRenames(io.Discard, aliases, true)
	if err != nil {
		t.Fatalf("migrateProjectRenames: %v", err)
	}
	if len(renamed) != 1 {
		t.Fatalf("renamed = %v, want one entry", renamed)
	}

	mod, err := os.ReadFile(filepath.Join(rskDir, "rsk.mod"))
	if err != nil {
		t.Fatalf("read rsk.mod: %v", err)
	}
	if !strings.Contains(string(mod), "hexagonal-arch") {
		t.Errorf("dry run rewrote rsk.mod:\n%s", mod)
	}
}

//nolint:paralleltest // t.Chdir forbids parallel tests
func TestMigrateProjectRenames(t *testing.T) {
	rskDir := writeProject(t)
	aliases := map[string]string{"hexagonal-arch": "hexagonal-architecture"}

	renamed, err := migrateProjectRenames(io.Discard, aliases, false)
	if err != nil {
		t.Fatalf("migrateProjectRenames: %v", err)
	}
	if len(renamed) != 1 || renamed[0].New != "hexagonal-architecture" {
		t.Fatalf("renamed = %v, want hexagonal-architecture", renamed)
	}

	m, err := manifest.ReadMod(rskDir)
	if err != nil {
		t.Fatalf("ReadMod: %v", err)
	}
	if _, stale := m.Skills["hexagonal-arch"]; stale {
		t.Error("rsk.mod still carries the old skill key")
	}
	if _, ok := m.Skills["hexagonal-architecture"]; !ok {
		t.Error("rsk.mod is missing the new skill key")
	}
	if m.Pinned[0] != "hexagonal-architecture" {
		t.Errorf("pinned[0] = %q, want the new name", m.Pinned[0])
	}
	if m.Skills["go-architect"] != "*" {
		t.Error("an unrelated skill changed")
	}

	lock, err := manifest.ReadLock(rskDir)
	if err != nil {
		t.Fatalf("ReadLock: %v", err)
	}
	if lock.Skills[0].Name != "hexagonal-architecture" {
		t.Errorf("lock name = %q, want the new name", lock.Skills[0].Name)
	}
	if lock.Skills[0].Version != "1.2.1" {
		t.Errorf("lock version = %q, want it preserved", lock.Skills[0].Version)
	}
}

//nolint:paralleltest // t.Chdir forbids parallel tests
func TestMigrateProjectRenamesNoProject(t *testing.T) {
	t.Chdir(t.TempDir())

	renamed, err := migrateProjectRenames(io.Discard, map[string]string{"old": "new"}, false)
	if err != nil {
		t.Fatalf("migrateProjectRenames: %v", err)
	}
	if len(renamed) != 0 {
		t.Errorf("renamed = %v, want none outside a project", renamed)
	}
}
