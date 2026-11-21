// Package trust records which declared sources the user has allowed rsk to
// contact. A grant is made only by `rsk source trust`, which needs a terminal,
// so an agent that can run shell commands cannot grant one. A grant is keyed by
// identity and URL: changing the URL of a trusted identity revokes it.
package trust

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
	"go.ralva.dev/ralvaskills/v3/internal/fsx"
	"go.ralva.dev/ralvaskills/v3/internal/schema"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

const (
	// FileName is the name of the grants file inside the config folder.
	FileName = "trust.json"

	privateFile = 0o600
	tempPattern = ".rsk-trust-*.tmp"
)

// ErrUntrusted is returned, wrapped in an UntrustedError, when a source has no
// grant.
var ErrUntrusted = errors.New("source is not trusted")

type (
	// GrantRecord is one trust decision.
	GrantRecord struct {
		ID        string    `json:"id"`
		URL       string    `json:"url"`
		GrantedAt time.Time `json:"granted_at"`
	}

	// File is the content of trust.json.
	File struct {
		SchemaVersion schema.Version `json:"version"`
		Grants        []GrantRecord  `json:"grants"`
	}

	// UntrustedError says which source needs a grant and the exact command that
	// makes one, so a refusal is always actionable.
	UntrustedError struct {
		ID  string
		URL string
	}
)

// Error implements error.
func (e *UntrustedError) Error() string {
	return fmt.Sprintf(
		"source %s (%s) is not trusted: run this yourself in a terminal to trust it: rsk source trust %s",
		e.ID, e.URL, e.ID,
	)
}

// Unwrap returns ErrUntrusted.
func (*UntrustedError) Unwrap() error {
	return ErrUntrusted
}

// DefaultPath returns where trust.json lives, inside the rsk config folder.
func DefaultPath() string {
	return filepath.Join(config.DefaultConfigFolderPath(), FileName)
}

// Granted reports whether id is trusted at exactly this URL.
func (f File) Granted(id, url string) bool {
	want := skill.CanonicalSourceID(id)
	for _, g := range f.Grants {
		if g.ID == want && g.URL == url {
			return true
		}
	}
	return false
}

// Grant records trust in id at url, replacing any earlier grant for id: an
// identity is trusted at one URL at a time.
func Grant(f File, id, url string, at time.Time) File {
	canonical := skill.CanonicalSourceID(id)
	grants := make([]GrantRecord, 0, len(f.Grants)+1)
	for _, g := range f.Grants {
		if g.ID != canonical {
			grants = append(grants, g)
		}
	}
	grants = append(grants, GrantRecord{ID: canonical, URL: url, GrantedAt: at.UTC()})
	f.Grants = grants
	return f
}

// Revoke removes the grant for id and reports whether there was one.
func Revoke(f File, id string) (File, bool) {
	canonical := skill.CanonicalSourceID(id)
	kept := make([]GrantRecord, 0, len(f.Grants))
	for _, g := range f.Grants {
		if g.ID != canonical {
			kept = append(kept, g)
		}
	}
	removed := len(kept) != len(f.Grants)
	f.Grants = kept
	return f, removed
}

// Require returns nil when id is trusted at url, and an UntrustedError that
// names the command to run otherwise. Every code path that contacts a git
// source calls it first.
func Require(f File, id, url string) error {
	if f.Granted(id, url) {
		return nil
	}
	return &UntrustedError{ID: skill.CanonicalSourceID(id), URL: url}
}

// Load reads trust.json. A missing file grants nothing.
func Load(path string) (File, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return File{SchemaVersion: schema.Trust}, nil
	}
	if err != nil {
		return File{}, fmt.Errorf("read %s: %w", path, err)
	}
	var f File
	if err = json.Unmarshal(data, &f); err != nil {
		return File{}, fmt.Errorf("parse %s: %w", path, err)
	}
	if err = schema.Check(f.SchemaVersion, schema.Trust); err != nil {
		return File{}, fmt.Errorf("%s: %w", path, err)
	}
	return f, nil
}

// Save writes trust.json atomically, readable by the user only.
func Save(path string, f File) error {
	f.SchemaVersion = schema.Trust
	if f.Grants == nil {
		f.Grants = []GrantRecord{}
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
