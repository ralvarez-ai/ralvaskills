package sources_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/schema"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
	"go.ralva.dev/ralvaskills/v3/internal/sources"
)

func TestNewEntryFromAGitHubIdentity(t *testing.T) {
	t.Parallel()
	e, err := sources.NewEntry("@Acme/Skills-Private", sources.Options{Rev: "main", Roots: []string{"skills"}})
	if err != nil {
		t.Fatal(err)
	}
	want := sources.Entry{
		ID: "@acme/skills-private", Type: source.KindGit,
		URL: "git@github.com:Acme/Skills-Private.git", Rev: "main", Roots: []string{"skills"},
	}
	if e.ID != want.ID || e.Type != want.Type || e.URL != want.URL || e.Rev != want.Rev || len(e.Roots) != 1 {
		t.Errorf("entry = %+v, want %+v", e, want)
	}
}

func TestNewEntryDerivesTheIdentityFromAGitHubURL(t *testing.T) {
	t.Parallel()
	for _, url := range []string{
		"git@github.com:acme/skills.git",
		"git@github.com:acme/skills",
		"ssh://git@github.com/acme/skills.git",
		"https://github.com/acme/skills.git",
		"https://github.com/acme/skills",
		"https://GitHub.com/Acme/Skills/",
	} {
		e, err := sources.NewEntry(url, sources.Options{})
		if err != nil {
			t.Errorf("NewEntry(%q): %v", url, err)
			continue
		}
		if e.ID != "@acme/skills" || e.Type != source.KindGit || e.URL != url {
			t.Errorf("NewEntry(%q) = %+v", url, e)
		}
	}
}

func TestNewEntryRefusesAnIdentityThatDoesNotMatchItsURL(t *testing.T) {
	t.Parallel()
	_, err := sources.NewEntry("https://github.com/evil/skills.git", sources.Options{ID: "@anthropics/skills"})
	if err == nil {
		t.Fatal("registered @anthropics/skills against another repository")
	}
	if _, err = sources.NewEntry(
		"https://github.com/acme/skills.git",
		sources.Options{ID: "@acme/skills"},
	); err != nil {
		t.Errorf("a matching explicit id was refused: %v", err)
	}
}

func TestNewEntryNonGitHubHostNeedsADeclaredIdentity(t *testing.T) {
	t.Parallel()
	url := "git@git.example.com:team/skills.git"
	if _, err := sources.NewEntry(url, sources.Options{}); err == nil {
		t.Error("a non-GitHub URL without an id was accepted")
	}
	e, err := sources.NewEntry(url, sources.Options{ID: "@team/skills"})
	if err != nil || e.ID != "@team/skills" || e.Type != source.KindGit {
		t.Errorf("entry = %+v, %v", e, err)
	}
}

func TestNewEntryRefusesReservedAndMisplacedIdentities(t *testing.T) {
	t.Parallel()
	url := "git@git.example.com:team/skills.git"
	for _, id := range []string{
		skill.DefaultRegistryID, "@ralvarezdev/ralvaskills", "@RALVAREZ-AI/ralvaskills", skill.OfficialID,
		"@local/anything", "not-an-identity", "",
	} {
		if _, err := sources.NewEntry(url, sources.Options{ID: id}); err == nil {
			t.Errorf("identity %q was accepted for a git source", id)
		}
	}
	// The same reserved names cannot be reached through the @owner/repo form.
	for _, target := range []string{"@anthropics/skills", "@ralvarezdev/ralvaskills", "@ralvarez-ai/ralvaskills"} {
		if _, err := sources.NewEntry(target, sources.Options{}); err == nil {
			t.Errorf("%s was registered as a declared source", target)
		}
	}
}

func TestNewEntryRefusesHostileInput(t *testing.T) {
	t.Parallel()
	for _, target := range []string{
		"ext::sh -c id", "file:///etc", "http://github.com/a/b", "--upload-pack=x",
		"https://user:ghp_abcdefghijklmnopqrstuvwxyz0123456789@github.com/a/b.git", "",
		"https://github.com/a", "https://github.com/a/b/c",
	} {
		if e, err := sources.NewEntry(target, sources.Options{ID: "@a/b"}); err == nil {
			t.Errorf("NewEntry(%q) accepted: %+v", target, e)
		}
	}
	base := "@acme/skills"
	for name, opts := range map[string]sources.Options{
		"rev option":     {Rev: "--upload-pack=x"},
		"rev colon":      {Rev: "a:b"},
		"root traversal": {Roots: []string{"../x"}},
		"root absolute":  {Roots: []string{"/etc"}},
		"root dash":      {Roots: []string{"-x"}},
		"empty exclude":  {Excludes: []string{""}},
		"nul exclude":    {Excludes: []string{"a\x00b"}},
		"exclude parent": {Excludes: []string{"../**"}},
		"empty root":     {Roots: []string{""}},
	} {
		if e, err := sources.NewEntry(base, opts); err == nil {
			t.Errorf("%s accepted: %+v", name, e)
		}
	}
}

func TestNewEntryPathSource(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "My Skills")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	e, err := sources.NewEntry(dir, sources.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if e.Type != source.KindPath || e.URL != dir || e.ID != "@local/my-skills" {
		t.Errorf("entry = %+v, want a path source named @local/my-skills", e)
	}

	e, err = sources.NewEntry(dir, sources.Options{ID: "@me/scratch"})
	if err != nil || e.ID != "@me/scratch" {
		t.Errorf("declared id: %+v, %v", e, err)
	}

	if _, err = sources.NewEntry(filepath.Join(dir, "missing"), sources.Options{}); err == nil {
		t.Error("a path that does not exist was accepted")
	}
	file := filepath.Join(dir, "f")
	if err = os.WriteFile(file, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err = sources.NewEntry(file, sources.Options{}); err == nil {
		t.Error("a file was accepted as a path source")
	}
	if _, err = sources.NewEntry(dir, sources.Options{ID: skill.OfficialID}); err == nil {
		t.Error("a path source took a reserved identity")
	}
}

func TestNewEntryRelativePathBecomesAbsolute(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	if err := os.Mkdir("skills-here", 0o755); err != nil {
		t.Fatal(err)
	}
	e, err := sources.NewEntry("./skills-here", sources.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !filepath.IsAbs(e.URL) || !strings.HasSuffix(e.URL, "skills-here") {
		t.Errorf("URL = %q, want an absolute path", e.URL)
	}
}

func TestAddRefusesDuplicates(t *testing.T) {
	t.Parallel()
	a, _ := sources.NewEntry("@acme/skills", sources.Options{})
	f, err := sources.Add(sources.File{}, a)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = sources.Add(f, a); err == nil {
		t.Error("the same identity was added twice")
	}
	other := a
	other.ID = "@acme/other"
	if _, err = sources.Add(f, other); err == nil {
		t.Error("two identities for one URL were accepted: the second would shadow the first")
	}
	again, _ := sources.NewEntry("https://github.com/ACME/Skills.git", sources.Options{})
	if _, err = sources.Add(f, again); err == nil {
		t.Error("a differently spelled URL for the same identity was accepted")
	}
}

func TestRemoveAndFind(t *testing.T) {
	t.Parallel()
	a, _ := sources.NewEntry("@acme/skills", sources.Options{})
	b, _ := sources.NewEntry("@acme/other", sources.Options{})
	f, _ := sources.Add(sources.File{}, a)
	f, _ = sources.Add(f, b)

	if got, ok := sources.Find(f, "@ACME/Skills"); !ok || got.ID != "@acme/skills" {
		t.Errorf("Find is not case-insensitive: %+v, %v", got, ok)
	}
	f, removed := sources.Remove(f, "@acme/skills")
	if !removed || len(f.Sources) != 1 || f.Sources[0].ID != "@acme/other" {
		t.Errorf("Remove = %+v, %v", f, removed)
	}
	if _, removed = sources.Remove(f, "@acme/skills"); removed {
		t.Error("removed a source twice")
	}
}

func TestSaveLoadRoundTripAndPermissions(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "rsk", "sources.json")
	if f, err := sources.Load(path); err != nil || len(f.Sources) != 0 {
		t.Fatalf("a missing file should load empty: %+v, %v", f, err)
	}

	a, _ := sources.NewEntry(
		"@acme/skills",
		sources.Options{Rev: "v1", Roots: []string{"skills"}, Excludes: []string{"**/testdata/**"}},
	)
	f, _ := sources.Add(sources.File{}, a)
	if err := sources.Save(path, f); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0o600 {
		t.Errorf("sources.json mode = %v, %v; want 0600", info.Mode().Perm(), err)
	}
	got, err := sources.Load(path)
	if err != nil || len(got.Sources) != 1 || got.Sources[0].Rev != "v1" ||
		got.Sources[0].Excludes[0] != "**/testdata/**" {
		t.Fatalf("round trip = %+v, %v", got, err)
	}
	if got.SchemaVersion != schema.Sources {
		t.Errorf("version = %d, want %d", got.SchemaVersion, schema.Sources)
	}
}

func TestLoadRefusesWhatItCannotTrust(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	write := func(name, body string) string {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
		return p
	}

	if _, err := sources.Load(
		write("newer.json", `{"version":2,"sources":[]}`),
	); !errors.Is(
		err,
		schema.ErrUnsupported,
	) {
		t.Errorf("a newer file err = %v, want ErrUnsupported", err)
	}
	if _, err := sources.Load(write("bad.json", `{not json`)); err == nil {
		t.Error("corrupt JSON was accepted")
	}
	// A hand-edited file goes through the same checks as `source add`.
	hostile := `{"version":1,"sources":[{"id":"@anthropics/skills","type":"git","url":"https://github.com/evil/skills.git"}]}`
	if _, err := sources.Load(write("hostile.json", hostile)); err == nil {
		t.Error("a hand-edited file registered @anthropics/skills against another repository")
	}
	ext := `{"version":1,"sources":[{"id":"@a/b","type":"git","url":"ext::sh -c id"}]}`
	if _, err := sources.Load(write("ext.json", ext)); err == nil {
		t.Error("a hand-edited file carried an ext:: URL")
	}
}

func TestEntrySpec(t *testing.T) {
	t.Parallel()
	e, _ := sources.NewEntry("@acme/skills", sources.Options{Rev: "main", Roots: []string{"skills"}})
	spec := e.Spec()
	if spec.ID != e.ID || spec.Type != source.KindGit || spec.URL != e.URL || spec.Rev != "main" ||
		spec.Roots[0] != "skills" {
		t.Errorf("spec = %+v", spec)
	}
}
