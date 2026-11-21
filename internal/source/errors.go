// Package source resolves skills from their respective origins.
package source

import "errors"

var (
	// ErrNotFound is returned when a skill cannot be located in a source.
	ErrNotFound = errors.New("skill not found")

	// ErrAmbiguous is returned when a bare name is provided by more than one
	// source. It is never resolved by order.
	ErrAmbiguous = errors.New("skill name is ambiguous")

	// ErrUnsupported is returned when a reference asks for something no source
	// in this build can provide yet.
	ErrUnsupported = errors.New("unsupported reference")

	// ErrUnknownSource is returned when a qualified reference names a source
	// that is not in the chain. It never falls back to another source.
	ErrUnknownSource = errors.New("unknown source")
)
