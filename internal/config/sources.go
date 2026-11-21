package config

import (
	"path/filepath"

	"go.ralva.dev/ralvaskills/v3/internal/source"
)

// RegistrySpec describes the hosted registry in use.
func (c Config) RegistrySpec() source.Spec {
	return source.Spec{ID: source.DefaultRegistryID, Type: source.KindRegistry, URL: c.EffectiveRegistryURL()}
}

// LocalSpec describes the local repo clone, the source in local mode.
func (c Config) LocalSpec() source.Spec {
	name := filepath.Base(c.RepoPath)
	if c.RepoPath == "" || name == "." || name == string(filepath.Separator) {
		name = "ralvaskills"
	}
	return source.Spec{ID: source.LocalScope + "/" + name, Type: source.KindPath, URL: c.RepoPath}
}

// OfficialSpec describes the anthropics/skills cache.
func (c Config) OfficialSpec() source.Spec {
	return source.Spec{ID: source.OfficialID, Type: source.KindOfficial, URL: source.OfficialSkillsURL}
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
	return source.Options{RegistryCache: c.RegistryCache(), OfficialCache: c.OfficialCache}
}

// Sources builds the chain for this machine's configuration. Every command
// that needs a skill source starts here.
func (c Config) Sources() (source.Chain, error) {
	return source.Build(c.SourceSpecs(), c.SourceOptions())
}
