package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/ui"
)

func bumpLockForProject(cmd *cobra.Command, args []string, cfg config.Config) error {
	ctx := cmd.Context()

	rskDir, ok := tryRskDir()
	if !ok {
		return nil
	}

	catalog, warn := config.LoadCatalog("")
	if warn != nil {
		ui.Warn(cmd.OutOrStdout(), fmt.Sprintf("user catalog: %v", warn))
	}

	names, err := skillNamesFromArgs(args, catalog)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		m, modOK := tryReadMod(rskDir)
		if !modOK {
			return nil
		}
		for name := range m.Skills {
			names = append(names, name)
		}
		if len(names) == 0 {
			return nil
		}
	}

	chain, err := cfg.Sources()
	if err != nil {
		return err
	}

	m, err := manifest.ReadMod(rskDir)
	if err != nil {
		return err
	}
	targets := projectSkillsDirs(filepath.Dir(rskDir), m)

	lock, err := manifest.ReadLock(rskDir)
	if err != nil {
		return err
	}
	lock, resolved, err := linkAndLock(lock, names, targets, func(name string) (skill.Skill, error) {
		return findSkillByName(ctx, name, chain)
	})
	if err != nil {
		return err
	}

	// A newer version may declare a permission the older one did not, or drop
	// one it no longer needs; reconcile before the lock is written so both the
	// client config and the record describe the same version.
	lock, result := reconcileLock(cmd, rskDir, lock, resolved)
	if !result.Empty() {
		reportPermissions(cmd.OutOrStdout(), result)
	}
	return manifest.WriteLock(rskDir, lock)
}

// linkAndLock links every resolvable skill into targets and updates its lock
// entry in place, keeping the version the skill resolved to.
func linkAndLock(
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
