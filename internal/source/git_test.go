package source

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go.ralva.dev/ralvaskills/v3/internal/git"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

const (
	testID     = "@acme/skills"
	testURL    = "git@github.com:acme/skills.git"
	testCommit = "1111111111111111111111111111111111111111"
	otherSHA   = "2222222222222222222222222222222222222222"
)

// fakeRepo is an in-memory stand-in for a bare clone, so the resolver's logic
// (index, freshness, fail-closed, extraction) is tested without a network.
type fakeRepo struct {
	entries    []git.TreeEntry
	blobs      map[string]string
	commit     string
	fetchErr   error
	archiveErr error
	missing    bool // the first Archive fails as if the commit's objects were gone

	refs  map[string]string // ref to commit, when a test needs tags; nil accepts any ref
	trees map[string]string // directory to tree id; missing means "t-<dir>"

	fetches  []string
	listed   [][]string
	archived []string
}

func (f *fakeRepo) TreeID(_ context.Context, _, dir string) (string, error) {
	if id, ok := f.trees[dir]; ok {
		return id, nil
	}
	return "t-" + dir, nil
}

func (f *fakeRepo) Fetch(_ context.Context, _ git.RemoteURL, ref string) (string, error) {
	f.fetches = append(f.fetches, ref)
	if f.fetchErr != nil {
		return "", f.fetchErr
	}
	f.missing = false
	if c, ok := f.refs[ref]; ok {
		return c, nil
	}
	if f.refs != nil && ref != "HEAD" {
		return "", fmt.Errorf("%w: %s", git.ErrRefNotFound, ref)
	}
	return f.commit, nil
}

func (f *fakeRepo) ListTree(_ context.Context, _ string, roots []string) ([]git.TreeEntry, error) {
	f.listed = append(f.listed, roots)
	return f.entries, nil
}

func (f *fakeRepo) ReadBlobs(_ context.Context, _ string, paths []string, maxBytes int64) (map[string]git.Blob, error) {
	out := make(map[string]git.Blob, len(paths))
	for _, p := range paths {
		body, ok := f.blobs[p]
		switch {
		case !ok:
			out[p] = git.Blob{Missing: true}
		case int64(len(body)) > maxBytes:
			out[p] = git.Blob{Size: int64(len(body))}
		default:
			out[p] = git.Blob{Data: []byte(body), Size: int64(len(body))}
		}
	}
	return out, nil
}

func (f *fakeRepo) Archive(_ context.Context, commit, dir, dest string, _ git.Limits) error {
	if f.archiveErr != nil {
		return f.archiveErr
	}
	if f.missing {
		return &git.Error{Op: "archive", Err: errors.New("exit status 128")}
	}
	f.archived = append(f.archived, commit+":"+dir)
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dest, "SKILL.md"), []byte(f.blobs[dir+"/SKILL.md"]), 0o644)
}

func skillMD(name, version string) string {
	return "---\nname: " + name + "\nversion: " + version + "\ndescription: does " + name + "\n---\n"
}

func tree(p string) git.TreeEntry {
	return git.TreeEntry{Mode: "040000", Type: "tree", ObjectID: "t-" + p, Path: p}
}

func blob(p string) git.TreeEntry {
	return git.TreeEntry{Mode: "100644", Type: "blob", ObjectID: "b-" + p, Path: p}
}

func baseRepo() *fakeRepo {
	return &fakeRepo{
		commit: testCommit,
		entries: []git.TreeEntry{
			tree("skills"), tree("skills/alpha"), blob("skills/alpha/SKILL.md"), blob("skills/alpha/notes.md"),
			tree("skills/beta"), blob("skills/beta/SKILL.md"),
		},
		blobs: map[string]string{
			"skills/alpha/SKILL.md": skillMD("alpha", "1.0.0"),
			"skills/beta/SKILL.md":  skillMD("beta", "2.0.0"),
		},
	}
}

type harness struct {
	g     *Git
	repo  *fakeRepo
	now   time.Time
	trust error
}

func newHarness(t *testing.T, mutate func(*Spec), repo *fakeRepo) *harness {
	t.Helper()
	h := &harness{repo: repo, now: time.Date(2026, 11, 21, 12, 0, 0, 0, time.UTC)}
	spec := Spec{ID: testID, Type: KindGit, URL: testURL}
	if mutate != nil {
		mutate(&spec)
	}
	g, err := newGit(spec, Options{
		CacheRoot: t.TempDir(),
		Gate:      func(string, string) error { return h.trust },
		Now:       func() time.Time { return h.now },
	}, repo)
	if err != nil {
		t.Fatal(err)
	}
	h.g = g
	return h
}

func TestGitAllListsTheIndexedSkills(t *testing.T) {
	t.Parallel()
	h := newHarness(t, func(s *Spec) { s.Roots = []string{"skills"} }, baseRepo())

	got, err := h.g.All(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Name != "alpha" || got[1].Name != "beta" {
		t.Fatalf("skills = %+v", got)
	}
	a := got[0]
	if a.Version != "1.0.0" || a.Description != "does alpha" || a.Source != skill.DeclaredSource(testID) {
		t.Errorf("alpha = %+v", a)
	}
	if len(h.repo.listed) != 1 || len(h.repo.listed[0]) != 1 || h.repo.listed[0][0] != "skills" {
		t.Errorf("roots passed to ListTree = %v", h.repo.listed)
	}
	if len(h.repo.fetches) != 1 || h.repo.fetches[0] != "HEAD" {
		t.Errorf("fetched %v, want the remote HEAD when no rev is set", h.repo.fetches)
	}
}

func TestGitFetchesTheDeclaredRev(t *testing.T) {
	t.Parallel()
	h := newHarness(t, func(s *Spec) { s.Rev = "v1.2" }, baseRepo())
	if _, err := h.g.All(t.Context()); err != nil {
		t.Fatal(err)
	}
	if h.repo.fetches[0] != "v1.2" {
		t.Errorf("fetched %v", h.repo.fetches)
	}
}

func TestGitIsGatedBeforeAnythingElse(t *testing.T) {
	t.Parallel()
	repo := baseRepo()
	h := newHarness(t, nil, repo)
	h.trust = errors.New("source is not trusted")

	if _, err := h.g.All(t.Context()); err == nil || errors.Is(err, ErrNotFound) {
		t.Errorf("All err = %v, want the gate's error", err)
	}
	if _, err := h.g.Find(t.Context(), "alpha"); err == nil || errors.Is(err, ErrNotFound) {
		t.Errorf("Find err = %v, want the gate's error, not a not-found", err)
	}
	if len(repo.fetches)+len(repo.listed)+len(repo.archived) != 0 {
		t.Errorf("an untrusted source reached git: fetches=%v listed=%v archived=%v",
			repo.fetches, repo.listed, repo.archived)
	}
}

func TestGitIndexBuildRules(t *testing.T) {
	t.Parallel()
	repo := &fakeRepo{
		commit: testCommit,
		entries: []git.TreeEntry{
			tree("skills/alpha"), blob("skills/alpha/SKILL.md"),
			// nested skill inside a skill: the outer one wins, as in skill.Walk
			tree("skills/alpha/inner"), blob("skills/alpha/inner/SKILL.md"),
			// a symlink inside a skill excludes the whole skill
			tree("skills/linked"), blob("skills/linked/SKILL.md"),
			{Mode: git.ModeSymlink, Type: "blob", ObjectID: "x", Path: "skills/linked/escape"},
			// a submodule inside a skill excludes it as well
			tree("skills/subbed"), blob("skills/subbed/SKILL.md"),
			{Mode: git.ModeSubmodule, Type: "commit", ObjectID: "y", Path: "skills/subbed/mod"},
			// SKILL.md that is itself a symlink is not a skill
			tree("skills/fake"), {Mode: git.ModeSymlink, Type: "blob", ObjectID: "z", Path: "skills/fake/SKILL.md"},
			// invalid directory name
			tree("skills/Bad_Name"), blob("skills/Bad_Name/SKILL.md"),
			// excluded by pattern
			tree("skills/node_modules/dep"), blob("skills/node_modules/dep/SKILL.md"),
			// a SKILL.md at the repository root has no name
			blob("SKILL.md"),
			// SKILL.md larger than the cap
			tree("skills/huge"), blob("skills/huge/SKILL.md"),
			tree("skills/keep"), blob("skills/keep/SKILL.md"),
		},
		blobs: map[string]string{
			"skills/alpha/SKILL.md": skillMD("alpha", "1.0.0"),
			"skills/keep/SKILL.md":  skillMD("keep", "3.0.0"),
			"skills/huge/SKILL.md":  strings.Repeat("x", maxSkillFile+1),
		},
	}
	h := newHarness(t, func(s *Spec) { s.Excludes = []string{"**/node_modules/**"} }, repo)

	got, err := h.g.All(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(got))
	for _, s := range got {
		names = append(names, s.Name)
	}
	if strings.Join(names, ",") != "alpha,keep" {
		t.Errorf("indexed %v, want only alpha and keep", names)
	}

	idx, _ := h.g.loadIndex()
	reasons := map[string]string{}
	for _, sk := range idx.Skipped {
		reasons[sk.Path] = sk.Reason
	}
	for path, want := range map[string]string{
		"skills/linked":   "symlink",
		"skills/subbed":   "submodule",
		"skills/fake":     "SKILL.md is not a regular file",
		"skills/Bad_Name": "name",
		"skills/huge":     "too large",
	} {
		if !strings.Contains(reasons[path], want) {
			t.Errorf("skipped[%s] = %q, want it to mention %q (all: %v)", path, reasons[path], want, reasons)
		}
	}
}

func TestGitDuplicateNamesAreAnError(t *testing.T) {
	t.Parallel()
	repo := baseRepo()
	repo.entries = append(repo.entries, tree("other/alpha"), blob("other/alpha/SKILL.md"))
	repo.blobs["other/alpha/SKILL.md"] = skillMD("alpha", "9.0.0")
	h := newHarness(t, nil, repo)

	_, err := h.g.All(t.Context())
	if !errors.Is(err, ErrDuplicate) {
		t.Fatalf("err = %v, want ErrDuplicate", err)
	}
	for _, want := range []string{"skills/alpha", "other/alpha"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not name %s", err, want)
		}
	}
}

func TestGitIndexIsFreshWithinTheTTLAndRefetchedAfterIt(t *testing.T) {
	t.Parallel()
	h := newHarness(t, nil, baseRepo())
	ctx := t.Context()

	if _, err := h.g.All(ctx); err != nil {
		t.Fatal(err)
	}
	h.now = h.now.Add(23 * time.Hour)
	if _, err := h.g.All(ctx); err != nil {
		t.Fatal(err)
	}
	if len(h.repo.fetches) != 1 {
		t.Errorf("fetched %d times inside the TTL, want 1", len(h.repo.fetches))
	}

	h.now = h.now.Add(2 * time.Hour)
	if _, err := h.g.All(ctx); err != nil {
		t.Fatal(err)
	}
	if len(h.repo.fetches) != 2 {
		t.Errorf("fetched %d times after the TTL, want 2", len(h.repo.fetches))
	}
}

func TestGitFailsClosedWhenTheRefreshFailsPastTheTTL(t *testing.T) {
	t.Parallel()
	h := newHarness(t, nil, baseRepo())
	ctx := t.Context()
	if _, err := h.g.All(ctx); err != nil {
		t.Fatal(err)
	}

	h.now = h.now.Add(25 * time.Hour)
	h.repo.fetchErr = errors.New("permission denied (publickey)")
	_, err := h.g.Find(ctx, "alpha")
	if err == nil || errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v: a stale index must not answer when the refresh fails", err)
	}
	if !strings.Contains(err.Error(), "publickey") {
		t.Errorf("error %q hides why the refresh failed", err)
	}
}

func TestGitOfflineUsesAnExistingIndexAndSaysWhenThereIsNone(t *testing.T) {
	t.Parallel()
	h := newHarness(t, nil, baseRepo())
	ctx := t.Context()

	h.g.offline = true
	if _, err := h.g.All(ctx); !errors.Is(err, ErrNeverFetched) {
		t.Fatalf("err = %v, want ErrNeverFetched when offline with no index", err)
	}
	if len(h.repo.fetches) != 0 {
		t.Error("--offline reached the network")
	}

	h.g.offline = false
	if _, err := h.g.All(ctx); err != nil {
		t.Fatal(err)
	}
	h.g.offline = true
	h.now = h.now.Add(30 * 24 * time.Hour)
	h.repo.fetchErr = errors.New("unreachable")
	if _, err := h.g.All(ctx); err != nil {
		t.Errorf("offline must accept an existing index however old: %v", err)
	}
}

func TestGitIndexIsBoundToItsSpec(t *testing.T) {
	t.Parallel()
	repo := baseRepo()
	h := newHarness(t, nil, repo)
	if _, err := h.g.All(t.Context()); err != nil {
		t.Fatal(err)
	}

	// Same cache, but the declaration now points at another URL: the cached
	// index must not answer for it, even offline.
	moved := *h.g
	moved.spec.URL = "git@evil.example:acme/skills.git"
	moved.offline = true
	if _, err := moved.All(t.Context()); !errors.Is(err, ErrNeverFetched) {
		t.Errorf("err = %v, want the index of another URL to be ignored", err)
	}
}

func TestGitFindExtractsOncePerCommit(t *testing.T) {
	t.Parallel()
	h := newHarness(t, nil, baseRepo())
	ctx := t.Context()

	got, err := h.g.Find(ctx, "alpha")
	if err != nil {
		t.Fatal(err)
	}
	wantSuffix := filepath.Join("sources", "acme", "skills", "alpha", testCommit)
	if !strings.HasSuffix(got.Path, wantSuffix) || got.Source != skill.DeclaredSource(testID) ||
		got.Version != "1.0.0" {
		t.Errorf("skill = %+v, want a path ending %s", got, wantSuffix)
	}
	if _, err = os.Stat(filepath.Join(got.Path, "SKILL.md")); err != nil {
		t.Errorf("the skill was not extracted: %v", err)
	}

	if _, err = h.g.Find(ctx, "alpha"); err != nil {
		t.Fatal(err)
	}
	if len(h.repo.archived) != 1 || h.repo.archived[0] != testCommit+":skills/alpha" {
		t.Errorf("archived %v, want one extraction of skills/alpha", h.repo.archived)
	}
}

func TestGitFindMissingNameIsNotFound(t *testing.T) {
	t.Parallel()
	h := newHarness(t, nil, baseRepo())
	_, err := h.g.Find(t.Context(), "ghost")
	if !errors.Is(err, ErrNotFound) || !strings.Contains(err.Error(), testID) {
		t.Errorf("err = %v, want ErrNotFound naming the source", err)
	}
}

func TestGitFindRefetchesACommitWhoseObjectsAreGone(t *testing.T) {
	t.Parallel()
	repo := baseRepo()
	repo.missing = true
	h := newHarness(t, nil, repo)
	ctx := t.Context()
	// Build the index first, then lose the objects.
	if _, err := h.g.All(ctx); err != nil {
		t.Fatal(err)
	}
	repo.missing = true
	repo.fetches = nil

	if _, err := h.g.Find(ctx, "alpha"); err != nil {
		t.Fatalf("Find: %v", err)
	}
	if len(repo.fetches) != 1 || repo.fetches[0] != testCommit {
		t.Errorf("fetches = %v, want a fetch of the locked commit", repo.fetches)
	}

	// Offline there is nothing to fetch it from.
	repo.missing = true
	repo.fetches = nil
	h.g.offline = true
	if err := os.RemoveAll(h.g.extractDir("alpha", testCommit)); err != nil {
		t.Fatal(err)
	}
	if _, err := h.g.Find(ctx, "alpha"); err == nil {
		t.Error("offline Find succeeded without the objects")
	}
	if len(repo.fetches) != 0 {
		t.Errorf("offline Find fetched %v", repo.fetches)
	}
}

func TestGitFindDoesNotRetryAnUnsafeArchive(t *testing.T) {
	t.Parallel()
	repo := baseRepo()
	h := newHarness(t, nil, repo)
	if _, err := h.g.All(t.Context()); err != nil {
		t.Fatal(err)
	}
	repo.archiveErr = fmt.Errorf("%w: link", git.ErrUnsafeEntry)
	repo.fetches = nil

	if _, err := h.g.Find(t.Context(), "alpha"); !errors.Is(err, git.ErrUnsafeEntry) {
		t.Errorf("err = %v, want ErrUnsafeEntry", err)
	}
	if len(repo.fetches) != 0 {
		t.Error("an unsafe archive triggered a refetch")
	}
}

func TestNewBuildsAGitSourceOnlyWithAGateAndACache(t *testing.T) {
	t.Parallel()
	spec := Spec{ID: testID, Type: KindGit, URL: testURL, Rev: "main", Roots: []string{"skills"}}
	gate := func(string, string) error { return nil }

	r, err := New(spec, Options{CacheRoot: t.TempDir(), Gate: gate})
	if err != nil || r.Spec().ID != testID {
		t.Fatalf("New = %v, %v", r, err)
	}
	if _, err = New(spec, Options{CacheRoot: t.TempDir()}); err == nil {
		t.Error("a git source was built without a trust gate")
	}
	if _, err = New(spec, Options{Gate: gate}); err == nil {
		t.Error("a git source was built without a cache location")
	}
	for name, bad := range map[string]Spec{
		"private":    {ID: testID, Type: KindGit, URL: testURL, Private: true},
		"traversal":  {ID: "@acme/..", Type: KindGit, URL: testURL},
		"ext url":    {ID: testID, Type: KindGit, URL: "ext::sh -c id"},
		"file url":   {ID: testID, Type: KindGit, URL: "file:///srv/x.git"},
		"bad root":   {ID: testID, Type: KindGit, URL: testURL, Roots: []string{"../x"}},
		"noncanon":   {ID: "@Acme/Skills", Type: KindGit, URL: testURL},
		"credential": {ID: testID, Type: KindGit, URL: "https://u:ghp_abcdefghijklmnopqrstuvwxyz0123456789@github.com/acme/skills.git"},
	} {
		if _, err = New(bad, Options{CacheRoot: t.TempDir(), Gate: gate}); err == nil {
			t.Errorf("%s: New accepted %+v", name, bad)
		}
	}
	// Other kinds still refuse what they cannot honor.
	if _, err = New(Spec{ID: "@a/a", Type: KindPath, URL: "/x", Rev: "main"}, Options{}); err == nil {
		t.Error("a path source accepted a rev")
	}
}

func TestGlobMatching(t *testing.T) {
	t.Parallel()
	tests := []struct {
		pattern, path string
		want          bool
	}{
		{"**/node_modules/**", "a/node_modules/b", true},
		{"**/node_modules/**", "node_modules/b", true},
		{"**/node_modules/**", "a/node_modules", true},
		{"**/node_modules/**", "a/mynode_modules/b", false},
		{"skills/*/draft", "skills/x/draft", true},
		{"skills/*/draft", "skills/x/y/draft", false},
		{"**/testdata/**", "pkg/testdata/skills/x", true},
		{"vendor/**", "vendor/a/b", true},
		{"vendor/**", "src/vendor/a", false},
		{"a?c", "abc", true},
		{"a?c", "a/c", false},
		{"[x]", "x", false},
		{"a.b", "a.b", true},
		{"a.b", "aXb", false},
	}
	for _, tc := range tests {
		if got := globMatch(tc.pattern, tc.path); got != tc.want {
			t.Errorf("globMatch(%q, %q) = %v, want %v", tc.pattern, tc.path, got, tc.want)
		}
	}
}
