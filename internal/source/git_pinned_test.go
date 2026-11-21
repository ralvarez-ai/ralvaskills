package source

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

func TestGitSkillsCarryWhereTheyCameFrom(t *testing.T) {
	t.Parallel()
	h := newHarness(t, func(s *Spec) { s.Rev = "main" }, baseRepo())

	all, err := h.g.All(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	want := skill.GitRef{URL: testURL, Rev: "main", Commit: testCommit, Tree: "t-skills/alpha", Path: "skills/alpha"}
	if all[0].Git == nil || *all[0].Git != want {
		t.Errorf("All()[0].Git = %+v, want %+v", all[0].Git, want)
	}
	found, err := h.g.Find(t.Context(), "alpha")
	if err != nil {
		t.Fatal(err)
	}
	if found.Git == nil || *found.Git != want {
		t.Errorf("Find().Git = %+v, want %+v", found.Git, want)
	}
}

func tagged(version string) *fakeRepo {
	repo := baseRepo()
	repo.refs = map[string]string{"alpha@v" + version: otherSHA}
	return repo
}

func TestGitFindVersionResolvesThroughATag(t *testing.T) {
	t.Parallel()
	repo := tagged("1.0.0")
	h := newHarness(t, nil, repo)

	got, err := h.g.FindVersion(t.Context(), "alpha", "1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if got.Git == nil || got.Git.Commit != otherSHA || got.Git.Rev != "alpha@v1.0.0" || got.Version != "1.0.0" {
		t.Errorf("skill = %+v, git = %+v; want the tag's commit and rev", got, got.Git)
	}
	if !strings.HasSuffix(got.Path, filepath.Join("alpha", otherSHA)) {
		t.Errorf("path = %s, want the extraction at the tag's commit", got.Path)
	}
	if repo.fetches[0] != "alpha@v1.0.0" {
		t.Errorf("fetched %v", repo.fetches)
	}

	// A leading v in the version names the same tag.
	if _, err = h.g.FindVersion(t.Context(), "alpha", "v1.0.0"); err != nil {
		t.Errorf("v1.0.0: %v", err)
	}
}

func TestGitFindVersionRefusesWhatItCannotTrust(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	h := newHarness(t, nil, tagged("1.0.0"))
	if _, err := h.g.FindVersion(ctx, "alpha", "9.9.9"); !errors.Is(err, ErrNotFound) ||
		!strings.Contains(err.Error(), "alpha@v9.9.9") {
		t.Errorf("a missing tag err = %v, want ErrNotFound naming the tag", err)
	}
	// The tag exists but the skill at it declares another version.
	mismatch := tagged("2.0.0")
	h = newHarness(t, nil, mismatch)
	if _, err := h.g.FindVersion(ctx, "alpha", "2.0.0"); err == nil || !strings.Contains(err.Error(), "1.0.0") {
		t.Errorf("a tag whose skill says 1.0.0 err = %v, want a mismatch error", err)
	}
	for _, bad := range []string{"", "1 0", "--x", "1.0.0;rm", "../x"} {
		if _, err := h.g.FindVersion(ctx, "alpha", bad); err == nil {
			t.Errorf("version %q accepted", bad)
		}
	}
	h.g.offline = true
	if _, err := h.g.FindVersion(ctx, "alpha", "2.0.0"); err == nil {
		t.Error("offline FindVersion succeeded without the network")
	}
	h.g.offline = false
	h.trust = errors.New("source is not trusted")
	repo := tagged("1.0.0")
	h2 := newHarness(t, nil, repo)
	h2.trust = h.trust
	if _, err := h2.g.FindVersion(ctx, "alpha", "1.0.0"); err == nil || len(repo.fetches) != 0 {
		t.Errorf("an untrusted source was fetched: %v %v", err, repo.fetches)
	}
}

func TestGitLockedRestoresTheRecordedCommitAndChecksItsTree(t *testing.T) {
	t.Parallel()
	repo := baseRepo()
	h := newHarness(t, nil, repo)
	locked := LockedSkill{Name: "alpha", Commit: otherSHA, Tree: "t-skills/alpha", Path: "skills/alpha"}

	got, err := h.g.Locked(t.Context(), locked)
	if err != nil {
		t.Fatal(err)
	}
	if got.Git == nil || got.Git.Commit != otherSHA || got.Git.Tree != "t-skills/alpha" ||
		!strings.HasSuffix(got.Path, filepath.Join("alpha", otherSHA)) {
		t.Errorf("skill = %+v git = %+v", got, got.Git)
	}
	if len(repo.listed)+len(repo.fetches) != 0 {
		t.Errorf(
			"restoring from a lock must not index or fetch the branch: listed=%v fetches=%v",
			repo.listed,
			repo.fetches,
		)
	}

	bad := locked
	bad.Tree = "somebody-elses-tree"
	h2 := newHarness(t, nil, baseRepo())
	if _, err = h2.g.Locked(t.Context(), bad); !errors.Is(err, ErrTreeMismatch) {
		t.Errorf("a tree that is not the recorded one err = %v, want ErrTreeMismatch", err)
	}
	if _, statErr := os.Stat(h2.g.extractDir("alpha", otherSHA)); statErr == nil {
		t.Error("an extraction that failed its integrity check was left in the cache")
	}
}

func TestGitLockedValidatesAndIsGated(t *testing.T) {
	t.Parallel()
	repo := baseRepo()
	h := newHarness(t, nil, repo)
	for name, l := range map[string]LockedSkill{
		"bad commit":   {Name: "alpha", Commit: "main", Path: "skills/alpha"},
		"bad path":     {Name: "alpha", Commit: otherSHA, Path: "../x"},
		"bad name":     {Name: "Alpha", Commit: otherSHA, Path: "skills/alpha"},
		"name vs path": {Name: "alpha", Commit: otherSHA, Path: "skills/beta"},
	} {
		if _, err := h.g.Locked(t.Context(), l); err == nil {
			t.Errorf("%s was accepted", name)
		}
	}

	h.trust = errors.New("source is not trusted")
	if _, err := h.g.Locked(
		t.Context(),
		LockedSkill{Name: "alpha", Commit: otherSHA, Path: "skills/alpha"},
	); err == nil {
		t.Error("an untrusted source was restored")
	}
	if len(repo.archived)+len(repo.fetches) != 0 {
		t.Errorf("an untrusted source reached git: %v %v", repo.archived, repo.fetches)
	}
}

func TestGitLockedFetchesAMissingCommitUnlessOffline(t *testing.T) {
	t.Parallel()
	repo := baseRepo()
	repo.missing = true
	h := newHarness(t, nil, repo)
	locked := LockedSkill{Name: "alpha", Commit: otherSHA, Path: "skills/alpha"}

	h.g.offline = true
	if _, err := h.g.Locked(t.Context(), locked); err == nil || len(repo.fetches) != 0 {
		t.Errorf("offline: err=%v fetches=%v, want a failure without a fetch", err, repo.fetches)
	}
	h.g.offline = false
	if _, err := h.g.Locked(t.Context(), locked); err != nil {
		t.Fatalf("Locked: %v", err)
	}
	if len(repo.fetches) != 1 || repo.fetches[0] != otherSHA {
		t.Errorf("fetches = %v, want a fetch of the locked commit", repo.fetches)
	}

	gone := baseRepo()
	gone.missing = true
	gone.fetchErr = errors.New("not our ref")
	h = newHarness(t, nil, gone)
	if _, err := h.g.Locked(t.Context(), locked); err == nil || !strings.Contains(err.Error(), otherSHA) {
		t.Errorf("an unreachable commit err = %v, want it to name the commit", err)
	}
}

func TestGitTreeAtReportsTheInstalledTree(t *testing.T) {
	t.Parallel()
	repo := baseRepo()
	repo.trees = map[string]string{"skills/alpha": "old-tree"}
	h := newHarness(t, nil, repo)

	got, err := h.g.TreeAt(t.Context(), "alpha", otherSHA)
	if err != nil || got != "old-tree" {
		t.Errorf("TreeAt = %q, %v", got, err)
	}
	if _, err = h.g.TreeAt(t.Context(), "ghost", otherSHA); !errors.Is(err, ErrNotFound) {
		t.Errorf("an unknown name err = %v, want ErrNotFound", err)
	}
}
