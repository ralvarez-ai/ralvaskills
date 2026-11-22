package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/install"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/ui"
	updatecheck "go.ralva.dev/ralvaskills/v3/internal/update"
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
	targets := install.ProjectSkillsDirs(filepath.Dir(rskDir), m)

	lock, err := manifest.ReadLock(rskDir)
	if err != nil {
		return err
	}
	lock, resolved, err := updatecheck.LinkAndLock(lock, names, targets, func(name string) (skill.Skill, error) {
		return install.FindSkill(ctx, name, chain)
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
