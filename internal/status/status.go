// Package status reads what is linked into the skill directories a project and
// the machine use: which skills, at which version, from which source, and in
// which bundles. It only reads; the caller renders what it returns.
package status

import (
	"fmt"
	"os"
	"path/filepath"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
)

// Linked describes one symlink found in a target directory.
type Linked struct {
	Name    string
	Version string
	Source  skill.Source
	Bundles []string
}

// Section groups the linked skills under a single target directory.
type Section struct {
	Title    string
	Subtitle string
	Dir      string
	Skills   []Linked
}

// BuildSections lists the directories a status report covers: the global
// targets unless projectOnly (or just forTool's when set), then the project
// directories unless globalOnly. The sections' skills are empty; Scan fills them.
func BuildSections(
	cfg config.Config, globalOnly, projectOnly bool, forTool tool.ID, projectDirs []string,
) []Section {
	var sections []Section

	if !projectOnly {
		if forTool != "" {
			if dir, ok := cfg.GlobalTargets[string(forTool)]; ok {
				sections = append(sections, Section{
					Title:    "Global — " + forTool.String(),
					Subtitle: dir,
					Dir:      dir,
				})
			}
		} else {
			for toolID, dir := range cfg.GlobalTargets {
				sections = append(sections, Section{
					Title:    "Global — " + toolID,
					Subtitle: dir,
					Dir:      dir,
				})
			}
		}
	}

	if !globalOnly {
		for _, dir := range projectDirs {
			sections = append(sections, Section{
				Title:    "Project",
				Subtitle: dir,
				Dir:      dir,
			})
		}
	}

	return sections
}

// Scan fills in each section's skills in place by scanning its target
// directory. A directory that can't be scanned becomes a warning rather than a
// failure.
func Scan(
	sections []Section, cfg config.Config, membership map[string][]string, includePersonal bool,
) []string {
	var warnings []string
	for i := range sections {
		entries, err := ScanLinked(
			sections[i].Dir,
			cfg.RepoPath, cfg.OfficialCache, cfg.RegistryCache(),
			membership, includePersonal,
		)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("scan %s: %v", sections[i].Dir, err))
		}
		sections[i].Skills = entries
	}
	return warnings
}

// ScanLinked reads the symlinks in targetDir and reports each linked skill's
// version and source. A missing directory is an empty result, not an error.
func ScanLinked(
	targetDir, repoPath, officialCachePath, registryCachePath string,
	membership map[string][]string,
	includePersonal bool,
) ([]Linked, error) {
	entries, err := os.ReadDir(targetDir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var result []Linked
	for _, e := range entries {
		linkPath := filepath.Join(targetDir, e.Name())
		fi, statErr := os.Lstat(linkPath)
		if statErr != nil || !skill.IsLink(fi) {
			continue
		}
		symlinkTarget, readErr := os.Readlink(linkPath)
		if readErr != nil {
			continue
		}
		if !includePersonal && skill.IsPersonalPath(symlinkTarget) {
			continue
		}
		skillVersion, readVersionErr := skill.ReadVersion(symlinkTarget)
		if readVersionErr != nil {
			skillVersion = ""
		}
		src := skill.DetectSource(symlinkTarget, repoPath, officialCachePath, registryCachePath)
		if id, _, _, isGit := skill.ParseGitLink(symlinkTarget, filepath.Dir(officialCachePath)); isGit {
			src = skill.DeclaredSource(id)
			if id == skill.OfficialID {
				src = skill.SourceOfficial
			}
		}
		result = append(result, Linked{
			Name:    e.Name(),
			Version: skillVersion,
			Source:  src,
			Bundles: membership[e.Name()],
		})
	}
	return result, nil
}

// BundleMembership returns a map from skill name to the names of the bundles
// that include it.
func BundleMembership(bundles []config.Bundle) map[string][]string {
	idx := make(map[string][]string)
	for _, b := range bundles {
		for _, ref := range b.Skills {
			idx[ref.Name] = append(idx[ref.Name], b.Name)
		}
	}
	return idx
}
