package install_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
	"go.ralva.dev/ralvaskills/v3/internal/install"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
)

func TestApplyProjectWritesManifestLockAndPin(t *testing.T) {
	t.Parallel()

	project := t.TempDir()
	rskDir := filepath.Join(project, ".rsk")
	if err := manifest.WriteMod(rskDir, manifest.Mod{Tools: []tool.ID{tool.ClaudeID}}); err != nil {
		t.Fatalf("seed manifest: %v", err)
	}

	src := filepath.Join(t.TempDir(), "go-architect")
	if err := os.MkdirAll(src, fsperm.Dir); err != nil {
		t.Fatalf("mkdir src: %v", err)
	}
	target := filepath.Join(project, ".claude", "skills")

	results, err := install.Apply([]skill.Skill{{
		Name:    "go-architect",
		Version: "1.2.3",
		Path:    src,
		Source:  skill.SourceLocal,
	}}, install.Options{
		Scope:       install.ScopeProject,
		Targets:     []string{target},
		RskDir:      rskDir,
		Pin:         true,
		Constraints: map[string]string{"go-architect": "1.2.3"},
	})
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if len(results) != 1 || results[0].Err != nil {
		t.Fatalf("results = %+v", results)
	}
	if _, linkErr := os.Lstat(filepath.Join(target, "go-architect")); linkErr != nil {
		t.Errorf("link missing: %v", linkErr)
	}

	m, err := manifest.ReadMod(rskDir)
	if err != nil {
		t.Fatalf("ReadMod: %v", err)
	}
	if m.Skills["go-architect"] != "1.2.3" {
		t.Errorf("skills = %+v, want go-architect pinned to 1.2.3", m.Skills)
	}
	if !slices.Contains(m.Pinned, "go-architect") {
		t.Errorf("pinned = %+v, want go-architect", m.Pinned)
	}

	lock, err := manifest.ReadLock(rskDir)
	if err != nil {
		t.Fatalf("ReadLock: %v", err)
	}
	if len(lock.Skills) != 1 || lock.Skills[0].Version != "1.2.3" {
		t.Errorf("lock = %+v, want one go-architect entry at 1.2.3", lock.Skills)
	}
}

func TestApplyProjectLinksAllNativeDirectories(t *testing.T) {
	t.Parallel()
	project := t.TempDir()
	rskDir := filepath.Join(project, ".rsk")
	m := manifest.Mod{Tools: []tool.ID{tool.ClaudeID, tool.CodexID, tool.OpenCodeID}}
	if err := manifest.WriteMod(rskDir, m); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(t.TempDir(), "example")
	if err := os.MkdirAll(src, fsperm.Dir); err != nil {
		t.Fatal(err)
	}
	results, err := install.Apply([]skill.Skill{{Name: "example", Path: src, Version: "1.0.0"}}, install.Options{
		Scope: install.ScopeProject, Targets: install.ProjectSkillsDirs(project, m), RskDir: rskDir,
	})
	if err != nil || len(results) != 1 || results[0].Err != nil {
		t.Fatalf("Apply: results=%+v error=%v", results, err)
	}
	for _, base := range []string{".claude", ".agents", ".opencode"} {
		path := filepath.Join(project, base, "skills", "example")
		if !skill.IsLinked("example", filepath.Dir(path)) {
			t.Errorf("missing skill link %s", path)
		}
	}
}

func TestApplyGlobalWritesNoManifest(t *testing.T) {
	t.Parallel()

	project := t.TempDir()
	target := filepath.Join(project, "global-skills")
	src := filepath.Join(t.TempDir(), "tdd")
	if err := os.MkdirAll(src, fsperm.Dir); err != nil {
		t.Fatalf("mkdir src: %v", err)
	}

	results, err := install.Apply([]skill.Skill{{
		Name:    "tdd",
		Version: "1.0.0",
		Path:    src,
	}}, install.Options{
		Scope:   install.ScopeGlobal,
		Targets: []string{target},
	})
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if results[0].Err != nil {
		t.Fatalf("results = %+v", results)
	}
	if _, linkErr := os.Lstat(filepath.Join(target, "tdd")); linkErr != nil {
		t.Errorf("link missing: %v", linkErr)
	}
	if _, statErr := os.Stat(filepath.Join(project, ".rsk")); !os.IsNotExist(statErr) {
		t.Error("global apply wrote a project manifest")
	}
}

func TestApplyProjectQualifiesAnOfficialSkillAndLeavesRegistryOnesBare(t *testing.T) {
	t.Parallel()

	project := t.TempDir()
	rskDir := filepath.Join(project, ".rsk")
	if err := manifest.WriteMod(rskDir, manifest.Mod{Tools: []tool.ID{tool.ClaudeID}}); err != nil {
		t.Fatal(err)
	}
	mk := func(name string) string {
		dir := filepath.Join(t.TempDir(), name)
		if err := os.MkdirAll(dir, fsperm.Dir); err != nil {
			t.Fatal(err)
		}
		return dir
	}

	_, err := install.Apply([]skill.Skill{
		{Name: "xlsx", Path: mk("xlsx"), Source: skill.SourceOfficial},
		{Name: "tdd", Version: "1.0.0", Path: mk("tdd"), Source: skill.SourceRegistry},
	}, install.Options{
		Scope: install.ScopeProject, Targets: []string{filepath.Join(project, ".claude", "skills")},
		RskDir: rskDir, Pin: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	modText, _ := os.ReadFile(manifest.ModPath(rskDir))
	for _, want := range []string{"version = 2", `"@anthropics/skills/xlsx"`, `tdd = "*"`} {
		if !strings.Contains(string(modText), want) {
			t.Errorf("rsk.mod lacks %q:\n%s", want, modText)
		}
	}
	m, err := manifest.ReadMod(rskDir)
	if err != nil || m.Skills["xlsx"] != "*" || m.Origin["xlsx"] != skill.OfficialID {
		t.Fatalf("read back %+v, %v", m, err)
	}
	lock, err := manifest.ReadLock(rskDir)
	if err != nil {
		t.Fatal(err)
	}
	if entry, ok := manifest.FindLockEntry(lock, "xlsx"); !ok || entry.Source != skill.SourceOfficial {
		t.Errorf("lock entry = %+v, %v", entry, ok)
	}
}

func TestApplyProjectRecordsAGitSkillInTheLock(t *testing.T) {
	t.Parallel()

	src := filepath.Join(t.TempDir(), "alpha")
	if err := os.MkdirAll(src, fsperm.Dir); err != nil {
		t.Fatal(err)
	}
	ref := &skill.GitRef{
		URL: "git@github.com:acme/skills.git", Rev: "main", Path: "skills/alpha",
		Commit: strings.Repeat("a", 40), Tree: strings.Repeat("b", 40),
	}
	gitSkill := skill.Skill{
		Name: "alpha", Version: "1.0.0", Path: src, Source: skill.DeclaredSource("@acme/skills"), Git: ref,
	}

	project := t.TempDir()
	rskDir := filepath.Join(project, ".rsk")
	if err := manifest.WriteMod(rskDir, manifest.Mod{Tools: []tool.ID{tool.ClaudeID}}); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(project, ".claude", "skills")
	if _, err := install.Apply([]skill.Skill{gitSkill}, install.Options{
		Scope:       install.ScopeProject,
		Targets:     []string{target},
		RskDir:      rskDir,
		Constraints: map[string]string{"alpha": "1.0.0"},
	}); err != nil {
		t.Fatal(err)
	}

	lock, err := manifest.ReadLock(rskDir)
	if err != nil {
		t.Fatal(err)
	}
	entry, ok := manifest.FindLockEntry(lock, "alpha")
	if !ok || entry.Source != gitSkill.Source || entry.Git == nil || *entry.Git != *ref {
		t.Errorf("lock entry = %+v, %v", entry, ok)
	}
	m, err := manifest.ReadMod(rskDir)
	if err != nil || m.Skills["alpha"] != "1.0.0" || m.Origin["alpha"] != "@acme/skills" {
		t.Errorf("mod = %+v, %v", m, err)
	}
	if _, statErr := os.Lstat(filepath.Join(target, "alpha")); statErr != nil {
		t.Errorf("the skill was not linked: %v", statErr)
	}

	// Reinstalling a registry skill of the same name drops the git provenance.
	reg := skill.Skill{Name: "alpha", Version: "3.0.0", Path: src, Source: skill.SourceRegistry}
	if _, err = install.Apply([]skill.Skill{reg}, install.Options{
		Scope: install.ScopeProject, Targets: []string{target}, RskDir: rskDir,
	}); err != nil {
		t.Fatal(err)
	}
	lock, _ = manifest.ReadLock(rskDir)
	entry, _ = manifest.FindLockEntry(lock, "alpha")
	if entry.Git != nil || entry.Source != skill.SourceRegistry {
		t.Errorf("a registry install kept git fields: %+v", entry)
	}
}

// A path source is a directory on one machine: a project manifest cannot
// restore it elsewhere, so it is refused before anything is linked or written.
func TestApplyProjectRefusesASkillFromAPathSource(t *testing.T) {
	t.Parallel()
	project := t.TempDir()
	rskDir := filepath.Join(project, ".rsk")
	if err := manifest.WriteMod(rskDir, manifest.Mod{Tools: []tool.ID{tool.ClaudeID}}); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(t.TempDir(), "scratch-skill")
	if err := os.MkdirAll(src, fsperm.Dir); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(project, ".claude", "skills")

	_, err := install.Apply(
		[]skill.Skill{{Name: "scratch-skill", Path: src, Source: skill.SourcePath}},
		install.Options{
			Scope: install.ScopeProject, Targets: []string{target}, RskDir: rskDir,
		},
	)
	if err == nil || !strings.Contains(err.Error(), "path source") {
		t.Fatalf("err = %v, want a refusal naming the path source", err)
	}
	if _, statErr := os.Lstat(filepath.Join(target, "scratch-skill")); statErr == nil {
		t.Error("the skill was linked before the refusal")
	}
	if m, readErr := manifest.ReadMod(rskDir); readErr != nil || len(m.Skills) != 0 {
		t.Errorf("manifest = %+v, %v", m, readErr)
	}
}

func TestApplyGlobalLinksASkillFromAPathSource(t *testing.T) {
	t.Parallel()
	src := filepath.Join(t.TempDir(), "scratch-skill")
	if err := os.MkdirAll(src, fsperm.Dir); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(t.TempDir(), "skills")

	results, err := install.Apply(
		[]skill.Skill{{Name: "scratch-skill", Path: src, Source: skill.SourcePath}},
		install.Options{
			Scope: install.ScopeGlobal, Targets: []string{target},
		},
	)
	if err != nil || len(results) != 1 || results[0].Err != nil {
		t.Fatalf("Apply = %+v, %v", results, err)
	}
}
