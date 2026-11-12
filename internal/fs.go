// Package internal provides internal utilities for the rsk CLI.
package internal

import (
	"path/filepath"
	"runtime"

	"github.com/adrg/xdg"
)

const (
	// ProjectFolderName is the name of the project-local rsk directory that
	// marks a directory as an rsk project root.
	ProjectFolderName = ".rsk"

	// ClaudeFolderName is Claude Code's configuration directory.
	ClaudeFolderName = ".claude"

	// AgentsFolderName contains agent skills shared with Codex.
	AgentsFolderName = ".agents"

	// OpenCodeProjectFolderName is OpenCode's project configuration directory.
	OpenCodeProjectFolderName = ".opencode"

	// TempFilePattern is the pattern for rsk's atomic-write temporary files.
	TempFilePattern = ".rsk-*.tmp"
)

// ConfigHome returns ~/.config on Windows instead of %APPDATA%, and falls back
// to configHome() on all other platforms.
func ConfigHome() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(xdg.Home, ".config")
	}
	return xdg.ConfigHome
}
