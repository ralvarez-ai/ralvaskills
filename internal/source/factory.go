package source

import (
	"errors"
	"fmt"
)

// New is the only way to build a Resolver. Every command, the MCP server and
// the update check construct their sources through it, so the rules about
// which sources exist and how they are consulted live in one place.
//
// It rejects a Spec it cannot fully honor rather than ignoring a field, and a
// git source it cannot gate.
//
//nolint:ireturn // returns interface so callers depend on the resolver, not the transport
func New(spec Spec, opts Options) (Resolver, error) {
	if spec.Private {
		return nil, fmt.Errorf("source %s: private is not supported yet", spec.ID)
	}
	if spec.Type != KindGit && spec.Rev != "" {
		return nil, fmt.Errorf("source %s: rev applies to git sources only", spec.ID)
	}
	walksRoots := spec.Type == KindGit || (spec.Type == KindPath && spec.Declared)
	if !walksRoots && (len(spec.Roots) > 0 || len(spec.Excludes) > 0) {
		return nil, fmt.Errorf("source %s: roots and excludes apply to git and declared path sources only", spec.ID)
	}
	switch spec.Type {
	case KindGit:
		return newGit(spec, opts)
	case KindRegistry:
		return newRegistry(spec, opts.RegistryCache), nil
	case KindPath:
		if spec.Declared {
			return newPath(spec)
		}
		return newLocal(spec), nil
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
