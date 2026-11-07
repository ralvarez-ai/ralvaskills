// Package config handles loading and saving the rsk configuration file.
package config

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"

	"github.com/ralvarezdev/ralvaskills/v3/internal"
)

const (
	// EnvConfigHome overrides the entire rsk config directory
	// (normally ~/.config/rsk) — config.json, catalog.toml, and all caches
	// live under it instead. Intended for sandboxing rsk during manual or
	// automated testing without touching the real machine config.
	EnvConfigHome = "RSK_CONFIG_HOME"

	// EnvRegistryURL overrides the registry base URL for this process, ahead of
	// registry_url in config.json and DefaultRegistryURL. It is never written back
	// to config.json, so a different host needs neither a release nor an edit.
	EnvRegistryURL = "RSK_REGISTRY_URL"

	// ConfigFolderName is the name of the folder under ~/.config/ where rsk
	// stores its configuration.
	ConfigFolderName = "rsk"

	// ConfigFileName is the filename of the main rsk configuration file.
	ConfigFileName = "config.json"

	// CatalogFileName is the filename of the user-overridable skill catalog.
	CatalogFileName = "catalog.toml"

	// RegistryFolderName is the subdirectory under the config folder used to
	// cache tarballs and the registry index.
	RegistryFolderName = "registry"

	// DefaultRegistryURL is the base URL of the hosted skill registry used when
	// no local repo clone is configured and neither RSK_REGISTRY_URL nor
	// registry_url is set. Released binaries carry it, so the host it names is
	// never retired.
	DefaultRegistryURL = "https://skills.ralva.dev"
)

// configHomeOverride is set by the --config flag; it wins over the environment.
var configHomeOverride string

// SetConfigHome overrides the config folder for this process (the --config
// flag). An empty dir clears the override.
func SetConfigHome(dir string) {
	configHomeOverride = dir
}

// DefaultConfigFolderPath returns the canonical path of the rsk config folder
// (~/.config/rsk), the --config override, or the RSK_CONFIG_HOME environment
// override, in that order.
func DefaultConfigFolderPath() string {
	if configHomeOverride != "" {
		return configHomeOverride
	}
	if dir := os.Getenv(EnvConfigHome); dir != "" {
		return dir
	}
	return filepath.Join(internal.ConfigHome(), ConfigFolderName)
}

// DefaultConfigFilePath returns the canonical path of the rsk config file
// (~/.config/rsk/config.json).
func DefaultConfigFilePath() string {
	return filepath.Join(DefaultConfigFolderPath(), ConfigFileName)
}

// DefaultCatalogPath returns the canonical path of the user skill catalog
// (~/.config/rsk/catalog.toml).
func DefaultCatalogPath() string {
	return filepath.Join(DefaultConfigFolderPath(), CatalogFileName)
}

// DefaultRegistryCachePath returns the directory used to cache skill downloads
// from the hosted registry. The cache is placed next to officialCache so both
// caches share the same parent directory.
//
// officialCache is required by config validation. The fallback branch is only
// reachable when called on an unvalidated Config.
func DefaultRegistryCachePath(officialCache string) string {
	if officialCache != "" {
		return filepath.Join(filepath.Dir(filepath.Clean(officialCache)), RegistryFolderName)
	}
	return filepath.Join(xdg.CacheHome, ConfigFolderName, RegistryFolderName)
}
