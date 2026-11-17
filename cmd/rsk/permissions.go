package main

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/spf13/cobra"

	"go.ralva.dev/lib/term"

	"go.ralva.dev/ralvaskills/v3/internal/cmdx"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/permission"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
	"go.ralva.dev/ralvaskills/v3/internal/ui"
)

// reconcilePermissions applies the skills' allowed-tools declarations to the
// project lock and to the clients that need a written rule, then writes the
// lock. It never fails the install or update that preceded it: a permission is
// best effort, and the skill is installed either way.
func reconcilePermissions(cmd *cobra.Command, rskDir string, skills []skill.Skill) {
	if len(skills) == 0 {
		return
	}
	lock, err := manifest.ReadLock(rskDir)
	if err != nil {
		ui.Warn(cmd.OutOrStdout(), fmt.Sprintf("tool permissions: %v", err))
		return
	}

	next, result := reconcileLock(cmd, rskDir, lock, skills)
	if result.Empty() {
		return
	}
	if err = manifest.WriteLock(rskDir, next); err != nil {
		ui.Warn(cmd.OutOrStdout(), fmt.Sprintf("tool permissions: %v", err))
		return
	}
	reportPermissions(cmd.OutOrStdout(), result)
}

// reconcileLock runs the reconciliation against an in-memory lock, for callers
// that write the lock themselves (rsk update).
func reconcileLock(
	cmd *cobra.Command, rskDir string, lock manifest.Lock, skills []skill.Skill,
) (manifest.Lock, permission.Result) {
	out := cmd.OutOrStdout()
	m, err := manifest.ReadMod(rskDir)
	if err != nil {
		ui.Warn(out, fmt.Sprintf("tool permissions: %v", err))
		return lock, permission.Result{}
	}

	yes := cmdx.Bool(cmd, cmdx.FlagYes)
	opts := permission.Options{
		ProjectRoot:   filepath.Dir(rskDir),
		Tools:         m.Tools,
		Skills:        skills,
		Mode:          permissionMode(cmd),
		RevokeRemoved: boolFlag(cmd, cmdx.FlagRevokeRemoved),
		Confirm: func(summary string) bool {
			return term.Confirm(cmd, yes, true, summary)
		},
	}

	next, result, err := permission.Reconcile(lock, opts)
	if err != nil {
		ui.Warn(out, fmt.Sprintf("tool permissions: %v", err))
		return lock, permission.Result{}
	}
	return next, result
}

// permissionMode decides whether this session may write tool permissions.
// Without a terminal there is nobody to ask, so the default is to skip: rsk
// never grants a permission unattended.
func permissionMode(cmd *cobra.Command) permission.Mode {
	switch {
	case boolFlag(cmd, cmdx.FlagApplyPermissions):
		return permission.ModeApply
	case boolFlag(cmd, cmdx.FlagSkipPermissions):
		return permission.ModeSkip
	case !term.StdioIsTerminal():
		return permission.ModeSkip
	default:
		return permission.ModeAsk
	}
}

// boolFlag reads an optional bool flag: commands that do not register it (and
// the legacy paths that share this code) read false instead of panicking.
func boolFlag(cmd *cobra.Command, name string) bool {
	if cmd.Flags().Lookup(name) == nil {
		return false
	}
	return cmdx.Bool(cmd, name)
}

// printPermissionsPreview lists the permissions the skills declare, before
// anything is written. It is what a dry run exists to show.
func printPermissionsPreview(out io.Writer, skills []skill.Skill) {
	var declared []permission.SkillRules
	for _, s := range skills {
		entries := skill.NormalizeAllowedTools(s.AllowedTools)
		if len(entries) == 0 {
			continue
		}
		declared = append(declared, permission.SkillRules{Skill: s.Name, Entries: entries})
	}
	if len(declared) == 0 {
		return
	}

	ui.Header(out, "Permissions declared:")
	for _, row := range declared {
		for _, entry := range row.Entries {
			fmt.Fprintf(out, "  %s  %s\n", ui.SkillName(row.Skill), entry)
		}
	}
	fmt.Fprintln(out)
}

// reportPermissions prints what rsk wrote, skipped, or left behind.
func reportPermissions(out io.Writer, result permission.Result) {
	for _, rule := range result.Added {
		ui.Success(out, fmt.Sprintf("allow  %s  %s", rule.Action, ui.MutedPath(rule.Resource)))
	}
	for _, warning := range result.Warnings {
		ui.Warn(out, warning)
	}
	if result.Skipped {
		ui.Warn(out, "tool permissions not applied — re-run with --apply-permissions, or add them by hand")
	}
	if result.Revoked {
		ui.Info(out, fmt.Sprintf("withdrew %d permission(s) no longer declared", len(result.Removed)))
		return
	}
	if result.Skipped || len(result.Added) > 0 {
		return
	}
	if len(result.Removed) > 0 {
		ui.Warn(out, fmt.Sprintf(
			"%d permission(s) are no longer declared; run 'rsk update --revoke-removed' to withdraw them",
			len(result.Removed)))
		return
	}
	for _, row := range result.Declared {
		for _, entry := range row.Entries {
			ui.Info(out, fmt.Sprintf("%s: %s (the client applies this itself)", ui.SkillName(row.Skill), entry))
		}
	}
}

// reportRevoked prints the permissions an uninstall withdrew.
func reportRevoked(out io.Writer, rules []tool.Rule) {
	for _, rule := range rules {
		ui.Success(out, fmt.Sprintf("withdrawn  %s  %s", rule.Action, ui.MutedPath(rule.Resource)))
	}
}
