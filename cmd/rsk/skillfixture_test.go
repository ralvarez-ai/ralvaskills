package main

import (
	"os"
	"path/filepath"
	"testing"
)

// writeCatalogSkill writes a minimal skill named name under dir.
func writeCatalogSkill(t *testing.T, dir, name string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, name), 0o750); err != nil {
		t.Fatal(err)
	}
	body := "---\nname: " + name + "\nversion: 1.0.0\n---\n"
	if err := os.WriteFile(filepath.Join(dir, name, "SKILL.md"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}
