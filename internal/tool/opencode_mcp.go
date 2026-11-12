package tool

import (
	"fmt"
	"path/filepath"
	"reflect"
)

const (
	// openCodeMCPKey is the opencode.json key holding local MCP servers.
	openCodeMCPKey = "mcp"

	// mcpLocalType is OpenCode's type for a local server.
	mcpLocalType = "local"
)

// openCodeMCPServer is the opencode MCP entry rsk registers.
func openCodeMCPServer() map[string]any {
	return map[string]any{
		mcpTypeKey:    mcpLocalType,
		mcpCommandKey: []any{rskMCPServerName, mcpSubcommand, mcpReadOnlyFlag},
		"enabled":     true,
	}
}

// legacyOpenCodeMCPServer is the pre-read-only entry consumers may still have.
func legacyOpenCodeMCPServer() map[string]any {
	return map[string]any{
		mcpTypeKey:    mcpLocalType,
		mcpCommandKey: []any{rskMCPServerName, mcpSubcommand},
		"enabled":     true,
	}
}

// openCodeV2MCPServer is the entry written into mcp.servers.
func openCodeV2MCPServer() map[string]any {
	return map[string]any{
		mcpTypeKey:    mcpLocalType,
		mcpCommandKey: []any{rskMCPServerName, mcpSubcommand, mcpReadOnlyFlag},
	}
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

// hasMCP reports an rsk server in either layout: the current mcp.servers map or
// the pre-v2 mcp.rsk entry.
func (*openCodeTool) hasMCP(projectDir string) (bool, error) {
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

// mcpConfigName returns opencode.json after refusing the alternate files
// OpenCode may resolve ahead of it.
func (*openCodeTool) mcpConfigName(projectDir string) (string, error) {
	return OpenCodeFileName, rejectAlternateOpenCodeConfigs(projectDir)
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
