package source_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
)

// fakeSource is an in-memory Resolver. A name in skills resolves; a name in
// failOn returns failErr; anything else is not found.
type fakeSource struct {
	spec    source.Spec
	skills  map[string]bool
	failErr error
	calls   []string
}

func newFake(id string, kind source.Kind, names ...string) *fakeSource {
	f := &fakeSource{spec: source.Spec{ID: id, Type: kind}, skills: map[string]bool{}}
	for _, n := range names {
		f.skills[n] = true
	}
	return f
}

func (f *fakeSource) Spec() source.Spec {
	return f.spec
}

func (f *fakeSource) All(context.Context) ([]skill.Skill, error) {
	return nil, nil
}

func (f *fakeSource) Find(_ context.Context, name string) (skill.Skill, error) {
	f.calls = append(f.calls, "Find:"+name)
	return f.lookup(name)
}

func (f *fakeSource) lookup(name string) (skill.Skill, error) {
	if f.failErr != nil {
		return skill.Skill{}, f.failErr
	}
	if f.skills[name] {
		return skill.Skill{Name: name, Path: f.spec.ID + "/" + name}, nil
	}
	return skill.Skill{}, fmt.Errorf("%w: %s has no %q", source.ErrNotFound, f.spec.ID, name)
}

// lazyFake is a source whose Find may fetch over the network (as Official
// clones on demand) and which can also answer from what it already holds.
type lazyFake struct{ *fakeSource }

func (l lazyFake) FindCached(_ context.Context, name string) (skill.Skill, error) {
	l.calls = append(l.calls, "FindCached:"+name)
	return l.lookup(name)
}

func TestChainFindBare(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	boom := errors.New("registry unreachable")

	tests := []struct {
		name     string
		chain    func() source.Chain
		wantPath string
		wantErr  error
		wantMsg  []string
	}{
		{
			name: "one source has it",
			chain: func() source.Chain {
				return source.Chain{newFake("@a/a", source.KindRegistry, "x"), newFake("@b/b", source.KindOfficial)}
			},
			wantPath: "@a/a/x",
		},
		{
			name: "found in the second source",
			chain: func() source.Chain {
				return source.Chain{newFake("@a/a", source.KindRegistry), newFake("@b/b", source.KindOfficial, "x")}
			},
			wantPath: "@b/b/x",
		},
		{
			name: "two sources have it is ambiguous and names both",
			chain: func() source.Chain {
				return source.Chain{
					newFake("@a/a", source.KindRegistry, "x"), newFake("@b/b", source.KindOfficial, "x"),
				}
			},
			wantErr: source.ErrAmbiguous,
			wantMsg: []string{`"x"`, "@a/a", "@b/b"},
		},
		{
			name: "no source has it",
			chain: func() source.Chain {
				return source.Chain{newFake("@a/a", source.KindRegistry), newFake("@b/b", source.KindOfficial)}
			},
			wantErr: source.ErrNotFound,
		},
		{
			name: "a failing source fails closed even when another has the name",
			chain: func() source.Chain {
				bad := newFake("@a/a", source.KindRegistry)
				bad.failErr = boom
				return source.Chain{bad, newFake("@b/b", source.KindOfficial, "x")}
			},
			wantErr: boom,
		},
		{
			name: "a failing source fails closed regardless of order",
			chain: func() source.Chain {
				bad := newFake("@b/b", source.KindOfficial)
				bad.failErr = boom
				return source.Chain{newFake("@a/a", source.KindRegistry, "x"), bad}
			},
			wantErr: boom,
		},
		{
			name:    "empty chain",
			chain:   func() source.Chain { return nil },
			wantErr: source.ErrNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := tc.chain().Find(ctx, skill.Ref{Name: "x"})
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("err = %v, want %v", err, tc.wantErr)
				}
				for _, m := range tc.wantMsg {
					if !strings.Contains(err.Error(), m) {
						t.Errorf("err %q does not mention %q", err, m)
					}
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Path != tc.wantPath {
				t.Errorf("path = %q, want %q", got.Path, tc.wantPath)
			}
		})
	}
}

func TestChainFindQualified(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	a := newFake("@a/a", source.KindRegistry, "x")
	b := newFake("@b/b", source.KindOfficial, "x")
	chain := source.Chain{a, b}

	got, err := chain.Find(ctx, skill.Ref{Source: "@b/b", Name: "x"})
	if err != nil || got.Path != "@b/b/x" {
		t.Fatalf("qualified find = %+v, %v; want @b/b/x", got, err)
	}
	if len(a.calls) != 0 {
		t.Errorf("a qualified reference consulted another source: %v", a.calls)
	}

	if _, err = chain.Find(ctx, skill.Ref{Source: "@c/c", Name: "x"}); !errors.Is(err, source.ErrUnknownSource) {
		t.Errorf("unknown source err = %v, want ErrUnknownSource, not a fallback", err)
	}

	_, err = chain.Find(ctx, skill.Ref{Source: "@a/a", Name: "missing"})
	if !errors.Is(err, source.ErrNotFound) {
		t.Errorf("qualified miss err = %v, want ErrNotFound", err)
	}
}

func TestChainFindLazySource(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("a hit elsewhere never triggers the fetching Find", func(t *testing.T) {
		t.Parallel()
		reg := newFake("@a/a", source.KindRegistry, "x")
		off := lazyFake{newFake("@b/b", source.KindOfficial)}
		if _, err := (source.Chain{reg, off}).Find(ctx, skill.Ref{Name: "x"}); err != nil {
			t.Fatal(err)
		}
		for _, c := range off.calls {
			if strings.HasPrefix(c, "Find:") {
				t.Errorf("lazy source was asked to fetch although another source had the name: %v", off.calls)
			}
		}
	})

	t.Run("a cached duplicate is still ambiguous", func(t *testing.T) {
		t.Parallel()
		reg := newFake("@a/a", source.KindRegistry, "x")
		off := lazyFake{newFake("@b/b", source.KindOfficial, "x")}
		_, err := (source.Chain{reg, off}).Find(ctx, skill.Ref{Name: "x"})
		if !errors.Is(err, source.ErrAmbiguous) {
			t.Fatalf("err = %v, want ErrAmbiguous", err)
		}
	})

	t.Run("a miss everywhere falls back to the fetching Find", func(t *testing.T) {
		t.Parallel()
		reg := newFake("@a/a", source.KindRegistry)
		off := &fetchingFake{fakeSource: newFake("@b/b", source.KindOfficial), fetched: "x"}
		got, err := (source.Chain{reg, off}).Find(ctx, skill.Ref{Name: "x"})
		if err != nil || got.Path != "@b/b/x" {
			t.Fatalf("got %+v, %v; want the skill found after the fetch", got, err)
		}
	})
}

// fetchingFake finds nothing in its cache but "fetches" name on Find, the way
// Official clones anthropics/skills on first use.
type fetchingFake struct {
	*fakeSource
	fetched string
}

func (f *fetchingFake) FindCached(context.Context, string) (skill.Skill, error) {
	return skill.Skill{}, fmt.Errorf("%w: cache empty", source.ErrNotFound)
}

func (f *fetchingFake) Find(_ context.Context, name string) (skill.Skill, error) {
	f.skills[f.fetched] = true
	f.calls = append(f.calls, "Find:"+name)
	return f.lookup(name)
}

func TestChainSelect(t *testing.T) {
	t.Parallel()
	reg := newFake("@a/a", source.KindRegistry, "x")
	off := newFake("@b/b", source.KindOfficial, "y")
	chain := source.Chain{reg, off}

	only := chain.Select(source.KindOfficial)
	if len(only) != 1 || only[0].Spec().ID != "@b/b" {
		t.Fatalf("Select(official) = %v", only)
	}
	both := chain.Select(source.KindRegistry, source.KindOfficial)
	if len(both) != 2 {
		t.Fatalf("Select(registry, official) kept %d sources, want 2", len(both))
	}
	if _, err := only.Find(context.Background(), skill.Ref{Name: "x"}); !errors.Is(err, source.ErrNotFound) {
		t.Errorf("a selected chain must not reach excluded sources, err = %v", err)
	}
}

func TestChainFindQualifiedNormalisesTheSourceAndRejectsPaths(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	reg := newFake(skill.DefaultRegistryID, source.KindRegistry, "x")
	off := newFake(skill.OfficialID, source.KindOfficial, "x")
	chain := source.Chain{reg, off}

	for _, qualified := range []string{
		"@ralvarez-ai/ralvaskills/x", "@ralvarezdev/ralvaskills/x", "@Ralvarez-AI/RalvaSkills/x",
	} {
		ref, _, err := skill.ParseRef(qualified)
		if err != nil {
			t.Fatal(err)
		}
		got, findErr := chain.Find(ctx, ref)
		if findErr != nil || got.Path != skill.DefaultRegistryID+"/x" {
			t.Errorf("Find(%s) = %+v, %v; want the registry's skill", qualified, got, findErr)
		}
	}

	withPath := skill.Ref{Source: skill.OfficialID, Name: "x", Path: "a/x"}
	if _, err := chain.Find(ctx, withPath); !errors.Is(err, source.ErrUnsupported) {
		t.Errorf("a path selection err = %v, want ErrUnsupported until a source can honor it", err)
	}
}
