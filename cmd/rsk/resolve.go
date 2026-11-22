package main

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/install"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
	"go.ralva.dev/ralvaskills/v3/internal/sources"
	"go.ralva.dev/ralvaskills/v3/internal/ui"
)

// flagOffline is the root flag that keeps git sources from contacting their
// remotes.
const flagOffline = "offline"

// sourceRepoFor lets tests supply the clone a git source reads, so none of them
// reaches a network. It is nil outside tests.
var sourceRepoFor func(source.Spec) source.RepoOps

// sourceChain builds the chain a command resolves skill names through: the
// built-in sources, then the declared ones, then any source offered for this run
// only (see sources.ChainOptions.Implicit). Commands that only restore
// or relink what a project already records use cfg.Sources() instead, so they
// never depend on a declared source.
func sourceChain(cmd *cobra.Command, cfg config.Config, implicit ...sources.Entry) (source.Chain, error) {
	offline, err := offlineRequested(cmd)
	if err != nil {
		return nil, err
	}
	chain, err := sources.Chain(cfg, sources.ChainOptions{Offline: offline, RepoFor: sourceRepoFor, Implicit: implicit})
	if err != nil {
		return nil, err
	}
	if offline && len(chain) > len(cfg.SourceSpecs()) {
		ui.Info(cmd.ErrOrStderr(), "offline: answering from cached source indexes")
	}
	return chain, nil
}

// offlineRequested reports whether --offline was given. A command without the
// flag is online.
func offlineRequested(cmd *cobra.Command) (bool, error) {
	if cmd.Flags().Lookup(flagOffline) == nil {
		return false, nil
	}
	return cmd.Flags().GetBool(flagOffline)
}

// isLinkedAnywhere reports whether a skill named name is symlinked into any of
// the given target directories.
func isLinkedAnywhere(name string, targets []string) bool {
	for _, t := range targets {
		if skill.IsLinked(name, t) {
			return true
		}
	}
	return false
}

// skillNamesFromArgs resolves bundle-or-skill args to a deduplicated list of
// skill names by reading the catalog. A name matching a bundle expands to that
// bundle's skill names; otherwise the name is treated as a skill name.
func skillNamesFromArgs(args []string, catalog []config.Bundle) ([]string, error) {
	seen := make(map[string]bool)
	var names []string
	for _, raw := range args {
		bareName, _, err := parseNameVersion(raw)
		if err != nil {
			return nil, err
		}
		if b, ok := config.FindBundle(catalog, bareName); ok {
			for _, ref := range b.Skills {
				if !seen[ref.Name] {
					seen[ref.Name] = true
					names = append(names, ref.Name)
				}
			}
			continue
		}
		if !seen[bareName] {
			seen[bareName] = true
			names = append(names, bareName)
		}
	}
	return names, nil
}

// allTargetDirs returns all configured target dirs (all global targets + project-local).
func allTargetDirs(cfg config.Config) []string {
	dirs := make([]string, 0, len(cfg.GlobalTargets)+1)
	for _, dir := range cfg.GlobalTargets {
		dirs = append(dirs, dir)
	}
	if rskDir, err := manifest.ProjectFolderPath(); err == nil {
		if m, modErr := manifest.ReadMod(rskDir); modErr == nil {
			dirs = append(dirs, install.ProjectSkillsDirs(filepath.Dir(rskDir), m)...)
		}
	}
	return dirs
}

func joinKeys(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

// parseNameVersion splits a "name" or "name@version" argument.
func parseNameVersion(arg string) (name, version string, err error) {
	return install.ParseNameVersion(arg)
}

// syncPinnedAllTools writes pinned skill entries for every tool in rsk.mod.
func syncPinnedAllTools(rskDir string, m manifest.Mod) error {
	return install.SyncPinnedAllTools(rskDir, m)
}
