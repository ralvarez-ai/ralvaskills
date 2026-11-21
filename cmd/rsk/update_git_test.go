package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
)

const (
	acmeID    = "@acme/skills"
	oldCommit = "1111111111111111111111111111111111111111"
	newCommit = "2222222222222222222222222222222222222222"
	oldTree   = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	newTree   = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
)

// fakeGitSource lists a fixed set of skills and answers TreeAt from a table.
type fakeGitSource struct {
	id      string
	skills  map[string]skill.Skill
	listErr error
	trees   map[string]string // commit to tree id, for TreeAt
	calls   int
}

func (f *fakeGitSource) Spec() source.Spec {
	return source.Spec{ID: f.id, Type: source.KindGit}
}

func (f *fakeGitSource) All(context.Context) ([]skill.Skill, error) {
	f.calls++
	if f.listErr != nil {
		return nil, f.listErr
	}
	out := make([]skill.Skill, 0, len(f.skills))
	for name := range f.skills {
		out = append(out, f.skills[name])
	}
	return out, nil
}

func (f *fakeGitSource) Find(_ context.Context, name string) (skill.Skill, error) {
	s, ok := f.skills[name]
	if !ok {
		return skill.Skill{}, source.ErrNotFound
	}
	return s, nil
}

func (f *fakeGitSource) TreeAt(_ context.Context, _, commit string) (string, error) {
	if tree, ok := f.trees[commit]; ok {
		return tree, nil
	}
	return "", errors.New("object gone")
}

func currentGitSkill(name, tree string) skill.Skill {
	return skill.Skill{
		Name:   name,
		Source: skill.DeclaredSource(acmeID),
		Git: &skill.GitRef{
			URL:    "git@github.com:acme/skills.git",
			Commit: newCommit,
			Tree:   tree,
			Path:   "skills/" + name,
		},
	}
}

func TestParseGitLink(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	dest := filepath.Join(root, "sources", "acme", "skills", "alpha", oldCommit)
	id, name, commit, ok := parseGitLink(dest, root)
	if !ok || id != "@acme/skills" || name != "alpha" || commit != oldCommit {
		t.Errorf("parseGitLink = %q %q %q %v", id, name, commit, ok)
	}
	for _, bad := range []string{
		filepath.Join(root, "registry", "alpha", "1.0.0"),
		filepath.Join(root, "sources", "acme", "skills", "alpha"),
		filepath.Join(root, "sources", "acme", "skills", "alpha", oldCommit, "extra"),
		"/elsewhere/alpha",
	} {
		if _, _, _, ok = parseGitLink(bad, root); ok {
			t.Errorf("parseGitLink(%s) accepted a path that is not a git extraction", bad)
		}
	}
}

func TestPlanGitUpdatesComparesTreesNotCommits(t *testing.T) {
	t.Parallel()
	src := &fakeGitSource{id: "@acme/skills", skills: map[string]skill.Skill{
		// Same tree at a newer commit: the repository moved, this skill did not.
		"unchanged": currentGitSkill("unchanged", oldTree),
		"changed":   currentGitSkill("changed", newTree),
	}}
	installs := []gitInstall{
		{name: "unchanged", id: "@acme/skills", commit: oldCommit, tree: oldTree},
		{name: "changed", id: "@acme/skills", commit: oldCommit, tree: oldTree},
		{name: "pinned", id: "@acme/skills", commit: oldCommit, tree: oldTree, constraint: "1.2.0"},
		{name: "atcommit", id: "@acme/skills", commit: oldCommit, tree: oldTree, constraint: oldCommit},
		{name: "gone", id: "@acme/skills", commit: oldCommit, tree: oldTree},
	}
	var out bytes.Buffer
	plan, failed := planGitUpdates(t.Context(), &out, source.Chain{src}, installs)

	if len(plan) != 1 || plan[0].install.name != "changed" || plan[0].next.Git.Commit != newCommit {
		t.Fatalf("plan = %+v, want only the skill whose own tree changed", plan)
	}
	if failed != 0 {
		t.Errorf("failed = %d; a skill that left the source is a warning, not a failure", failed)
	}
	text := out.String()
	for _, want := range []string{
		"pinned is pinned to 1.2.0", "atcommit is pinned to 1111111", "gone is no longer in @acme/skills",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("output lacks %q:\n%s", want, text)
		}
	}
	if src.calls != 1 {
		t.Errorf("the source was listed %d times, want once per source", src.calls)
	}
}

func TestPlanGitUpdatesReportsASourceThatCannotBeConsulted(t *testing.T) {
	t.Parallel()
	broken := &fakeGitSource{id: "@acme/skills", listErr: errors.New("permission denied (publickey)")}
	installs := []gitInstall{
		{name: "a", id: "@acme/skills", commit: oldCommit, tree: oldTree},
		{name: "b", id: "@acme/skills", commit: oldCommit, tree: oldTree},
		{name: "c", id: "@other/repo", commit: oldCommit, tree: oldTree},
	}
	var out bytes.Buffer
	plan, failed := planGitUpdates(t.Context(), &out, source.Chain{broken}, installs)

	if len(plan) != 0 || failed != 3 {
		t.Errorf(
			"plan=%v failed=%d; a source that fails must fail every skill from it, and an undeclared one too",
			plan,
			failed,
		)
	}
	for _, want := range []string{"permission denied", "@other/repo is not declared"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output lacks %q:\n%s", want, out.String())
		}
	}
	if broken.calls != 1 {
		t.Errorf("a failing source was retried %d times", broken.calls)
	}
}

func TestPlanGitUpdatesMeasuresAGlobalInstallThroughItsCommit(t *testing.T) {
	t.Parallel()
	src := &fakeGitSource{
		id:     "@acme/skills",
		skills: map[string]skill.Skill{"alpha": currentGitSkill("alpha", newTree)},
		trees:  map[string]string{oldCommit: newTree},
	}
	// No lock tree: the install's own commit says the directory already has the
	// current tree, so nothing changed even though the repository moved.
	same := []gitInstall{{name: "alpha", id: "@acme/skills", commit: oldCommit}}
	if plan, _ := planGitUpdates(t.Context(), &bytes.Buffer{}, source.Chain{src}, same); len(plan) != 0 {
		t.Errorf("plan = %+v, want none", plan)
	}

	src.trees[oldCommit] = oldTree
	if plan, _ := planGitUpdates(t.Context(), &bytes.Buffer{}, source.Chain{src}, same); len(plan) != 1 {
		t.Errorf("plan = %+v, want the changed skill", plan)
	}

	// When the old commit's objects are gone the tree is unknown, which counts as changed.
	delete(src.trees, oldCommit)
	if plan, _ := planGitUpdates(t.Context(), &bytes.Buffer{}, source.Chain{src}, same); len(plan) != 1 {
		t.Errorf("an unknown installed tree must update, got %+v", plan)
	}
}

func TestProjectAndGlobalGitInstalls(t *testing.T) {
	t.Parallel()
	lock := manifest.Lock{Skills: []manifest.LockEntry{
		{Name: "alpha", Source: skill.DeclaredSource(acmeID), Git: &skill.GitRef{Commit: oldCommit, Tree: oldTree}},
		{Name: "tdd", Source: skill.SourceRegistry, Version: "1.0.0"},
		{Name: "beta", Source: skill.DeclaredSource(acmeID), Git: &skill.GitRef{Commit: oldCommit, Tree: oldTree}},
	}}
	m := manifest.Mod{Skills: map[string]string{"alpha": "*", "beta": "2.0.0", "tdd": "*"}}

	got := projectGitInstalls(lock, m, nil)
	if len(got) != 2 || got[0].name != "alpha" || got[1].constraint != "2.0.0" || got[0].tree != oldTree {
		t.Errorf("project installs = %+v", got)
	}
	if only := projectGitInstalls(lock, m, map[string]bool{"beta": true}); len(only) != 1 || only[0].name != "beta" {
		t.Errorf("filtered = %+v", only)
	}

	root := t.TempDir()
	target := t.TempDir()
	extracted := filepath.Join(root, "sources", "acme", "skills", "alpha", oldCommit)
	if err := os.MkdirAll(extracted, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(extracted, filepath.Join(target, "alpha")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("/somewhere/else", filepath.Join(target, "other")); err != nil {
		t.Fatal(err)
	}
	globals := globalGitInstalls([]string{target, filepath.Join(target, "missing")}, root, "", nil)
	if len(globals) != 1 || globals[0].id != "@acme/skills" || globals[0].commit != oldCommit {
		t.Errorf("global installs = %+v", globals)
	}
}

func TestRelinkGitRepointsLinksAndLockGitUpdatesTheEntry(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	oldDir, newDir := filepath.Join(dir, "old"), filepath.Join(dir, "new")
	for _, d := range []string{oldDir, newDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	target := filepath.Join(dir, "skills")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(oldDir, filepath.Join(target, "alpha")); err != nil {
		t.Fatal(err)
	}

	next := currentGitSkill("alpha", newTree)
	next.Path, next.Version = newDir, "1.1.0"
	src := &fakeGitSource{id: "@acme/skills", skills: map[string]skill.Skill{"alpha": next}}
	plan := []gitUpdate{{install: gitInstall{name: "alpha", id: "@acme/skills", commit: oldCommit}, next: next}}

	var out, errOut bytes.Buffer
	resolved, failed := relinkGit(t.Context(), &out, &errOut, source.Chain{src}, plan, []string{target})
	if failed != 0 || len(resolved) != 1 {
		t.Fatalf("resolved=%v failed=%d err=%s", resolved, failed, errOut.String())
	}
	if got, _ := os.Readlink(filepath.Join(target, "alpha")); got != newDir {
		t.Errorf("link points at %q, want %q", got, newDir)
	}

	lock := manifest.Lock{Skills: []manifest.LockEntry{{
		Name: "alpha", Source: skill.DeclaredSource(acmeID), Version: "1.0.0",
		Git: &skill.GitRef{Commit: oldCommit, Tree: oldTree}, Permissions: []string{"Read"},
	}}}
	lock = lockGitUpdates(lock, resolved)
	e := lock.Skills[0]
	if e.Git == nil || e.Git.Commit != newCommit || e.Git.Tree != newTree || e.Version != "1.1.0" ||
		e.Git.Path != "skills/alpha" || len(e.Permissions) != 1 {
		t.Errorf("entry = %+v, want the new commit and tree with the permissions record kept", e)
	}
}

func TestGitLockDriftFindsALinkThatIsNotTheLockedCommit(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	target := t.TempDir()
	extract := func(name, commit string) string {
		d := filepath.Join(root, "sources", "acme", "skills", name, commit)
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
		return d
	}
	link := func(name, dest string) {
		if err := os.Symlink(dest, filepath.Join(target, name)); err != nil {
			t.Fatal(err)
		}
	}
	link("ok", extract("ok", oldCommit))
	link("moved", extract("moved", newCommit)) // linked at a commit other than the locked one
	extract("missing", oldCommit)              // locked, extracted, but not linked in this target

	lock := manifest.Lock{Skills: []manifest.LockEntry{
		{Name: "ok", Source: skill.DeclaredSource(acmeID), Git: &skill.GitRef{Commit: oldCommit}},
		{Name: "moved", Source: skill.DeclaredSource(acmeID), Git: &skill.GitRef{Commit: oldCommit}},
		{Name: "missing", Source: skill.DeclaredSource(acmeID), Git: &skill.GitRef{Commit: oldCommit}},
		{Name: "tdd", Source: skill.SourceRegistry},
	}}
	got := gitLockDrift(lock, []string{target}, root)
	if len(got) != 2 {
		t.Fatalf("drift = %v, want the moved link and the missing one", got)
	}
	joined := strings.Join(got, "\n")
	for _, want := range []string{"moved", "missing"} {
		if !strings.Contains(joined, want) {
			t.Errorf("drift %q does not mention %s", joined, want)
		}
	}
	if strings.Contains(joined, "ok:") || strings.Contains(joined, "tdd") {
		t.Errorf("drift %q reports a skill that matches its lock or is not a git skill", joined)
	}
}

func TestOfficialSkillsJoinTheGitUpdatePass(t *testing.T) {
	t.Parallel()
	lock := manifest.Lock{Skills: []manifest.LockEntry{
		// Written before the git transport: no git record at all.
		{Name: "xlsx", Source: skill.SourceOfficial},
		{Name: "pdf", Source: skill.SourceOfficial, Git: &skill.GitRef{Commit: oldCommit, Tree: oldTree}},
		{Name: "tdd", Source: skill.SourceRegistry, Version: "1.0.0"},
	}}
	got := projectGitInstalls(lock, manifest.Mod{Skills: map[string]string{"xlsx": "*", "pdf": "*"}}, nil)
	if len(got) != 2 || got[0].id != skill.OfficialID || got[0].tree != "" || got[1].tree != oldTree ||
		got[1].commit != oldCommit {
		t.Errorf("installs = %+v, want both official entries under the official identity", got)
	}
}

func TestGlobalGitInstallsFindsLinksIntoTheOldOfficialClone(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	officialCache := filepath.Join(root, "anthropic")
	target := t.TempDir()
	// The old layout: the repository was cloned into <official_cache>/skills.
	oldClone := filepath.Join(officialCache, "skills", "skills", "xlsx")
	if err := os.MkdirAll(oldClone, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(oldClone, filepath.Join(target, "xlsx")); err != nil {
		t.Fatal(err)
	}
	// And one already in the new layout.
	fresh := filepath.Join(root, "sources", "anthropics", "skills", "pdf", newCommit)
	if err := os.MkdirAll(fresh, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(fresh, filepath.Join(target, "pdf")); err != nil {
		t.Fatal(err)
	}

	got := globalGitInstalls([]string{target}, root, officialCache, nil)
	if len(got) != 2 {
		t.Fatalf("installs = %+v, want both official skills", got)
	}
	byName := map[string]gitInstall{}
	for _, in := range got {
		byName[in.name] = in
	}
	if in := byName["xlsx"]; in.id != skill.OfficialID || in.commit != "" {
		t.Errorf("legacy link = %+v, want the official identity and no known commit", in)
	}
	if in := byName["pdf"]; in.id != skill.OfficialID || in.commit != newCommit {
		t.Errorf("new link = %+v", in)
	}
}

func TestPlanGitUpdatesMigratesALegacyOfficialInstall(t *testing.T) {
	t.Parallel()
	src := &fakeGitSource{id: skill.OfficialID, skills: map[string]skill.Skill{
		"xlsx": {
			Name: "xlsx", Source: skill.SourceOfficial,
			Git: &skill.GitRef{Commit: newCommit, Tree: newTree, Path: "skills/xlsx"},
		},
	}}
	legacy := []gitInstall{{name: "xlsx", id: skill.OfficialID}}
	plan, failed := planGitUpdates(t.Context(), &bytes.Buffer{}, source.Chain{src}, legacy)
	if failed != 0 || len(plan) != 1 || plan[0].next.Git.Commit != newCommit {
		t.Errorf("plan = %+v failed=%d: an install with no recorded tree must be re-pointed", plan, failed)
	}
}

func TestPlanGitUpdatesRefreshesEachSourceOnceAndReportsAFailure(t *testing.T) {
	t.Parallel()
	ok := &refreshingSource{}
	ok.id = "@acme/skills"
	ok.skills = map[string]skill.Skill{"a": currentGitSkill("a", newTree), "b": currentGitSkill("b", newTree)}
	installs := []gitInstall{
		{name: "a", id: "@acme/skills", commit: oldCommit, tree: newTree},
		{name: "b", id: "@acme/skills", commit: oldCommit, tree: newTree},
	}
	if _, failed := planGitUpdates(t.Context(), &bytes.Buffer{}, source.Chain{ok}, installs); failed != 0 {
		t.Errorf("failed = %d", failed)
	}
	if ok.refreshes != 1 {
		t.Errorf("refreshed %d times: an update must see the source as it is now, once", ok.refreshes)
	}

	bad := &refreshingSource{refreshErr: errors.New("network is unreachable")}
	bad.id = "@acme/skills"
	bad.skills = map[string]skill.Skill{"a": currentGitSkill("a", newTree)}
	var out bytes.Buffer
	plan, failed := planGitUpdates(t.Context(), &out, source.Chain{bad}, installs[:1])
	if len(plan) != 0 || failed != 1 || !strings.Contains(out.String(), "network is unreachable") {
		t.Errorf("plan=%v failed=%d out=%q; a source that cannot be refreshed must fail, not answer from its old index",
			plan, failed, out.String())
	}
}

// refreshingSource is a fakeGitSource that can be refreshed.
type refreshingSource struct {
	fakeGitSource
	refreshes  int
	refreshErr error
}

func (r *refreshingSource) Refresh(context.Context) error {
	r.refreshes++
	return r.refreshErr
}

func TestGitLockDriftCoversOfficialEntries(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	lock := manifest.Lock{Skills: []manifest.LockEntry{
		{Name: "xlsx", Source: skill.SourceOfficial, Git: &skill.GitRef{Commit: oldCommit}},
		{Name: "legacy", Source: skill.SourceOfficial}, // no record: nothing to compare
	}}
	got := gitLockDrift(lock, []string{t.TempDir()}, root)
	if len(got) != 1 || !strings.Contains(got[0], "xlsx") {
		t.Errorf("drift = %v, want only the official entry that has a record", got)
	}
}
