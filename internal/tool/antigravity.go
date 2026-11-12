package tool

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.ralva.dev/ralvaskills/v3/internal"
	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

const (
	// AntigravityID identifies the Google Antigravity skill target shared by the
	// Antigravity App, IDE and CLI.
	AntigravityID ID = "antigravity"

	// AntigravityMCPFileName is the project file Antigravity reads for MCP
	// servers. It lives inside the shared .agents/ workspace folder.
	AntigravityMCPFileName = "mcp_config.json"

	// antigravityHomeFolder is the home-relative folder Antigravity reads for
	// personal customizations. It ignores XDG_CONFIG_HOME, so this path must not
	// go through internal.ConfigHome().
	antigravityHomeFolder = ".gemini"
	antigravityConfigDir  = "config"
)

// antigravityTool implements Tool for Google Antigravity.
type antigravityTool struct{}

// ID returns AntigravityID.
func (*antigravityTool) ID() ID {
	return AntigravityID
}

// SkillsDir returns the global skills directory read by the Antigravity App,
// IDE and CLI (~/.gemini/config/skills).
func (*antigravityTool) SkillsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(
		home, antigravityHomeFolder, antigravityConfigDir, skill.SkillsFolderName,
	)
}

// ProjectSkillsDir returns Antigravity's workspace skills directory
// (.agents/skills). Antigravity scans it one level deep, so each skill folder
// sits directly inside it, and Codex discovers the same directory.
func (*antigravityTool) ProjectSkillsDir(projectRoot string) string {
	return filepath.Join(projectRoot, internal.AgentsFolderName, skill.SkillsFolderName)
}

// SyncPinned is a no-op: Antigravity discovers skills on demand and reads
// AGENTS.md itself, so there is no generated import for rsk to write.
func (*antigravityTool) SyncPinned(string, []string) error {
	return nil
}

// RemovePinned is a no-op: rsk never writes an Antigravity skill import.
func (*antigravityTool) RemovePinned(string) error {
	return nil
}

// antigravityMCPServer is the Antigravity MCP entry rsk registers. The
// read-only server omits install_skills until Antigravity's per-tool approval
// behavior is verified, as with Codex and OpenCode.
func antigravityMCPServer() map[string]any {
	return map[string]any{
		mcpCommandKey: rskMCPServerName,
		"args":        []any{"mcp", mcpReadOnlyFlag},
	}
}

// antigravityMCPPath returns the project MCP config Antigravity reads.
func antigravityMCPPath(projectDir string) string {
	return filepath.Join(projectDir, internal.AgentsFolderName, AntigravityMCPFileName)
}

// RegisterMCP writes the rsk server into the project's
// .agents/mcp_config.json, preserving other servers and keys. It is idempotent
// and returns the path written.
func (*antigravityTool) RegisterMCP(projectDir string) (string, error) {
	dir := filepath.Join(projectDir, internal.AgentsFolderName)
	if err := refuseIfLinkDir(dir); err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, fsperm.Dir); err != nil {
		return "", err
	}
	return registerMCPServer(
		antigravityMCPPath(projectDir), "mcpServers",
		antigravityMCPServer(), readClaudeConfig, writeClaudeConfig,
	)
}

// UnregisterMCP removes the rsk server from .agents/mcp_config.json, leaving
// other servers intact.
func (*antigravityTool) UnregisterMCP(projectDir string) error {
	if err := refuseIfLinkDir(filepath.Join(projectDir, internal.AgentsFolderName)); err != nil {
		return err
	}
	return removeMCPServer(
		antigravityMCPPath(projectDir), "mcpServers",
		antigravityMCPServer(), readClaudeConfig, writeClaudeConfig,
	)
}

// refuseIfLinkDir rejects a symlinked configuration directory so rsk does not
// write the MCP config outside the project.
func refuseIfLinkDir(dir string) error {
	info, err := os.Lstat(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("stat %s: %w", dir, err)
	}
	if skill.IsLink(info) {
		return fmt.Errorf("%s is a symlink (→ %s): refusing to write the MCP config through it",
			dir, linkTarget(info))
	}
	return nil
}
