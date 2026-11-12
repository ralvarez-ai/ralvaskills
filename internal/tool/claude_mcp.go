package tool

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal/fsx"
)

const (
	// MCPFileName is the project file Claude Code reads for MCP servers.
	MCPFileName = ".mcp.json"

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
