package skill

import (
	"errors"
	"fmt"
	"path"
	"regexp"
	"strings"
)

const (
	// DefaultRegistryID is the canonical identity of the hosted registry.
	DefaultRegistryID = "@ralvarez-ai/ralvaskills"

	// OfficialID is the identity of the anthropics/skills cache.
	OfficialID = "@anthropics/skills"

	// LocalScope is the identity scope reserved for path sources.
	LocalScope = "@local"
)

var (
	ownerPattern = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)
	// repoPattern admits what GitHub allows in a repository name, except a
	// name made only of dots: "." and ".." would escape a cache directory.
	repoPattern = regexp.MustCompile(`^[a-z0-9_-][a-z0-9._-]*$|^\.[a-z0-9._-]*[a-z0-9_-][a-z0-9._-]*$`)

	// builtinAliases maps an alternative spelling of a built-in source to its
	// canonical identity. It is the only place an alias can come from: never
	// from fetched data.
	builtinAliases = map[string]string{
		"@ralvarezdev/ralvaskills": DefaultRegistryID,
	}
)

// Ref is a parsed skill reference.
//
//	name                              bare, resolved across every source
//	@owner/repo/name                  qualified: resolved in exactly that source
//	@owner/repo//path/to/name         qualified with a repo-relative path
//
// A version, when given, follows the last "@" and is returned beside the Ref,
// not inside it.
type Ref struct {
	// Source is the canonical identity of the source, or empty for a bare name.
	Source string

	// Name is the skill name: the basename of its directory.
	Name string

	// Path is the repo-relative directory selected with the "//" form.
	Path string
}

// ParseRef parses a reference with an optional "@version" suffix.
func ParseRef(arg string) (ref Ref, version string, err error) {
	if !strings.HasPrefix(arg, "@") {
		name, bareVersion, _ := strings.Cut(arg, "@")
		if name == "" {
			return Ref{}, "", fmt.Errorf("skill name must not be empty (got %q)", arg)
		}
		if strings.Contains(name, "/") {
			return Ref{}, "", fmt.Errorf(
				"skill name %q must not contain \"/\"; a qualified reference starts with \"@\"", name,
			)
		}
		return Ref{Name: name}, bareVersion, nil
	}

	body := arg[1:]
	if i := strings.LastIndex(body, "@"); i >= 0 {
		body, version = body[:i], body[i+1:]
	}
	ref, err = parseQualified(body)
	if err != nil {
		return Ref{}, "", fmt.Errorf("reference %q: %w", arg, err)
	}
	return ref, version, nil
}

// ParseKey parses a key as stored in rsk.mod and rsk.lock: a reference with no
// version and no path.
func ParseKey(key string) (Ref, error) {
	ref, version, err := ParseRef(key)
	if err != nil {
		return Ref{}, err
	}
	if version != "" || strings.Contains(strings.TrimPrefix(key, "@"), "@") {
		return Ref{}, fmt.Errorf("key %q must not carry a version", key)
	}
	if ref.Path != "" {
		return Ref{}, fmt.Errorf("key %q must not carry a path", key)
	}
	return ref, nil
}

// ParseSourceID validates a source identity ("@owner/repo") and returns its
// canonical form: case-folded, with a built-in alias mapped to its target.
func ParseSourceID(id string) (string, error) {
	body, ok := strings.CutPrefix(id, "@")
	owner, repo, hasRepo := strings.Cut(body, "/")
	if !ok || !hasRepo || strings.Contains(repo, "/") ||
		!ownerPattern.MatchString(strings.ToLower(owner)) || !repoPattern.MatchString(strings.ToLower(repo)) {
		return "", fmt.Errorf("source identity %q must look like @owner/repo", id)
	}
	return CanonicalSourceID(id), nil
}

func parseQualified(body string) (Ref, error) {
	owner, rest, _ := strings.Cut(body, "/")
	repo, tail, _ := strings.Cut(rest, "/")
	if !ownerPattern.MatchString(strings.ToLower(owner)) || !repoPattern.MatchString(strings.ToLower(repo)) {
		return Ref{}, errors.New("expected @owner/repo/name")
	}
	ref := Ref{Source: CanonicalSourceID("@" + owner + "/" + repo)}

	if dir, ok := strings.CutPrefix(tail, "/"); ok {
		if err := validateRepoPath(dir); err != nil {
			return Ref{}, err
		}
		ref.Path = dir
		ref.Name = path.Base(dir)
	} else {
		ref.Name = tail
	}
	if err := ValidateName(ref.Name); err != nil {
		return Ref{}, err
	}
	return ref, nil
}

func validateRepoPath(p string) error {
	if p == "" {
		return errors.New("path after \"//\" is empty")
	}
	for seg := range strings.SplitSeq(p, "/") {
		if seg == "" || seg == "." || seg == ".." {
			return fmt.Errorf("path %q must be relative, with no empty, \".\" or \"..\" segment", p)
		}
	}
	return nil
}

// CanonicalSourceID folds case and maps a built-in alias to its canonical
// identity, so two spellings of one source never become two entries.
func CanonicalSourceID(id string) string {
	id = strings.ToLower(id)
	if canonical, ok := builtinAliases[id]; ok {
		return canonical
	}
	return id
}

// Key is the form a skill takes in rsk.mod and rsk.lock: the bare name for a
// skill from the default registry, the qualified form for any other source.
// It depends on the skill and not on the machine, so one committed file reads
// the same everywhere.
func (r Ref) Key() string {
	if r.Source == "" || r.Source == DefaultRegistryID {
		return r.Name
	}
	return r.Source + "/" + r.Name
}

// String returns the reference as written, with its path when it has one.
func (r Ref) String() string {
	switch {
	case r.Source == "":
		return r.Name
	case r.Path != "":
		return r.Source + "//" + r.Path
	default:
		return r.Source + "/" + r.Name
	}
}
