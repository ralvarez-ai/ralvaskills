package source_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
)

const pathSkillBody = "---\nname: %s\nversion: 1.2.0\ndescription: A test skill.\n---\nBody\n"

func mkSkill(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o750); err != nil {
		t.Fatal(err)
	}
	body := strings.Replace(pathSkillBody, "%s", filepath.Base(dir), 1)
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

//nolint:ireturn // the test reaches the resolver the way production does, through the factory
func newPathSource(t *testing.T, root string, mod func(*source.Spec)) source.Resolver {
	t.Helper()
	spec := source.Spec{ID: "@me/scratch", Type: source.KindPath, URL: root, Declared: true}
	if mod != nil {
		mod(&spec)
	}
	r, err := source.New(spec, source.Options{})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return r
}

func skillNames(skills []skill.Skill) []string {
	out := make([]string, len(skills))
	for i := range skills {
		out[i] = skills[i].Name
	}
	return out
}

func TestADeclaredPathSourceIsWalkedAsGivenNotUnderSkills(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	mkSkill(t, filepath.Join(root, "alpha"))
	mkSkill(t, filepath.Join(root, "group", "beta"))
	mkSkill(t, filepath.Join(root, "skills", "gamma"))

	got, err := newPathSource(t, root, nil).All(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if want := "alpha beta gamma"; strings.Join(skillNames(got), " ") != want {
		t.Fatalf("skills = %v, want %s", skillNames(got), want)
	}
	for i := range got {
		if got[i].Source != skill.SourcePath || got[i].Version != "1.2.0" || got[i].Description == "" ||
			!filepath.IsAbs(got[i].Path) {
			t.Errorf("skill %+v", got[i])
		}
	}
}

func TestAPathSourceHonoursRootsAndExcludes(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	mkSkill(t, filepath.Join(root, "a", "one"))
	mkSkill(t, filepath.Join(root, "b", "two"))
	mkSkill(t, filepath.Join(root, "b", "drafts", "three"))

	got, err := newPathSource(t, root, func(s *source.Spec) {
		s.Roots, s.Excludes = []string{"b"}, []string{"**/drafts/**"}
	}).All(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(skillNames(got), " ") != "two" {
		t.Errorf("skills = %v, want only two", skillNames(got))
	}
}

func TestAPathSourceNeverFollowsOrKeepsASymlink(t *testing.T) {
	t.Parallel()
	root, outside := t.TempDir(), t.TempDir()
	mkSkill(t, filepath.Join(root, "good"))
	mkSkill(t, filepath.Join(outside, "escaped"))
	mkSkill(t, filepath.Join(root, "tainted"))
	if err := os.Symlink(outside, filepath.Join(root, "tainted", "link")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(outside, "escaped"), filepath.Join(root, "linked-dir")); err != nil {
		t.Fatal(err)
	}
	mkSkill(t, filepath.Join(root, "md-link"))
	if err := os.Remove(filepath.Join(root, "md-link", "SKILL.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(
		filepath.Join(outside, "escaped", "SKILL.md"),
		filepath.Join(root, "md-link", "SKILL.md"),
	); err != nil {
		t.Fatal(err)
	}

	src := newPathSource(t, root, nil)
	got, err := src.All(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(skillNames(got), " ") != "good" {
		t.Fatalf("skills = %v, want only good", skillNames(got))
	}
	scanner, ok := src.(interface {
		Skipped(context.Context) ([]source.Skipped, error)
	})
	if !ok {
		t.Fatal("a path source does not report what it skipped")
	}
	skipped, err := scanner.Skipped(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	reasons := map[string]string{}
	for _, s := range skipped {
		reasons[s.Path] = s.Reason
	}
	if !strings.Contains(reasons["tainted"], "symlink") || !strings.Contains(reasons["md-link"], "regular file") {
		t.Errorf("skipped = %v", reasons)
	}
	if _, found := reasons["linked-dir"]; found {
		t.Errorf("a symlinked directory was inspected: %v", reasons)
	}
}

func TestAPathSourceRefusesARootThatEscapesOrIsALink(t *testing.T) {
	t.Parallel()
	root, outside := t.TempDir(), t.TempDir()
	mkSkill(t, filepath.Join(outside, "escaped"))
	if err := os.Symlink(outside, filepath.Join(root, "door")); err != nil {
		t.Fatal(err)
	}

	for name, roots := range map[string][]string{
		"dot dot":      {"../x"},
		"absolute":     {"/etc"},
		"through link": {"door"},
	} {
		spec := source.Spec{ID: "@me/scratch", Type: source.KindPath, URL: root, Declared: true, Roots: roots}
		r, err := source.New(spec, source.Options{})
		if err == nil {
			_, err = r.All(context.Background())
		}
		if err == nil {
			t.Errorf("%s: root %v was accepted", name, roots)
		}
	}
}

func TestAPathSourceFailsOnTwoSkillsWithOneName(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	mkSkill(t, filepath.Join(root, "a", "same"))
	mkSkill(t, filepath.Join(root, "b", "same"))

	_, err := newPathSource(t, root, nil).All(context.Background())
	if !errors.Is(err, source.ErrDuplicate) || !strings.Contains(err.Error(), "a/same") ||
		!strings.Contains(err.Error(), "b/same") {
		t.Errorf("err = %v, want ErrDuplicate naming both", err)
	}
}

func TestAPathSourceSkipsABadNameAndFindsByName(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	mkSkill(t, filepath.Join(root, "Not_Kebab"))
	mkSkill(t, filepath.Join(root, "fine"))

	src := newPathSource(t, root, nil)
	got, err := src.All(context.Background())
	if err != nil || strings.Join(skillNames(got), " ") != "fine" {
		t.Fatalf("skills = %v, %v", skillNames(got), err)
	}
	if s, findErr := src.Find(context.Background(), "fine"); findErr != nil || s.Name != "fine" {
		t.Errorf("Find fine = %+v, %v", s, findErr)
	}
	if _, findErr := src.Find(context.Background(), "absent"); !errors.Is(findErr, source.ErrNotFound) {
		t.Errorf("Find absent = %v, want ErrNotFound", findErr)
	}
}

func TestOnlyADeclaredPathSourceMayCarryRootsAndExcludes(t *testing.T) {
	t.Parallel()
	for name, mod := range map[string]func(*source.Spec){
		"roots":    func(s *source.Spec) { s.Roots = []string{"s"} },
		"excludes": func(s *source.Spec) { s.Excludes = []string{"x"} },
	} {
		spec := source.Spec{ID: "@local/clone", Type: source.KindPath, URL: t.TempDir()}
		mod(&spec)
		if _, err := source.New(spec, source.Options{}); err == nil {
			t.Errorf("the local clone accepted %s", name)
		}
	}
	spec := source.Spec{ID: "@me/scratch", Type: source.KindPath, URL: t.TempDir(), Declared: true, Rev: "main"}
	if _, err := source.New(spec, source.Options{}); err == nil {
		t.Error("a path source accepted a rev")
	}
}

func TestScanSkillDirReadsOneSkillDirectoryUnderTheSameRules(t *testing.T) {
	t.Parallel()
	root, outside := t.TempDir(), t.TempDir()
	mkSkill(t, filepath.Join(root, "good-one"))
	mkSkill(t, filepath.Join(root, "tainted"))
	if err := os.Symlink(outside, filepath.Join(root, "tainted", "link")); err != nil {
		t.Fatal(err)
	}
	mkSkill(t, filepath.Join(root, "Not_Kebab"))
	if err := os.Symlink(filepath.Join(root, "good-one"), filepath.Join(root, "alias")); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "group", "inner"), 0o750); err != nil {
		t.Fatal(err)
	}
	mkSkill(t, filepath.Join(root, "group", "inner", "deep"))

	got, err := source.ScanSkillDir(filepath.Join(root, "good-one"))
	if err != nil || got.Name != "good-one" || got.Version != "1.2.0" || got.Source != skill.SourcePath ||
		got.Path != filepath.Join(root, "good-one") {
		t.Fatalf("skill = %+v, %v", got, err)
	}
	for name, dir := range map[string]string{
		"a symlink inside":      filepath.Join(root, "tainted"),
		"a bad name":            filepath.Join(root, "Not_Kebab"),
		"a symlinked directory": filepath.Join(root, "alias"),
		"a directory of skills": filepath.Join(root, "group"),
		"no such directory":     filepath.Join(root, "absent"),
		"a file":                filepath.Join(root, "good-one", "SKILL.md"),
	} {
		if s, scanErr := source.ScanSkillDir(dir); scanErr == nil {
			t.Errorf("%s was accepted: %+v", name, s)
		}
	}
	if _, err = source.ScanSkillDir(
		filepath.Join(root, "group"),
	); err == nil ||
		!strings.Contains(err.Error(), "rsk source add") {
		t.Errorf("a directory of skills should point at 'rsk source add': %v", err)
	}
}
