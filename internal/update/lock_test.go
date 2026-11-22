package update_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/update"
)

func TestLinkAndLockLinksResolvedSkillsAndSkipsTheRest(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "src", "alpha")
	if err := os.MkdirAll(srcDir, 0o750); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "skills")
	if err := os.MkdirAll(target, 0o750); err != nil {
		t.Fatal(err)
	}
	alpha := skill.Skill{Name: "alpha", Version: "2.0.0", Source: skill.SourceRegistry, Path: srcDir}
	resolve := func(name string) (skill.Skill, error) {
		if name == "alpha" {
			return alpha, nil
		}
		return skill.Skill{}, errors.New("not found")
	}
	lock := manifest.Lock{Skills: []manifest.LockEntry{{Name: "alpha", Version: "1.0.0"}}}

	got, resolved, err := update.LinkAndLock(lock, []string{"missing", "alpha"}, []string{target}, resolve)
	if err != nil {
		t.Fatal(err)
	}
	if len(resolved) != 1 || resolved[0].Name != "alpha" {
		t.Fatalf("resolved = %v, want only alpha", resolved)
	}
	entry, ok := manifest.FindLockEntry(got, "alpha")
	if !ok || entry.Version != "2.0.0" {
		t.Errorf("entry = %+v, want the resolved version recorded", entry)
	}
	if dest, linkErr := os.Readlink(filepath.Join(target, "alpha")); linkErr != nil || dest != srcDir {
		t.Errorf("link = %q, %v; want it pointing at %q", dest, linkErr, srcDir)
	}
	if _, found := manifest.FindLockEntry(got, "missing"); found {
		t.Error("a skill that did not resolve must not get a lock entry")
	}
}
