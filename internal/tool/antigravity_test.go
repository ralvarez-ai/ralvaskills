package tool

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func antigravityReadConfig(t *testing.T, path string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	cfg := make(map[string]any)
	if err = json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	return cfg
}

func TestAntigravityPinnedSyncIsNoOp(t *testing.T) {
	t.Parallel()
	project := t.TempDir()
	client := &antigravityTool{}
	if err := client.SyncPinned(project, []string{"go-architect"}); err != nil {
		t.Fatal(err)
	}
	if err := client.RemovePinned(project); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(project)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("pinned sync wrote %d entries, want none", len(entries))
	}
}

func TestAntigravityMCPPreservesConfigAndIsIdempotent(t *testing.T) {
	t.Parallel()
	project := t.TempDir()
	dir := filepath.Join(project, ".agents")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, AntigravityMCPFileName)
	original := map[string]any{
		"mcpServers": map[string]any{"other": map[string]any{"command": "other"}},
		"theme":      "dark",
	}
	raw, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, raw, 0o640); err != nil {
		t.Fatal(err)
	}

	client := &antigravityTool{}
	for range 2 {
		got, registerErr := client.RegisterMCP(project)
		if registerErr != nil {
			t.Fatal(registerErr)
		}
		if got != path {
			t.Fatalf("RegisterMCP path = %q, want %q", got, path)
		}
	}

	cfg := antigravityReadConfig(t, path)
	if cfg["theme"] != "dark" {
		t.Errorf("unrelated key lost: %v", cfg)
	}
	servers, ok := cfg["mcpServers"].(map[string]any)
	if !ok {
		t.Fatalf("mcpServers = %T, want object", cfg["mcpServers"])
	}
	if _, exists := servers["other"]; !exists {
		t.Errorf("unrelated server lost: %v", servers)
	}
	want := map[string]any{"command": "rsk", "args": []any{"mcp", mcpReadOnlyFlag}}
	if !reflect.DeepEqual(servers["rsk"], want) {
		t.Errorf("rsk server = %v, want %v", servers["rsk"], want)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o640 {
		t.Fatalf("mode = %v", info.Mode().Perm())
	}

	if err = client.UnregisterMCP(project); err != nil {
		t.Fatal(err)
	}
	if got := antigravityReadConfig(t, path); !reflect.DeepEqual(got, original) {
		t.Fatalf("restored config = %v, want %v", got, original)
	}
}

func TestAntigravityMCPRefusesConflictingServer(t *testing.T) {
	t.Parallel()
	project := t.TempDir()
	dir := filepath.Join(project, ".agents")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, AntigravityMCPFileName)
	original := `{"mcpServers":{"rsk":{"command":"other"}}}`
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	client := &antigravityTool{}
	if _, err := client.RegisterMCP(project); err == nil {
		t.Fatal("expected conflict")
	}
	if err := client.UnregisterMCP(project); err == nil {
		t.Fatal("expected refusal to remove a modified rsk server")
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != original {
		t.Fatalf("config = %q, err=%v", data, err)
	}
}

func TestAntigravityMCPRefusesSymlinkedAgentsDir(t *testing.T) {
	t.Parallel()
	project := t.TempDir()
	target := t.TempDir()
	if err := os.Symlink(target, filepath.Join(project, ".agents")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}

	client := &antigravityTool{}
	if _, err := client.RegisterMCP(project); err == nil {
		t.Fatal("expected symlinked .agents to be refused")
	}
	if _, err := os.Lstat(filepath.Join(target, AntigravityMCPFileName)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("wrote through the symlink: %v", err)
	}
}

func TestRefuseIfLinkDirMissingDirIsAllowed(t *testing.T) {
	t.Parallel()
	if err := refuseIfLinkDir(filepath.Join(t.TempDir(), "absent")); err != nil {
		t.Fatalf("missing directory: %v", err)
	}
}
