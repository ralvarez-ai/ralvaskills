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

	// ClaudeLocalFileName is Claude Code's personal, uncommitted instruction file
	// at the project root. rsk keeps one line in it, the import of the digest of
	// the pins from a private source, and leaves everything else alone.
	ClaudeLocalFileName = "CLAUDE.local.md"

	// ClaudeFileName is the name of the generated file inside .rsk/ that lists
	// pinned skills for Claude Code to import.
	ClaudeFileName = "CLAUDE.md"

	// Settings.json configuration keys
	claudeSettingsFileName = "settings.json"
	claudeSettingsPermKey  = "permissions"
	claudeSettingsAllowKey = "allow"
	claudeSettingsDenyKey  = "deny"

	claudeLocalImport   = "@" + internal.ProjectFolderName + "/" + digest.LocalIndexFileName
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
	if err := removeLocalLine(filepath.Join(projectDir, ClaudeLocalFileName)); err != nil {
		return err
	}
	return removeClaudeImport(filepath.Join(projectDir, ClaudeFileName))
}

// SyncLocal keeps the import of .rsk/SKILLS.local.md in CLAUDE.local.md while
// there are pins from a private source, and removes it, and the file if rsk's
// line was all it held, when there are none.
func (*ClaudeTool) SyncLocal(projectDir string, has bool) error {
	path := filepath.Join(projectDir, ClaudeLocalFileName)
	if !has {
		return removeLocalLine(path)
	}
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if hasLine(string(data), claudeLocalImport) {
		return nil
	}
	if err = refuseIfLink(path); err != nil {
		return err
	}
	text := string(data)
	if text != "" && !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	return fsx.WriteAtomic(path, rskTempPattern, func(w io.Writer) error {
		_, writeErr := io.WriteString(w, text+claudeLocalImport+"\n")
		return writeErr
	})
}

func hasLine(text, line string) bool {
	for l := range strings.SplitSeq(text, "\n") {
		if strings.TrimSpace(l) == line {
			return true
		}
	}
	return false
}

// removeLocalLine deletes rsk's import from CLAUDE.local.md, and the file when
// nothing else is left in it. A missing file is not an error.
func removeLocalLine(path string) error {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if !hasLine(string(data), claudeLocalImport) {
		return nil
	}
	if err = refuseIfLink(path); err != nil {
		return err
	}
	var kept []string
	for line := range strings.SplitSeq(string(data), "\n") {
		if strings.TrimSpace(line) != claudeLocalImport {
			kept = append(kept, line)
		}
	}
	rest := strings.Join(kept, "\n")
	if strings.TrimSpace(rest) == "" {
		return os.Remove(path)
	}
	return fsx.WriteAtomic(path, rskTempPattern, func(w io.Writer) error {
		_, writeErr := io.WriteString(w, rest)
		return writeErr
	})
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
