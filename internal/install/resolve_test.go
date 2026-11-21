package install_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/install"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
)

func writeSkillDir(t *testing.T, root, name string) {
	t.Helper()
	dir := filepath.Join(root, skill.SkillsFolderName, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nname: " + name + "\ndescription: test\nmetadata:\n  version: 1.0.0\n---\n"
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// localModeChain builds the real chain for a machine in local mode: a repo
// clone holding the local names and an official cache holding the official ones.
func localModeChain(t *testing.T, local, official []string) source.Chain {
	t.Helper()
	repo, cache := t.TempDir(), t.TempDir()
	for _, n := range local {
		writeSkillDir(t, repo, n)
	}
	for _, n := range official {
		writeSkillDir(t, cache, n)
	}
	chain, err := config.Config{RepoPath: repo, OfficialCache: cache}.Sources()
	if err != nil {
		t.Fatal(err)
	}
	return chain
}

func TestFindSkillSearchesEverySource(t *testing.T) {
	t.Parallel()
	chain := localModeChain(t, []string{"only-local", "both"}, []string{"only-official", "both"})

	for name, wantSource := range map[string]skill.Source{
		"only-local":    skill.SourceLocal,
		"only-official": skill.SourceOfficial,
	} {
		got, err := install.FindSkill(t.Context(), name, chain)
		if err != nil {
			t.Fatalf("FindSkill(%s): %v", name, err)
		}
		if got.Source != wantSource {
			t.Errorf("FindSkill(%s).Source = %s, want %s", name, got.Source, wantSource)
		}
	}

	_, err := install.FindSkill(t.Context(), "both", chain)
	if !errors.Is(err, source.ErrAmbiguous) {
		t.Fatalf("a name in two sources: err = %v, want ErrAmbiguous, not the first match", err)
	}
}

func TestFindSkillNotFoundKeepsTheMessage(t *testing.T) {
	t.Parallel()
	chain := localModeChain(t, []string{"other"}, nil)

	_, err := install.FindSkill(t.Context(), "nope", chain)
	if err == nil {
		t.Fatal("expected an error")
	}
	want := "skill \"nope\" not found in local repo or official cache\n  Run 'rsk catalog' to browse available skills"
	if err.Error() != want {
		t.Errorf("error = %q, want %q", err.Error(), want)
	}
}

func TestResolveBundleMembersStayInTheirSource(t *testing.T) {
	t.Parallel()
	// "dup" exists in both sources: a bundle member names its source, so it is
	// never ambiguous there, while the bare name would be.
	chain := localModeChain(t, []string{"dup", "mine"}, []string{"dup", "theirs"})
	bundles := []config.Bundle{{Name: "b", Skills: []config.SkillRef{
		{Name: "dup", Source: skill.SourceOfficial},
		{Name: "mine", Source: skill.SourceLocal},
		{Name: "theirs", Source: skill.SourceOfficial},
		{Name: "ghost", Source: skill.SourceLocal},
		{Name: "mine", Source: skill.SourceOfficial},
	}}}

	skills, warnings, err := install.Resolve(t.Context(), []string{"b"}, bundles, chain)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]skill.Source{}
	for _, s := range skills {
		got[s.Name] = s.Source
	}
	if got["dup"] != skill.SourceOfficial || got["mine"] != skill.SourceLocal || got["theirs"] != skill.SourceOfficial {
		t.Errorf("resolved = %v", got)
	}
	if len(warnings) != 2 {
		t.Fatalf("warnings = %q, want one for ghost and one for mine-in-official", warnings)
	}
	if !strings.Contains(warnings[0], "ghost is not yet available (planned)") ||
		!strings.Contains(warnings[1], "mine (official) not in cache") {
		t.Errorf("warnings = %q", warnings)
	}
}

func TestResolveQualifiedReferenceNamesItsSource(t *testing.T) {
	t.Parallel()
	// "both" is ambiguous bare, so the qualified form is the way out.
	chain := localModeChain(t, []string{"both"}, []string{"both"})

	got, _, err := install.Resolve(t.Context(), []string{"@anthropics/skills/both"}, nil, chain)
	if err != nil || len(got) != 1 || got[0].Source != skill.SourceOfficial {
		t.Fatalf("Resolve = %+v, %v; want the official skill", got, err)
	}

	_, _, err = install.Resolve(t.Context(), []string{"@acme/elsewhere/both"}, nil, chain)
	if !errors.Is(err, source.ErrUnknownSource) {
		t.Errorf("unknown source err = %v, want ErrUnknownSource", err)
	}

	_, _, err = install.Resolve(t.Context(), []string{"@anthropics/skills/ghost"}, nil, chain)
	if err == nil || !strings.Contains(err.Error(), `skill "ghost" not found in @anthropics/skills`) {
		t.Errorf("qualified miss err = %v", err)
	}
}

func TestParseNameVersionAndConstraintsHandleQualifiedReferences(t *testing.T) {
	t.Parallel()
	name, version, err := install.ParseNameVersion("@anthropics/skills/xlsx@1.2.0")
	if err != nil || name != "xlsx" || version != "1.2.0" {
		t.Errorf("ParseNameVersion = %q, %q, %v; want xlsx, 1.2.0", name, version, err)
	}

	got, err := install.ConstraintsFromArgs([]string{"@anthropics/skills/xlsx@1.2.0", "go-grpc", "docs@2.0.0"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got["xlsx"] != "1.2.0" || got["docs"] != "2.0.0" {
		t.Errorf("constraints = %v", got)
	}
}
