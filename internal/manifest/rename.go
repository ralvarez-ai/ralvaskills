package manifest

import (
	"cmp"
	"slices"
)

// Rename is one skill rename applied to a manifest or a lock.
type Rename struct {
	Old string
	New string
}

// RenameMod rewrites the skill keys and the pinned list of m through aliases and
// reports what moved, sorted by old name. When the new name is already present,
// its constraint wins and the old key is dropped: the new name is authoritative.
func RenameMod(m Mod, aliases map[string]string) (Mod, []Rename) {
	var renamed []Rename
	for old, next := range aliases {
		constraint, ok := m.Skills[old]
		if !ok {
			continue
		}
		if _, taken := m.Skills[next]; !taken {
			m.Skills[next] = constraint
		}
		delete(m.Skills, old)
		renamed = append(renamed, Rename{Old: old, New: next})
	}

	for i, name := range m.Pinned {
		if next, ok := aliases[name]; ok {
			m.Pinned[i] = next
		}
	}

	slices.SortFunc(renamed, func(a, b Rename) int {
		return cmp.Compare(a.Old, b.Old)
	})
	return m, renamed
}

// RenameLock rewrites the skill names of l through aliases and reports what
// moved, sorted by old name. An old entry is dropped when the new name already
// has one, so a lock never carries the same name twice.
func RenameLock(l Lock, aliases map[string]string) (Lock, []Rename) {
	present := make(map[string]bool, len(l.Skills))
	for _, entry := range l.Skills {
		present[entry.Name] = true
	}

	var renamed []Rename
	kept := make([]LockEntry, 0, len(l.Skills))
	for _, entry := range l.Skills {
		next, ok := aliases[entry.Name]
		if !ok {
			kept = append(kept, entry)
			continue
		}
		renamed = append(renamed, Rename{Old: entry.Name, New: next})
		if present[next] {
			continue
		}
		entry.Name = next
		present[next] = true
		kept = append(kept, entry)
	}

	slices.SortFunc(renamed, func(a, b Rename) int {
		return cmp.Compare(a.Old, b.Old)
	})
	l.Skills = kept
	return l, renamed
}
