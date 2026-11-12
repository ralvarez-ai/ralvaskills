package tool_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
)

func TestClaudeMCPRegisterIsIdempotentAndPreservesOthers(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, tool.MCPFileName)
	writeFile(t, path, `{"mcpServers":{"other":{"command":"other"}}}`)

	claude := mustTool(t, tool.ClaudeID)

	if _, err := claude.RegisterMCP(dir); err != nil {
		t.Fatalf("RegisterMCP: %v", err)
	}
	first := readFile(t, path)
	if _, err := claude.RegisterMCP(dir); err != nil {
		t.Fatalf("RegisterMCP (second): %v", err)
	}
	if second := readFile(t, path); first != second {
		t.Error("RegisterMCP is not idempotent")
	}

	servers := mcpServers(t, path, "mcpServers")
	if _, ok := servers["other"]; !ok {
		t.Error("RegisterMCP dropped the unrelated server")
	}
	rsk, ok := servers["rsk"].(map[string]any)
	if !ok {
		t.Fatalf("rsk server missing: %+v", servers)
	}
	if rsk["command"] != "rsk" {
		t.Errorf("rsk command = %v, want rsk", rsk["command"])
	}

	if err := claude.UnregisterMCP(dir); err != nil {
		t.Fatalf("UnregisterMCP: %v", err)
	}
	servers = mcpServers(t, path, "mcpServers")
	if _, present := servers["rsk"]; present {
		t.Error("UnregisterMCP left the rsk server behind")
	}
	if _, present := servers["other"]; !present {
		t.Error("UnregisterMCP dropped the unrelated server")
	}
}

func TestOpenCodeMCPRegisterPreservesInstructions(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, tool.OpenCodeFileName)
	writeFile(t, path, `{"instructions":[".claude/skills/tdd/SKILL.md"]}`)

	opencode := mustTool(t, tool.OpenCodeID)
	if _, err := opencode.RegisterMCP(dir); err != nil {
		t.Fatalf("RegisterMCP: %v", err)
	}

	cfg := readJSON(t, path)
	if _, ok := cfg["instructions"]; !ok {
		t.Error("RegisterMCP dropped the instructions key")
	}
	servers, _ := cfg["mcp"].(map[string]any)
	rsk, ok := servers["rsk"].(map[string]any)
	if !ok {
		t.Fatalf("mcp.rsk missing: %+v", cfg)
	}
	if rsk["type"] != "local" {
		t.Errorf("rsk type = %v, want local", rsk["type"])
	}

	if err := opencode.UnregisterMCP(dir); err != nil {
		t.Fatalf("UnregisterMCP: %v", err)
	}
	if cfg = readJSON(t, path); cfg["mcp"] != nil {
		t.Errorf("UnregisterMCP left mcp behind: %+v", cfg["mcp"])
	}
}

func TestOpenCodeV2RegistrationPreservesOtherServers(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, tool.OpenCodeFileName)
	writeFile(t, path, `{"mcp":{"servers":{"other":{"type":"local","command":["other"]}}}}`)
	client := mustTool(t, tool.OpenCodeID)
	for range 2 {
		if _, err := client.RegisterMCP(dir); err != nil {
			t.Fatal(err)
		}
	}
	cfg := readJSON(t, path)
	mcp, ok := cfg["mcp"].(map[string]any)
	if !ok {
		t.Fatal("missing mcp")
	}
	servers, ok := mcp["servers"].(map[string]any)
	if !ok || servers["other"] == nil {
		t.Fatal("missing other server")
	}
	rsk, ok := servers["rsk"].(map[string]any)
	if !ok || !strings.Contains(readFile(t, path), "--read-only") || rsk["type"] != "local" {
		t.Fatalf("invalid rsk server: %+v", servers["rsk"])
	}
	if err := client.UnregisterMCP(dir); err != nil {
		t.Fatal(err)
	}
	cfg = readJSON(t, path)
	mcp, ok = cfg["mcp"].(map[string]any)
	if !ok {
		t.Fatal("missing mcp after removal")
	}
	servers, ok = mcp["servers"].(map[string]any)
	if !ok {
		t.Fatal("missing servers after removal")
	}
	if _, exists := servers["rsk"]; exists || servers["other"] == nil {
		t.Fatalf("servers after removal = %+v", servers)
	}
}

func TestMCPRegistrationRefusesConflictingConfig(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		id      tool.ID
		file    string
		content string
	}{
		{tool.ClaudeID, tool.MCPFileName, `{"mcpServers":{"rsk":{"command":"other"}}}`},
		{tool.ClaudeID, tool.MCPFileName, `{"mcpServers":"invalid"}`},
		{tool.OpenCodeID, tool.OpenCodeFileName, `{"mcp":{"servers":{"rsk":{"type":"local","command":["other"]}}}}`},
	} {
		t.Run(string(tc.id)+tc.content, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			path := filepath.Join(dir, tc.file)
			writeFile(t, path, tc.content)
			if _, err := mustTool(t, tc.id).RegisterMCP(dir); err == nil {
				t.Fatal("expected conflict")
			}
			if got := readFile(t, path); got != tc.content {
				t.Fatalf("config changed: %s", got)
			}
		})
	}
}

func TestMCPPreflightDoesNotWriteBeforeConflict(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, tool.OpenCodeFileName)
	original := `{"mcp":{"rsk":{"type":"local","command":["other"]}}}`
	writeFile(t, path, original)
	if err := tool.PreflightMCP(dir, []tool.ID{tool.ClaudeID, tool.CodexID, tool.OpenCodeID}); err == nil {
		t.Fatal("expected conflicting OpenCode server to fail preflight")
	}
	if got := readFile(t, path); got != original {
		t.Fatalf("config changed: %s", got)
	}
	for _, name := range []string{tool.MCPFileName, ".codex"} {
		if _, err := os.Lstat(filepath.Join(dir, name)); !os.IsNotExist(err) {
			t.Fatalf("preflight wrote %s: %v", name, err)
		}
	}
}

func TestWorkflowPointerIsIdempotentAndRemovable(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "CLAUDE.md")
	writeFile(t, path, "# My project\n\nKeep me.\n")

	if err := tool.WriteMCPPointer(dir); err != nil {
		t.Fatalf("WriteMCPPointer: %v", err)
	}
	first := readFile(t, path)
	if err := tool.WriteMCPPointer(dir); err != nil {
		t.Fatalf("WriteMCPPointer (second): %v", err)
	}
	second := readFile(t, path)
	if first != second {
		t.Error("WriteMCPPointer is not idempotent")
	}
	if strings.Count(second, "<!-- rsk-mcp:start -->") != 1 {
		t.Errorf("pointer block count = %d, want 1", strings.Count(second, "<!-- rsk-mcp:start -->"))
	}
	if !strings.Contains(second, "Keep me.") {
		t.Error("WriteMCPPointer dropped existing content")
	}

	if err := tool.RemoveMCPPointer(dir); err != nil {
		t.Fatalf("RemoveMCPPointer: %v", err)
	}
	after := readFile(t, path)
	if strings.Contains(after, "rsk-mcp") {
		t.Errorf("RemoveMCPPointer left the block: %q", after)
	}
	if !strings.Contains(after, "Keep me.") {
		t.Error("RemoveMCPPointer dropped existing content")
	}
}

//
//nolint:ireturn // returns the registry interface the test needs to call
func mustTool(t *testing.T, id tool.ID) tool.Tool {
	t.Helper()

	tt, ok := tool.Get(id)
	if !ok {
		t.Fatalf("tool %q not registered", id)
	}
	return tt
}

func mcpServers(t *testing.T, path, key string) map[string]any {
	t.Helper()

	cfg := readJSON(t, path)
	servers, ok := cfg[key].(map[string]any)
	if !ok {
		t.Fatalf("%s: %s is not an object: %+v", path, key, cfg)
	}
	return servers
}

func readJSON(t *testing.T, path string) map[string]any {
	t.Helper()

	var cfg map[string]any
	if err := json.Unmarshal([]byte(readFile(t, path)), &cfg); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	return cfg
}

func readFile(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), fsperm.File); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
