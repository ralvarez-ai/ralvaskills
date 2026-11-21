package source

import (
	"context"

	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

// Resolver is one source in a Chain. Commands depend on this interface, and
// on Chain, rather than on concrete types, so that source selection stays in
// this package.
type Resolver interface {
	// Spec returns the descriptor the source was built from, including its
	// canonical identity.
	Spec() Spec

	// All returns every skill available from this source.
	All(ctx context.Context) ([]skill.Skill, error)

	// Find returns the skill with the given name. Returns an error wrapping
	// ErrNotFound when the skill is not present in this source.
	Find(ctx context.Context, name string) (skill.Skill, error)
}
