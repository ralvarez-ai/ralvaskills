package permission

import (
	"slices"

	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
)

// RevokeSkills withdraws the rules rsk wrote for names and drops their lock
// entries. A rule another skill still relies on stays: it belongs to more than
// one skill only while both are installed.
func RevokeSkills(
	lock manifest.Lock, projectRoot string, tools []tool.ID, names []string,
) (manifest.Lock, []tool.Rule, error) {
	if len(names) == 0 {
		return lock, nil, nil
	}

	candidates := make(map[string]struct{})
	for _, name := range names {
		entry, ok := manifest.FindLockEntry(lock, name)
		if !ok {
			continue
		}
		for _, key := range entry.Granted {
			candidates[key] = struct{}{}
		}
	}

	next := lock
	for _, name := range names {
		next = manifest.RemoveLockEntry(next, name)
	}

	rules := tool.ParseRules(withoutOwned(candidates, next))
	if len(rules) == 0 {
		return next, nil, nil
	}

	writers := permissionWriters(tools)
	for _, w := range writers {
		if err := w.RevokePermissions(projectRoot, rules); err != nil {
			return lock, nil, err
		}
	}
	return next, rules, nil
}

// withoutOwned drops the keys another lock entry still owns.
func withoutOwned(keys map[string]struct{}, lock manifest.Lock) []string {
	kept := make([]string, 0, len(keys))
	for key := range keys {
		if !lockOwns(lock, key) {
			kept = append(kept, key)
		}
	}
	return kept
}

// lockOwns reports whether any entry in lock grants key.
func lockOwns(lock manifest.Lock, key string) bool {
	for _, entry := range lock.Skills {
		if slices.Contains(entry.Granted, key) {
			return true
		}
	}
	return false
}
