package update

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal/git"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
)

const shortCommit = 7

// Level is how a Notice is shown.
type Level int

const (
	// Info is a neutral message.
	Info Level = iota
	// Warn is a skipped or degraded step that does not fail the command.
	Warn
)

// Notice is one message for the user, in the order it was produced.
type Notice struct {
	Level Level
	Text  string
}

type (
	// GitInstall is an installed git skill and what is known of its content. A
	// project install has the tree its lock recorded; a global one has none and
	// is measured against its commit.
	GitInstall struct {
		Name, ID, Commit, Tree, Constraint string
	}

	// GitUpdate is an installed skill whose directory changed at the source.
	GitUpdate struct {
		Install GitInstall
		Next    skill.Skill
	}

	// Refresher is implemented by a source that can fetch its branch now.
	Refresher interface {
		Refresh(ctx context.Context) error
	}

	// treeAter is implemented by a source that can say what tree a skill's
	// directory had at an earlier commit.
	treeAter interface {
		TreeAt(ctx context.Context, name, commit string) (string, error)
	}
)

// PinLabel shortens a commit id used as a pin, and leaves a version as it is.
func PinLabel(constraint string) string {
	if git.ValidateCommit(constraint) == nil {
		return Short(constraint)
	}
	return constraint
}

// Short shortens a commit id for display.
func Short(commit string) string {
	if len(commit) > shortCommit {
		return commit[:shortCommit]
	}
	return commit
}

// GlobalGitInstalls finds the git skills linked into targets, by reading where
// each link points.
func GlobalGitInstalls(targets []string, cacheRoot, officialCache string, only map[string]bool) []GitInstall {
	var installs []GitInstall
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
			id, name, commit, ok := skill.ParseGitLink(dest, cacheRoot)
			if !ok && isLegacyOfficialLink(dest, officialCache) {
				// Linked into the working clone the official source used before the
				// git transport. Its commit is unknown, so it will be re-pointed.
				id, name, commit, ok = skill.OfficialID, e.Name(), "", true
			}
			if !ok || name != e.Name() || seen[id+"/"+name] || (len(only) > 0 && !only[name]) {
				continue
			}
			seen[id+"/"+name] = true
			installs = append(installs, GitInstall{Name: name, ID: id, Commit: commit})
		}
	}
	return installs
}

// isLegacyOfficialLink reports whether dest points into the working clone the
// official source used before it moved to the git transport.
func isLegacyOfficialLink(dest, officialCache string) bool {
	if officialCache == "" {
		return false
	}
	prefix := filepath.Join(officialCache, skill.SkillsFolderName) + string(filepath.Separator)
	return strings.HasPrefix(filepath.Clean(dest), prefix)
}

// ProjectGitInstalls reads the git entries of a project's lock.
func ProjectGitInstalls(lock manifest.Lock, m manifest.Mod, only map[string]bool) []GitInstall {
	var installs []GitInstall
	for i := range lock.Skills {
		e := &lock.Skills[i]
		if !e.Source.UsesGit() || (len(only) > 0 && !only[e.Name]) {
			continue
		}
		in := GitInstall{Name: e.Name, ID: e.Source.Identity(), Constraint: m.Skills[e.Name]}
		if e.Git != nil {
			in.Commit, in.Tree = e.Git.Commit, e.Git.Tree
		}
		installs = append(installs, in)
	}
	return installs
}

// gitLister lists each git source at most once per run and remembers which
// ones failed, so a source that cannot be consulted is reported once and counts
// against every skill installed from it.
type gitLister struct {
	notices   []Notice
	resolvers map[string]source.Resolver
	listed    map[string]map[string]skill.Skill
	broken    map[string]bool
}

func newGitLister(chain source.Chain) *gitLister {
	l := &gitLister{
		resolvers: map[string]source.Resolver{},
		listed:    map[string]map[string]skill.Skill{}, broken: map[string]bool{},
	}
	for _, r := range chain.Select(source.KindGit) {
		l.resolvers[r.Spec().ID] = r
	}
	return l
}

func (l *gitLister) warn(text string) {
	l.notices = append(l.notices, Notice{Warn, text})
}

// current returns the skill of an install as its source has it now, and the
// tree the installed copy has. ok is false when it could not be had; failed says
// whether that counts as a failure (an undeclared or unreachable source) or only
// a warning (the skill left the source).
func (l *gitLister) current(ctx context.Context, in *GitInstall) (cur skill.Skill, tree string, ok, failed bool) {
	r, declared := l.resolvers[in.ID]
	if !declared {
		l.warn(fmt.Sprintf("%s: source %s is not declared on this machine, skipped", in.Name, in.ID))
		return skill.Skill{}, "", false, true
	}
	if l.broken[in.ID] {
		return skill.Skill{}, "", false, true
	}
	if l.listed[in.ID] == nil {
		// An update must see the source as it is now, not as it was a day ago.
		var err error
		if rf, canRefresh := r.(Refresher); canRefresh {
			err = rf.Refresh(ctx)
		}
		var all []skill.Skill
		if err == nil {
			all, err = r.All(ctx)
		}
		if err != nil {
			l.warn(fmt.Sprintf("source %s: %v", in.ID, err))
			l.broken[in.ID] = true
			return skill.Skill{}, "", false, true
		}
		l.listed[in.ID] = make(map[string]skill.Skill, len(all))
		for i := range all {
			l.listed[in.ID][all[i].Name] = all[i]
		}
	}
	cur, found := l.listed[in.ID][in.Name]
	if !found || cur.Git == nil {
		l.warn(fmt.Sprintf("%s is no longer in %s", in.Name, in.ID))
		return skill.Skill{}, "", false, false
	}
	return cur, installedTree(ctx, r, in), true, false
}

// GitPlan is what PlanGit found: the skills to update, how many installs could
// not be checked, and the messages produced on the way, in order.
type GitPlan struct {
	Updates []GitUpdate
	Failed  int
	Notices []Notice
}

// PlanGit compares each install with its source as it is now. A skill is
// outdated when its own directory changed, not when the repository did: in a
// monorepo an unrelated commit moves the commit id of every skill.
func PlanGit(ctx context.Context, chain source.Chain, installs []GitInstall) GitPlan {
	lister := newGitLister(chain)
	var plan GitPlan
	for i := range installs {
		in := &installs[i]
		if in.Constraint != "" && in.Constraint != "*" {
			lister.notices = append(lister.notices,
				Notice{Info, fmt.Sprintf("%s is pinned to %s: not updated", in.Name, PinLabel(in.Constraint))})
			continue
		}
		cur, tree, ok, didFail := lister.current(ctx, in)
		if didFail {
			plan.Failed++
		}
		if !ok || tree == cur.Git.Tree {
			continue
		}
		plan.Updates = append(plan.Updates, GitUpdate{Install: *in, Next: cur})
	}
	plan.Notices = lister.notices
	return plan
}

// installedTree is the tree the installed copy has: the lock's, or for a global
// install the one its commit gives the skill's directory. Empty when unknown,
// which counts as changed.
func installedTree(ctx context.Context, r source.Resolver, in *GitInstall) string {
	if in.Tree != "" {
		return in.Tree
	}
	if ta, ok := r.(treeAter); ok {
		if tree, err := ta.TreeAt(ctx, in.Name, in.Commit); err == nil {
			return tree
		}
	}
	return ""
}

// GitRelink is the outcome of updating one skill: the skill as resolved, and
// what went wrong. FindErr is set when the new version could not be extracted;
// LinkErrs holds one error per link that could not be re-pointed. A skill with
// neither error is updated.
type GitRelink struct {
	Update   GitUpdate
	Skill    skill.Skill
	FindErr  error
	LinkErrs []error
}

// Failed reports whether the skill was not updated.
func (r *GitRelink) Failed() bool {
	return r.FindErr != nil || len(r.LinkErrs) > 0
}

// RelinkGit extracts each outdated skill at its new commit and re-points the
// links that exist. It reports every skill's outcome, so one failure does not
// stop the rest.
func RelinkGit(ctx context.Context, chain source.Chain, updates []GitUpdate, targets []string) []GitRelink {
	results := make([]GitRelink, 0, len(updates))
	for i := range updates {
		u := &updates[i]
		res := GitRelink{Update: *u}
		s, err := chain.Find(ctx, skill.Ref{Source: u.Install.ID, Name: u.Install.Name})
		if err != nil {
			res.FindErr = err
			results = append(results, res)
			continue
		}
		res.Skill = s
		for _, target := range targets {
			if !skill.IsLinked(s.Name, target) {
				continue
			}
			if linkErr := skill.Link(s, target); linkErr != nil {
				res.LinkErrs = append(res.LinkErrs, linkErr)
			}
		}
		results = append(results, res)
	}
	return results
}

// LockGitUpdates records the new commit and tree of each updated skill.
func LockGitUpdates(lock manifest.Lock, resolved []skill.Skill) manifest.Lock {
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

// GitLockDrift compares each git entry of the lock with the project's links,
// offline: a link that does not point at the extraction of the locked commit, or
// is missing from every target, means the install is not what the lock records.
func GitLockDrift(lock manifest.Lock, targets []string, cacheRoot string) []string {
	var drift []string
	for i := range lock.Skills {
		e := &lock.Skills[i]
		if !e.Source.UsesGit() || e.Git == nil {
			continue
		}
		owner, repo, _ := strings.Cut(strings.TrimPrefix(e.Source.Identity(), "@"), "/")
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
					fmt.Sprintf("%s: linked at %s, the lock says commit %s", e.Name, dest, Short(e.Git.Commit)),
				)
			}
		}
		if !linked {
			drift = append(drift, fmt.Sprintf("%s: locked at commit %s but not linked", e.Name, Short(e.Git.Commit)))
		}
	}
	return drift
}
