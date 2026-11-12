package main

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
	"go.ralva.dev/ralvaskills/v3/internal/ui"
)

// applyRenames runs the rename migration that matches the scope and re-links
// what moved. Failures warn rather than abort: a migration must not block the
// update that already succeeded.
func applyRenames(
	cmd *cobra.Command, cfg config.Config, reg *source.Registry, index source.Index,
	targets []string, opts updateOpts, out, errOut io.Writer,
) {
	if opts.global {
		migrateErr := migrateGlobalRenames(cmd.Context(), out, reg, index.Aliases, targets, opts.dryRun)
		if migrateErr != nil {
			ui.Warn(out, fmt.Sprintf("global rename migration: %v", migrateErr))
		}
		return
	}

	renamed, migrateErr := migrateProjectRenames(out, index.Aliases, opts.dryRun)
	if migrateErr != nil {
		ui.Warn(out, fmt.Sprintf("rename migration: %v", migrateErr))
	}
	if len(renamed) == 0 || opts.dryRun {
		return
	}
	if lockErr := bumpLockForProject(cmd, nil, cfg); lockErr != nil {
		ui.Warn(out, fmt.Sprintf("re-link renamed skills: %v", lockErr))
	}
	refreshPinnedDigest(out, errOut)
}

// migrateProjectRenames rewrites rsk.mod and rsk.lock through the registry's
// alias table and removes the links the old names left behind. The new links
// are created by the caller's normal resolve-and-link pass; nothing is written
// in dry-run mode.
func migrateProjectRenames(out io.Writer, aliases map[string]string, dryRun bool) ([]manifest.Rename, error) {
	if len(aliases) == 0 {
		return nil, nil
	}
	rskDir, ok := tryRskDir()
	if !ok {
		return nil, nil
	}

	m, err := manifest.ReadMod(rskDir)
	if err != nil {
		return nil, err
	}
	lock, err := manifest.ReadLock(rskDir)
	if err != nil {
		return nil, err
	}

	m, fromMod := manifest.RenameMod(m, aliases)
	lock, fromLock := manifest.RenameLock(lock, aliases)
	renamed := mergeRenames(fromMod, fromLock)
	if len(renamed) == 0 {
		return nil, nil
	}

	for _, r := range renamed {
		if dryRun {
			ui.Info(out, fmt.Sprintf("would rename %s → %s (rsk.mod, rsk.lock, .rsk/skills)", r.Old, r.New))
			continue
		}
		ui.Info(out, fmt.Sprintf("renamed %s → %s", r.Old, r.New))
	}
	if dryRun {
		return renamed, nil
	}

	if err = manifest.WriteMod(rskDir, m); err != nil {
		return nil, err
	}
	if err = manifest.WriteLock(rskDir, lock); err != nil {
		return nil, err
	}
	for _, target := range projectSkillsDirs(filepath.Dir(rskDir), m) {
		for _, r := range renamed {
			if unlinkErr := skill.Unlink(r.Old, target); unlinkErr != nil {
				return nil, unlinkErr
			}
		}
	}
	return renamed, nil
}

// migrateGlobalRenames points a machine-wide link of a retired name at the
// skill that replaced it.
func migrateGlobalRenames(
	ctx context.Context, out io.Writer, reg *source.Registry,
	aliases map[string]string, targets []string, dryRun bool,
) error {
	for _, target := range targets {
		for old, next := range aliases {
			if !skill.IsLinked(old, target) {
				continue
			}
			if dryRun {
				ui.Info(out, fmt.Sprintf("would rename %s → %s in %s", old, next, target))
				continue
			}
			s, findErr := reg.FindVersion(ctx, next, "")
			if findErr != nil {
				ui.Warn(out, fmt.Sprintf("rename %s: %v", old, findErr))
				continue
			}
			if linkErr := skill.Link(s, target); linkErr != nil {
				return linkErr
			}
			if unlinkErr := skill.Unlink(old, target); unlinkErr != nil {
				return unlinkErr
			}
			ui.Success(out, fmt.Sprintf("%s → %s", old, next))
		}
	}
	return nil
}

// mergeRenames unions the renames reported by the mod and the lock, which
// usually agree, and sorts them by old name.
func mergeRenames(mod, lock []manifest.Rename) []manifest.Rename {
	seen := make(map[string]bool, len(mod)+len(lock))
	merged := make([]manifest.Rename, 0, len(mod)+len(lock))
	for _, list := range [][]manifest.Rename{mod, lock} {
		for _, r := range list {
			if seen[r.Old] {
				continue
			}
			seen[r.Old] = true
			merged = append(merged, r)
		}
	}
	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Old < merged[j].Old
	})
	return merged
}

// localAliases reads the retired-name table from the local ralvaskills clone.
func localAliases(cfg config.Config) (map[string]string, error) {
	path := filepath.Join(cfg.RepoPath, skill.RetirementsPath)
	retired, err := skill.LoadRetirements(path)
	if err != nil {
		return nil, err
	}
	return retired, nil
}
