package source

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal/git"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

// FindVersion returns the skill at an exact version, which a git source
// resolves in one of two ways. A full commit id selects that commit directly,
// for a repository that has no tags. Anything else is a version, resolved
// through the tag <name>@v<version>: a missing tag is an error, and so is a tag
// whose skill declares another version. Both need the network, so neither is
// available offline.
func (g *Git) FindVersion(ctx context.Context, name, version string) (skill.Skill, error) {
	if git.ValidateCommit(version) == nil {
		return g.pin(ctx, name, version, "")
	}
	ver := strings.TrimPrefix(version, "v")
	if !versionPattern.MatchString(ver) {
		return skill.Skill{}, fmt.Errorf("source %s: %q is not a version a tag can name", g.spec.ID, version)
	}
	return g.pin(ctx, name, name+"@v"+ver, ver)
}

// pin fetches ref (a tag or a commit id), finds the skill in the tree it gives,
// and extracts it. wantVersion, when set, must be what the skill declares there.
func (g *Git) pin(ctx context.Context, name, ref, wantVersion string) (skill.Skill, error) {
	if err := git.ValidateRef(ref); err != nil {
		return skill.Skill{}, fmt.Errorf("source %s: %w", g.spec.ID, err)
	}
	if err := g.gate(g.spec.ID, g.spec.URL); err != nil {
		return skill.Skill{}, err
	}
	if g.offline {
		return skill.Skill{}, fmt.Errorf(
			"source %s: resolving %s needs the network; run without --offline", g.spec.ID, ref,
		)
	}

	commit, err := g.repo.Fetch(ctx, g.remote, ref)
	if errors.Is(err, git.ErrRefNotFound) {
		return skill.Skill{}, fmt.Errorf("%w: source %s has no tag or commit %s", ErrNotFound, g.spec.ID, ref)
	}
	if err != nil {
		return skill.Skill{}, fmt.Errorf("source %s: fetch %s: %w", g.spec.ID, ref, err)
	}
	entries, err := g.repo.ListTree(ctx, commit, g.spec.Roots)
	if err != nil {
		return skill.Skill{}, err
	}
	skills, _, err := g.buildSkills(ctx, commit, entries)
	if err != nil {
		return skill.Skill{}, err
	}
	for i := range skills {
		if skills[i].Name != name {
			continue
		}
		if wantVersion != "" && skills[i].Version != "" && skills[i].Version != wantVersion {
			return skill.Skill{}, fmt.Errorf(
				"source %s: tag %s holds %s at version %s", g.spec.ID, ref, name, skills[i].Version,
			)
		}
		dest, extractErr := g.materializeAt(ctx, name, skills[i].Path, commit)
		if extractErr != nil {
			return skill.Skill{}, fmt.Errorf("source %s: extract %q: %w", g.spec.ID, name, extractErr)
		}
		return g.toSkill(&skills[i], dest, commit, ref), nil
	}
	return skill.Skill{}, fmt.Errorf("%w: source %s has no skill %q at %s", ErrNotFound, g.spec.ID, name, ref)
}

// Locked restores a skill at the commit a lock recorded, without consulting the
// branch or the name index. The extraction is checked against the recorded tree
// id, so content that is not what was locked is refused and removed.
func (g *Git) Locked(ctx context.Context, l LockedSkill) (skill.Skill, error) {
	if err := validateLocked(l); err != nil {
		return skill.Skill{}, fmt.Errorf("source %s: %w", g.spec.ID, err)
	}
	if err := g.gate(g.spec.ID, g.spec.URL); err != nil {
		return skill.Skill{}, err
	}
	dest, err := g.materializeAt(ctx, l.Name, l.Path, l.Commit)
	if err != nil {
		return skill.Skill{}, fmt.Errorf("source %s: restore %q at %s: %w", g.spec.ID, l.Name, l.Commit, err)
	}
	tree, err := g.repo.TreeID(ctx, l.Commit, l.Path)
	if err == nil && l.Tree != "" && tree != l.Tree {
		err = fmt.Errorf("%w: %s at %s is %s, the lock says %s", ErrTreeMismatch, l.Path, l.Commit, tree, l.Tree)
	}
	if err != nil {
		_ = os.RemoveAll(dest)
		return skill.Skill{}, err
	}

	fm := readMeta(dest)
	return skill.Skill{
		Name: l.Name, Version: fm.Version, Description: fm.Description, Path: dest,
		Source: skill.DeclaredSource(g.spec.ID), AllowedTools: skill.NormalizeAllowedTools(fm.AllowedTools),
		Git: &skill.GitRef{URL: g.spec.URL, Rev: l.Rev, Commit: l.Commit, Tree: tree, Path: l.Path},
	}, nil
}

func validateLocked(l LockedSkill) error {
	if err := skill.ValidateName(l.Name); err != nil {
		return err
	}
	if err := git.ValidateCommit(l.Commit); err != nil {
		return err
	}
	if err := git.ValidateRepoPath(l.Path); err != nil {
		return err
	}
	if path.Base(l.Path) != l.Name {
		return fmt.Errorf("locked path %s does not end in the skill name %s", l.Path, l.Name)
	}
	return nil
}

func readMeta(dir string) skill.Frontmatter {
	f, err := os.Open(filepath.Join(dir, skillFile))
	if err != nil {
		return skill.Frontmatter{}
	}
	defer func() { _ = f.Close() }()
	return skill.ParseFrontmatter(f)
}

// TreeAt returns the tree id name's directory has at commit, using the path the
// current index gives it. It is how an install with no lock entry (a global one)
// finds out whether the skill changed since the commit it was extracted at.
func (g *Git) TreeAt(ctx context.Context, name, commit string) (string, error) {
	idx, err := g.index(ctx)
	if err != nil {
		return "", err
	}
	for i := range idx.Skills {
		if idx.Skills[i].Name == name {
			return g.repo.TreeID(ctx, commit, idx.Skills[i].Path)
		}
	}
	return "", fmt.Errorf("%w: source %s has no skill %q", ErrNotFound, g.spec.ID, name)
}
