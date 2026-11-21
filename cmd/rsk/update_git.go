package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"go.ralva.dev/lib/term"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/install"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
	"go.ralva.dev/ralvaskills/v3/internal/ui"
)

const (
	shortCommit  = 7
	gitLinkParts = 4
)

type (
	// gitInstall is an installed git skill and what is known of its content. A
	// project install has the tree its lock recorded; a global one has none and
	// is measured against its commit.
	gitInstall struct {
		name, id, commit, tree, constraint string
	}

	// gitUpdate is an installed skill whose directory changed at the source.
	gitUpdate struct {
		install gitInstall
		next    skill.Skill
	}

	// treeAter is implemented by a source that can say what tree a skill's
	// directory had at an earlier commit.
	treeAter interface {
		TreeAt(ctx context.Context, name, commit string) (string, error)
	}
)

func short(commit string) string {
	if len(commit) > shortCommit {
		return commit[:shortCommit]
	}
	return commit
}

// parseGitLink reads the identity, name and commit out of the target of a link
// into the git cache: <root>/sources/<owner>/<repo>/<name>/<commit>.
func parseGitLink(dest, cacheRoot string) (id, name, commit string, ok bool) {
	prefix := filepath.Join(cacheRoot, "sources") + string(filepath.Separator)
	rest, found := strings.CutPrefix(filepath.Clean(dest), prefix)
	if !found {
		return "", "", "", false
	}
	parts := strings.Split(rest, string(filepath.Separator))
	if len(parts) != gitLinkParts {
		return "", "", "", false
	}
	return "@" + parts[0] + "/" + parts[1], parts[2], parts[3], true
}

// globalGitInstalls finds the git skills linked into targets, by reading where
// each link points.
func globalGitInstalls(targets []string, cacheRoot string, only map[string]bool) []gitInstall {
	var installs []gitInstall
	seen := map[string]bool{}
	for _, target := range targets {
		entries, err := os.ReadDir(target)
		if err != nil {
			continue
		}
		for _, e := range entries {
			dest, linkErr := os.Readlink(filepath.Join(target, e.Name()))
			if linkErr != nil {
				continue
			}
			id, name, commit, ok := parseGitLink(dest, cacheRoot)
			if !ok || name != e.Name() || seen[id+"/"+name] || (len(only) > 0 && !only[name]) {
				continue
			}
			seen[id+"/"+name] = true
			installs = append(installs, gitInstall{name: name, id: id, commit: commit})
		}
	}
	return installs
}

// projectGitInstalls reads the git entries of a project's lock.
func projectGitInstalls(lock manifest.Lock, m manifest.Mod, only map[string]bool) []gitInstall {
	var installs []gitInstall
	for i := range lock.Skills {
		e := &lock.Skills[i]
		if !e.Source.IsDeclared() || (len(only) > 0 && !only[e.Name]) {
			continue
		}
		in := gitInstall{name: e.Name, id: string(e.Source), constraint: m.Skills[e.Name]}
		if e.Git != nil {
			in.commit, in.tree = e.Git.Commit, e.Git.Tree
		}
		installs = append(installs, in)
	}
	return installs
}

// gitLister lists each git source at most once per run and remembers which
// ones failed, so a source that cannot be consulted is reported once and counts
// against every skill installed from it.
type gitLister struct {
	out       io.Writer
	resolvers map[string]source.Resolver
	listed    map[string]map[string]skill.Skill
	broken    map[string]bool
}

func newGitLister(out io.Writer, chain source.Chain) *gitLister {
	l := &gitLister{
		out: out, resolvers: map[string]source.Resolver{},
		listed: map[string]map[string]skill.Skill{}, broken: map[string]bool{},
	}
	for _, r := range chain.Select(source.KindGit) {
		l.resolvers[r.Spec().ID] = r
	}
	return l
}

// current returns the skill of an install as its source has it now, and the
// tree the installed copy has. ok is false when it could not be had; failed says
// whether that counts as a failure (an undeclared or unreachable source) or only
// a warning (the skill left the source).
func (l *gitLister) current(ctx context.Context, in *gitInstall) (cur skill.Skill, tree string, ok, failed bool) {
	r, declared := l.resolvers[in.id]
	if !declared {
		ui.Warn(l.out, fmt.Sprintf("%s: source %s is not declared on this machine, skipped", in.name, in.id))
		return skill.Skill{}, "", false, true
	}
	if l.broken[in.id] {
		return skill.Skill{}, "", false, true
	}
	if l.listed[in.id] == nil {
		all, err := r.All(ctx)
		if err != nil {
			ui.Warn(l.out, fmt.Sprintf("source %s: %v", in.id, err))
			l.broken[in.id] = true
			return skill.Skill{}, "", false, true
		}
		l.listed[in.id] = make(map[string]skill.Skill, len(all))
		for i := range all {
			l.listed[in.id][all[i].Name] = all[i]
		}
	}
	cur, found := l.listed[in.id][in.name]
	if !found || cur.Git == nil {
		ui.Warn(l.out, fmt.Sprintf("%s is no longer in %s", in.name, in.id))
		return skill.Skill{}, "", false, false
	}
	return cur, installedTree(ctx, r, in), true, false
}

// planGitUpdates compares each install with its source as it is now. A skill is
// outdated when its own directory changed, not when the repository did: in a
// monorepo an unrelated commit moves the commit id of every skill.
func planGitUpdates(
	ctx context.Context, out io.Writer, chain source.Chain, installs []gitInstall,
) (plan []gitUpdate, failed int) {
	lister := newGitLister(out, chain)
	for i := range installs {
		in := &installs[i]
		if in.constraint != "" && in.constraint != "*" {
			ui.Info(out, fmt.Sprintf("%s is pinned to %s: not updated", in.name, in.constraint))
			continue
		}
		cur, tree, ok, didFail := lister.current(ctx, in)
		if didFail {
			failed++
		}
		if !ok || tree == cur.Git.Tree {
			continue
		}
		plan = append(plan, gitUpdate{install: *in, next: cur})
	}
	return plan, failed
}

// installedTree is the tree the installed copy has: the lock's, or for a global
// install the one its commit gives the skill's directory. Empty when unknown,
// which counts as changed.
func installedTree(ctx context.Context, r source.Resolver, in *gitInstall) string {
	if in.tree != "" {
		return in.tree
	}
	if ta, ok := r.(treeAter); ok {
		if tree, err := ta.TreeAt(ctx, in.name, in.commit); err == nil {
			return tree
		}
	}
	return ""
}

func printGitUpdatePlan(out io.Writer, plan []gitUpdate) {
	fmt.Fprintln(out)
	ui.Header(out, "Git sources")
	for i := range plan {
		u := &plan[i]
		fmt.Fprintf(out, "  %s  %s  %s  %s  (%s)\n", ui.SkillName(u.install.name),
			ui.SkillVersion(short(u.install.commit)), ui.Arrow, ui.SkillVersion(short(u.next.Git.Commit)), u.install.id)
	}
	fmt.Fprintln(out)
}

// relinkGit extracts each outdated skill at its new commit and re-points the
// links that exist. It returns what was resolved and what failed.
func relinkGit(
	ctx context.Context, out, errOut io.Writer, chain source.Chain, plan []gitUpdate, targets []string,
) (resolved []skill.Skill, failed int) {
	for i := range plan {
		u := &plan[i]
		s, err := chain.Find(ctx, skill.Ref{Source: u.install.id, Name: u.install.name})
		if err != nil {
			ui.Failf(errOut, "update %s: %v", u.install.name, err)
			failed++
			continue
		}
		linkFailed := false
		for _, target := range targets {
			if !skill.IsLinked(s.Name, target) {
				continue
			}
			if linkErr := skill.Link(s, target); linkErr != nil {
				ui.Failf(errOut, "re-link %s: %v", s.Name, linkErr)
				linkFailed = true
			}
		}
		if linkFailed {
			failed++
			continue
		}
		ui.Success(out, fmt.Sprintf("%s  %s  %s  %s", ui.SkillName(s.Name),
			ui.SkillVersion(short(u.install.commit)), ui.Arrow, ui.SkillVersion(short(s.Git.Commit))))
		resolved = append(resolved, s)
	}
	return resolved, failed
}

// lockGitUpdates records the new commit and tree of each updated skill.
func lockGitUpdates(lock manifest.Lock, resolved []skill.Skill) manifest.Lock {
	for i := range resolved {
		s := &resolved[i]
		entry, found := manifest.FindLockEntry(lock, s.Name)
		if !found {
			continue
		}
		entry.Version = s.Version
		entry.Git = s.Git
		lock = manifest.UpsertLockEntry(lock, entry)
	}
	return lock
}

// gitProject is the project a git update acts on: its manifest directory and
// its lock, which is rewritten when a skill moves.
type gitProject struct {
	rskDir string
	lock   manifest.Lock
}

// collectGitInstalls finds the git skills installed in the scope `rsk update`
// runs in: from the project's lock, or from where the global links point.
func collectGitInstalls(
	cfg config.Config, opts updateOpts, args, targets []string,
) (installs []gitInstall, project *gitProject, err error) {
	only := map[string]bool{}
	for _, raw := range args {
		name, _, parseErr := install.ParseNameVersion(raw)
		if parseErr != nil {
			return nil, nil, parseErr
		}
		only[name] = true
	}
	if opts.global {
		return globalGitInstalls(targets, cfg.SourceOptions().CacheRoot, only), nil, nil
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
	return projectGitInstalls(lock, m, only), &gitProject{rskDir: dir, lock: lock}, nil
}

// runGitUpdate updates the git skills installed in the scope `rsk update` runs
// in. It does nothing, and prints nothing, when there are none.
func runGitUpdate(cmd *cobra.Command, cfg config.Config, opts updateOpts, args []string) error {
	targets, err := resolveTargetDirs(cfg, opts.global, opts.forTool)
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
	plan, failed := planGitUpdates(ctx, out, chain, installs)
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
	lock := lockGitUpdates(project.lock, resolved)
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

// gitLockDrift compares each git entry of the lock with the project's links,
// offline: a link that does not point at the extraction of the locked commit, or
// is missing from every target, means the install is not what the lock records.
func gitLockDrift(lock manifest.Lock, targets []string, cacheRoot string) []string {
	var drift []string
	for i := range lock.Skills {
		e := &lock.Skills[i]
		if !e.Source.IsDeclared() {
			continue
		}
		if e.Git == nil {
			continue
		}
		owner, repo, _ := strings.Cut(strings.TrimPrefix(string(e.Source), "@"), "/")
		want := filepath.Join(cacheRoot, "sources", owner, repo, e.Name, e.Git.Commit)
		linked := false
		for _, target := range targets {
			dest, err := os.Readlink(filepath.Join(target, e.Name))
			if err != nil {
				continue
			}
			linked = true
			if filepath.Clean(dest) != want {
				drift = append(
					drift,
					fmt.Sprintf("%s: linked at %s, the lock says commit %s", e.Name, dest, short(e.Git.Commit)),
				)
			}
		}
		if !linked {
			drift = append(drift, fmt.Sprintf("%s: locked at commit %s but not linked", e.Name, short(e.Git.Commit)))
		}
	}
	return drift
}
