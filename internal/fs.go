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
)

// ConfigHome returns ~/.config on Windows instead of %APPDATA%, and falls back
// to configHome() on all other platforms.
func ConfigHome() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(xdg.Home, ".config")
	}
	return xdg.ConfigHome
}
