package main

import (
	"path/filepath"
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source/sourcetest"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
)

// TestPlanSyncSupportsVersionlessOfficialSkills covers the official-source
// lock entry: those skills carry no version in their frontmatter, so the lock
// records an empty one and sync must not reject it. An entry written before the
// git transport has no git record and resolves to what the source has now.
//
//nolint:paralleltest // the official source seam is process-wide
func TestPlanSyncSupportsVersionlessOfficialSkills(t *testing.T) {
	withOfficialSkills(t, map[string]string{"frontend-design": ""})
	project := t.TempDir()
	cfg := config.Config{OfficialCache: filepath.Join(t.TempDir(), "anthropic")}
	m := manifest.Mod{Tools: []tool.ID{tool.ClaudeID}, Skills: map[string]string{"frontend-design": "*"}}
	lock := manifest.Lock{
		Skills: []manifest.LockEntry{{Name: "frontend-design", Source: skill.SourceOfficial}},
	}
	links, err := planSyncLinks(t.Context(), project, m, lock, cfg)
	if err != nil || len(links) != 1 {
		t.Fatalf("links=%+v err=%v", links, err)
	}
	want := filepath.Join(cfg.SourceOptions().CacheRoot, "sources", "anthropics", "skills", "frontend-design",
		strings.Repeat("c", 40))
	if links[0].skill.Path != want || links[0].skill.Source != skill.SourceOfficial {
		t.Fatalf("skill = %+v, want an official skill extracted at %s", links[0].skill, want)
	}
}

// An official entry that carries a git record is restored at its locked commit.
//
//nolint:paralleltest // the official source seam is process-wide
func TestPlanSyncRestoresAnOfficialEntryAtItsLockedCommit(t *testing.T) {
	repo := withOfficialSkills(t, map[string]string{"xlsx": ""})
	project := t.TempDir()
	cfg := config.Config{OfficialCache: filepath.Join(t.TempDir(), "anthropic")}
	m := manifest.Mod{Tools: []tool.ID{tool.ClaudeID}, Skills: map[string]string{"xlsx": "*"}}
	locked := strings.Repeat("a", 40)
	lock := manifest.Lock{Skills: []manifest.LockEntry{{
		Name: "xlsx", Source: skill.SourceOfficial,
		Git: &skill.GitRef{
			URL:    "https://github.com/anthropics/skills",
			Commit: locked,
			Tree:   sourcetest.TreeOf("skills/xlsx"),
			Path:   "skills/xlsx",
		},
	}}}

	links, err := planSyncLinks(t.Context(), project, m, lock, cfg)
	if err != nil || len(links) != 1 {
		t.Fatalf("links=%+v err=%v", links, err)
	}
	if !strings.HasSuffix(links[0].skill.Path, filepath.Join("xlsx", locked)) {
		t.Errorf("path = %s, want the extraction at the locked commit", links[0].skill.Path)
	}
	if len(repo.Fetches) != 0 {
		t.Errorf("restoring from a lock fetched %v: it must not consult the branch", repo.Fetches)
	}

	lock.Skills[0].Git.URL = "https://github.com/evil/skills"
	if _, err = planSyncLinks(t.Context(), project, m, lock, cfg); err == nil {
		t.Error("a lock that fetches the official source from another URL was followed")
	}
}

//nolint:paralleltest // the config folder override is process-wide
func TestPlanSyncRestoresGitEntriesOnlyThroughADeclaredTrustedSource(t *testing.T) {
	config.SetConfigHome(t.TempDir())
	t.Cleanup(func() { config.SetConfigHome("") })

	project := t.TempDir()
	cfg := config.Config{OfficialCache: filepath.Join(t.TempDir(), "anthropic")}
	m := manifest.Mod{Tools: []tool.ID{tool.ClaudeID}, Skills: map[string]string{"alpha": "*"}}
	lock := manifest.Lock{Skills: []manifest.LockEntry{{
		Name: "alpha", Version: "1.0.0", Source: skill.DeclaredSource("@acme/skills"),
		Git: &skill.GitRef{
			URL: "git@github.com:acme/skills.git", Path: "skills/alpha",
			Commit: strings.Repeat("a", 40), Tree: strings.Repeat("b", 40),
		},
	}}}

	_, err := planSyncLinks(t.Context(), project, m, lock, cfg)
	if err == nil || !strings.Contains(err.Error(), "rsk source add @acme/skills") {
		t.Errorf("a lock naming an undeclared source err = %v, want the command to add it", err)
	}
}

//nolint:paralleltest // the official source seam is process-wide
func TestPlanSyncRefusesAnOfficialSkillTheSourceLacks(t *testing.T) {
	withOfficialSkills(t, map[string]string{"other": ""})
	cfg := config.Config{OfficialCache: filepath.Join(t.TempDir(), "anthropic")}
	m := manifest.Mod{Tools: []tool.ID{tool.ClaudeID}, Skills: map[string]string{"example": "*"}}
	lock := manifest.Lock{Skills: []manifest.LockEntry{{Name: "example", Source: skill.SourceOfficial}}}

	if _, err := planSyncLinks(t.Context(), t.TempDir(), m, lock, cfg); err == nil ||
		!strings.Contains(err.Error(), `no skill "example"`) {
		t.Errorf("err = %v, want a refusal because the official source has no such skill", err)
	}
}
