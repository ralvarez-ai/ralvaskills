package tool

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// jsonTools is every adapter whose MCP config is a JSON server map. Adding one
// to the registry and to this table is all a new JSON client needs.
func jsonTools(t *testing.T) []*jsonAdapter {
	t.Helper()
	ids := []ID{AntigravityID, KimiID, DroidID, GeminiID, QwenID, CursorID, AmpID}
	adapters := make([]*jsonAdapter, 0, len(ids))
	for _, id := range ids {
		client, ok := Get(id)
		if !ok {
			t.Fatalf("%s not registered", id)
		}
		adapter, ok := client.(*jsonAdapter)
		if !ok {
			t.Fatalf("%s is %T, want *jsonAdapter", id, client)
		}
		adapters = append(adapters, adapter)
	}
	return adapters
}

func readJSONObject(t *testing.T, path string) map[string]any {
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

func TestJSONAdaptersShareProjectSkillsDir(t *testing.T) {
	t.Parallel()
	project := t.TempDir()
	for _, adapter := range jsonTools(t) {
		if got, want := adapter.ProjectSkillsDir(project),
			filepath.Join(project, ".agents", "skills"); got != want {
			t.Errorf("%s project dir = %q, want %q", adapter.ID(), got, want)
		}
		if adapter.SkillsDir() == "" {
			t.Errorf("%s global dir is empty", adapter.ID())
		}
	}
}

func TestJSONAdaptersMCPRoundTrip(t *testing.T) {
	t.Parallel()
	for _, adapter := range jsonTools(t) {
		t.Run(string(adapter.ID()), func(t *testing.T) {
			t.Parallel()
			cfg := adapter.JSONMCP()
			project := t.TempDir()
			path := filepath.Join(project, filepath.FromSlash(cfg.file))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			original := map[string]any{
				cfg.key:     map[string]any{"other": map[string]any{"command": "other"}},
				"unrelated": "keep",
			}
			raw, err := json.Marshal(original)
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(path, raw, 0o640); err != nil {
				t.Fatal(err)
			}

			for range 2 {
				got, registerErr := adapter.RegisterMCP(project)
				if registerErr != nil {
					t.Fatal(registerErr)
				}
				if got != path {
					t.Fatalf("path = %q, want %q", got, path)
				}
			}

			got := readJSONObject(t, path)
			servers, ok := got[cfg.key].(map[string]any)
			if !ok {
				t.Fatalf("%s = %T, want object", cfg.key, got[cfg.key])
			}
			if _, exists := servers["other"]; !exists {
				t.Errorf("unrelated server lost: %v", servers)
			}
			if !reflect.DeepEqual(servers[rskMCPServerName], cfg.server()) {
				t.Errorf("rsk server = %v, want %v", servers[rskMCPServerName], cfg.server())
			}
			if got["unrelated"] != "keep" {
				t.Errorf("unrelated key lost: %v", got)
			}
			info, statErr := os.Stat(path)
			if statErr != nil {
				t.Fatal(statErr)
			}
			if info.Mode().Perm() != 0o640 {
				t.Fatalf("mode = %v", info.Mode().Perm())
			}

			if err = adapter.UnregisterMCP(project); err != nil {
				t.Fatal(err)
			}
			if restored := readJSONObject(t, path); !reflect.DeepEqual(restored, original) {
				t.Fatalf("restored = %v, want %v", restored, original)
			}
		})
	}
}

func TestJSONAdaptersRefuseConflictingServer(t *testing.T) {
	t.Parallel()
	for _, adapter := range jsonTools(t) {
		t.Run(string(adapter.ID()), func(t *testing.T) {
			t.Parallel()
			cfg := adapter.JSONMCP()
			project := t.TempDir()
			path := filepath.Join(project, filepath.FromSlash(cfg.file))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			original := `{"` + cfg.key + `":{"rsk":{"command":"other"}}}`
			if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
				t.Fatal(err)
			}

			if _, err := adapter.RegisterMCP(project); err == nil {
				t.Fatal("expected conflict")
			}
			if err := adapter.UnregisterMCP(project); err == nil {
				t.Fatal("expected refusal to remove a modified rsk server")
			}
			data, err := os.ReadFile(path)
			if err != nil || string(data) != original {
				t.Fatalf("config = %q, err=%v", data, err)
			}
		})
	}
}

func TestJSONAdaptersPinnedSyncIsNoOp(t *testing.T) {
	t.Parallel()
	for _, adapter := range jsonTools(t) {
		t.Run(string(adapter.ID()), func(t *testing.T) {
			t.Parallel()
			project := t.TempDir()
			if err := adapter.SyncPinned(project, []string{"go-architect"}); err != nil {
				t.Fatal(err)
			}
			if err := adapter.RemovePinned(project); err != nil {
				t.Fatal(err)
			}
			entries, err := os.ReadDir(project)
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 0 {
				t.Fatalf("pinned sync wrote %d entries, want none", len(entries))
			}
		})
	}
}

func TestJSONAdaptersRefuseSymlinkedDir(t *testing.T) {
	t.Parallel()
	for _, adapter := range jsonTools(t) {
		t.Run(string(adapter.ID()), func(t *testing.T) {
			t.Parallel()
			cfg := adapter.JSONMCP()
			relDir := filepath.Dir(filepath.FromSlash(cfg.file))
			if relDir == "." {
				t.Skip("config lives in the project root")
			}
			project := t.TempDir()
			target := t.TempDir()
			if err := os.Symlink(target, filepath.Join(project, relDir)); err != nil {
				t.Skipf("symlink unavailable: %v", err)
			}

			if _, err := adapter.RegisterMCP(project); err == nil {
				t.Fatal("expected symlinked config directory to be refused")
			}
			entries, err := os.ReadDir(target)
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 0 {
				t.Fatalf("wrote through the symlink: %v", entries)
			}
		})
	}
}

func TestKimiSkillsDirHonorsHomeOverride(t *testing.T) {
	t.Setenv(kimiHomeEnv, "/tmp/isolated-kimi")
	client, _ := Get(KimiID)
	if got, want := client.SkillsDir(), "/tmp/isolated-kimi/skills"; got != want {
		t.Fatalf("SkillsDir = %q, want %q", got, want)
	}
}

func TestCursorServerCarriesStdioType(t *testing.T) {
	t.Parallel()
	server := readOnlyMCPServer()
	if _, hasType := server[mcpTypeKey]; hasType {
		t.Fatal("the shared entry should not carry a transport type")
	}
	cursor := cursorMCPServer()
	if cursor[mcpTypeKey] != mcpStdioType {
		t.Fatalf("cursor type = %v, want %q", cursor[mcpTypeKey], mcpStdioType)
	}
	if !reflect.DeepEqual(cursor["args"], server["args"]) {
		t.Fatalf("cursor args = %v, want %v", cursor["args"], server["args"])
	}
}

func TestAmpUsesDottedMCPKeyAndRejectsJSONC(t *testing.T) {
	t.Parallel()
	project := t.TempDir()
	client, _ := Get(AmpID)
	path, err := client.RegisterMCP(project)
	if err != nil {
		t.Fatal(err)
	}
	if path != filepath.Join(project, ".amp", "settings.json") {
		t.Fatalf("path = %q", path)
	}
	cfg := readJSONObject(t, path)
	if _, ok := cfg[ampMCPKey].(map[string]any); !ok {
		t.Fatalf("%s = %T, want object", ampMCPKey, cfg[ampMCPKey])
	}
	if _, nested := cfg["amp"]; nested {
		t.Fatalf("wrote a nested amp object: %v", cfg)
	}

	other := t.TempDir()
	if mkErr := os.MkdirAll(filepath.Join(other, ".amp"), 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}
	if writeErr := os.WriteFile(filepath.Join(other, ampSettingsJSONC), []byte("{}"), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}
	if _, registerErr := client.RegisterMCP(other); registerErr == nil {
		t.Fatal("expected the .jsonc settings file to be refused")
	}
}

func TestProjectMCPFileCoversEveryJSONAdapter(t *testing.T) {
	t.Parallel()
	for _, adapter := range jsonTools(t) {
		file, ok := ProjectMCPFile(adapter.ID())
		if !ok || file != adapter.JSONMCP().file {
			t.Errorf("%s ProjectMCPFile = %q, %v", adapter.ID(), file, ok)
		}
	}
	if _, ok := ProjectMCPFile(CodexID); ok {
		t.Error("Codex is a TOML client and has no JSON MCP file")
	}
	if _, ok := ProjectMCPFile(ID("unknown")); ok {
		t.Error("unknown tool reported a file")
	}
}

func TestJSONAdapterUnregisterMissingFileIsNoOp(t *testing.T) {
	t.Parallel()
	for _, adapter := range jsonTools(t) {
		if err := adapter.UnregisterMCP(t.TempDir()); err != nil {
			t.Errorf("%s: %v", adapter.ID(), err)
		}
	}
}

func TestRefuseIfLinkDirMissingDirIsAllowed(t *testing.T) {
	t.Parallel()
	if err := refuseIfLinkDir(filepath.Join(t.TempDir(), "absent")); err != nil {
		t.Fatalf("missing directory: %v", err)
	}
}

func TestTrimSlash(t *testing.T) {
	t.Parallel()
	for in, want := range map[string]string{
		"/home/u/skills/": "/home/u/skills",
		"  /home/u  ":     "/home/u",
		"/":               "",
	} {
		if got := trimSlash(in); got != want {
			t.Errorf("trimSlash(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestJSONAdapterAlternateErrorsAreActionable(t *testing.T) {
	t.Parallel()
	client, _ := Get(AmpID)
	project := t.TempDir()
	if err := os.MkdirAll(filepath.Join(project, ".amp"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, ampSettingsJSONC), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, registerErr := client.RegisterMCP(project)
	if registerErr == nil || !strings.Contains(registerErr.Error(), ampSettingsJSONC) {
		t.Fatalf("err = %v, want one naming %s", registerErr, ampSettingsJSONC)
	}
	if errors.Is(registerErr, os.ErrNotExist) {
		t.Fatal("unexpected missing-file error")
	}
}
