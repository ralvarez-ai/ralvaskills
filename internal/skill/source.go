package skill

import (
	"fmt"
	"strings"
)

// Source identifies where an installed skill originates.
type Source string

const (
	// SourceLocal is a skill symlinked from the local ralvaskills repo clone.
	SourceLocal Source = "local"

	// SourceOfficial is a skill fetched from the official anthropic/skills cache.
	SourceOfficial Source = "official"

	// SourceRegistry is a skill downloaded from the hosted registry.
	SourceRegistry Source = "registry"

	// SourceGit is the kind of a skill read from a declared git repository.
	// Such a skill's Source is the repository's identity itself ("@owner/repo",
	// see DeclaredSource), so one value says both what kind of source it is and
	// which one; Kind folds it back to SourceGit. ParseSource does not accept it.
	SourceGit Source = "git"
)

// ParseSource converts a string to its Source value. Returns an error for
// unrecognized values.
func ParseSource(s string) (Source, error) {
	switch Source(s) {
	case SourceLocal, SourceOfficial, SourceRegistry:
		return Source(s), nil
	default:
		return "", fmt.Errorf("unknown skill source %q", s)
	}
}

// DeclaredSource is the Source of a skill read from the declared git source
// with the given canonical identity.
func DeclaredSource(id string) Source {
	return Source(id)
}

// IsDeclared reports whether s names a declared git source by its identity.
func (s Source) IsDeclared() bool {
	return strings.HasPrefix(string(s), "@")
}

// Kind returns s itself for the built-in sources, and SourceGit for a declared
// one.
func (s Source) Kind() Source {
	if s.IsDeclared() {
		return SourceGit
	}
	return s
}

// Identity returns the canonical identity of the source a skill kind comes
// from. A local clone has none (neither does an unknown kind): its identity is the name of a directory on one
// machine, which a committed file must not carry.
func (s Source) Identity() string {
	switch s {
	case SourceRegistry:
		return DefaultRegistryID
	case SourceOfficial:
		return OfficialID
	default:
		if s.IsDeclared() {
			return string(s)
		}
		return ""
	}
}

// SourceFromIdentity maps a source identity, or one of the older kind names
// ("registry", "official", "local") that rsk.lock v1 wrote, to its Source. An
// identity this rsk has no transport for is an error, never a guess.
func SourceFromIdentity(id string) (Source, error) {
	if src, err := ParseSource(id); err == nil {
		return src, nil
	}
	switch canonical := CanonicalSourceID(id); {
	case canonical == DefaultRegistryID:
		return SourceRegistry, nil
	case canonical == OfficialID:
		return SourceOfficial, nil
	case strings.HasPrefix(canonical, LocalScope+"/"):
		return SourceLocal, nil
	default:
		return "", fmt.Errorf("source %q is not supported by this rsk", id)
	}
}

// String returns the string representation of the source, identical to its
// wire encoding.
func (s Source) String() string {
	return string(s)
}

// Label returns the short display label used in table output.
func (s Source) Label() string {
	switch s.Kind() {
	case SourceLocal:
		return "ralva"
	case SourceOfficial:
		return "anthr"
	case SourceRegistry:
		return "reg  "
	case SourceGit:
		return "git  "
	default:
		return "?????"
	}
}

// MarshalText implements encoding.TextMarshaler so Source encodes as its
// string value in TOML, JSON, and other text-based formats.
func (s Source) MarshalText() ([]byte, error) {
	return []byte(s), nil
}

// UnmarshalText implements encoding.TextUnmarshaler so Source decodes from
// its string value in TOML, JSON, and other text-based formats.
func (s *Source) UnmarshalText(b []byte) error {
	src, err := ParseSource(string(b))
	if err != nil {
		return err
	}
	*s = src
	return nil
}
