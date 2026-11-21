package source

import (
	"errors"
	"fmt"
)

// New is the only way to build a Resolver. Every command, the MCP server and
// the update check construct their sources through it, so the rules about
// which sources exist and how they are consulted live in one place.
//
// It rejects a Spec it cannot fully honor rather than ignoring a field.
//
//nolint:ireturn // returns interface so callers depend on the resolver, not the transport
func New(spec Spec, opts Options) (Resolver, error) {
	if spec.Rev != "" || len(spec.Roots) > 0 || len(spec.Excludes) > 0 || spec.Private {
		return nil, fmt.Errorf("source %s: rev, roots, excludes and private are not supported yet", spec.ID)
	}
	switch spec.Type {
	case KindRegistry:
		return newRegistry(spec, opts.RegistryCache), nil
	case KindPath:
		return newLocal(spec), nil
	case KindOfficial:
		return newOfficial(spec, opts.OfficialCache), nil
	default:
		return nil, fmt.Errorf("source %s: unsupported kind %q", spec.ID, spec.Type)
	}
}

// Build constructs a Chain from specs, keeping declaration order.
func Build(specs []Spec, opts Options) (Chain, error) {
	chain := make(Chain, 0, len(specs))
	var errs []error
	for _, spec := range specs {
		r, err := New(spec, opts)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		chain = append(chain, r)
	}
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}
	return chain, nil
}
