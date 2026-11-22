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

	// Local names the skills that live in rsk.local.mod, the layer for skills
	// from a private source. A skill is in exactly one layer: the reader marks
	// the ones it found in the local file, and the writer sends them back there,
	// so nothing private reaches the committed rsk.mod.
	Local map[string]bool `toml:"-"`
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

// ReadMod decodes rsk.mod from rskDir and merges rsk.local.mod, when there is
// one, marking the skills it holds as local.
func ReadMod(rskDir string) (Mod, error) {
	m := Mod{
		Pinned: []string{}, Skills: map[string]string{}, Origin: map[string]string{}, Local: map[string]bool{},
	}
	f, err := decodeMod(ModPath(rskDir))
	if err != nil {
		return Mod{}, err
	}
	m.SchemaVersion, m.Tools = f.SchemaVersion, f.Tools
	if err = m.merge(&f, false); err != nil {
		return Mod{}, err
	}
	if _, statErr := os.Stat(LocalModPath(rskDir)); statErr == nil {
		local, localErr := decodeMod(LocalModPath(rskDir))
		if localErr != nil {
			return Mod{}, localErr
		}
		if err = m.merge(&local, true); err != nil {
			return Mod{}, err
		}
	}
	m.normalize()
	return m, nil
}

func decodeMod(path string) (modFile, error) {
	var f modFile
	if _, err := toml.DecodeFile(path, &f); err != nil {
		return modFile{}, fmt.Errorf("read %s: %w", path, err)
	}
	if err := schema.Check(f.SchemaVersion, schema.Mod); err != nil {
		return modFile{}, fmt.Errorf("%s: %w", filepath.Base(path), err)
	}
	return f, nil
}

// merge adds a layer's skills and pins to m, marking them local when asked. The
// same name in two layers is an error: a skill is unique by name in one scope.
func (m *Mod) merge(f *modFile, local bool) error {
	for key, constraint := range f.Skills {
		ref, err := skill.ParseKey(key)
		if err != nil {
			return fmt.Errorf("rsk.mod: %w", err)
		}
		if _, dup := m.Skills[ref.Name]; dup {
			return fmt.Errorf("rsk.mod: two skills named %q in one scope", ref.Name)
		}
		m.Skills[ref.Name] = constraint
		if ref.Source != "" && ref.Source != skill.DefaultRegistryID {
			m.Origin[ref.Name] = ref.Source
		}
		if local {
			m.Local[ref.Name] = true
		}
	}
	for _, key := range f.Pinned {
		ref, err := skill.ParseKey(key)
		if err != nil {
			return fmt.Errorf("rsk.mod: %w", err)
		}
		m.Pinned = append(m.Pinned, ref.Name)
	}
	return nil
}

// WriteMod encodes m as TOML to rskDir/rsk.mod atomically, creating parent
// directories as needed. It writes version 2 only when a skill in m has an
// origin to qualify; otherwise the file stays version 1. The skills marked
// local go to rsk.local.mod instead, which is kept out of git and removed when
// no skill is left in it.
func WriteMod(rskDir string, m Mod) error {
	m.normalize()
	if err := writeModFile(ModPath(rskDir), &m, func(name string) bool { return !m.Local[name] }, true); err != nil {
		return err
	}
	var local []string
	for name := range m.Skills {
		if m.Local[name] {
			local = append(local, name)
		}
	}
	if len(local) == 0 {
		return removeIfPresent(LocalModPath(rskDir))
	}
	if err := ExcludeLocalFiles(filepath.Dir(rskDir)); err != nil {
		return err
	}
	return writeModFile(LocalModPath(rskDir), &m, func(name string) bool { return m.Local[name] }, false)
}

// writeModFile writes the skills and pins of m that keep selects. Only the
// committed file carries the tools list.
func writeModFile(path string, m *Mod, keep func(name string) bool, withTools bool) error {
	f := modFile{
		SchemaVersion: schema.Base,
		Pinned:        []string{},
		Skills:        map[string]string{},
	}
	if withTools {
		f.Tools = m.Tools
	}
	key := func(name string) string {
		if origin := m.Origin[name]; origin != "" {
			f.SchemaVersion = schema.Mod
			return skill.Ref{Source: origin, Name: name}.Key()
		}
		return name
	}
	for name, constraint := range m.Skills {
		if keep(name) {
			f.Skills[key(name)] = constraint
		}
	}
	for _, name := range m.Pinned {
		if keep(name) {
			f.Pinned = append(f.Pinned, key(name))
		}
	}

	if err := os.MkdirAll(filepath.Dir(path), fsperm.Dir); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}
	return fsx.WriteAtomic(path, TempFilePattern, func(w io.Writer) error {
		return toml.NewEncoder(w).Encode(f)
	})
}

// removeIfPresent deletes path and ignores a file that is already gone.
func removeIfPresent(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove %s: %w", path, err)
	}
	return nil
}
