package git

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
}

func output(t *testing.T, dir string, args ...string) string {
	t.Helper()
	c := exec.CommandContext(t.Context(), "git", args...)
	c.Dir = dir
	out, err := c.Output()
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
	return strings.TrimSpace(string(out))
}

// fixture is a source repository with one skill, a script, and (optionally)
// the entries a hostile repository would carry.
type fixture struct {
	dir    string
	remote RemoteURL
}

func newFixture(t *testing.T) fixture {
	t.Helper()
	requireGit(t)
	ctx := t.Context()
	dir := t.TempDir()
	run(ctx, t, dir, "init", "-b", "main")
	for path, body := range map[string]string{
		"skills/alpha/SKILL.md":            "alpha",
		"skills/alpha/refs/notes.md":       "notes",
		"skills/beta/SKILL.md":             "beta",
		"packages/x/skills/gamma/SKILL.md": "gamma",
		"README.md":                        "readme",
	} {
		full := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	script := filepath.Join(dir, "skills", "alpha", "run.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	run(ctx, t, dir, "add", ".")
	run(ctx, t, dir, "commit", "-m", "seed")
	return fixture{dir: dir, remote: RemoteURL{raw: dir}}
}

func (f fixture) commit(t *testing.T) string {
	t.Helper()
	return output(t, f.dir, "rev-parse", "HEAD")
}

func newRepo(t *testing.T) Repo {
	t.Helper()
	return Repo{Dir: filepath.Join(t.TempDir(), "git", "acme", "skills.git"), allowFile: true}
}

func TestFetchBranchTagAndSHA(t *testing.T) {
	t.Parallel()
	f := newFixture(t)
	ctx := t.Context()
	want := f.commit(t)
	run(ctx, t, f.dir, "-c", "tag.gpgSign=false", "tag", "-a", "-m", "release", "alpha@v1.0.0")

	for _, ref := range []string{"main", "alpha@v1.0.0", want} {
		repo := newRepo(t)
		got, err := repo.Fetch(ctx, f.remote, ref)
		if err != nil {
			t.Fatalf("Fetch(%q): %v", ref, err)
		}
		if got != want {
			t.Errorf("Fetch(%q) = %s, want %s", ref, got, want)
		}
	}

	repo := newRepo(t)
	if _, err := repo.Fetch(ctx, f.remote, "no-such-branch"); !errors.Is(err, ErrRefNotFound) {
		t.Errorf("fetching a missing ref err = %v, want ErrRefNotFound", err)
	}
}

func TestFetchKeepsTheCommitReachable(t *testing.T) {
	t.Parallel()
	f := newFixture(t)
	repo := newRepo(t)
	commit, err := repo.Fetch(t.Context(), f.remote, "main")
	if err != nil {
		t.Fatal(err)
	}
	pinned := output(t, repo.Dir, "for-each-ref", "--format=%(objectname)", "refs/rsk/")
	if !strings.Contains(pinned, commit) {
		t.Errorf("commit %s is not pinned under refs/rsk/ and could be garbage collected; refs: %q", commit, pinned)
	}
}

func TestFetchRefusesWhatCouldInjectAnOption(t *testing.T) {
	t.Parallel()
	f := newFixture(t)
	for _, ref := range []string{
		"", "--upload-pack=touch /tmp/x", "-x", "a b", "a:b", "a*", "a..b", "a\nb", "main^", "a~1", "x[",
	} {
		if _, err := newRepo(t).Fetch(t.Context(), f.remote, ref); !errors.Is(err, ErrInvalidRef) {
			t.Errorf("Fetch(%q) err = %v, want ErrInvalidRef", ref, err)
		}
	}
}

func TestListTreeFiltersByRootsAndFlagsDangerousModes(t *testing.T) {
	t.Parallel()
	f := newFixture(t)
	ctx := t.Context()
	if err := os.Symlink("/etc/passwd", filepath.Join(f.dir, "skills", "alpha", "link")); err != nil {
		t.Fatal(err)
	}
	run(ctx, t, f.dir, "add", ".")
	run(ctx, t, f.dir, "update-index", "--add", "--cacheinfo", "160000,"+strings.Repeat("a", 40)+",skills/sub")
	run(ctx, t, f.dir, "commit", "-m", "hostile")

	repo := newRepo(t)
	commit, err := repo.Fetch(ctx, f.remote, "main")
	if err != nil {
		t.Fatal(err)
	}

	entries, err := repo.ListTree(ctx, commit, []string{"skills"})
	if err != nil {
		t.Fatal(err)
	}
	modes := map[string]string{}
	for _, e := range entries {
		modes[e.Path] = e.Mode
		if e.Path != "skills" && !strings.HasPrefix(e.Path, "skills/") {
			t.Errorf("%s is outside the requested root", e.Path)
		}
	}
	if modes["skills/alpha"] != "040000" {
		t.Errorf("directory entries are missing: skills/alpha mode = %q", modes["skills/alpha"])
	}
	if modes["skills/alpha/SKILL.md"] != "100644" || modes["skills/alpha/run.sh"] != "100755" {
		t.Errorf("modes = %v", modes)
	}
	if modes["skills/alpha/link"] != ModeSymlink {
		t.Errorf("symlink mode = %q, want %q", modes["skills/alpha/link"], ModeSymlink)
	}
	if modes["skills/sub"] != ModeSubmodule {
		t.Errorf("submodule mode = %q, want %q", modes["skills/sub"], ModeSubmodule)
	}

	all, err := repo.ListTree(ctx, commit, nil)
	if err != nil {
		t.Fatal(err)
	}
	var sawReadme bool
	for _, e := range all {
		sawReadme = sawReadme || e.Path == "README.md"
	}
	if !sawReadme {
		t.Error("no roots should list the whole tree")
	}
}

func TestListTreeRefusesBadInput(t *testing.T) {
	t.Parallel()
	f := newFixture(t)
	repo := newRepo(t)
	commit, err := repo.Fetch(t.Context(), f.remote, "main")
	if err != nil {
		t.Fatal(err)
	}
	for _, roots := range [][]string{{"../x"}, {"/abs"}, {"a/../../b"}, {"-x"}, {""}, {"a\x00b"}} {
		if _, err = repo.ListTree(t.Context(), commit, roots); !errors.Is(err, ErrInvalidPath) {
			t.Errorf("ListTree roots %q err = %v, want ErrInvalidPath", roots, err)
		}
	}
	for _, bad := range []string{"", "main", "--all", "abc", strings.Repeat("g", 40)} {
		if _, err = repo.ListTree(t.Context(), bad, nil); !errors.Is(err, ErrInvalidCommit) {
			t.Errorf("ListTree commit %q err = %v, want ErrInvalidCommit", bad, err)
		}
	}
}

func TestTreeIDOfASkillIgnoresUnrelatedCommits(t *testing.T) {
	t.Parallel()
	f := newFixture(t)
	ctx := t.Context()
	repo := newRepo(t)
	first, err := repo.Fetch(ctx, f.remote, "main")
	if err != nil {
		t.Fatal(err)
	}
	before, err := repo.TreeID(ctx, first, "skills/alpha")
	if err != nil {
		t.Fatal(err)
	}
	if want := output(t, f.dir, "rev-parse", first+":skills/alpha"); before != want {
		t.Errorf("TreeID = %s, want %s", before, want)
	}

	if err = os.WriteFile(filepath.Join(f.dir, "README.md"), []byte("changed"), 0o644); err != nil {
		t.Fatal(err)
	}
	run(ctx, t, f.dir, "commit", "-am", "unrelated")
	second, err := repo.Fetch(ctx, f.remote, "main")
	if err != nil {
		t.Fatal(err)
	}
	after, err := repo.TreeID(ctx, second, "skills/alpha")
	if err != nil {
		t.Fatal(err)
	}
	if second == first || after != before {
		t.Errorf("commits %s/%s tree %s/%s: an unrelated commit moved the skill's tree", first, second, before, after)
	}

	if _, err = repo.TreeID(ctx, second, "skills/ghost"); err == nil {
		t.Error("TreeID of a missing path succeeded")
	}
}

func TestArchiveExtractsOneSkillDirectory(t *testing.T) {
	t.Parallel()
	f := newFixture(t)
	repo := newRepo(t)
	commit, err := repo.Fetch(t.Context(), f.remote, "main")
	if err != nil {
		t.Fatal(err)
	}

	dest := filepath.Join(t.TempDir(), "extracted", "alpha")
	if err = repo.Archive(t.Context(), commit, "skills/alpha", dest, Limits{}); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(filepath.Join(dest, "SKILL.md")); string(got) != "alpha" {
		t.Errorf("SKILL.md = %q", got)
	}
	if got, _ := os.ReadFile(filepath.Join(dest, "refs", "notes.md")); string(got) != "notes" {
		t.Errorf("refs/notes.md = %q", got)
	}
	if info, statErr := os.Stat(filepath.Join(dest, "run.sh")); statErr != nil || info.Mode().Perm()&0o100 == 0 {
		t.Errorf("run.sh lost its exec bit: %v %v", info, statErr)
	}
	if _, statErr := os.Stat(filepath.Join(dest, "skills")); statErr == nil {
		t.Error("the directory prefix was not stripped")
	}
	if _, statErr := os.Stat(filepath.Join(dest, "..", "beta")); statErr == nil {
		t.Error("a sibling skill was extracted")
	}
}

func TestArchiveRefusesASkillThatContainsASymlink(t *testing.T) {
	t.Parallel()
	f := newFixture(t)
	ctx := t.Context()
	if err := os.Symlink("/etc/passwd", filepath.Join(f.dir, "skills", "alpha", "link")); err != nil {
		t.Fatal(err)
	}
	run(ctx, t, f.dir, "add", ".")
	run(ctx, t, f.dir, "commit", "-m", "link")

	repo := newRepo(t)
	commit, err := repo.Fetch(ctx, f.remote, "main")
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	err = repo.Archive(ctx, commit, "skills/alpha", filepath.Join(root, "alpha"), Limits{})
	if !errors.Is(err, ErrUnsafeEntry) {
		t.Fatalf("err = %v, want ErrUnsafeEntry", err)
	}
	if entries, _ := os.ReadDir(root); len(entries) != 0 {
		t.Errorf("a refused skill left %v behind", entries)
	}
}

func TestGitFailureDoesNotLeakACredential(t *testing.T) {
	t.Parallel()
	requireGit(t)
	// A URL with credentials never reaches git: the parser refuses it.
	if _, err := ParseRemoteURL("https://user:ghp_secret0123456789abcdefghijkl@github.com/a/b.git"); err == nil {
		t.Fatal("credentials in a URL were accepted")
	}
}

func TestVersionIsAtLeastTheMinimum(t *testing.T) {
	t.Parallel()
	requireGit(t)
	v, err := Version(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if !versionAtLeast(v, MinVersion) {
		t.Skipf("git %v is older than the minimum %v the new transports need", v, MinVersion)
	}
}

func TestProductionRunnerRefusesTheFileProtocol(t *testing.T) {
	t.Parallel()
	f := newFixture(t)
	repo := Repo{Dir: filepath.Join(t.TempDir(), "x.git")} // allowFile left false
	if _, err := repo.Fetch(t.Context(), f.remote, "main"); err == nil {
		t.Fatal("the production runner fetched from a local path; protocol.file.allow=never is not in effect")
	}
}

func TestTreeIDRefusesAFile(t *testing.T) {
	t.Parallel()
	f := newFixture(t)
	repo := newRepo(t)
	commit, err := repo.Fetch(t.Context(), f.remote, "main")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.TreeID(t.Context(), commit, "skills/alpha/SKILL.md"); !errors.Is(err, ErrInvalidPath) {
		t.Errorf("TreeID of a file err = %v, want ErrInvalidPath", err)
	}
}

func TestLsRemoteResolvesARefWithoutCloning(t *testing.T) {
	t.Parallel()
	f := newFixture(t)
	want := f.commit(t)
	run := runner{allowFile: true}

	got, err := run.lsRemote(t.Context(), f.remote, "main")
	if err != nil || got != want {
		t.Fatalf("lsRemote(main) = %q, %v; want %q", got, err, want)
	}
	if got, err = run.lsRemote(t.Context(), f.remote, ""); err != nil || got != want {
		t.Errorf("lsRemote(HEAD) = %q, %v; want %q", got, err, want)
	}
	if _, err = run.lsRemote(t.Context(), f.remote, "nope"); !errors.Is(err, ErrRefNotFound) {
		t.Errorf("a missing ref err = %v, want ErrRefNotFound", err)
	}
	if _, err = run.lsRemote(t.Context(), f.remote, "--upload-pack=x"); !errors.Is(err, ErrInvalidRef) {
		t.Errorf("an option-looking ref err = %v, want ErrInvalidRef", err)
	}
	if _, err = (runner{}).lsRemote(t.Context(), f.remote, "main"); err == nil {
		t.Error("the production runner reached a local path")
	}
}

func TestListTreeReportsTheTreeIDOfEachDirectory(t *testing.T) {
	t.Parallel()
	f := newFixture(t)
	repo := newRepo(t)
	commit, err := repo.Fetch(t.Context(), f.remote, "main")
	if err != nil {
		t.Fatal(err)
	}
	entries, err := repo.ListTree(t.Context(), commit, []string{"skills"})
	if err != nil {
		t.Fatal(err)
	}
	want, err := repo.TreeID(t.Context(), commit, "skills/alpha")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Path == "skills/alpha" && (e.Type != "tree" || e.ObjectID != want) {
			t.Errorf("entry = %+v, want tree %s", e, want)
		}
	}
}

func TestReadBlobs(t *testing.T) {
	t.Parallel()
	f := newFixture(t)
	repo := newRepo(t)
	commit, err := repo.Fetch(t.Context(), f.remote, "main")
	if err != nil {
		t.Fatal(err)
	}

	paths := []string{"skills/alpha/SKILL.md", "README.md", "skills/ghost/SKILL.md", "skills/beta/SKILL.md"}
	got, err := repo.ReadBlobs(t.Context(), commit, paths, 1024)
	if err != nil {
		t.Fatal(err)
	}
	if string(got["skills/alpha/SKILL.md"].Data) != "alpha" || string(got["README.md"].Data) != "readme" ||
		string(got["skills/beta/SKILL.md"].Data) != "beta" {
		t.Errorf("blobs = %+v", got)
	}
	if !got["skills/ghost/SKILL.md"].Missing {
		t.Error("a missing file was not reported as missing")
	}

	small, err := repo.ReadBlobs(t.Context(), commit, []string{"README.md", "skills/beta/SKILL.md"}, 4)
	if err != nil {
		t.Fatal(err)
	}
	if small["README.md"].Data != nil || small["README.md"].Size != 6 {
		t.Errorf("an oversized file must report its size and no data: %+v", small["README.md"])
	}
	if string(small["skills/beta/SKILL.md"].Data) != "beta" {
		t.Errorf("a file under the limit lost its data after an oversized one: %+v", small["skills/beta/SKILL.md"])
	}

	if _, err = repo.ReadBlobs(t.Context(), commit, []string{"../x"}, 10); !errors.Is(err, ErrInvalidPath) {
		t.Errorf("a traversing path err = %v, want ErrInvalidPath", err)
	}
	if got, err = repo.ReadBlobs(t.Context(), commit, nil, 10); err != nil || len(got) != 0 {
		t.Errorf("no paths: %v, %v", got, err)
	}
}

func TestArchiveLeavesNothingWhenGitFails(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	dest := filepath.Join(root, "out", "alpha")
	// A clone that does not exist makes git archive fail before it writes a byte;
	// an empty stream is a valid tar, so extraction itself succeeds.
	repo := Repo{Dir: filepath.Join(root, "no-such-clone.git")}

	err := repo.Archive(t.Context(), strings.Repeat("a", 40), "skills/alpha", dest, Limits{})
	if err == nil {
		t.Fatal("archiving from a missing clone succeeded")
	}
	if _, statErr := os.Stat(dest); statErr == nil {
		t.Error("a failed archive left an (empty) destination, which blocks the retry")
	}
	if retryErr := repo.Archive(
		t.Context(),
		strings.Repeat("a", 40),
		"skills/alpha",
		dest,
		Limits{},
	); errors.Is(
		retryErr,
		os.ErrExist,
	) {
		t.Errorf("the retry was refused as if the destination existed: %v", retryErr)
	}
}

func TestArchiveNeverTouchesAnExistingDestination(t *testing.T) {
	t.Parallel()
	f := newFixture(t)
	repo := newRepo(t)
	commit, err := repo.Fetch(t.Context(), f.remote, "main")
	if err != nil {
		t.Fatal(err)
	}
	dest := t.TempDir()
	keep := filepath.Join(dest, "keep.txt")
	if err = os.WriteFile(keep, []byte("mine"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err = repo.Archive(t.Context(), commit, "skills/alpha", dest, Limits{}); err == nil {
		t.Fatal("archiving over an existing directory succeeded")
	}
	if _, statErr := os.Stat(keep); statErr != nil {
		t.Errorf("a refused archive deleted what was already there: %v", statErr)
	}
}
