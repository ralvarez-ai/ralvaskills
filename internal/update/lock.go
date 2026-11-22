package update

import (
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

// LinkAndLock links every resolvable skill into targets and updates its lock
// entry in place, keeping the version the skill resolved to.
func LinkAndLock(
	lock manifest.Lock, names, targets []string, resolve func(string) (skill.Skill, error),
) (manifest.Lock, []skill.Skill, error) {
	resolved := make([]skill.Skill, 0, len(names))
	for _, name := range names {
		s, err := resolve(name)
		if err != nil {
			continue
		}
		for _, target := range targets {
			if linkErr := skill.Link(s, target); linkErr != nil {
				return lock, nil, linkErr
			}
		}
		existing, found := manifest.FindLockEntry(lock, s.Name)
		if !found {
			existing = manifest.LockEntry{Name: s.Name}
		}
		existing.Version = s.Version
		existing.Source = s.Source
		existing.Git = s.Git
		lock = manifest.UpsertLockEntry(lock, existing)
		resolved = append(resolved, s)
	}
	return lock, resolved, nil
}

// refreshOfficialCache clones or pulls the official anthropics/skills cache.
