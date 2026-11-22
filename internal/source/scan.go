package source

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

type (
	// ScanResult is what a dry-run discovery found in a source: the skills that
	// would resolve, the directories left out and why, and the names that more
	// than one directory claims. A scan reports a clash where resolving the
	// source would refuse it.
	ScanResult struct {
		Skills  []ScannedSkill `json:"skills"`
		Skipped []Skipped      `json:"skipped,omitempty"`
		Clashes []Clash        `json:"clashes,omitempty"`
	}

	// ScannedSkill is one skill a scan found, with the problems in its SKILL.md.
	ScannedSkill struct {
		Name        string   `json:"name"`
		Path        string   `json:"path"`
		Version     string   `json:"version,omitempty"`
		Description string   `json:"description,omitempty"`
		Issues      []string `json:"issues,omitempty"`
	}

	// Clash is a skill name that more than one directory provides. None of them
	// resolves until one is excluded or renamed.
	Clash struct {
		Name  string   `json:"name"`
		Paths []string `json:"paths"`
	}
)

// issuesOf lists what is wrong with a skill's frontmatter, for a scan to show:
// the name must be its directory's, and a skill should carry a version and a
// description. None of these stops a skill resolving; they are what a repository
// owner fixes before others install from it.
func issuesOf(dirName string, fm *skill.Frontmatter) []string {
	var issues []string
	if fm.Name != "" && fm.Name != dirName {
		issues = append(issues, fmt.Sprintf("frontmatter name %q differs from the directory name %q", fm.Name, dirName))
	}
	if fm.Version == "" {
		issues = append(issues, "no version in the frontmatter")
	}
	if fm.Description == "" {
		issues = append(issues, "no description in the frontmatter")
	}
	return issues
}

// clashesOf returns the names with more than one directory, in name order.
func clashesOf(paths map[string][]string) []Clash {
	var out []Clash
	for name, dirs := range paths {
		if len(dirs) > 1 {
			out = append(out, Clash{Name: name, Paths: dirs})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// firstClash is the error a source returns when it holds a clash and cannot be
// resolved: it names the first two directories.
func firstClash(clashes []Clash) error {
	if len(clashes) == 0 {
		return nil
	}
	c := clashes[0]
	return fmt.Errorf("%w: %q is in %s and %s", ErrDuplicate, c.Name, c.Paths[0], c.Paths[1])
}

// frontmatterAt reads the frontmatter of the SKILL.md in dir.
func frontmatterAt(dir string) skill.Frontmatter {
	f, err := os.Open(filepath.Join(dir, skillFile))
	if err != nil {
		return skill.Frontmatter{}
	}
	defer func() { _ = f.Close() }()
	return skill.ParseFrontmatter(f)
}
