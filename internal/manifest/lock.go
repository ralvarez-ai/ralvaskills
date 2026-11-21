package manifest

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
	"go.ralva.dev/ralvaskills/v3/internal/fsx"
	"go.ralva.dev/ralvaskills/v3/internal/git"
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

		// Git says where a git skill was read, and is nil for every other
		// source. Its Commit and Tree are the integrity record: a restore must
		// find that tree at that commit.
		Git *skill.GitRef `toml:"-"`
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
const transportGit = "git"

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
		URL         string   `toml:"url,omitempty"`
		Rev         string   `toml:"rev,omitempty"`
		Commit      string   `toml:"commit,omitempty"`
		Tree        string   `toml:"tree,omitempty"`
		Path        string   `toml:"path,omitempty"`
	}
)

// transportOf names how a skill kind is fetched.
func transportOf(src skill.Source) string {
	switch src {
	case skill.SourceRegistry:
		return "registry"
	case skill.SourceOfficial:
		return transportGit
	case skill.SourceLocal:
		return "path"
	default:
		if src.IsDeclared() {
			return transportGit
		}
		return ""
	}
}

// ReadLock decodes rsk.lock from rskDir. Returns an empty Lock if the file
// does not exist.
func ReadLock(rskDir string) (Lock, error) {
	file := LockPath(rskDir)
	if _, statErr := os.Stat(file); errors.Is(statErr, os.ErrNotExist) {
		return Lock{}, nil
	}

	var f lockFile
	if _, err := toml.DecodeFile(file, &f); err != nil {
		return Lock{}, fmt.Errorf("read %s: %w", file, err)
	}
	if err := schema.Check(f.SchemaVersion, schema.Lock); err != nil {
		return Lock{}, fmt.Errorf("rsk.lock: %w", err)
	}

	l := Lock{SchemaVersion: f.SchemaVersion}
	for i := range f.Skills {
		entry, err := entryFromFile(&f.Skills[i])
		if err != nil {
			return Lock{}, fmt.Errorf("rsk.lock: %w", err)
		}
		l.Skills = append(l.Skills, entry)
	}
	return l, nil
}

func entryFromFile(e *lockEntryFile) (LockEntry, error) {
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
	entry := LockEntry{
		Name: ref.Name, Version: e.Version, Source: src, Permissions: e.Permissions, Granted: e.Granted,
	}
	if !src.IsDeclared() {
		// Older locks recorded a machine cache path under "path". It means
		// nothing outside a git entry, so it is dropped rather than refused.
		e.Path = ""
	}
	if e.URL != "" || e.Rev != "" || e.Commit != "" || e.Tree != "" || e.Path != "" {
		entry.Git = &skill.GitRef{URL: e.URL, Rev: e.Rev, Commit: e.Commit, Tree: e.Tree, Path: e.Path}
	}
	if err = checkGitFields(&entry, e.Transport); err != nil {
		return LockEntry{}, fmt.Errorf("%s: %w", e.Name, err)
	}
	return entry, nil
}

// checkGitFields holds a lock entry to its source: a git entry must carry a
// well-formed URL, commit, tree and path, and no other entry may carry any. A
// lock file is untrusted input, since restoring from it fetches a URL it names.
func checkGitFields(e *LockEntry, transport string) error {
	if !e.Source.IsDeclared() {
		if e.Git != nil {
			return errors.New("git fields on an entry that is not from a git source")
		}
		return nil
	}
	if transport != transportGit {
		return fmt.Errorf("a git source needs transport %q, got %q", transportGit, transport)
	}
	if e.Git == nil {
		return errors.New("a git entry needs url, commit, tree and path")
	}
	if _, err := git.ParseRemoteURL(e.Git.URL); err != nil {
		return err
	}
	for _, check := range []error{
		git.ValidateCommit(e.Git.Commit), git.ValidateCommit(e.Git.Tree), git.ValidateRepoPath(e.Git.Path),
	} {
		if check != nil {
			return check
		}
	}
	if e.Git.Rev != "" {
		if err := git.ValidateRef(e.Git.Rev); err != nil {
			return err
		}
	}
	if path.Base(e.Git.Path) != e.Name {
		return fmt.Errorf("path %s does not end in the skill name", e.Git.Path)
	}
	return nil
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
			if e.Git != nil {
				out.URL, out.Rev, out.Commit, out.Tree, out.Path = e.Git.URL, e.Git.Rev, e.Git.Commit, e.Git.Tree, e.Git.Path
			}
		}
		f.Skills[i] = out
	}

	file := LockPath(rskDir)
	if err := os.MkdirAll(filepath.Dir(file), fsperm.Dir); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}
	return fsx.WriteAtomic(file, TempFilePattern, func(w io.Writer) error {
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
