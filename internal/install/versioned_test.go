package install_test

import (
	"context"
	"errors"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/install"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
)

// gitLike is a source with a latest version and one pinned through a tag.
type gitLike struct {
	id       string
	versions map[string]string
}

func (g gitLike) Spec() source.Spec {
	return source.Spec{ID: g.id, Type: source.KindGit}
}

func (g gitLike) All(context.Context) ([]skill.Skill, error) {
	return nil, nil
}

func (g gitLike) Find(_ context.Context, name string) (skill.Skill, error) {
	return skill.Skill{Name: name, Version: "2.0.0", Source: skill.DeclaredSource(g.id)}, nil
}

func (g gitLike) FindVersion(_ context.Context, name, version string) (skill.Skill, error) {
	if _, ok := g.versions[version]; !ok {
		return skill.Skill{}, errors.New("no such tag")
	}
	return skill.Skill{Name: name, Version: version, Source: skill.DeclaredSource(g.id)}, nil
}

func TestResolvePinsAGitSkillThroughItsVersion(t *testing.T) {
	t.Parallel()
	chain := source.Chain{gitLike{id: "@acme/skills", versions: map[string]string{"1.0.0": ""}}}

	got, _, err := install.Resolve(t.Context(), []string{"@acme/skills/alpha@1.0.0"}, nil, chain)
	if err != nil || len(got) != 1 || got[0].Version != "1.0.0" {
		t.Fatalf("pinned = %+v, %v; want version 1.0.0", got, err)
	}
	got, _, err = install.Resolve(t.Context(), []string{"@acme/skills/alpha"}, nil, chain)
	if err != nil || got[0].Version != "2.0.0" {
		t.Errorf("unpinned = %+v, %v; want the latest", got, err)
	}
	if _, _, err = install.Resolve(t.Context(), []string{"@acme/skills/alpha@9.9.9"}, nil, chain); err == nil {
		t.Error("an unknown tag resolved")
	}
}

func TestResolveKeepsIgnoringTheVersionForOtherSources(t *testing.T) {
	t.Parallel()
	chain := localModeChain(t, []string{"mine"}, nil)
	got, _, err := install.Resolve(t.Context(), []string{"mine@9.9.9"}, nil, chain)
	if err != nil || len(got) != 1 || got[0].Name != "mine" {
		t.Errorf("a built-in source must resolve as before whatever the version: %+v, %v", got, err)
	}
}
