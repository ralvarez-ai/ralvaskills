package main

import (
	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/mcp"
)

// offlineAliases returns the retired-name table `rsk status` marks against. It
// never touches the network: the local ralvaskills clone's alias file wins, the
// cached registry index of any age is the fallback, and with neither it returns
// nil so status stays silent.
func offlineAliases(cfg config.Config) map[string]string {
	if cfg.RepoPath != "" {
		if aliases, err := localAliases(cfg); err == nil && len(aliases) > 0 {
			return aliases
		}
	}
	return mcp.CachedAliases(cfg.RegistryCache())
}
