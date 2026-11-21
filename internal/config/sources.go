package config

import (
	"fmt"
	"path/filepath"

	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
)

// RegistrySpec describes the hosted registry in use.
func (c Config) RegistrySpec() source.Spec {
	return source.Spec{ID: skill.DefaultRegistryID, Type: source.KindRegistry, URL: c.EffectiveRegistryURL()}
}

// LocalSpec describes the local repo clone, the source in local mode.
func (c Config) LocalSpec() source.Spec {
	name := filepath.Base(c.RepoPath)
	if c.RepoPath == "" || name == "." || name == string(filepath.Separator) {
		name = "ralvaskills"
	}
	return source.Spec{ID: skill.LocalScope + "/" + name, Type: source.KindPath, URL: c.RepoPath}
}

// BuiltinTrusted reports whether id at url is the one git source rsk ships, the
// official anthropics/skills repository at its own URL, which needs no grant. A
// different URL under the same identity is not trusted.
func BuiltinTrusted(id, url string) bool {
	return id == skill.OfficialID && url == source.OfficialSkillsURL
}

// builtinGate is the trust gate for a chain that has no declared sources: it
// admits the built-in official source and nothing else.
func builtinGate(id, url string) error {
	if BuiltinTrusted(id, url) {
		return nil
	}
	return fmt.Errorf("source %s is not a built-in source", id)
}

// OfficialSpec describes the official anthropics/skills source, read through the
// git transport like any other git source.
func (c Config) OfficialSpec() source.Spec {
	return source.Spec{ID: skill.OfficialID, Type: source.KindGit, URL: source.OfficialSkillsURL}
}

// SourceSpecs returns the sources in declaration order: the local clone in
// local mode or the hosted registry otherwise, then the official cache.
func (c Config) SourceSpecs() []source.Spec {
	primary := c.RegistrySpec()
	if c.LocalMode() {
		primary = c.LocalSpec()
	}
	return []source.Spec{primary, c.OfficialSpec()}
}

// SourceOptions returns the cache locations the sources are built with.
func (c Config) SourceOptions() source.Options {
	return source.Options{
		RegistryCache: c.RegistryCache(),
		CacheRoot:     filepath.Dir(c.OfficialCache),
		Gate:          builtinGate,
	}
}

// Sources builds the chain for this machine's configuration. Every command
// that needs a skill source starts here.
func (c Config) Sources() (source.Chain, error) {
	return source.Build(c.SourceSpecs(), c.SourceOptions())
}
