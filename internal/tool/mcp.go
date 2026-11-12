package tool

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
	"go.ralva.dev/ralvaskills/v3/internal/fsx"
)

const (
	// MCPFileName is the project file Claude Code reads for MCP servers.
	MCPFileName = ".mcp.json"

	// rskMCPServerName is the key rsk registers its server under.
	rskMCPServerName = "rsk"

	// openCodeMCPKey is the opencode.json key holding local MCP servers.
	openCodeMCPKey = "mcp"

	// mcpSubcommand is the argv every client passes to rsk.
	mcpSubcommand = "mcp"
	mcpCommandKey = "command"
	mcpTypeKey    = "type"
	mcpLocalType  = "local"
	mcpStdioType  = "stdio"

	// mcpServersKey is the JSON key most clients use for their server map.
	mcpServersKey = "mcpServers"

	// mcpReadOnlyFlag keeps a registration to the read-only server, which omits
	// install_skills until the client's per-install approval behavior is verified.
	mcpReadOnlyFlag = "--read-only"

	mcpPointerStart = "<!-- rsk-mcp:start -->"
	mcpPointerEnd   = "<!-- rsk-mcp:end -->"
)

// claudeMCPServer is the Claude Code MCP entry rsk registers.
func claudeMCPServer() map[string]any {
	return map[string]any{mcpCommandKey: rskMCPServerName, "args": []any{mcpSubcommand}}
}

// JSONMCP reports where Claude Code keeps its MCP servers, so preflight and
// detection handle it like the other JSON clients.
func (*ClaudeTool) JSONMCP() jsonMCPConfig {
	return jsonMCPConfig{
		file:   MCPFileName,
		key:    mcpServersKey,
		server: claudeMCPServer,
	}
}

// openCodeMCPServer is the opencode MCP entry rsk registers.
func openCodeMCPServer() map[string]any {
	return map[string]any{
		mcpTypeKey:    mcpLocalType,
		mcpCommandKey: []any{rskMCPServerName, mcpSubcommand, mcpReadOnlyFlag},
		"enabled":     true,
	}
}

func legacyOpenCodeMCPServer() map[string]any {
	return map[string]any{
		mcpTypeKey:    mcpLocalType,
		mcpCommandKey: []any{rskMCPServerName, mcpSubcommand},
		"enabled":     true,
	}
}

// RegisterMCP writes the rsk server into the project's .mcp.json, preserving
// other servers and keys. It is idempotent: a matching entry is left untouched.
// It returns the path written.
func (c *ClaudeTool) RegisterMCP(projectDir string) (string, error) {
	return c.JSONMCP().Register(projectDir)
}

// UnregisterMCP removes the rsk server from .mcp.json, leaving others intact.
func (c *ClaudeTool) UnregisterMCP(projectDir string) error {
	return c.JSONMCP().Unregister(projectDir)
}

// RegisterMCP writes the rsk server into the project's opencode.json, preserving
// other servers and keys. It is idempotent and returns the path written.
func (*openCodeTool) RegisterMCP(projectDir string) (string, error) {
	if err := rejectAlternateOpenCodeConfigs(projectDir); err != nil {
		return "", err
	}
	path := filepath.Join(projectDir, OpenCodeFileName)
	cfg, err := readJSONConfig(path)
	if err != nil {
		return "", err
	}
	if mcp, ok := cfg[openCodeMCPKey].(map[string]any); ok {
		if _, exists := mcp["servers"]; exists {
			return registerOpenCodeV2(path, cfg, mcp)
		}
	}
	return registerMCPServer(
		path, openCodeMCPKey,
		openCodeMCPServer(), []map[string]any{legacyOpenCodeMCPServer()},
		readJSONConfig, writeJSONConfig,
	)
}

// UnregisterMCP removes the rsk server from opencode.json, leaving others intact.
func (*openCodeTool) UnregisterMCP(projectDir string) error {
	if err := rejectAlternateOpenCodeConfigs(projectDir); err != nil {
		return err
	}
	path := filepath.Join(projectDir, OpenCodeFileName)
	cfg, err := readJSONConfig(path)
	if err != nil {
		return err
	}
	if mcp, ok := cfg[openCodeMCPKey].(map[string]any); ok {
		if _, exists := mcp["servers"]; exists {
			return removeOpenCodeV2(path, cfg, mcp)
		}
	}
	return removeMCPServer(
		path, openCodeMCPKey,
		openCodeMCPServer(), []map[string]any{legacyOpenCodeMCPServer()},
		readJSONConfig, writeJSONConfig,
	)
}

// openCodeAlternateConfigs lists the files OpenCode reads before or beside
// opencode.json. rsk refuses them instead of writing a second source of truth.
var openCodeAlternateConfigs = []string{
	"opencode.jsonc",
	".opencode/opencode.json",
	".opencode/opencode.jsonc",
}

func rejectAlternateOpenCodeConfigs(projectDir string) error {
	return rejectAlternates(projectDir, OpenCodeFileName, openCodeAlternateConfigs)
}

func openCodeV2MCPServer() map[string]any {
	return map[string]any{
		mcpTypeKey:    mcpLocalType,
		mcpCommandKey: []any{rskMCPServerName, mcpSubcommand, mcpReadOnlyFlag},
	}
}

// ProjectHasMCP reports whether a selected client already has an rsk server.
func ProjectHasMCP(projectDir string, ids []ID) (bool, error) {
	for _, id := range ids {
		found, err := clientHasMCP(projectDir, id)
		if err != nil {
			return false, err
		}
		if found {
			return true, nil
		}
	}
	return false, nil
}

func clientHasMCP(projectDir string, id ID) (bool, error) {
	switch id {
	case CodexID:
		data, err := readCodexConfig(filepath.Join(projectDir, codexFolder, codexConfig))
		return strings.Contains(string(data), "[mcp_servers.rsk]"), err
	case OpenCodeID:
		return openCodeHasMCP(projectDir)
	default:
	}
	provider, ok := jsonMCPProviderFor(id)
	if !ok {
		return false, fmt.Errorf("unknown tool %q", id)
	}
	cfg, err := readJSONConfig(filepath.Join(projectDir, filepath.FromSlash(provider.JSONMCP().file)))
	if err != nil {
		return false, err
	}
	servers, valid := cfg[provider.JSONMCP().key].(map[string]any)
	if !valid {
		return false, nil
	}
	_, exists := servers[rskMCPServerName]
	return exists, nil
}

// openCodeHasMCP covers both OpenCode layouts: the current mcp.servers map and
// the pre-v2 mcp.rsk entry.
func openCodeHasMCP(projectDir string) (bool, error) {
	config, err := readJSONConfig(filepath.Join(projectDir, OpenCodeFileName))
	if err != nil {
		return false, err
	}
	servers, valid := config[openCodeMCPKey].(map[string]any)
	if !valid {
		return false, nil
	}
	if _, exists := servers[rskMCPServerName]; exists {
		return true, nil
	}
	nested, valid := servers["servers"].(map[string]any)
	if !valid {
		return false, nil
	}
	_, exists := nested[rskMCPServerName]
	return exists, nil
}

// jsonMCPProviderFor returns the JSON MCP provider for id.
//
//nolint:ireturn // returns the interface the caller checks
func jsonMCPProviderFor(id ID) (jsonMCPProvider, bool) {
	t, ok := Get(id)
	if !ok {
		return nil, false
	}
	provider, ok := t.(jsonMCPProvider)
	return provider, ok
}

// PreflightMCP checks all selected client configurations before any registration.
// It stages their existing files in a temporary project and invokes the same
// registration logic used for the real project.
func PreflightMCP(projectDir string, ids []ID) error {
	stage, err := os.MkdirTemp("", "rsk-mcp-preflight-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage)
	for _, id := range ids {
		if err = preflightClientMCP(projectDir, stage, id); err != nil {
			return err
		}
	}
	return nil
}

func preflightClientMCP(projectDir, stage string, id ID) error {
	client, ok := Get(id)
	if !ok {
		return fmt.Errorf("unknown tool %q", id)
	}
	name, err := mcpConfigName(projectDir, id)
	if err != nil {
		return err
	}
	source := filepath.Join(projectDir, name)
	if err = refuseIfLink(source); err != nil {
		return err
	}
	data, err := os.ReadFile(source)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		target := filepath.Join(stage, name)
		if err = os.MkdirAll(filepath.Dir(target), fsperm.Dir); err != nil {
			return err
		}
		if err = os.WriteFile(target, data, fsperm.File); err != nil {
			return err
		}
	}
	if _, err = client.RegisterMCP(stage); err != nil {
		return fmt.Errorf("%s MCP preflight: %w", id, err)
	}
	return nil
}

func mcpConfigName(projectDir string, id ID) (string, error) {
	switch id {
	case CodexID:
		name := filepath.Join(codexFolder, codexConfig)
		_, err := readCodexConfig(filepath.Join(projectDir, name))
		return name, err
	case OpenCodeID:
		return OpenCodeFileName, rejectAlternateOpenCodeConfigs(projectDir)
	default:
	}
	provider, ok := jsonMCPProviderFor(id)
	if !ok {
		return "", fmt.Errorf("MCP preflight is not implemented for %q", id)
	}
	cfg := provider.JSONMCP()
	if err := cfg.checkAlternates(projectDir); err != nil {
		return "", err
	}
	return cfg.file, nil
}

func registerOpenCodeV2(path string, cfg, mcp map[string]any) (string, error) {
	if err := refuseIfLink(path); err != nil {
		return "", err
	}
	servers, ok := mcp["servers"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("%s: mcp.servers is not an object", path)
	}
	if _, exists := mcp[rskMCPServerName]; exists {
		return "", fmt.Errorf("%s: legacy mcp.rsk conflicts with v2; migrate it manually", path)
	}
	server := openCodeV2MCPServer()
	if existing, exists := servers[rskMCPServerName]; exists {
		if reflect.DeepEqual(existing, server) {
			return path, nil
		}
		return "", fmt.Errorf("%s: mcp.servers.rsk already exists with different settings", path)
	}
	servers[rskMCPServerName] = server
	return path, writeJSONConfig(path, cfg)
}

func removeOpenCodeV2(path string, cfg, mcp map[string]any) error {
	if err := refuseIfLink(path); err != nil {
		return err
	}
	servers, ok := mcp["servers"].(map[string]any)
	if !ok {
		return fmt.Errorf("%s: mcp.servers is not an object", path)
	}
	if _, exists := mcp[rskMCPServerName]; exists {
		return fmt.Errorf("%s: legacy mcp.rsk conflicts with v2; resolve manually", path)
	}
	existing, exists := servers[rskMCPServerName]
	if !exists {
		return nil
	}
	if !reflect.DeepEqual(existing, openCodeV2MCPServer()) {
		return fmt.Errorf("%s: refusing to remove a modified rsk MCP server", path)
	}
	delete(servers, rskMCPServerName)
	if len(servers) == 0 {
		delete(mcp, "servers")
	}
	if len(mcp) == 0 {
		delete(cfg, openCodeMCPKey)
	}
	return writeJSONConfig(path, cfg)
}

// matchesAny reports whether v equals one of the accepted legacy shapes. A
// client that changed its server entry passes the old shape here, so the
// generic writer stays ignorant of which client needed it.
func matchesAny(v any, candidates []map[string]any) bool {
	for _, candidate := range candidates {
		if reflect.DeepEqual(v, candidate) {
			return true
		}
	}
	return false
}

func registerMCPServer(
	path, key string,
	server map[string]any,
	legacy []map[string]any,
	read func(string) (map[string]any, error),
	write func(string, map[string]any) error,
) (string, error) {
	if err := refuseIfLink(path); err != nil {
		return "", err
	}
	cfg, err := read(path)
	if err != nil {
		return "", err
	}

	servers, ok := cfg[key].(map[string]any)
	if !ok {
		if _, exists := cfg[key]; exists {
			return "", fmt.Errorf("%s: %s is not an object", path, key)
		}
		servers = make(map[string]any)
	}
	if reflect.DeepEqual(servers[rskMCPServerName], server) {
		return path, nil
	}
	if existing, exists := servers[rskMCPServerName]; exists && !matchesAny(existing, legacy) {
		return "", fmt.Errorf("%s: MCP server rsk already exists with different settings", path)
	}

	servers[rskMCPServerName] = server
	cfg[key] = servers
	if writeErr := write(path, cfg); writeErr != nil {
		return "", writeErr
	}
	return path, nil
}

func removeMCPServer(
	path, key string,
	managed map[string]any,
	legacy []map[string]any,
	read func(string) (map[string]any, error),
	write func(string, map[string]any) error,
) error {
	if err := refuseIfLink(path); err != nil {
		return err
	}
	cfg, err := read(path)
	if err != nil {
		return err
	}

	servers, ok := cfg[key].(map[string]any)
	if !ok {
		return nil
	}
	existing, present := servers[rskMCPServerName]
	if !present {
		return nil
	}
	if !reflect.DeepEqual(existing, managed) && !matchesAny(existing, legacy) {
		return fmt.Errorf("%s: refusing to remove a modified rsk MCP server", path)
	}

	delete(servers, rskMCPServerName)
	if len(servers) == 0 {
		delete(cfg, key)
	} else {
		cfg[key] = servers
	}
	return write(path, cfg)
}

const mcpPointerBody = "## Project standards\n\n" +
	"Before scaffolding or a non-trivial change, call the rsk MCP `project_profile` " +
	"tool to learn which skills apply here. Installed skills appear automatically " +
	"in your skill list — invoke them with the Skill tool. Never install without " +
	"asking.\n"

const mcpPointerBlock = mcpPointerStart + "\n" + mcpPointerBody + mcpPointerEnd + "\n"

// WriteMCPPointer adds the rsk workflow pointer to the project's ./CLAUDE.md,
// replacing any previous block so repeated calls are idempotent. Other content
// is preserved.
func WriteMCPPointer(projectDir string) error {
	path := filepath.Join(projectDir, ClaudeFileName)
	if err := refuseIfLink(path); err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read %s: %w", path, err)
	}

	content := strings.TrimRight(removeMCPPointerBlock(string(data)), "\n")
	if content != "" {
		content += "\n\n"
	}
	return writeMCPPointerFile(path, content+mcpPointerBlock)
}

// RemoveMCPPointer removes the rsk workflow pointer block from ./CLAUDE.md.
func RemoveMCPPointer(projectDir string) error {
	path := filepath.Join(projectDir, ClaudeFileName)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if linkErr := refuseIfLink(path); linkErr != nil {
		return linkErr
	}
	return writeMCPPointerFile(path, removeMCPPointerBlock(string(data)))
}

func removeMCPPointerBlock(s string) string {
	start := strings.Index(s, mcpPointerStart)
	if start < 0 {
		return s
	}
	end := strings.Index(s[start:], mcpPointerEnd)
	if end < 0 {
		return s
	}
	end = start + end + len(mcpPointerEnd)
	if end < len(s) && s[end] == '\n' {
		end++
	}
	return s[:start] + s[end:]
}

func writeMCPPointerFile(path, content string) error {
	return fsx.WriteAtomic(path, rskTempPattern, func(w io.Writer) error {
		_, err := io.WriteString(w, content)
		return err
	})
}
