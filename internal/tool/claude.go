package tool

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"

	"go.ralva.dev/ralvaskills/v3/internal"
	"go.ralva.dev/ralvaskills/v3/internal/digest"
	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
	"go.ralva.dev/ralvaskills/v3/internal/fsx"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

const (
	// ClaudeID is the tool identifier for Claude Code.
	ClaudeID ID = "claude-code"

	// ClaudeFolderName is the name of the Claude Code configuration directory
	// inside the user's home directory or inside the project.
	ClaudeFolderName = internal.ClaudeFolderName

	// ClaudeFileName is the name of the generated file inside .rsk/ that lists
	// pinned skills for Claude Code to import.
	ClaudeFileName = "CLAUDE.md"

	// Settings.json configuration keys
	claudeSettingsFileName = "settings.json"
	claudeSettingsPermKey  = "permissions"
	claudeSettingsAllowKey = "allow"
	claudeSettingsDenyKey  = "deny"

	claudeImportLine    = "@" + internal.ProjectFolderName + "/" + ClaudeFileName
	claudeFileOpenFlags = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	rskTempPattern      = internal.TempFilePattern

	// claudeSkillsIndexReference is the single import line in .rsk/CLAUDE.md. The
	// path is relative to .rsk/ and points at the generated SKILLS.md, which
	// lives inside the project: Claude Code expands an @-import only when its
	// target resolves inside the project, and the old per-skill SKILL.md imports
	// went through symlinks into the machine cache, so they were skipped.
	claudeSkillsIndexReference = "@" + digest.IndexFileName + "\n"
)

var (
	// ProjectClaudeFilePath is the path of CLAUDE.md relative to the project
	// root (.rsk/CLAUDE.md).
	ProjectClaudeFilePath = filepath.Join(internal.ProjectFolderName, ClaudeFileName)
)

// ClaudeTool implements Tool for Claude Code.
type ClaudeTool struct{}

// ID returns ClaudeID.
func (*ClaudeTool) ID() ID {
	return ClaudeID
}

// SkillsDir returns the canonical path to the Claude Code global skills
// directory (~/.claude/skills).
func (*ClaudeTool) SkillsDir() string {
	return filepath.Join(xdg.Home, ClaudeFolderName, skill.SkillsFolderName)
}

// ProjectSkillsDir returns the project-local skills directory (.claude/skills/).
// Both Claude Code and OpenCode discover skills from this path.
func (*ClaudeTool) ProjectSkillsDir(projectRoot string) string {
	return filepath.Join(projectRoot, ClaudeFolderName, skill.SkillsFolderName)
}

// SyncPinned writes .rsk/CLAUDE.md atomically with a single import of the
// generated .rsk/SKILLS.md (empty when nothing is pinned; the callers write
// SKILLS.md itself), bootstraps the project's AGENTS.md, and imports it from CLAUDE.md
// ahead of the generated pin file (idempotent). projectDir is the project root
// (parent of .rsk/).
func (*ClaudeTool) SyncPinned(projectDir string, pinnedNames []string) error {
	if err := ensureAgentsFile(projectDir); err != nil {
		return fmt.Errorf("sync agent guide: %w", err)
	}

	rskDir := filepath.Join(projectDir, internal.ProjectFolderName)
	if err := writePinnedClaude(rskDir, pinnedNames); err != nil {
		return fmt.Errorf("sync claude-code pins: %w", err)
	}

	claudeMD := filepath.Join(projectDir, ClaudeFileName)
	if err := ensureAgentsImport(claudeMD); err != nil {
		return fmt.Errorf("sync agent guide import: %w", err)
	}
	if err := appendClaudeImport(claudeMD); err != nil {
		return fmt.Errorf("sync claude-code import: %w", err)
	}
	return nil
}

// RemovePinned removes the import marker from the project CLAUDE.md.
// Returns nil if the file does not exist.
func (*ClaudeTool) RemovePinned(projectDir string) error {
	return removeClaudeImport(filepath.Join(projectDir, ClaudeFileName))
}

func writePinnedClaude(rskDir string, pinnedNames []string) error {
	path := filepath.Join(rskDir, ClaudeFileName)
	if err := refuseIfLink(path); err != nil {
		return err
	}
	return fsx.WriteAtomic(path, rskTempPattern, func(w io.Writer) error {
		if len(pinnedNames) == 0 {
			return nil
		}
		_, err := io.WriteString(w, claudeSkillsIndexReference)
		return err
	})
}

func appendClaudeImport(claudeMDPath string) error {
	data, err := os.ReadFile(claudeMDPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read %s: %w", claudeMDPath, err)
	}
	if strings.Contains(string(data), claudeImportLine) {
		return nil
	}

	if err = refuseIfLink(claudeMDPath); err != nil {
		return err
	}

	f, err := os.OpenFile(claudeMDPath, claudeFileOpenFlags, fsperm.File)
	if err != nil {
		return fmt.Errorf("open %s: %w", claudeMDPath, err)
	}

	// A file that does not end in a newline would otherwise get the import
	// glued onto its last line.
	prefix := ""
	if len(data) > 0 && !strings.HasSuffix(string(data), "\n") {
		prefix = "\n"
	}

	_, writeErr := fmt.Fprintf(f, "%s%s\n", prefix, claudeImportLine)
	closeErr := f.Close()
	if writeErr != nil {
		return writeErr
	}
	return closeErr
}

// refuseIfLink fails when path is a symlink or junction rather than a regular
// file.
//
// rsk's two write paths disagree about links, and neither behavior is safe.
// appendClaudeImport opens with O_APPEND, which resolves the link and writes
// the import line into an unrelated target file. The WriteAtomic callers rename
// over the link, silently replacing it with a regular file. A link that points
// at rsk's own generated .rsk/CLAUDE.md is worse than either: the append lands
// in the generated file, which then imports itself.
//
// Failing loudly is deliberate. Someone who symlinked CLAUDE.md did it on
// purpose, and rsk cannot know whether the target wants an rsk import line, so
// it must not guess in either direction.
func refuseIfLink(path string) error {
	fi, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat %s: %w", path, err)
	}
	if skill.IsLink(fi) {
		return fmt.Errorf("%s is a symlink (→ %s): rsk will not write through it, "+
			"because that would modify the target file or replace the link. "+
			"Remove the link, or point CLAUDE.md at a real file, then retry",
			path, linkTarget(fi))
	}
	if !fi.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file (mode %s): refusing to write",
			path, fi.Mode())
	}
	return nil
}

// linkTarget returns a symlink's target for the error message, or "" when the
// mode bit came from a Windows junction, which has no readable target.
func linkTarget(fi os.FileInfo) string {
	target, err := os.Readlink(fi.Name())
	if err != nil {
		return "unknown"
	}
	return target
}

func removeClaudeImport(claudeMDPath string) error {
	data, err := os.ReadFile(claudeMDPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read %s: %w", claudeMDPath, err)
	}

	if err = refuseIfLink(claudeMDPath); err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	kept := lines[:0]
	for _, line := range lines {
		if strings.TrimSpace(line) != claudeImportLine {
			kept = append(kept, line)
		}
	}

	payload := []byte(strings.Join(kept, "\n"))
	return fsx.WriteAtomic(claudeMDPath, rskTempPattern, func(w io.Writer) error {
		_, writeErr := w.Write(payload)
		return writeErr
	})
}

// SettingsPath returns the path to the project's .claude/settings.json file.
func (*ClaudeTool) SettingsPath(projectDir string) string {
	return filepath.Join(projectDir, ClaudeFolderName, claudeSettingsFileName)
}

// ReadPermissions reads the tool permissions from .claude/settings.json.
// Returns empty slices if the file does not exist. Permissions is a map with
// "allow" and "deny" keys, each mapping to a []string of rules.
func (ct *ClaudeTool) ReadPermissions(projectDir string) (allow, deny []string, err error) {
	path := ct.SettingsPath(projectDir)
	cfg, err := readJSONConfig(path)
	if err != nil {
		return nil, nil, err
	}

	perms, ok := cfg[claudeSettingsPermKey].(map[string]any)
	if !ok {
		return nil, nil, nil
	}

	allow = toStringSlice(perms[claudeSettingsAllowKey])
	deny = toStringSlice(perms[claudeSettingsDenyKey])
	return allow, deny, nil
}

// WritePermissions writes tool permissions to .claude/settings.json, preserving
// other keys in the file. Creates the file if it does not exist.
func (ct *ClaudeTool) WritePermissions(projectDir string, allow, deny []string) error {
	path := ct.SettingsPath(projectDir)
	cfg, err := readJSONConfig(path)
	if err != nil {
		return err
	}

	perms := map[string]any{}
	if len(allow) > 0 {
		perms[claudeSettingsAllowKey] = allow
	}
	if len(deny) > 0 {
		perms[claudeSettingsDenyKey] = deny
	}

	if len(perms) == 0 {
		delete(cfg, claudeSettingsPermKey)
	} else {
		cfg[claudeSettingsPermKey] = perms
	}

	return writeJSONConfig(path, cfg)
}

func toStringSlice(v any) []string {
	if v == nil {
		return nil
	}
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(arr))
	for _, item := range arr {
		if s, stringOK := item.(string); stringOK {
			result = append(result, s)
		}
	}
	return result
}
