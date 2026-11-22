package catalog_test

import (
	"os"
	"path/filepath"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/catalog"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
)

func writeCatalogSkill(t *testing.T, dir, name string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, name), 0o750); err != nil {
		t.Fatal(err)
	}
	body := "---\nname: " + name + "\nversion: 1.0.0\n---\n"
	if err := os.WriteFile(filepath.Join(dir, name, "SKILL.md"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

// A declared path source is a declared source: the local-clone walk must leave
// it alone, so `--source local` never lists it and the declared walk does.
func TestListIncludesADeclaredPathSourceWithTheDeclaredOnes(t *testing.T) {
	t.Parallel()
	clone, scratch := t.TempDir(), t.TempDir()
	writeCatalogSkill(t, filepath.Join(clone, "skills"), "from-clone")
	writeCatalogSkill(t, scratch, "from-path")
	chain, err := source.Build([]source.Spec{
		{ID: "@local/clone", Type: source.KindPath, URL: clone},
		{ID: "@me/scratch", Type: source.KindPath, URL: scratch, Declared: true},
	}, source.Options{})
	if err != nil {
		t.Fatal(err)
	}

	list := func(filter skill.Source) []string {
		res, listErr := catalog.List(t.Context(), chain, catalog.ListOptions{Source: filter, Personal: true})
		if listErr != nil {
			t.Fatal(listErr)
		}
		names := make([]string, len(res.Skills))
		for i := range res.Skills {
			names[i] = res.Skills[i].Name + ":" + string(res.Skills[i].Source)
		}
		if len(res.Warnings) != 0 {
			t.Errorf("warnings: %v", res.Warnings)
		}
		return names
	}
	if got := list(skill.SourceLocal); len(got) != 1 || got[0] != "from-clone:local" {
		t.Errorf("--source local = %v", got)
	}
	if got := list(skill.SourceGit); len(got) != 1 || got[0] != "from-path:path" {
		t.Errorf("--source git = %v", got)
	}
	if got := list(""); len(got) != 2 {
		t.Errorf("all = %v", got)
	}
}

// A private source is registered, trusted by being declared, and still left out
// of the catalog unless the listing is asked to include it.
func TestListLeavesAPrivateSourceOutUnlessIncluded(t *testing.T) {
	t.Parallel()
	open, hidden := t.TempDir(), t.TempDir()
	writeCatalogSkill(t, open, "public-one")
	writeCatalogSkill(t, hidden, "private-one")
	chain, err := source.Build([]source.Spec{
		{ID: "@me/open", Type: source.KindPath, URL: open, Declared: true},
		{ID: "@me/hidden", Type: source.KindPath, URL: hidden, Declared: true, Private: true},
	}, source.Options{})
	if err != nil {
		t.Fatal(err)
	}
	names := func(includePrivate bool) []string {
		res, listErr := catalog.List(t.Context(), chain, catalog.ListOptions{Private: includePrivate, Personal: true})
		if listErr != nil {
			t.Fatal(listErr)
		}
		out := make([]string, len(res.Skills))
		for i := range res.Skills {
			out[i] = res.Skills[i].Name
		}
		return out
	}
	if got := names(false); len(got) != 1 || got[0] != "public-one" {
		t.Errorf("default listing = %v", got)
	}
	if got := names(true); len(got) != 2 {
		t.Errorf("--include private listing = %v", got)
	}
}
