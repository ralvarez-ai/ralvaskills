package tool

import (
	"path/filepath"
	"slices"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal"
	"go.ralva.dev/ralvaskills/v3/internal/digest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

const (
	// OpenCodeID is the tool identifier for OpenCode.
	OpenCodeID ID = "opencode"

	// OpenCodeFolderName is the name of the OpenCode configuration directory
	// inside the XDG config home.
	OpenCodeFolderName = "opencode"

	// OpenCodeFileName is the name of the OpenCode configuration file written
	// inside the project root by rsk.
	OpenCodeFileName = "opencode.json"

	// OpenCodeInstructionsKey is the key in opencode.json whose value is a list
	// of paths to SKILL.md files that OpenCode should import as instructions.
	OpenCodeInstructionsKey = "instructions"

	// OpenCodePermissionsKey is the key in opencode.json that holds the
	// permission rules, an array of {action, resource, effect} objects.
	OpenCodePermissionsKey = "permissions"

	// OpenCodeSkillsPrefix is the forward-slash path prefix of the per-skill
	// entries earlier versions wrote to opencode.json. They are no longer
	// written, only recognized so a re-sync removes them. Uses "/" rather than
	// filepath.Separator because opencode.json stores Unix-style paths.
	OpenCodeSkillsPrefix = ClaudeFolderName + "/" + skill.SkillsFolderName + "/"

	// OpenCodeSkillsIndexEntry is the one rsk-managed entry in the instructions
	// list: the generated .rsk/SKILLS.md, relative to the project root.
	OpenCodeSkillsIndexEntry = internal.ProjectFolderName + "/" + digest.IndexFileName
)

// openCodeTool implements Tool for OpenCode.
type openCodeTool struct{}

// ID returns OpenCodeID.
func (*openCodeTool) ID() ID {
	return OpenCodeID
}

// SkillsDir returns the canonical path to the OpenCode global skills directory
// (~/.config/opencode/skills).
func (*openCodeTool) SkillsDir() string {
	return filepath.Join(internal.ConfigHome(), OpenCodeFolderName, skill.SkillsFolderName)
}

// ProjectSkillsDir returns OpenCode's native project-local skills directory.
func (*openCodeTool) ProjectSkillsDir(projectRoot string) string {
	return filepath.Join(projectRoot, internal.OpenCodeProjectFolderName, skill.SkillsFolderName)
}

// SyncPinned updates opencode.json in projectDir so that the
// .claude/skills/<name>/SKILL.md entries match pinnedNames exactly. Non-rsk
// entries are preserved. Creates the file if it does not exist.
func (*openCodeTool) SyncPinned(projectDir string, pinnedNames []string) error {
	return updateOCInstructions(filepath.Join(projectDir, OpenCodeFileName), pinnedNames)
}

// RemovePinned removes all rsk-managed .claude/skills/ entries from opencode.json
// in projectDir. Returns nil if the file does not exist.
func (*openCodeTool) RemovePinned(projectDir string) error {
	return updateOCInstructions(filepath.Join(projectDir, OpenCodeFileName), nil)
}

// updateOCInstructions reads opencode.json, replaces the rsk-managed
// .claude/skills/ entries with pinnedNames, and writes the file back atomically.
// Pass a nil or empty pinnedNames to remove all rsk entries.
func updateOCInstructions(path string, pinnedNames []string) error {
	cfg, err := readJSONConfig(path)
	if err != nil {
		return err
	}

	kept := filterOCEntries(cfg[OpenCodeInstructionsKey])
	if len(pinnedNames) > 0 {
		kept = append(kept, OpenCodeSkillsIndexEntry)
	}

	if len(kept) == 0 {
		delete(cfg, OpenCodeInstructionsKey)
	} else {
		cfg[OpenCodeInstructionsKey] = kept
	}

	return writeJSONConfig(path, cfg)
}

// filterOCEntries returns the subset of raw that are not rsk-managed: neither
// the generated SKILLS.md entry nor the per-skill .claude/skills/ entries that
// earlier versions wrote. raw must be cfg[openCodeInstructionsKey].
func filterOCEntries(raw any) []any {
	existing, ok := raw.([]any)
	if !ok {
		return nil
	}
	kept := make([]any, 0, len(existing))
	for _, e := range existing {
		if s, isStr := e.(string); isStr &&
			(strings.HasPrefix(s, OpenCodeSkillsPrefix) || s == OpenCodeSkillsIndexEntry) {
			continue
		}
		kept = append(kept, e)
	}
	return kept
}

// GrantPermissions appends the rules opencode.json does not already carry and
// returns the ones it wrote. OpenCode applies the last matching rule, so the
// new entries go at the end, after any broad "ask" the user configured.
func (*openCodeTool) GrantPermissions(projectDir string, rules []Rule) ([]Rule, error) {
	path := filepath.Join(projectDir, OpenCodeFileName)
	cfg, err := readJSONConfig(path)
	if err != nil {
		return nil, err
	}

	existing := ocPermissionEntries(cfg[OpenCodePermissionsKey])
	added := make([]Rule, 0, len(rules))
	for _, rule := range rules {
		if ocHasRule(existing, rule) {
			continue
		}
		existing = append(existing, ocRuleEntry(rule))
		added = append(added, rule)
	}
	if len(added) == 0 {
		return nil, nil
	}

	cfg[OpenCodePermissionsKey] = existing
	if err = writeJSONConfig(path, cfg); err != nil {
		return nil, err
	}
	return added, nil
}

// RevokePermissions removes the rules from opencode.json, leaving every other
// entry in place. Nothing is written when no entry matches.
func (*openCodeTool) RevokePermissions(projectDir string, rules []Rule) error {
	if len(rules) == 0 {
		return nil
	}
	path := filepath.Join(projectDir, OpenCodeFileName)
	cfg, err := readJSONConfig(path)
	if err != nil {
		return err
	}

	existing := ocPermissionEntries(cfg[OpenCodePermissionsKey])
	kept := make([]any, 0, len(existing))
	removed := false
	for _, entry := range existing {
		if rule, ok := ocRuleFromEntry(entry); ok && slices.Contains(rules, rule) {
			removed = true
			continue
		}
		kept = append(kept, entry)
	}
	if !removed {
		return nil
	}

	if len(kept) == 0 {
		delete(cfg, OpenCodePermissionsKey)
	} else {
		cfg[OpenCodePermissionsKey] = kept
	}
	return writeJSONConfig(path, cfg)
}

// ocPermissionEntries returns the permission entries of raw as a list of
// objects, dropping anything malformed instead of failing the whole write.
func ocPermissionEntries(raw any) []any {
	existing, ok := raw.([]any)
	if !ok {
		return nil
	}
	kept := make([]any, 0, len(existing))
	for _, entry := range existing {
		if _, isMap := entry.(map[string]any); isMap {
			kept = append(kept, entry)
		}
	}
	return kept
}

// ocRuleEntry renders one rule as the object OpenCode reads.
func ocRuleEntry(rule Rule) map[string]any {
	return map[string]any{
		"action":   rule.Action,
		"resource": rule.Resource,
		"effect":   effectAllow,
	}
}

// ocRuleFromEntry reads an allow entry back into a rule. A deny or ask entry is
// not a rule rsk owns, so it does not match.
func ocRuleFromEntry(entry any) (Rule, bool) {
	obj, ok := entry.(map[string]any)
	if !ok {
		return Rule{}, false
	}
	action := ocStringField(obj, "action")
	resource := ocStringField(obj, "resource")
	effect := ocStringField(obj, "effect")
	if action == "" || resource == "" || effect != effectAllow {
		return Rule{}, false
	}
	return Rule{Action: action, Resource: resource}, true
}

// ocStringField reads a string field, returning "" when it is absent or not a
// string, so a malformed entry is skipped instead of failing the write.
func ocStringField(obj map[string]any, key string) string {
	value, ok := obj[key].(string)
	if !ok {
		return ""
	}
	return value
}

// ocHasRule reports whether the entries already allow rule.
func ocHasRule(entries []any, rule Rule) bool {
	for _, entry := range entries {
		if existing, ok := ocRuleFromEntry(entry); ok && existing == rule {
			return true
		}
	}
	return false
}
