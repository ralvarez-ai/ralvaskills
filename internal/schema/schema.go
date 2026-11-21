// Package schema defines version constants for every rsk-managed file format
// and provides helpers for forward-compatibility validation and migration.
//
// Every persisted file (config.json, rsk.mod, rsk.lock, catalog.toml) carries
// a top-level "version" integer field. When rsk reads a file it calls Check to
// confirm the stored version is supported. When rsk writes a file it stamps the
// current version. Files that predate versioning carry version 0, which is
// treated as equivalent to version 1 (the initial stable format).
package schema

import (
	"errors"
	"fmt"
)

// Version is the schema version number embedded in persisted rsk files.
type Version int

const (
	// Config is the current schema version for ~/.config/rsk/config.json.
	Config Version = 1

	// Base is the version written whenever nothing in a file needs a newer
	// one, so a project keeps working for teammates on an older rsk until it
	// adopts a feature the older format cannot express.
	Base Version = 1

	// Mod is the newest schema version of .rsk/rsk.mod this build reads and
	// writes. Version 2 carries qualified skill keys (@owner/repo/name).
	Mod Version = 2

	// Lock is the newest schema version of .rsk/rsk.lock this build reads and
	// writes. Version 2 carries qualified keys, source identities and the
	// transport of each skill.
	Lock Version = 2

	// Catalog is the current schema version for ~/.config/rsk/catalog.toml
	// and any user-override catalog.
	Catalog Version = 1
)

// ErrUnsupported is returned by Check when a file carries a version that this
// build of rsk cannot read.
var ErrUnsupported = errors.New("unsupported schema version")

// Check returns an error wrapping ErrUnsupported if stored is newer than
// current, the newest version this build reads. Every version from 0 up to
// current is accepted, so a file written before a bump keeps loading. Version 0
// is legacy (pre-versioning) and equivalent to version 1.
func Check(stored, current Version) error {
	if stored >= 0 && stored <= current {
		return nil
	}
	return fmt.Errorf("%w: file version %d, rsk supports %d — upgrade rsk or re-run 'rsk init'",
		ErrUnsupported, stored, current)
}
