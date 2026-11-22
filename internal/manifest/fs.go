package manifest

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.ralva.dev/ralvaskills/v3/internal"
)

const (
	// ModFileName is the filename of the project skill manifest.
	ModFileName = "rsk.mod"

	// LockFileName is the filename of the resolved skill lockfile.
	LockFileName = "rsk.lock"

	// LocalModFileName holds the skills from a private source. It is never
	// committed: rsk keeps it out of git and merges it with rsk.mod on read.
	LocalModFileName = "rsk.local.mod"

	// LocalLockFileName is the lock for those skills, kept the same way.
	LocalLockFileName = "rsk.local.lock"

	// TempFilePattern is the glob pattern for temporary files created during
	// atomic write operations inside the .rsk directory.
	TempFilePattern = internal.TempFilePattern
)

// ProjectFolderPath returns the .rsk directory path for the current working
// directory. Returns an error if rsk.mod is not found, indicating the directory
// has not been initialized as an rsk project.
func ProjectFolderPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	rskDir := filepath.Join(cwd, internal.ProjectFolderName)
	if _, statErr := os.Stat(ModPath(rskDir)); os.IsNotExist(statErr) {
		return "", errors.New("no rsk.mod found — run 'rsk new' first")
	}
	return rskDir, nil
}

// LockPath returns the canonical path of rsk.lock inside rskDir.
func LockPath(rskDir string) string {
	return filepath.Join(rskDir, LockFileName)
}

// ModPath returns the canonical path of rsk.mod inside rskDir.
func ModPath(rskDir string) string {
	return filepath.Join(rskDir, ModFileName)
}

// LocalLockPath returns the canonical path of rsk.local.lock inside rskDir.
func LocalLockPath(rskDir string) string {
	return filepath.Join(rskDir, LocalLockFileName)
}

// LocalModPath returns the canonical path of rsk.local.mod inside rskDir.
func LocalModPath(rskDir string) string {
	return filepath.Join(rskDir, LocalModFileName)
}
