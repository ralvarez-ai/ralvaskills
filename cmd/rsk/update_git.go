package main

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"go.ralva.dev/lib/term"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/install"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
	"go.ralva.dev/ralvaskills/v3/internal/ui"
	updatecheck "go.ralva.dev/ralvaskills/v3/internal/update"
)

// gitProject is the project a git update acts on: its manifest directory and
// its lock, which is rewritten when a skill moves.
type gitProject struct {
	rskDir string
	lock   manifest.Lock
}

// printUpdateNotices shows the messages an update step produced, in order.
func printUpdateNotices(out io.Writer, notices []updatecheck.Notice) {
	for _, n := range notices {
		if n.Level == updatecheck.Warn {
			ui.Warn(out, n.Text)
			continue
		}
		ui.Info(out, n.Text)
	}
}

func printGitUpdatePlan(out io.Writer, plan []updatecheck.GitUpdate) {
	fmt.Fprintln(out)
	ui.Header(out, "Git sources")
	for i := range plan {
		u := &plan[i]
		fmt.Fprintf(out, "  %s  %s  %s  %s  (%s)\n", ui.SkillName(u.Install.Name),
			ui.SkillVersion(updatecheck.Short(u.Install.Commit)), ui.Arrow,
			ui.SkillVersion(updatecheck.Short(u.Next.Git.Commit)), u.Install.ID)
	}
	fmt.Fprintln(out)
}

// relinkGit updates each outdated skill and reports the outcome of each one as
// it would have been printed while updating. It returns what was resolved and
// how many failed.
func relinkGit(
	ctx context.Context, out, errOut io.Writer, chain source.Chain, plan []updatecheck.GitUpdate, targets []string,
) (resolved []skill.Skill, failed int) {
	results := updatecheck.RelinkGit(ctx, chain, plan, targets)
	for i := range results {
		r := &results[i]
		if r.FindErr != nil {
			ui.Failf(errOut, "update %s: %v", r.Update.Install.Name, r.FindErr)
			failed++
			continue
		}
		for _, linkErr := range r.LinkErrs {
			ui.Failf(errOut, "re-link %s: %v", r.Skill.Name, linkErr)
		}
		if r.Failed() {
			failed++
			continue
		}
		ui.Success(out, fmt.Sprintf("%s  %s  %s  %s", ui.SkillName(r.Skill.Name),
			ui.SkillVersion(updatecheck.Short(r.Update.Install.Commit)), ui.Arrow,
			ui.SkillVersion(updatecheck.Short(r.Skill.Git.Commit))))
		resolved = append(resolved, r.Skill)
	}
	return resolved, failed
}

// collectGitInstalls finds the git skills installed in the scope `rsk update`
// runs in: from the project's lock, or from where the global links point.
func collectGitInstalls(
	cfg config.Config, opts updateOpts, args, targets []string,
) (installs []updatecheck.GitInstall, project *gitProject, err error) {
	only := map[string]bool{}
	for _, raw := range args {
		name, _, parseErr := install.ParseNameVersion(raw)
		if parseErr != nil {
			return nil, nil, parseErr
		}
		only[name] = true
	}
	if opts.global {
		return updatecheck.GlobalGitInstalls(targets, cfg.SourceOptions().CacheRoot, cfg.OfficialCache, only), nil, nil
	}
	dir, ok := tryRskDir()
	if !ok {
		return nil, nil, nil
	}
	// A project whose manifest or lock cannot be read is left to the commands
	// that report that; this pass only adds work for git entries it can read.
	m, modErr := manifest.ReadMod(dir)
	lock, lockErr := manifest.ReadLock(dir)
	if modErr != nil || lockErr != nil {
		return nil, nil, nil //nolint:nilerr // see above
	}
	return updatecheck.ProjectGitInstalls(lock, m, only), &gitProject{rskDir: dir, lock: lock}, nil
}

// runGitUpdate updates the git skills installed in the scope `rsk update` runs
// in. It does nothing, and prints nothing, when there are none.
func runGitUpdate(cmd *cobra.Command, cfg config.Config, opts updateOpts, args []string) error {
	targets, err := install.Targets(cfg, opts.global, opts.forTool)
	if err != nil {
		return err
	}
	installs, project, err := collectGitInstalls(cfg, opts, args, targets)
	if err != nil || len(installs) == 0 {
		return err
	}

	out, ctx := cmd.OutOrStdout(), cmd.Context()
	chain, err := sourceChain(cmd, cfg)
	if err != nil {
		return err
	}
	planned := updatecheck.PlanGit(ctx, chain, installs)
	printUpdateNotices(out, planned.Notices)
	plan, failed := planned.Updates, planned.Failed
	if len(plan) == 0 {
		return gitUpdateResult(failed)
	}

	printGitUpdatePlan(out, plan)
	if opts.dryRun {
		return gitUpdateResult(failed)
	}
	if !term.Confirm(cmd, opts.yes, term.StdioIsTerminal(), "Update git skills?") {
		fmt.Fprintln(out, messageAborted)
		return gitUpdateResult(failed)
	}
	resolved, relinkFailed := relinkGit(ctx, out, cmd.ErrOrStderr(), chain, plan, targets)
	if project != nil && len(resolved) > 0 {
		if err = recordGitUpdates(cmd, project, resolved); err != nil {
			return err
		}
	}
	return gitUpdateResult(failed + relinkFailed)
}

// recordGitUpdates rewrites the project's lock for the skills that moved,
// reconciling the tool permissions their new version declares.
func recordGitUpdates(cmd *cobra.Command, project *gitProject, resolved []skill.Skill) error {
	lock := updatecheck.LockGitUpdates(project.lock, resolved)
	lock, result := reconcileLock(cmd, project.rskDir, lock, resolved)
	if !result.Empty() {
		reportPermissions(cmd.OutOrStdout(), result)
	}
	if err := manifest.WriteLock(project.rskDir, lock); err != nil {
		return err
	}
	refreshPinnedDigest(cmd.OutOrStdout(), cmd.ErrOrStderr())
	return nil
}

func gitUpdateResult(failed int) error {
	if failed > 0 {
		return fmt.Errorf("%d git skill(s) could not be checked or updated", failed)
	}
	return nil
}
