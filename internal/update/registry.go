package update

import (
	"context"
	"fmt"
	"os"

	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
)

// Pair is a skill with a newer version available in the registry.
type Pair struct {
	Name      string
	Installed string
	Latest    string
	Skill     skill.Skill
}

// Names extracts the skill names from pairs, for caching the update state when
// nothing was applied (a dry run or an aborted confirm).
func Names(pairs []Pair) []string {
	names := make([]string, len(pairs))
	for i := range pairs {
		names[i] = pairs[i].Name
	}
	return names
}

// LockedVersions reads the project's rsk.lock and returns each tracked skill's
// recorded version. It returns nil when the lock can't be read: the lock
// describes one project's resolved state, so there is nothing to compare.
func LockedVersions(rskDir string) map[string]string {
	lock, err := manifest.ReadLock(rskDir)
	if err != nil {
		return nil
	}
	versions := make(map[string]string, len(lock.Skills))
	for _, entry := range lock.Skills {
		versions[entry.Name] = entry.Version
	}
	return versions
}

// LinkedNames collects every skill name currently linked in any target dir.
func LinkedNames(targets []string) []string {
	seen := make(map[string]bool)
	var names []string
	for _, target := range targets {
		entries, err := os.ReadDir(target)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !seen[e.Name()] {
				seen[e.Name()] = true
				names = append(names, e.Name())
			}
		}
	}
	return names
}

// WithLockedNames appends the skills the project's lock records but no target
// dir links, so an entry whose link went missing is repaired as well.
func WithLockedNames(names []string, locked map[string]string) []string {
	if len(locked) == 0 {
		return names
	}
	seen := make(map[string]bool, len(names))
	for _, name := range names {
		seen[name] = true
	}
	for name := range locked {
		if !seen[name] {
			names = append(names, name)
		}
	}
	return names
}

// NeedsUpdate reports whether a skill is behind: the version linked in the
// target dirs, or the version the project's lock records, disagrees with the
// registry's latest. A skill absent from the lock is judged by its link alone,
// since a project need not track everything it links.
func NeedsUpdate(installed, locked, latest string) bool {
	if installed != latest {
		return true
	}
	return locked != "" && locked != latest
}

// FindOutdated checks each name against index, skipping any that aren't in the
// registry, aren't actually behind, or fail to resolve at their latest version
// (a warning, not an error).
func FindOutdated(
	ctx context.Context, reg *source.Registry, index map[string]*source.IndexEntry,
	names, targets []string, registryCache string, locked map[string]string,
) (pairs []Pair, warnings []string) {
	for _, name := range names {
		entry, ok := index[name]
		if !ok {
			continue
		}
		// Determine installed version from symlink target path.
		installedVer := InstalledVersionFromTargets(name, targets, registryCache)
		lockedVer := locked[name]
		if !NeedsUpdate(installedVer, lockedVer, entry.Latest) {
			continue
		}
		s, err := reg.FindVersion(ctx, name, entry.Latest)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("skip %s: %v", name, err))
			continue
		}
		// When the link is already current, the lock is what is behind; report
		// that version so the plan shows what will be rewritten.
		reported := installedVer
		if reported == entry.Latest {
			reported = lockedVer
		}
		pairs = append(pairs, Pair{Name: name, Installed: reported, Latest: entry.Latest, Skill: s})
	}
	return pairs, warnings
}

// RelinkResult is one attempt to re-link a pair into one target directory.
type RelinkResult struct {
	Pair
	Target string
	Err    error
}

// Relink re-links each pair wherever it's currently linked, and reports every
// attempt, so a failure does not stop the rest of the batch.
func Relink(pairs []Pair, targets []string) []RelinkResult {
	var results []RelinkResult
	for i := range pairs {
		p := &pairs[i]
		for _, target := range targets {
			if !skill.IsLinked(p.Name, target) {
				continue
			}
			results = append(results, RelinkResult{Pair: *p, Target: target, Err: skill.Link(p.Skill, target)})
		}
	}
	return results
}
