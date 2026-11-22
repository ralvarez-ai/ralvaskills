// Package sources reads and writes the machine's declared skill sources,
// ~/.config/rsk/sources.json. The file is separate from config.json on purpose:
// an older rsk pinned in some project neither fails on it nor rewrites it, so a
// declared source can never be dropped silently. It holds no credential; git
// authenticates a private source with the user's own SSH key or credential
// helper.
package sources

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
	"go.ralva.dev/ralvaskills/v3/internal/fsx"
	"go.ralva.dev/ralvaskills/v3/internal/git"
	"go.ralva.dev/ralvaskills/v3/internal/schema"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
)

const (
	// FileName is the name of the declarations file inside the config folder.
	FileName = "sources.json"

	githubHost  = "github.com"
	privateFile = 0o600
	tempPattern = ".rsk-sources-*.tmp"
)

var slugPattern = regexp.MustCompile(`[^a-z0-9._-]+`)

type (
	// Entry is one declared source.
	Entry struct {
		ID       string      `json:"id"`
		Type     source.Kind `json:"type"`
		URL      string      `json:"url"`
		Rev      string      `json:"rev,omitempty"`
		Roots    []string    `json:"roots,omitempty"`
		Excludes []string    `json:"excludes,omitempty"`
		Private  bool        `json:"private,omitempty"`
	}

	// File is the content of sources.json.
	File struct {
		SchemaVersion schema.Version `json:"version"`
		Sources       []Entry        `json:"sources"`
	}

	// Options are what `rsk source add` accepts besides the target.
	Options struct {
		Private  bool
		ID       string
		Rev      string
		Roots    []string
		Excludes []string
	}
)

// DefaultPath returns where sources.json lives: inside the rsk config folder, so
// --config and RSK_CONFIG_HOME apply to it like to every other rsk file.
func DefaultPath() string {
	return filepath.Join(config.DefaultConfigFolderPath(), FileName)
}

// NewEntry builds a declaration from what the user typed: an identity
// (@owner/repo, a GitHub repository over SSH), a git URL, or a directory. It
// applies every rule a declaration must satisfy, so a hand-edited file and a
// command line go through the same checks.
func NewEntry(target string, o Options) (Entry, error) {
	var e Entry
	var err error
	switch {
	case strings.HasPrefix(target, "@"):
		e, err = fromIdentity(target, o)
	case isPath(target):
		e, err = fromPath(target, o)
	default:
		e, err = fromURL(target, o)
	}
	if err != nil {
		return Entry{}, err
	}
	e.Rev, e.Roots, e.Excludes, e.Private = o.Rev, o.Roots, o.Excludes, o.Private
	if err = Validate(e); err != nil {
		return Entry{}, err
	}
	return e, nil
}

// Spec converts the declaration to a source descriptor.
func (e Entry) Spec() source.Spec {
	return source.Spec{
		ID: e.ID, Type: e.Type, URL: e.URL, Rev: e.Rev, Roots: e.Roots, Excludes: e.Excludes,
		Declared: e.Type == source.KindPath, Private: e.Private,
	}
}

func isPath(target string) bool {
	if strings.HasPrefix(target, "/") || strings.HasPrefix(target, ".") || strings.HasPrefix(target, "~") {
		return true
	}
	_, err := os.Stat(target)
	return err == nil
}

func fromIdentity(target string, o Options) (Entry, error) {
	id, err := skill.ParseSourceID(target)
	if err != nil {
		return Entry{}, err
	}
	if o.ID != "" && skill.CanonicalSourceID(o.ID) != id {
		return Entry{}, fmt.Errorf("--id %s does not match %s", o.ID, target)
	}
	return Entry{
		ID:   id,
		Type: source.KindGit,
		URL:  "git@" + githubHost + ":" + strings.TrimPrefix(target, "@") + ".git",
	}, nil
}

func fromURL(target string, o Options) (Entry, error) {
	if _, err := git.ParseRemoteURL(target); err != nil {
		return Entry{}, err
	}
	derived, onGitHub, err := githubIdentity(target)
	if err != nil {
		return Entry{}, err
	}
	switch {
	case onGitHub && o.ID != "" && skill.CanonicalSourceID(o.ID) != derived:
		return Entry{}, fmt.Errorf("--id %s does not match %s, which is %s", o.ID, target, derived)
	case onGitHub:
		return Entry{ID: derived, Type: source.KindGit, URL: target}, nil
	case o.ID == "":
		return Entry{}, fmt.Errorf(
			"%s is not on GitHub, so its identity cannot be derived: pass --id @owner/repo",
			target,
		)
	}
	id, err := skill.ParseSourceID(o.ID)
	if err != nil {
		return Entry{}, err
	}
	return Entry{ID: id, Type: source.KindGit, URL: target}, nil
}

func fromPath(target string, o Options) (Entry, error) {
	if o.Rev != "" {
		return Entry{}, errors.New("--rev applies to git sources only")
	}
	if strings.HasPrefix(target, "~") {
		return Entry{}, errors.New("expand ~ before passing a path")
	}
	abs, err := filepath.Abs(target)
	if err != nil {
		return Entry{}, fmt.Errorf("resolve %s: %w", target, err)
	}
	info, err := os.Stat(abs)
	if err != nil || !info.IsDir() {
		return Entry{}, fmt.Errorf("%s is not a directory", abs)
	}

	id := o.ID
	if id == "" {
		slug := strings.Trim(slugPattern.ReplaceAllString(strings.ToLower(filepath.Base(abs)), "-"), "-.")
		id = skill.LocalScope + "/" + slug
	}
	id, err = skill.ParseSourceID(id)
	if err != nil {
		return Entry{}, err
	}
	return Entry{ID: id, Type: source.KindPath, URL: abs}, nil
}

// githubIdentity derives @owner/repo from a GitHub URL in any of the three
// accepted spellings. onGitHub is false for every other host.
func githubIdentity(raw string) (id string, onGitHub bool, err error) {
	var host, path string
	if strings.Contains(raw, "://") {
		u, parseErr := url.Parse(raw)
		if parseErr != nil {
			return "", false, fmt.Errorf("parse %s: %w", raw, parseErr)
		}
		host, path = u.Hostname(), u.Path
	} else {
		hostPart, rest, _ := strings.Cut(raw, ":")
		_, h, _ := strings.Cut(hostPart, "@")
		host, path = h, rest
	}
	if !strings.EqualFold(host, githubHost) {
		return "", false, nil
	}
	repoPath := strings.TrimSuffix(strings.Trim(path, "/"), ".git")
	owner, repo, ok := strings.Cut(repoPath, "/")
	if !ok || owner == "" || repo == "" || strings.Contains(repo, "/") {
		return "", true, fmt.Errorf("%s must name one GitHub repository, owner/repo", raw)
	}
	id, err = skill.ParseSourceID("@" + owner + "/" + repo)
	return id, true, err
}

func reserved(id string) bool {
	canonical := skill.CanonicalSourceID(id)
	return canonical == skill.DefaultRegistryID || canonical == skill.OfficialID
}

// Validate reports why e cannot be a declared source. It checks shape only, so
// a path source whose directory is temporarily missing still loads.
func Validate(e Entry) error {
	id, err := skill.ParseSourceID(e.ID)
	if err != nil {
		return err
	}
	if id != e.ID {
		return fmt.Errorf("source %s: the identity must be canonical (%s)", e.ID, id)
	}
	if reserved(id) {
		return fmt.Errorf("source %s is built in and cannot be declared", id)
	}
	if err = validateKind(e, id); err != nil {
		return fmt.Errorf("source %s: %w", e.ID, err)
	}
	if e.Rev != "" {
		if err = git.ValidateRef(e.Rev); err != nil {
			return fmt.Errorf("source %s: %w", e.ID, err)
		}
	}
	for _, root := range e.Roots {
		if err = git.ValidateRepoPath(root); err != nil {
			return fmt.Errorf("source %s: root: %w", e.ID, err)
		}
	}
	for _, pattern := range e.Excludes {
		if pattern == "" || strings.HasPrefix(pattern, "/") || strings.ContainsAny(pattern, "\x00\n") ||
			slices.Contains(strings.Split(pattern, "/"), "..") {
			return fmt.Errorf("source %s: exclude %q must be a relative pattern", e.ID, pattern)
		}
	}
	return nil
}

func validateKind(e Entry, id string) error {
	switch e.Type {
	case source.KindGit:
		if _, err := git.ParseRemoteURL(e.URL); err != nil {
			return err
		}
		if strings.HasPrefix(id, skill.LocalScope+"/") {
			return fmt.Errorf("the %s scope is reserved for path sources", skill.LocalScope)
		}
		derived, onGitHub, err := githubIdentity(e.URL)
		if err != nil {
			return err
		}
		if onGitHub && derived != id {
			return fmt.Errorf("identity %s does not match %s, which is %s", id, e.URL, derived)
		}
		if !onGitHub && !strings.Contains(id, "/") {
			return errors.New("a declared identity is needed")
		}
	case source.KindPath:
		if !filepath.IsAbs(e.URL) {
			return fmt.Errorf("path %q must be absolute", e.URL)
		}
		if e.Rev != "" {
			return errors.New("a path source has no rev")
		}
	default:
		return fmt.Errorf("type %q is not a declarable source", e.Type)
	}
	return nil
}

// Find returns the declaration with the given identity (case-insensitive).
func Find(f File, id string) (Entry, bool) {
	want := skill.CanonicalSourceID(id)
	for _, e := range f.Sources {
		if e.ID == want {
			return e, true
		}
	}
	return Entry{}, false
}

// Add appends e to f. An identity or a URL that is already declared is an
// error: two entries for one repository would let the second shadow the first.
func Add(f File, e Entry) (File, error) {
	for _, have := range f.Sources {
		if have.ID == e.ID {
			return f, fmt.Errorf("source %s is already declared", e.ID)
		}
		if have.URL == e.URL {
			return f, fmt.Errorf("%s is already declared as %s", e.URL, have.ID)
		}
	}
	f.Sources = append(append([]Entry(nil), f.Sources...), e)
	return f, nil
}

// Remove drops the declaration with the given identity and reports whether
// there was one.
func Remove(f File, id string) (File, bool) {
	want := skill.CanonicalSourceID(id)
	kept := make([]Entry, 0, len(f.Sources))
	for _, e := range f.Sources {
		if e.ID != want {
			kept = append(kept, e)
		}
	}
	removed := len(kept) != len(f.Sources)
	f.Sources = kept
	return f, removed
}

// Load reads sources.json. A missing file is an empty declaration list. Every
// entry is validated, so a hand-edited file cannot register what `rsk source
// add` would refuse.
func Load(path string) (File, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return File{SchemaVersion: schema.Sources}, nil
	}
	if err != nil {
		return File{}, fmt.Errorf("read %s: %w", path, err)
	}
	var f File
	if err = json.Unmarshal(data, &f); err != nil {
		return File{}, fmt.Errorf("parse %s: %w", path, err)
	}
	if err = schema.Check(f.SchemaVersion, schema.Sources); err != nil {
		return File{}, fmt.Errorf("%s: %w", path, err)
	}
	if err = validateAll(f); err != nil {
		return File{}, fmt.Errorf("%s: %w", path, err)
	}
	return f, nil
}

func validateAll(f File) error {
	var checked File
	for _, e := range f.Sources {
		if err := Validate(e); err != nil {
			return err
		}
		var err error
		if checked, err = Add(checked, e); err != nil {
			return err
		}
	}
	return nil
}

// Save writes sources.json atomically, readable by the user only: the names of
// private sources are not for other accounts on the machine.
func Save(path string, f File) error {
	if err := validateAll(f); err != nil {
		return err
	}
	f.SchemaVersion = schema.Sources
	if f.Sources == nil {
		f.Sources = []Entry{}
	}
	if err := os.MkdirAll(filepath.Dir(path), fsperm.Dir); err != nil {
		return fmt.Errorf("create %s: %w", filepath.Dir(path), err)
	}
	return fsx.WriteAtomicPerm(path, tempPattern, privateFile, func(w io.Writer) error {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(f)
	})
}
