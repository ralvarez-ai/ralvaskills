package manifest_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
)

const (
	privateID   = "@acme/hidden"
	privateURL  = "git@github.com:acme/hidden.git"
	secretName  = "secret-skill"
	commitSHA   = "dddddddddddddddddddddddddddddddddddddddd"
	treeSHA     = "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	secretToken = "acme/hidden"
)

// project is a git repository with an .rsk folder.
func project(t *testing.T) (root, rskDir string) {
	t.Helper()
	root = t.TempDir()
	if out, err := exec.CommandContext(t.Context(), "git", "-C", root, "init", "-q").CombinedOutput(); err != nil {
		t.Skipf("git is not available: %v %s", err, out)
	}
	return root, filepath.Join(root, ".rsk")
}

func privateEntry() manifest.LockEntry {
	return manifest.LockEntry{
		Name: secretName, Version: "1.0.0", Source: skill.DeclaredSource(privateID), Local: true,
		Git: &skill.GitRef{URL: privateURL, Commit: commitSHA, Tree: treeSHA, Path: "skills/" + secretName},
	}
}

func TestAPrivateSkillIsWrittenOnlyToTheLocalFiles(t *testing.T) {
	t.Parallel()
	root, rsk := project(t)
	m := manifest.Mod{
		Tools: []tool.ID{tool.ClaudeID}, Skills: map[string]string{"public-one": "*", secretName: "*"},
		Pinned: []string{"public-one", secretName}, Local: map[string]bool{secretName: true},
	}
	m.SetOrigin(secretName, skill.DeclaredSource(privateID))
	lock := manifest.Lock{Skills: []manifest.LockEntry{
		{Name: "public-one", Version: "2.0.0", Source: skill.SourceRegistry}, privateEntry(),
	}}
	if err := manifest.WriteMod(rsk, m); err != nil {
		t.Fatal(err)
	}
	if err := manifest.WriteLock(rsk, lock); err != nil {
		t.Fatal(err)
	}

	for _, committed := range []string{"rsk.mod", "rsk.lock"} {
		data, err := os.ReadFile(filepath.Join(rsk, committed))
		if err != nil {
			t.Fatal(err)
		}
		for _, leak := range []string{secretName, secretToken, privateURL, commitSHA} {
			if strings.Contains(string(data), leak) {
				t.Errorf("%s carries %q:\n%s", committed, leak, data)
			}
		}
		if !strings.Contains(string(data), "public-one") {
			t.Errorf("%s lost the public skill:\n%s", committed, data)
		}
	}
	for _, local := range []string{"rsk.local.mod", "rsk.local.lock"} {
		data, err := os.ReadFile(filepath.Join(rsk, local))
		if err != nil || !strings.Contains(string(data), secretName) {
			t.Errorf("%s = %q, %v", local, data, err)
		}
	}

	exclude, err := os.ReadFile(filepath.Join(root, ".git", "info", "exclude"))
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range []string{"/.rsk/rsk.local.mod", "/.rsk/rsk.local.lock", "/.rsk/SKILLS.local.md", "/CLAUDE.local.md"} {
		if !slices.Contains(strings.Split(string(exclude), "\n"), line) {
			t.Errorf(".git/info/exclude lacks %q:\n%s", line, exclude)
		}
	}
	if out, gitErr := exec.CommandContext(t.Context(), "git", "-C", root, "status", "--porcelain", "--ignored=no").
		Output(); gitErr == nil &&
		strings.Contains(string(out), "local") {
		t.Errorf("git sees a local file as untracked:\n%s", out)
	}
}

func TestReadMergesBothLayersAndMarksTheLocalOnes(t *testing.T) {
	t.Parallel()
	_, rsk := project(t)
	m := manifest.Mod{
		Tools: []tool.ID{tool.ClaudeID}, Skills: map[string]string{"public-one": "*", secretName: "1.0.0"},
		Pinned: []string{secretName}, Local: map[string]bool{secretName: true},
	}
	m.SetOrigin(secretName, skill.DeclaredSource(privateID))
	if err := manifest.WriteMod(rsk, m); err != nil {
		t.Fatal(err)
	}
	if err := manifest.WriteLock(rsk, manifest.Lock{Skills: []manifest.LockEntry{privateEntry()}}); err != nil {
		t.Fatal(err)
	}

	got, err := manifest.ReadMod(rsk)
	if err != nil || got.Skills[secretName] != "1.0.0" || !got.Local[secretName] || got.Local["public-one"] ||
		got.Origin[secretName] != privateID || !slices.Contains(got.Pinned, secretName) {
		t.Fatalf("mod = %+v, %v", got, err)
	}
	lock, err := manifest.ReadLock(rsk)
	if err != nil || len(lock.Skills) != 1 || !lock.Skills[0].Local || lock.Skills[0].Git == nil {
		t.Fatalf("lock = %+v, %v", lock, err)
	}

	// Writing back what was read keeps every skill in the layer it came from.
	if err = manifest.WriteMod(rsk, got); err != nil || manifest.WriteLock(rsk, lock) != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(filepath.Join(rsk, "rsk.mod"))
	if strings.Contains(string(data), secretName) {
		t.Errorf("a round trip moved the private skill into rsk.mod:\n%s", data)
	}
}

func TestATeammateWithoutTheLocalFilesSeesACleanProject(t *testing.T) {
	t.Parallel()
	_, rsk := project(t)
	m := manifest.Mod{Tools: []tool.ID{tool.ClaudeID}, Skills: map[string]string{"public-one": "*"}}
	if err := manifest.WriteMod(rsk, m); err != nil {
		t.Fatal(err)
	}
	got, err := manifest.ReadMod(rsk)
	if err != nil || len(got.Skills) != 1 || len(got.Local) != 0 {
		t.Fatalf("mod = %+v, %v", got, err)
	}
	if lock, lockErr := manifest.ReadLock(rsk); lockErr != nil || len(lock.Skills) != 0 {
		t.Errorf("lock = %+v, %v", lock, lockErr)
	}
	for _, name := range []string{"rsk.local.mod", "rsk.local.lock"} {
		if _, statErr := os.Stat(filepath.Join(rsk, name)); statErr == nil {
			t.Errorf("%s was written for a project with no private skill", name)
		}
	}
}

func TestDroppingTheLastPrivateSkillRemovesTheLocalFiles(t *testing.T) {
	t.Parallel()
	_, rsk := project(t)
	m := manifest.Mod{
		Tools: []tool.ID{
			tool.ClaudeID,
		},
		Skills: map[string]string{secretName: "*"},
		Local:  map[string]bool{secretName: true},
	}
	m.SetOrigin(secretName, skill.DeclaredSource(privateID))
	if err := manifest.WriteMod(
		rsk,
		m,
	); err != nil ||
		manifest.WriteLock(rsk, manifest.Lock{Skills: []manifest.LockEntry{privateEntry()}}) != nil {
		t.Fatal(err)
	}
	m.Skills, m.Local = map[string]string{}, nil
	if err := manifest.WriteMod(rsk, m); err != nil || manifest.WriteLock(rsk, manifest.Lock{}) != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"rsk.local.mod", "rsk.local.lock"} {
		if _, err := os.Stat(filepath.Join(rsk, name)); err == nil {
			t.Errorf("%s is left behind with no private skill", name)
		}
	}
}

func TestANameInBothLayersIsAnError(t *testing.T) {
	t.Parallel()
	_, rsk := project(t)
	if err := os.MkdirAll(rsk, 0o750); err != nil {
		t.Fatal(err)
	}
	both := "version = 2\ntools = [\"claude-code\"]\npinned = []\n[skills]\n  \"" + privateID + "/" + secretName + "\" = \"*\"\n"
	for _, name := range []string{"rsk.mod", "rsk.local.mod"} {
		if err := os.WriteFile(filepath.Join(rsk, name), []byte(both), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := manifest.ReadMod(rsk); err == nil {
		t.Error("a skill in both layers was accepted")
	}
}

func TestWritingOutsideAGitRepositoryStillWorks(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	rsk := filepath.Join(root, ".rsk")
	m := manifest.Mod{
		Tools: []tool.ID{
			tool.ClaudeID,
		},
		Skills: map[string]string{secretName: "*"},
		Local:  map[string]bool{secretName: true},
	}
	m.SetOrigin(secretName, skill.DeclaredSource(privateID))
	if err := manifest.WriteMod(rsk, m); err != nil {
		t.Fatalf("a project that is not a git repository could not keep a private skill: %v", err)
	}
	if _, err := os.Stat(filepath.Join(rsk, "rsk.local.mod")); err != nil {
		t.Error(err)
	}
}
