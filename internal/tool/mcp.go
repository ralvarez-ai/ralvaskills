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
	mcpCommandKey  = "command"
	mcpTypeKey     = "type"
	mcpLocalType   = "local"

	mcpPointerStart = "<!-- rsk-mcp:start -->"
	mcpPointerEnd   = "<!-- rsk-mcp:end -->"
)

// claudeMCPServer is the Claude Code MCP entry rsk registers.
func claudeMCPServer() map[string]any {
	return map[string]any{mcpCommandKey: rskMCPServerName, "args": []any{openCodeMCPKey}}
}

// openCodeMCPServer is the opencode MCP entry rsk registers.
func openCodeMCPServer() map[string]any {
	return map[string]any{
		mcpTypeKey:    mcpLocalType,
		mcpCommandKey: []any{rskMCPServerName, openCodeMCPKey, "--read-only"},
		"enabled":     true,
	}
}

func legacyOpenCodeMCPServer() map[string]any {
	return map[string]any{
		mcpTypeKey:    mcpLocalType,
		mcpCommandKey: []any{rskMCPServerName, openCodeMCPKey},
		"enabled":     true,
	}
}

// RegisterMCP writes the rsk server into the project's .mcp.json, preserving
// other servers and keys. It is idempotent: a matching entry is left untouched.
// It returns the path written.
func (*ClaudeTool) RegisterMCP(projectDir string) (string, error) {
	return registerMCPServer(
		filepath.Join(projectDir, MCPFileName), "mcpServers",
		claudeMCPServer(), readClaudeConfig, writeClaudeConfig,
	)
}

// UnregisterMCP removes the rsk server from .mcp.json, leaving others intact.
func (*ClaudeTool) UnregisterMCP(projectDir string) error {
	return removeMCPServer(
		filepath.Join(projectDir, MCPFileName), "mcpServers",
		claudeMCPServer(), readClaudeConfig, writeClaudeConfig,
	)
}

// RegisterMCP writes the rsk server into the project's opencode.json, preserving
// other servers and keys. It is idempotent and returns the path written.
func (*openCodeTool) RegisterMCP(projectDir string) (string, error) {
	if err := rejectAlternateOpenCodeConfigs(projectDir); err != nil {
		return "", err
	}
	path := filepath.Join(projectDir, OpenCodeFileName)
	cfg, err := readOCConfig(path)
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
		openCodeMCPServer(), readOCConfig, writeOCConfig,
	)
}

// UnregisterMCP removes the rsk server from opencode.json, leaving others intact.
func (*openCodeTool) UnregisterMCP(projectDir string) error {
	if err := rejectAlternateOpenCodeConfigs(projectDir); err != nil {
		return err
	}
	path := filepath.Join(projectDir, OpenCodeFileName)
	cfg, err := readOCConfig(path)
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
		openCodeMCPServer(), readOCConfig, writeOCConfig,
	)
}

func rejectAlternateOpenCodeConfigs(projectDir string) error {
	for _, name := range []string{"opencode.jsonc", ".opencode/opencode.json", ".opencode/opencode.jsonc"} {
		path := filepath.Join(projectDir, filepath.FromSlash(name))
		if _, err := os.Lstat(path); err == nil {
			return fmt.Errorf("%s: alternate OpenCode config requires manual MCP setup", path)
		} else if !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func openCodeV2MCPServer() map[string]any {
	return map[string]any{
		mcpTypeKey:    mcpLocalType,
		mcpCommandKey: []any{rskMCPServerName, openCodeMCPKey, "--read-only"},
	}
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
	case ClaudeID:
		return MCPFileName, nil
	case CodexID:
		name := filepath.Join(codexFolder, codexConfig)
		_, err := readCodexConfig(filepath.Join(projectDir, name))
		return name, err
	case OpenCodeID:
		return OpenCodeFileName, rejectAlternateOpenCodeConfigs(projectDir)
	default:
		return "", fmt.Errorf("MCP preflight is not implemented for %q", id)
	}
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
	return path, writeOCConfig(path, cfg)
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
	return writeOCConfig(path, cfg)
}

func registerMCPServer(
	path, key string,
	server map[string]any,
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
	if existing, exists := servers[rskMCPServerName]; exists &&
		(key != openCodeMCPKey || !reflect.DeepEqual(existing, legacyOpenCodeMCPServer())) {
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
	if !reflect.DeepEqual(existing, managed) &&
		(key != openCodeMCPKey || !reflect.DeepEqual(existing, legacyOpenCodeMCPServer())) {
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
