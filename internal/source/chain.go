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
func (c Chain) Find(ctx context.Context, ref skill.Ref) (skill.Skill, error) {
	if ref.Path != "" {
		return skill.Skill{}, fmt.Errorf("%w: selecting %q by path needs a git source", ErrUnsupported, ref.String())
	}
	if ref.Source != "" {
		return c.findIn(ctx, ref)
	}

	hits, err := c.collect(func(r Resolver) (skill.Skill, error) {
		return r.Find(ctx, ref.Name)
	})
	if err != nil {
		return skill.Skill{}, err
	}
	return pickOne(ref.Name, hits)
}

// Answer is one source's reply to a name lookup. Err is nil when the source has
// the skill, wraps ErrNotFound when it does not, and is any other error when
// the source could not be consulted.
type Answer struct {
	Spec  Spec
	Skill skill.Skill
	Err   error
}

// Survey asks every source for name and keeps each answer, where Find stops at
// the first source it cannot consult. It exists to explain a resolution, so
// `rsk resolve` can say who has a name and who could not be asked.
func (c Chain) Survey(ctx context.Context, name string) []Answer {
	answers := make([]Answer, len(c))
	for i, r := range c {
		answers[i].Spec = r.Spec()
		answers[i].Skill, answers[i].Err = r.Find(ctx, name)
	}
	return answers
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

// versionFinder is implemented by a source that can resolve an exact version.
type versionFinder interface {
	FindVersion(ctx context.Context, name, version string) (skill.Skill, error)
}

// FindVersion resolves name at an exact version in the source with identity id.
// It never searches: the caller already knows which source the skill is from.
func (c Chain) FindVersion(ctx context.Context, id, name, version string) (skill.Skill, error) {
	for _, r := range c {
		if skill.CanonicalSourceID(r.Spec().ID) != skill.CanonicalSourceID(id) {
			continue
		}
		vf, ok := r.(versionFinder)
		if !ok {
			return skill.Skill{}, fmt.Errorf("%w: source %s cannot resolve a version", ErrUnsupported, id)
		}
		return vf.FindVersion(ctx, name, version)
	}
	return skill.Skill{}, fmt.Errorf("%w: %q for skill %q", ErrUnknownSource, id, name)
}

// ByID returns the sources with the given identities, in chain order.
func (c Chain) ByID(ids ...string) Chain {
	var out Chain
	for _, r := range c {
		for _, id := range ids {
			if skill.CanonicalSourceID(r.Spec().ID) == skill.CanonicalSourceID(id) {
				out = append(out, r)
				break
			}
		}
	}
	return out
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

// Primary returns the hosted registry or the local clone: the sources a bundle
// member is read from. A path source the user declared is not one of them.
func (c Chain) Primary() Chain {
	var out Chain
	for _, r := range c {
		spec := r.Spec()
		if spec.Type == KindRegistry || (spec.Type == KindPath && !spec.Declared) {
			out = append(out, r)
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
