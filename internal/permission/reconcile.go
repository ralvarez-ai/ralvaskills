// Package permission reconciles the tool permissions a skill declares in its
// allowed-tools frontmatter with the rules rsk has written into each client's
// config, and records both in rsk.lock.
package permission

import (
	"fmt"
	"slices"

	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
)

// Mode decides whether rsk writes permissions without asking.
type Mode int

const (
	// ModeAsk writes only when Confirm answers yes. This is the default.
	ModeAsk Mode = iota

	// ModeApply writes without asking, for --apply-permissions.
	ModeApply

	// ModeSkip reports what the skill declares and writes nothing.
	ModeSkip
)

// Options configures Reconcile.
type Options struct {
	ProjectRoot string
	Tools       []tool.ID
	Skills      []skill.Skill
	Mode        Mode

	// RevokeRemoved also withdraws the rules a newer version no longer
	// declares. Off by default: a withdrawal is never silent.
	RevokeRemoved bool

	// Confirm asks the user. A nil Confirm behaves like a decline, which is
	// what a session without a terminal needs.
	Confirm func(summary string) bool
}

// SkillRules is one skill's declaration, as reported in the install preview.
type SkillRules struct {
	Skill   string
	Entries []string
}

// Result reports what Reconcile found and did.
type Result struct {
	// Declared lists every skill that declares allowed-tools, in input order.
	Declared []SkillRules

	// Added is the rules rsk wrote into a client config.
	Added []tool.Rule

	// Removed is the rules a newer version no longer declares. They are
	// revoked only when RevokeRemoved is set, and reported otherwise.
	Removed []tool.Rule

	// Revoked is true when Removed was actually withdrawn from a client.
	Revoked bool

	// Skipped is true when a declaration was not written because the user
	// declined or there was nobody to ask.
	Skipped bool

	// Warnings lists entries no client could express, such as MCP tool names.
	Warnings []string

	// Changed is true when the lock records something different, which is what
	// the caller needs in order to write it.
	Changed bool
}

// Empty reports whether Reconcile found nothing to record or report. A
// withdrawal offer counts as something to report even when the lock is already
// in its final state.
func (r Result) Empty() bool {
	return !r.Changed && !r.Skipped && len(r.Removed) == 0 && len(r.Warnings) == 0
}

// Reconcile compares each skill's allowed-tools declaration with the lock,
// writes the additions into every client config that needs them, optionally
// revokes what a newer version dropped, and returns the updated lock. The
// caller writes the lock; nothing else in the project is touched.
func Reconcile(lock manifest.Lock, opts Options) (manifest.Lock, Result, error) {
	var result Result

	declared := make(map[string][]string, len(opts.Skills))
	for _, s := range opts.Skills {
		if err := skill.ValidateAllowedTools(s.AllowedTools); err != nil {
			return lock, result, fmt.Errorf("%s: allowed-tools: %w", s.Name, err)
		}
		// An empty declaration still enters the diff: that is how a newer
		// version's withdrawal is noticed.
		entries := skill.NormalizeAllowedTools(s.AllowedTools)
		declared[s.Name] = entries
		if len(entries) > 0 {
			result.Declared = append(result.Declared, SkillRules{Skill: s.Name, Entries: entries})
		}
	}
	if len(opts.Skills) == 0 {
		return lock, result, nil
	}

	addedRules, removedEntries := diffLock(lock, declared)
	result.Added = addedRules

	writers := permissionWriters(opts.Tools)
	if len(writers) == 0 {
		// Nobody to write for: Claude Code applies the frontmatter itself, so
		// only the record changes and a withdrawal has no writer to act on.
		result.Added = nil
		next := recordDeclarations(lock, declared, nil, nil)
		result.Changed = changedDeclarations(lock, next)
		return next, result, nil
	}

	var written []tool.Rule
	if len(addedRules) > 0 {
		if !consent(opts, result) {
			result.Skipped = true
			return lock, result, nil
		}
		for _, w := range writers {
			got, err := w.GrantPermissions(opts.ProjectRoot, addedRules)
			if err != nil {
				return lock, result, err
			}
			written = append(written, got...)
		}
		result.Added = tool.SortRules(written)
	}

	revoked, err := revokeRemoved(opts, writers, removedEntries, declared, lock, &result)
	if err != nil {
		return lock, result, err
	}
	next := recordDeclarations(lock, declared, written, revoked)
	result.Changed = len(written) > 0 || len(revoked) > 0 || changedDeclarations(lock, next)
	return next, result, nil
}

// revokeWithdrawals withdraws the rules a skill no longer declares. It runs
// only when explicitly asked, and skips a rule another skill still declares.
func revokeRemoved(
	opts Options, writers []tool.PermissionWriter, removed []string,
	declared map[string][]string, lock manifest.Lock, result *Result,
) ([]tool.Rule, error) {
	if len(removed) == 0 {
		return nil, nil
	}
	rules := tool.SortRules(translateAll(removed, result))
	rules = withoutRules(rules, stillDeclaredKeys(declared, lock))
	if len(rules) == 0 {
		return nil, nil
	}
	result.Removed = rules
	if !opts.RevokeRemoved {
		return nil, nil
	}
	for _, w := range writers {
		if err := w.RevokePermissions(opts.ProjectRoot, rules); err != nil {
			return nil, err
		}
	}
	result.Revoked = true
	return rules, nil
}

// consent decides whether to write. ModeAsk with no Confirm is a decline, so a
// piped or headless session never changes permissions.
func consent(opts Options, result Result) bool {
	if opts.Mode == ModeApply {
		return true
	}
	if opts.Mode == ModeSkip || opts.Confirm == nil {
		return false
	}
	summary := fmt.Sprintf("Apply %d tool permission(s)?", len(result.Added))
	return opts.Confirm(summary)
}

// diffLock returns the rules to write and the entries that appeared or
// disappeared since the locked version.
func diffLock(lock manifest.Lock, declared map[string][]string) (addedRules []tool.Rule, removed []string) {
	for name, entries := range declared {
		entry, ok := manifest.FindLockEntry(lock, name)
		var previous []string
		if ok {
			previous = entry.Permissions
		}
		for _, newEntry := range entries {
			if !slices.Contains(previous, newEntry) {
				addedRules = append(addedRules, translateAll([]string{newEntry}, nil)...)
			}
		}
		for _, oldEntry := range previous {
			if !slices.Contains(entries, oldEntry) {
				removed = append(removed, oldEntry)
			}
		}
	}
	return tool.SortRules(addedRules), removed
}

// recordDeclarations stores what each skill declares and which rules rsk wrote,
// so the next update can diff and a later uninstall can withdraw exactly those.
func recordDeclarations(
	lock manifest.Lock, declared map[string][]string, written, revoked []tool.Rule,
) manifest.Lock {
	for name, entries := range declared {
		entry, _ := manifest.FindLockEntry(lock, name)
		entry.Name = name
		entry.Permissions = entries
		entry.Granted = mergeGranted(entry.Granted, written, revoked)
		lock = manifest.UpsertLockEntry(lock, entry)
	}
	return lock
}

// changedDeclarations reports whether before and after disagree on what any
// skill declares or what rsk granted it.
func changedDeclarations(before, after manifest.Lock) bool {
	for _, entry := range after.Skills {
		old, ok := manifest.FindLockEntry(before, entry.Name)
		if !ok || !slices.Equal(old.Permissions, entry.Permissions) || !slices.Equal(old.Granted, entry.Granted) {
			return true
		}
	}
	return len(before.Skills) != len(after.Skills)
}

// mergeGranted adds written rule keys and drops revoked ones.
func mergeGranted(existing []string, written, revoked []tool.Rule) []string {
	kept := make(map[string]struct{}, len(existing)+len(written))
	for _, key := range existing {
		kept[key] = struct{}{}
	}
	for _, rule := range written {
		kept[rule.Key()] = struct{}{}
	}
	for _, rule := range revoked {
		delete(kept, rule.Key())
	}
	if len(kept) == 0 {
		return nil
	}
	keys := make([]string, 0, len(kept))
	for key := range kept {
		keys = append(keys, key)
	}
	return tool.RuleKeys(tool.ParseRules(keys))
}

// stillDeclaredKeys returns the rule keys that remain in someone's hands: every
// declaration in this run, plus the grants of the skills the run does not touch.
func stillDeclaredKeys(declared map[string][]string, lock manifest.Lock) []string {
	keys := declaredRuleKeys(declared)
	for _, entry := range lock.Skills {
		if _, inRun := declared[entry.Name]; inRun {
			continue
		}
		keys = append(keys, entry.Granted...)
	}
	return keys
}

// declaredRuleKeys returns the rule keys every declaration in declared maps to.
func declaredRuleKeys(declared map[string][]string) []string {
	keys := make([]string, 0, len(declared))
	for _, entries := range declared {
		keys = append(keys, tool.RuleKeys(translateAll(entries, nil))...)
	}
	return keys
}

// withoutRules drops the rules whose key is in keys.
func withoutRules(rules []tool.Rule, keys []string) []tool.Rule {
	if len(keys) == 0 {
		return rules
	}
	kept := make([]tool.Rule, 0, len(rules))
	for _, rule := range rules {
		if !slices.Contains(keys, rule.Key()) {
			kept = append(kept, rule)
		}
	}
	return kept
}

// translateAll converts every entry, recording the ones the client cannot
// express in result when it is not nil.
func translateAll(entries []string, result *Result) []tool.Rule {
	rules := make([]tool.Rule, 0, len(entries))
	for _, entry := range entries {
		rule, ok := tool.Translate(entry)
		if !ok {
			if result != nil {
				result.Warnings = append(result.Warnings,
					entry+": no equivalent in this client, skipped")
			}
			continue
		}
		rules = append(rules, rule)
	}
	return rules
}

// permissionWriters returns the writers for opts.Tools, in their configured
// order, skipping clients that apply the frontmatter themselves.
func permissionWriters(ids []tool.ID) []tool.PermissionWriter {
	writers := make([]tool.PermissionWriter, 0, len(ids))
	for _, id := range ids {
		if writer, ok := tool.WriterFor(id); ok {
			writers = append(writers, writer)
		}
	}
	return writers
}
