package tool

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNativeSkillDirectories(t *testing.T) {
	t.Parallel()
	project := t.TempDir()
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		id              ID
		project, global string
	}{
		{CodexID, filepath.Join(project, ".agents", "skills"), filepath.Join(home, ".agents", "skills")},
		{OpenCodeID, filepath.Join(project, ".opencode", "skills"), ""},
		{ClaudeID, filepath.Join(project, ".claude", "skills"), ""},
	} {
		adapter, ok := Get(tc.id)
		if !ok {
			t.Fatalf("%s not registered", tc.id)
		}
		if got := adapter.ProjectSkillsDir(project); got != tc.project {
			t.Errorf("%s project directory = %q, want %q", tc.id, got, tc.project)
		}
		if tc.global != "" && adapter.SkillsDir() != tc.global {
			t.Errorf("%s global directory = %q, want %q", tc.id, adapter.SkillsDir(), tc.global)
		}
	}
}

func TestCodexMCPPreservesConfigAndRemovesOnlyManagedBlock(t *testing.T) {
	t.Parallel()
	project := t.TempDir()
	dir := filepath.Join(project, ".codex")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "config.toml")
	original := "# keep this comment\nmodel = \"example\"\n\n[mcp_servers.other]\ncommand = \"other\"\n"
	if err := os.WriteFile(path, []byte(original), 0o640); err != nil {
		t.Fatal(err)
	}
	client := &codexTool{}
	for range 2 {
		if _, err := client.RegisterMCP(project); err != nil {
			t.Fatal(err)
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), original) || strings.Count(string(data), codexMCPStart) != 1 ||
		!strings.Contains(string(data), "--read-only") {
		t.Fatalf("config = %s", data)
	}
	if err = client.UnregisterMCP(project); err != nil {
		t.Fatal(err)
	}
	data, err = os.ReadFile(path)
	if err != nil || string(data) != original {
		t.Fatalf("restored config = %q, err=%v", data, err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o640 {
		t.Fatalf("mode = %v", info.Mode().Perm())
	}
}

func TestCodexMCPRefusesConflictingServer(t *testing.T) {
	t.Parallel()
	project := t.TempDir()
	dir := filepath.Join(project, ".codex")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "config.toml")
	original := "[mcp_servers.rsk]\ncommand = \"other\"\n"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}
	client := &codexTool{}
	if _, err := client.RegisterMCP(project); err == nil {
		t.Fatal("expected conflict")
	}
	if err := client.UnregisterMCP(project); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != original {
		t.Fatalf("config = %q, err=%v", data, err)
	}
}
