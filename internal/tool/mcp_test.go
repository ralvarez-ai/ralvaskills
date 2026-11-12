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

func TestProjectHasMCPDetectsConfiguredClients(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writeFile(
		t,
		filepath.Join(dir, tool.OpenCodeFileName),
		`{"mcp":{"servers":{"rsk":{"type":"local","command":["rsk","mcp"]}}}}`,
	)
	found, err := tool.ProjectHasMCP(dir, []tool.ID{tool.OpenCodeID})
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	found, err = tool.ProjectHasMCP(dir, []tool.ID{tool.ClaudeID})
	if err != nil || found {
		t.Fatalf("unconfigured client found=%v err=%v", found, err)
	}
}

// TestProjectHasMCPRecognizesLegacyConfigurations covers the configurations
// written by earlier releases: Claude's unchanged .mcp.json entry, OpenCode's
// pre-v2 mcp.rsk entry, OpenCode's nested mcp.servers entry, and Codex's
// managed TOML block.
func TestProjectHasMCPRecognizesLegacyConfigurations(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		id    tool.ID
		setup func(*testing.T, string)
	}{
		{
			name: "claude",
			id:   tool.ClaudeID,
			setup: func(t *testing.T, dir string) {
				t.Helper()
				writeFile(t, filepath.Join(dir, tool.MCPFileName),
					`{"mcpServers":{"rsk":{"command":"rsk","args":["mcp"]}}}`)
			},
		},
		{
			name: "opencode pre-v2",
			id:   tool.OpenCodeID,
			setup: func(t *testing.T, dir string) {
				t.Helper()
				writeFile(t, filepath.Join(dir, tool.OpenCodeFileName),
					`{"mcp":{"rsk":{"type":"local","command":["rsk","mcp"],"enabled":true}}}`)
			},
		},
		{
			name: "opencode v2",
			id:   tool.OpenCodeID,
			setup: func(t *testing.T, dir string) {
				t.Helper()
				writeFile(t, filepath.Join(dir, tool.OpenCodeFileName),
					`{"mcp":{"servers":{"rsk":{"type":"local","command":["rsk","mcp","--read-only"]}}}}`)
			},
		},
		{
			name: "codex",
			id:   tool.CodexID,
			setup: func(t *testing.T, dir string) {
				t.Helper()
				if _, err := mustTool(t, tool.CodexID).RegisterMCP(dir); err != nil {
					t.Fatalf("RegisterMCP: %v", err)
				}
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			tc.setup(t, dir)
			found, err := tool.ProjectHasMCP(dir, []tool.ID{tc.id})
			if err != nil || !found {
				t.Fatalf("found=%v err=%v", found, err)
			}
		})
	}
}

func TestProjectHasMCPRejectsMalformedAndUnknownClients(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, tool.OpenCodeFileName), `{`)
	if _, err := tool.ProjectHasMCP(dir, []tool.ID{tool.OpenCodeID}); err == nil {
		t.Fatal("expected a parse error for malformed opencode.json")
	}
	if _, err := tool.ProjectHasMCP(dir, []tool.ID{tool.ID("unknown")}); err == nil {
		t.Fatal("expected an error for an unknown tool")
	}
}

// TestLegacyOpenCodeMCPConfigUpgradesWithoutLosingSettings guards consumers
// whose opencode.json still carries the pre-v2 command without --read-only:
// registration must upgrade that one entry in place and leave everything else.
func TestLegacyOpenCodeMCPConfigUpgradesWithoutLosingSettings(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, tool.OpenCodeFileName)
	writeFile(t, path, `{"instructions":["docs/keep.md"],"mcp":{`+
		`"rsk":{"type":"local","command":["rsk","mcp"],"enabled":true},`+
		`"other":{"type":"local","command":["other"]}}}`)

	client := mustTool(t, tool.OpenCodeID)
	if found, err := tool.ProjectHasMCP(dir, []tool.ID{tool.OpenCodeID}); err != nil || !found {
		t.Fatalf("legacy entry not detected: found=%v err=%v", found, err)
	}
	if _, err := client.RegisterMCP(dir); err != nil {
		t.Fatalf("RegisterMCP over legacy config: %v", err)
	}
	first := readFile(t, path)
	if _, err := client.RegisterMCP(dir); err != nil {
		t.Fatalf("RegisterMCP (second): %v", err)
	}
	if second := readFile(t, path); first != second {
		t.Fatal("RegisterMCP is not idempotent after upgrading a legacy entry")
	}
	if !strings.Contains(first, `"--read-only"`) {
		t.Fatalf("legacy command was not upgraded: %s", first)
	}

	cfg := readJSON(t, path)
	instructions, ok := cfg["instructions"].([]any)
	if !ok || len(instructions) != 1 || instructions[0] != "docs/keep.md" {
		t.Fatalf("instructions not preserved: %+v", cfg["instructions"])
	}
	servers, ok := cfg["mcp"].(map[string]any)
	if !ok || servers["other"] == nil {
		t.Fatalf("other MCP server not preserved: %+v", cfg["mcp"])
	}
}

func TestOpenCodeMCPRejectsAlternateConfigFiles(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"opencode.jsonc", ".opencode/opencode.json", ".opencode/opencode.jsonc"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			path := filepath.Join(dir, filepath.FromSlash(name))
			if err := os.MkdirAll(filepath.Dir(path), fsperm.Dir); err != nil {
				t.Fatal(err)
			}
			original := `{"mcp":{}}`
			writeFile(t, path, original)
			if _, err := mustTool(t, tool.OpenCodeID).RegisterMCP(dir); err == nil {
				t.Fatalf("expected %s to require manual MCP setup", name)
			}
			if got := readFile(t, path); got != original {
				t.Fatalf("%s changed: %s", name, got)
			}
		})
	}
}

func TestMCPPreflightLeavesLegacyConfigUntouched(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, tool.OpenCodeFileName)
	original := `{"mcp":{"rsk":{"type":"local","command":["rsk","mcp"],"enabled":true}}}`
	writeFile(t, path, original)
	if err := tool.PreflightMCP(dir, []tool.ID{tool.OpenCodeID}); err != nil {
		t.Fatalf("preflight over legacy config: %v", err)
	}
	if got := readFile(t, path); got != original {
		t.Fatalf("preflight changed the real config: %s", got)
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

// TestAntigravityMCPPreflightAndDetection covers the Antigravity adapter in the
// shared MCP paths: preflight reads .agents/mcp_config.json without writing it,
// and ProjectHasMCP reads the mcpServers map Antigravity uses.
func TestAntigravityMCPPreflightAndDetection(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	agents := filepath.Join(dir, ".agents")
	if err := os.MkdirAll(agents, fsperm.Dir); err != nil {
		t.Fatal(err)
	}

	if err := tool.PreflightMCP(dir, []tool.ID{tool.AntigravityID}); err != nil {
		t.Fatalf("preflight: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(agents, tool.AntigravityMCPFileName)); !os.IsNotExist(err) {
		t.Fatalf("preflight wrote the real config: %v", err)
	}

	if found, err := tool.ProjectHasMCP(dir, []tool.ID{tool.AntigravityID}); err != nil || found {
		t.Fatalf("unconfigured found=%v err=%v", found, err)
	}

	path, err := mustTool(t, tool.AntigravityID).RegisterMCP(dir)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(agents, tool.AntigravityMCPFileName); path != want {
		t.Fatalf("path = %q, want %q", path, want)
	}
	if found, foundErr := tool.ProjectHasMCP(dir, []tool.ID{tool.AntigravityID}); foundErr != nil || !found {
		t.Fatalf("configured found=%v err=%v", found, foundErr)
	}
	if _, exists := mcpServers(t, path, "mcpServers")["rsk"]; !exists {
		t.Fatalf("rsk server missing from %s", path)
	}
}
