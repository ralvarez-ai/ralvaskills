package source_test

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/source"
)

func TestNewBuildsEachKind(t *testing.T) {
	t.Parallel()
	opts := source.Options{
		RegistryCache: t.TempDir(), CacheRoot: t.TempDir(), Gate: func(string, string) error { return nil },
	}

	tests := []struct {
		name string
		spec source.Spec
	}{
		{"registry", source.Spec{
			ID: "@ralvarez-ai/ralvaskills", Type: source.KindRegistry, URL: "https://skills.example",
		}},
		{"path", source.Spec{ID: "@local/ralvaskills", Type: source.KindPath, URL: t.TempDir()}},
		{"official", source.Spec{ID: "@anthropics/skills", Type: source.KindGit, URL: source.OfficialSkillsURL}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r, err := source.New(tc.spec, opts)
			if err != nil {
				t.Fatalf("New: %v", err)
			}
			if r.Spec().ID != tc.spec.ID || r.Spec().Type != tc.spec.Type {
				t.Errorf("resolver = %s/%s, want %s/%s", r.Spec().ID, r.Spec().Type, tc.spec.ID, tc.spec.Type)
			}
		})
	}
}

func TestNewRejectsWhatItCannotHonour(t *testing.T) {
	t.Parallel()
	opts := source.Options{}
	path := func(mod func(*source.Spec)) source.Spec {
		spec := source.Spec{ID: "@a/a", Type: source.KindPath, URL: "/x"}
		mod(&spec)
		return spec
	}
	tests := []struct {
		name string
		spec source.Spec
	}{
		{"unknown kind", source.Spec{ID: "@a/a", Type: "bogus", URL: "git@host:a/a.git"}},
		{"git without a gate or a cache", source.Spec{ID: "@a/a", Type: source.KindGit, URL: "git@host:a/a.git"}},
		{"empty kind", source.Spec{ID: "@a/a"}},
		{"private is not enforced yet", path(func(s *source.Spec) { s.Private = true })},
		{"rev is not honored yet", path(func(s *source.Spec) { s.Rev = "main" })},
		{"roots are not honored yet", path(func(s *source.Spec) { s.Roots = []string{"s"} })},
		{"excludes are not honored yet", path(func(s *source.Spec) { s.Excludes = []string{"**"} })},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if _, err := source.New(tc.spec, opts); err == nil {
				t.Fatal("New accepted a spec it would silently ignore")
			}
		})
	}
}

func TestBuildKeepsDeclarationOrderAndRegistryAccessor(t *testing.T) {
	t.Parallel()
	opts := source.Options{
		RegistryCache: t.TempDir(), CacheRoot: t.TempDir(), Gate: func(string, string) error { return nil },
	}
	chain, err := source.Build([]source.Spec{
		{ID: "@ralvarez-ai/ralvaskills", Type: source.KindRegistry, URL: "https://skills.example"},
		{ID: "@anthropics/skills", Type: source.KindGit, URL: source.OfficialSkillsURL},
	}, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(chain) != 2 {
		t.Fatalf("chain has %d sources, want 2", len(chain))
	}
	if chain[0].Spec().ID != "@ralvarez-ai/ralvaskills" || chain[1].Spec().ID != "@anthropics/skills" {
		t.Fatalf("chain order = %s, %s", chain[0].Spec().ID, chain[1].Spec().ID)
	}
	if _, ok := chain.Registry(); !ok {
		t.Error("Registry() found no registry in a chain that has one")
	}
	if _, ok := chain.Select(source.KindGit).Registry(); ok {
		t.Error("Registry() found a registry in a chain without one")
	}
	if _, err = source.Build([]source.Spec{{ID: "@a/a", Type: "bogus"}}, opts); err == nil {
		t.Error("Build accepted an unknown kind")
	}
}

func TestPathSourceFindsASkill(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeSkill(t, root, "skills/x/SKILL.md")
	r, err := source.New(source.Spec{ID: "@local/r", Type: source.KindPath, URL: root}, source.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = r.Find(context.Background(), "x"); err != nil {
		t.Fatalf("Find: %v", err)
	}
}

// TestOnlyTheFactoryConstructsSources keeps the concrete constructors
// unexported, so nothing outside this package can build a source except
// through source.New (the compiler enforces it; this guards against someone
// exporting one again).
func TestOnlyTheFactoryConstructsSources(t *testing.T) {
	t.Parallel()
	fset := token.NewFileSet()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	banned := map[string]bool{"NewOfficial": true, "NewRegistry": true, "NewLocal": true}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		f, parseErr := parser.ParseFile(fset, e.Name(), nil, 0)
		if parseErr != nil {
			t.Fatal(parseErr)
		}
		for _, d := range f.Decls {
			if fn, ok := d.(*ast.FuncDecl); ok && fn.Recv == nil && banned[fn.Name.Name] {
				t.Errorf("%s: %s must stay unexported; construct sources through source.New",
					fset.Position(fn.Pos()), fn.Name.Name)
			}
		}
	}
}

func writeSkill(t *testing.T, root, rel string) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nname: x\ndescription: test skill\nmetadata:\n  version: 1.0.0\n---\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
