// Package manifest handles rsk.mod and rsk.lock read/write for project-local installs.
package manifest

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
	"go.ralva.dev/ralvaskills/v3/internal/fsx"
	"go.ralva.dev/ralvaskills/v3/internal/schema"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
)

// Mod is the in-memory representation of rsk.mod. Skill names are bare here
// whatever the file holds: a skill is unique by name within one scope, and Origin
// records which source a name came from when it is not the default registry.
type Mod struct {
	// SchemaVersion identifies the file format version. Version 0 is accepted
	// as legacy (pre-versioning) and treated as version 1.
	SchemaVersion schema.Version    `toml:"version,omitempty"`
	Tools         []tool.ID         `toml:"tools"`
	Pinned        []string          `toml:"pinned"`
	Skills        map[string]string `toml:"skills"`

	// Origin maps a bare skill name to the identity of its source, for skills
	// that do not come from the default registry. It is what makes a file
	// qualify those keys on write; it is never written itself.
	Origin map[string]string `toml:"-"`
}

// modFile is the on-disk shape of rsk.mod: keys are the qualified form for a
// skill from any source other than the default registry.
type modFile struct {
	SchemaVersion schema.Version    `toml:"version,omitempty"`
	Tools         []tool.ID         `toml:"tools"`
	Pinned        []string          `toml:"pinned"`
	Skills        map[string]string `toml:"skills"`
}

// SetOrigin records which kind of source name came from, so the next write
// qualifies its key. A skill from the default registry, or a local clone whose
// identity is machine-specific, has no origin and keeps a bare key.
func (m *Mod) SetOrigin(name string, src skill.Source) {
	m.normalize()
	if id := src.Identity(); id != "" && id != skill.DefaultRegistryID {
		m.Origin[name] = id
		return
	}
	delete(m.Origin, name)
}

// normalize fills nil/empty fields with safe zero-value defaults.
func (m *Mod) normalize() {
	if len(m.Tools) == 0 {
		m.Tools = []tool.ID{tool.ClaudeID}
	}
	if m.Skills == nil {
		m.Skills = make(map[string]string)
	}
	if m.Pinned == nil {
		m.Pinned = []string{}
	}
	if m.Origin == nil {
		m.Origin = make(map[string]string)
	}
}

// ReadMod decodes rsk.mod from rskDir.
func ReadMod(rskDir string) (Mod, error) {
	path := ModPath(rskDir)
	var f modFile
	if _, err := toml.DecodeFile(path, &f); err != nil {
		return Mod{}, fmt.Errorf("read %s: %w", path, err)
	}
	if err := schema.Check(f.SchemaVersion, schema.Mod); err != nil {
		return Mod{}, fmt.Errorf("rsk.mod: %w", err)
	}
	m := Mod{
		SchemaVersion: f.SchemaVersion,
		Tools:         f.Tools,
		Pinned:        make([]string, 0, len(f.Pinned)),
		Skills:        make(map[string]string, len(f.Skills)),
		Origin:        make(map[string]string),
	}
	for key, constraint := range f.Skills {
		ref, err := skill.ParseKey(key)
		if err != nil {
			return Mod{}, fmt.Errorf("rsk.mod: %w", err)
		}
		if _, dup := m.Skills[ref.Name]; dup {
			return Mod{}, fmt.Errorf("rsk.mod: two skills named %q in one scope", ref.Name)
		}
		m.Skills[ref.Name] = constraint
		if ref.Source != "" && ref.Source != skill.DefaultRegistryID {
			m.Origin[ref.Name] = ref.Source
		}
	}
	for _, key := range f.Pinned {
		ref, err := skill.ParseKey(key)
		if err != nil {
			return Mod{}, fmt.Errorf("rsk.mod: %w", err)
		}
		m.Pinned = append(m.Pinned, ref.Name)
	}
	m.normalize()
	return m, nil
}

// WriteMod encodes m as TOML to rskDir/rsk.mod atomically, creating parent
// directories as needed. It writes version 2 only when a skill in m has an
// origin to qualify; otherwise the file stays version 1.
func WriteMod(rskDir string, m Mod) error {
	m.normalize()
	f := modFile{
		SchemaVersion: schema.Base,
		Tools:         m.Tools,
		Pinned:        make([]string, len(m.Pinned)),
		Skills:        make(map[string]string, len(m.Skills)),
	}
	key := func(name string) string {
		if origin := m.Origin[name]; origin != "" {
			f.SchemaVersion = schema.Mod
			return skill.Ref{Source: origin, Name: name}.Key()
		}
		return name
	}
	for name, constraint := range m.Skills {
		f.Skills[key(name)] = constraint
	}
	for i, name := range m.Pinned {
		f.Pinned[i] = key(name)
	}

	path := ModPath(rskDir)
	if err := os.MkdirAll(filepath.Dir(path), fsperm.Dir); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}
	return fsx.WriteAtomic(path, TempFilePattern, func(w io.Writer) error {
		return toml.NewEncoder(w).Encode(f)
	})
}
