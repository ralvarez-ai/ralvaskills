package source

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

// Chain is an ordered set of sources. The order is the declaration order and
// only keeps listings stable: it never decides a resolution.
type Chain []Resolver

// cachedFinder is implemented by a source whose Find may reach the network to
// fill its own cache (Official clones anthropics/skills on first use).
// FindCached answers from what the source already holds and never fetches.
type cachedFinder interface {
	FindCached(ctx context.Context, name string) (skill.Skill, error)
}

// Find resolves ref.
//
// A qualified reference resolves in exactly the named source; an unknown
// source is an error, not a fallback.
//
// A bare name asks every source, and resolves only if exactly one has it. Two
// sources providing it is an error naming both. A source that cannot be
// consulted fails the whole lookup, never letting another source answer in its
// place: otherwise an unreachable private source would silently resolve the
// name from a public one.
//
// A source that can fetch on demand is first asked only for what it already
// holds, so a name another source provides never triggers a fetch. Only when no
// source has the name is it asked to fetch.
func (c Chain) Find(ctx context.Context, ref skill.Ref) (skill.Skill, error) {
	if ref.Path != "" {
		return skill.Skill{}, fmt.Errorf("%w: selecting %q by path needs a git source", ErrUnsupported, ref.String())
	}
	if ref.Source != "" {
		return c.findIn(ctx, ref)
	}

	hits, err := c.collect(func(r Resolver) (skill.Skill, error) {
		if cf, ok := r.(cachedFinder); ok {
			return cf.FindCached(ctx, ref.Name)
		}
		return r.Find(ctx, ref.Name)
	})
	if err != nil {
		return skill.Skill{}, err
	}
	if len(hits) == 0 {
		hits, err = c.collect(func(r Resolver) (skill.Skill, error) {
			if _, ok := r.(cachedFinder); !ok {
				return skill.Skill{}, ErrNotFound
			}
			return r.Find(ctx, ref.Name)
		})
		if err != nil {
			return skill.Skill{}, err
		}
	}
	return pickOne(ref.Name, hits)
}

func pickOne(name string, hits []hit) (skill.Skill, error) {
	switch len(hits) {
	case 0:
		return skill.Skill{}, fmt.Errorf("%w: %q is in no source", ErrNotFound, name)
	case 1:
		return hits[0].skill, nil
	default:
		ids := make([]string, len(hits))
		for i := range hits {
			ids[i] = hits[i].source
		}
		return skill.Skill{}, fmt.Errorf("%w: %q is provided by %s", ErrAmbiguous, name, strings.Join(ids, " and "))
	}
}

// Select returns the sources of the given kinds, in chain order.
func (c Chain) Select(kinds ...Kind) Chain {
	var out Chain
	for _, r := range c {
		for _, k := range kinds {
			if r.Spec().Type == k {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

// Registry returns the chain's hosted registry, for the callers that need its
// index rather than a single skill. It reports false when there is none.
func (c Chain) Registry() (*Registry, bool) {
	for _, r := range c {
		if reg, ok := r.(*Registry); ok {
			return reg, true
		}
	}
	return nil, false
}

func (c Chain) findIn(ctx context.Context, ref skill.Ref) (skill.Skill, error) {
	for _, r := range c {
		if skill.CanonicalSourceID(r.Spec().ID) == skill.CanonicalSourceID(ref.Source) {
			return r.Find(ctx, ref.Name)
		}
	}
	return skill.Skill{}, fmt.Errorf("%w: %q for skill %q", ErrUnknownSource, ref.Source, ref.Name)
}

type hit struct {
	source string
	skill  skill.Skill
}

// collect asks every source through find. A not-found answer skips the source;
// any other error aborts the lookup (fail closed).
func (c Chain) collect(find func(Resolver) (skill.Skill, error)) ([]hit, error) {
	var hits []hit
	for _, r := range c {
		s, err := find(r)
		switch {
		case err == nil:
			hits = append(hits, hit{source: r.Spec().ID, skill: s})
		case errors.Is(err, ErrNotFound):
		default:
			return nil, err
		}
	}
	return hits, nil
}
