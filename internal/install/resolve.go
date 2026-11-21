// Package install resolves and installs skills into a project or global scope.
// It is the single implementation shared by `rsk install` and the MCP server,
// so the CLI and the agent cannot diverge on what gets written.
package install

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal/cmdx"
	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/digest"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
)

// FindSkill returns the skill named name from the chain. A bare name is looked
// up in every source; a name two sources provide is an error, never resolved by
// order, and a source that cannot be consulted fails the lookup.
func FindSkill(ctx context.Context, name string, chain source.Chain) (skill.Skill, error) {
	return FindRef(ctx, skill.Ref{Name: name}, chain)
}

// FindRef resolves ref: a qualified reference in exactly the source it names, a
// bare name across the whole chain.
func FindRef(ctx context.Context, ref skill.Ref, chain source.Chain) (skill.Skill, error) {
	s, err := chain.Find(ctx, ref)
	if err == nil {
		return s, nil
	}
	if errors.Is(err, source.ErrNotFound) {
		if ref.Source != "" {
			return skill.Skill{}, fmt.Errorf("skill %q not found in %s", ref.Name, ref.Source)
		}
		return skill.Skill{}, fmt.Errorf(
			"skill %q not found in local repo or official cache\n  Run 'rsk catalog' to browse available skills",
			ref.Name,
		)
	}
	return skill.Skill{}, err
}

// Resolve expands bundle-or-skill names into concrete, deduplicated skills,
// along with soft warnings for bundle members that could not be resolved.
func Resolve(
	ctx context.Context,
	names []string,
	catalog []config.Bundle,
	chain source.Chain,
) ([]skill.Skill, []string, error) {
	var skills []skill.Skill
	var warnings []string

	for _, raw := range names {
		ref, version, err := skill.ParseRef(raw)
		if err != nil {
			return nil, nil, err
		}
		if bundle, ok := config.FindBundle(catalog, ref.Name); ok && ref.Source == "" {
			bundleSkills, bundleWarnings, resolveErr := resolveBundleSkills(
				ctx, bundle, chain,
			)
			if resolveErr != nil {
				return nil, nil, resolveErr
			}
			skills = append(skills, bundleSkills...)
			warnings = append(warnings, bundleWarnings...)
			continue
		}
		s, findErr := FindRef(ctx, ref, chain)
		if findErr != nil {
			return nil, nil, findErr
		}
		// A git source resolves an exact version through a tag. The other
		// sources keep ignoring it at resolution time and only record it as a
		// constraint, as they always have.
		if version != "" && s.Source.IsDeclared() {
			if s, findErr = chain.FindVersion(ctx, string(s.Source), s.Name, version); findErr != nil {
				return nil, nil, findErr
			}
		}
		skills = append(skills, s)
	}

	return DedupSkills(skills), warnings, nil
}

func resolveBundleSkills(
	ctx context.Context, bundle config.Bundle, chain source.Chain,
) ([]skill.Skill, []string, error) {
	var skills []skill.Skill
	var warnings []string

	for _, ref := range bundle.Skills {
		src, notFoundMsg := sourceForRef(ref.Source, chain)
		s, err := src.Find(ctx, skill.Ref{Name: ref.Name})
		if err == nil {
			skills = append(skills, s)
			continue
		}
		if !errors.Is(err, source.ErrNotFound) {
			return nil, nil, fmt.Errorf("resolve %s skill %q: %w", ref.Source, ref.Name, err)
		}
		warnings = append(warnings, fmt.Sprintf(notFoundMsg, ref.Name))
	}
	return skills, warnings, nil
}

// sourceForRef narrows the chain to the sources a bundle member names: the
// official cache for an official member, the configured local source (a clone or
// the hosted registry) otherwise.
func sourceForRef(src skill.Source, chain source.Chain) (sub source.Chain, notFoundMsg string) {
	if src == skill.SourceOfficial {
		return chain.Select(source.KindOfficial),
			"%s (official) not in cache — run 'rsk update --include official' to fetch it"
	}
	return chain.Select(source.KindRegistry, source.KindPath), "%s is not yet available (planned) — skipped"
}

// DedupSkills returns skills with duplicates removed, keeping first occurrences.
func DedupSkills(skills []skill.Skill) []skill.Skill {
	seen := make(map[string]bool, len(skills))
	result := make([]skill.Skill, 0, len(skills))
	for _, s := range skills {
		if !seen[s.Name] {
			seen[s.Name] = true
			result = append(result, s)
		}
	}
	return result
}

// ParseNameVersion splits "name[@version]" or "@owner/repo/name[@version]" into
// the skill's bare name and the version. The source of a qualified reference is
// dropped: use skill.ParseRef when it matters.
func ParseNameVersion(arg string) (name, version string, err error) {
	ref, version, err := skill.ParseRef(arg)
	if err != nil {
		return "", "", err
	}
	return ref.Name, version, nil
}

// ConstraintsFromArgs returns name → version constraints from "name@version"
// and "@owner/repo/name@version" args.
func ConstraintsFromArgs(args []string) (map[string]string, error) {
	constraints := make(map[string]string, len(args))
	for _, raw := range args {
		name, version, err := ParseNameVersion(raw)
		if err != nil {
			return nil, err
		}
		if version != "" {
			constraints[name] = version
		}
	}
	return constraints, nil
}

// Targets returns the skill directories an install should act on. Without
// global it returns the project's per-tool skills dirs; with global it returns
// the configured global dirs.
func Targets(cfg config.Config, global bool, forTool tool.ID) ([]string, error) {
	if !global {
		rskDir, err := manifest.ProjectFolderPath()
		if err != nil {
			return nil, err
		}
		m, err := manifest.ReadMod(rskDir)
		if err != nil {
			return nil, err
		}
		return ProjectSkillsDirs(filepath.Dir(rskDir), m), nil
	}

	if forTool != "" {
		dir, ok := cfg.GlobalTargets[string(forTool)]
		if !ok {
			return nil, fmt.Errorf(
				"--for: tool %q is not configured — configured tools: %s",
				forTool, joinKeys(cfg.GlobalTargets),
			)
		}
		return []string{dir}, nil
	}

	scope := cfg.DefaultTargetScope
	if scope == cmdx.ForAll {
		dirs := make([]string, 0, len(cfg.GlobalTargets))
		for _, dir := range cfg.GlobalTargets {
			dirs = append(dirs, dir)
		}
		return dirs, nil
	}

	dir, ok := cfg.GlobalTargets[scope]
	if !ok {
		return nil, fmt.Errorf(
			"default_target_scope %q does not match any configured tool — run 'rsk init' to fix",
			scope,
		)
	}
	return []string{dir}, nil
}

// ProjectSkillsDirs returns the deduplicated project-local skills dirs for the
// tools listed in m.
func ProjectSkillsDirs(projectRoot string, m manifest.Mod) []string {
	seen := make(map[string]bool, len(m.Tools))
	dirs := make([]string, 0, len(m.Tools))
	for _, id := range m.Tools {
		t, ok := tool.Get(id)
		if !ok {
			continue
		}
		dir := t.ProjectSkillsDir(projectRoot)
		if !seen[dir] {
			seen[dir] = true
			dirs = append(dirs, dir)
		}
	}
	return dirs
}

// SyncPinnedAllTools rewrites each configured tool's project config from m's
// pinned list.
func SyncPinnedAllTools(rskDir string, m manifest.Mod) error {
	projectDir := filepath.Dir(rskDir)
	if _, err := digest.Sync(projectDir, m.Pinned); err != nil {
		return err
	}
	for _, id := range m.Tools {
		t, ok := tool.Get(id)
		if !ok {
			return fmt.Errorf("unknown tool %q in rsk.mod", id)
		}
		if err := t.SyncPinned(projectDir, m.Pinned); err != nil {
			return err
		}
	}
	return nil
}

func joinKeys(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}
