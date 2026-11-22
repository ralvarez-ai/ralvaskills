package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

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
func TestCatalogListsADeclaredPathSourceWithTheDeclaredOnes(t *testing.T) {
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
		var warn bytes.Buffer
		listed := loadSkillsBySource(t.Context(), &warn, nil, chain, filter, false)
		names := make([]string, len(listed))
		for i := range listed {
			names[i] = listed[i].Name + ":" + string(listed[i].Source)
		}
		if warn.Len() != 0 {
			t.Errorf("warnings: %s", warn.String())
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
func TestCatalogLeavesAPrivateSourceOutUnlessIncluded(t *testing.T) {
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
		var warn bytes.Buffer
		listed := loadSkillsBySource(t.Context(), &warn, nil, chain, "", includePrivate)
		out := make([]string, len(listed))
		for i := range listed {
			out[i] = listed[i].Name
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
