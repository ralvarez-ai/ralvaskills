package manifest

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
	"go.ralva.dev/ralvaskills/v3/internal/fsx"
	"go.ralva.dev/ralvaskills/v3/internal/schema"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

type (
	// LockEntry records the resolved state of one installed skill. It
	// deliberately keeps no cache path: the path is machine state, the
	// version and source are what makes the lock reproducible.
	LockEntry struct {
		Name    string       `toml:"name"`
		Version string       `toml:"version"`
		Source  skill.Source `toml:"source"`

		// Permissions is the skill's allowed-tools declaration at the locked
		// version, normalized. Update diffs the new declaration against this.
		Permissions []string `toml:"permissions,omitempty"`

		// Granted is the client rules rsk wrote for this skill. Revocation
		// removes exactly these, never a rule the user added by hand.
		Granted []string `toml:"granted,omitempty"`
	}

	// Lock is the in-memory representation of rsk.lock.
	Lock struct {
		// SchemaVersion identifies the file format version. Version 0 is
		// accepted as legacy (pre-versioning) and treated as version 1.
		SchemaVersion schema.Version `toml:"version,omitempty"`
		Skills        []LockEntry    `toml:"skills"`
	}
)

// lockFile and lockEntryFile are the on-disk shape of rsk.lock. In memory an
// entry's name is bare and its Source a kind; on disk a version 2 file carries
// the qualified key, the source identity and the transport.
type (
	lockFile struct {
		SchemaVersion schema.Version  `toml:"version,omitempty"`
		Skills        []lockEntryFile `toml:"skills"`
	}

	lockEntryFile struct {
		Name        string   `toml:"name"`
		Version     string   `toml:"version"`
		Source      string   `toml:"source"`
		Permissions []string `toml:"permissions,omitempty"`
		Granted     []string `toml:"granted,omitempty"`
		Transport   string   `toml:"transport,omitempty"`
	}
)

// transportOf names how a skill kind is fetched.
func transportOf(src skill.Source) string {
	switch src {
	case skill.SourceRegistry:
		return "registry"
	case skill.SourceOfficial:
		return "git"
	case skill.SourceLocal:
		return "path"
	default:
		return ""
	}
}

// ReadLock decodes rsk.lock from rskDir. Returns an empty Lock if the file
// does not exist.
func ReadLock(rskDir string) (Lock, error) {
	path := LockPath(rskDir)
	if _, statErr := os.Stat(path); errors.Is(statErr, os.ErrNotExist) {
		return Lock{}, nil
	}

	var f lockFile
	if _, err := toml.DecodeFile(path, &f); err != nil {
		return Lock{}, fmt.Errorf("read %s: %w", path, err)
	}
	if err := schema.Check(f.SchemaVersion, schema.Lock); err != nil {
		return Lock{}, fmt.Errorf("rsk.lock: %w", err)
	}

	l := Lock{SchemaVersion: f.SchemaVersion}
	for _, e := range f.Skills {
		entry, err := entryFromFile(e)
		if err != nil {
			return Lock{}, fmt.Errorf("rsk.lock: %w", err)
		}
		l.Skills = append(l.Skills, entry)
	}
	return l, nil
}

func entryFromFile(e lockEntryFile) (LockEntry, error) {
	ref, err := skill.ParseKey(e.Name)
	if err != nil {
		return LockEntry{}, err
	}
	src, err := skill.SourceFromIdentity(e.Source)
	if err != nil {
		return LockEntry{}, fmt.Errorf("%s: %w", e.Name, err)
	}
	if ref.Source != "" {
		keySrc, keyErr := skill.SourceFromIdentity(ref.Source)
		if keyErr != nil {
			return LockEntry{}, fmt.Errorf("%s: %w", e.Name, keyErr)
		}
		if keySrc != src {
			return LockEntry{}, fmt.Errorf("%s: key names source %s but source is %q", e.Name, ref.Source, e.Source)
		}
	}
	return LockEntry{
		Name: ref.Name, Version: e.Version, Source: src, Permissions: e.Permissions, Granted: e.Granted,
	}, nil
}

// WriteLock encodes l as TOML to rskDir/rsk.lock atomically, creating parent
// directories as needed. It writes version 2 only when an entry comes from a
// source other than the default registry; otherwise the file stays version 1.
func WriteLock(rskDir string, l Lock) error {
	f := lockFile{SchemaVersion: schema.Base, Skills: make([]lockEntryFile, len(l.Skills))}
	for _, e := range l.Skills {
		if id := e.Source.Identity(); id != "" && id != skill.DefaultRegistryID {
			f.SchemaVersion = schema.Lock
			break
		}
	}
	v2 := f.SchemaVersion == schema.Lock
	for i, e := range l.Skills {
		out := lockEntryFile{
			Name: e.Name, Version: e.Version, Source: e.Source.String(),
			Permissions: e.Permissions, Granted: e.Granted,
		}
		if v2 {
			if id := e.Source.Identity(); id != "" {
				out.Name = skill.Ref{Source: id, Name: e.Name}.Key()
				out.Source = id
			}
			out.Transport = transportOf(e.Source)
		}
		f.Skills[i] = out
	}

	path := LockPath(rskDir)
	if err := os.MkdirAll(filepath.Dir(path), fsperm.Dir); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}
	return fsx.WriteAtomic(path, TempFilePattern, func(w io.Writer) error {
		return toml.NewEncoder(w).Encode(f)
	})
}

// UpsertLockEntry adds entry to l if no entry with the same name exists,
// or replaces the existing entry in place.
func UpsertLockEntry(l Lock, entry LockEntry) Lock {
	for i, e := range l.Skills {
		if e.Name == entry.Name {
			l.Skills[i] = entry
			return l
		}
	}
	l.Skills = append(l.Skills, entry)
	return l
}

// RemoveLockEntry removes the entry with the given name from l. No-op if
// the name is not present.
func RemoveLockEntry(l Lock, name string) Lock {
	kept := make([]LockEntry, 0, len(l.Skills))
	for _, e := range l.Skills {
		if e.Name != name {
			kept = append(kept, e)
		}
	}
	l.Skills = kept
	return l
}

// FindLockEntry returns the entry for name and whether it exists.
func FindLockEntry(l Lock, name string) (LockEntry, bool) {
	for _, e := range l.Skills {
		if e.Name == name {
			return e, true
		}
	}
	return LockEntry{}, false
}
