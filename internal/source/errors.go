// Package source resolves skills from their respective origins.
package source

import "errors"

var (
	// ErrNotFound is returned when a skill cannot be located in a source.
	ErrNotFound = errors.New("skill not found")

	// ErrAmbiguous is returned when a bare name is provided by more than one
	// source. It is never resolved by order.
	ErrAmbiguous = errors.New("skill name is ambiguous")

	// ErrDuplicate is returned when one source holds two skills with the same
	// name. A source is never resolved by order, so it is an error naming both.
	ErrDuplicate = errors.New("duplicate skill name in a source")

	// ErrTreeMismatch is returned when a skill restored at a locked commit does
	// not have the tree id the lock recorded: the content is not what was locked.
	ErrTreeMismatch = errors.New("skill content does not match the locked tree")

	// ErrNeverFetched is returned when --offline asks for a source that has no
	// index on this machine yet.
	ErrNeverFetched = errors.New("source has never been fetched")

	// ErrUnsupported is returned when a reference asks for something no source
	// in this build can provide yet.
	ErrUnsupported = errors.New("unsupported reference")

	// ErrUnknownSource is returned when a qualified reference names a source
	// that is not in the chain. It never falls back to another source.
	ErrUnknownSource = errors.New("unknown source")
)
