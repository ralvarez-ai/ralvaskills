package skill

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	// MaxNameLen is the longest a skill name may be.
	MaxNameLen = 64

	// RetirementsPath is where the repository records retired skill names,
	// relative to the repository root. Every skill rename adds an entry here.
	RetirementsPath = "registry/aliases.toml"
)

var namePattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Retirements maps a retired skill name to the name that replaced it. No skill
// may take a key of this map again, and the registry publishes it as the alias
// table that keeps the old name resolving.
type Retirements map[string]string

// retirementsFile is the on-disk shape of registry/aliases.toml.
type retirementsFile struct {
	Aliases Retirements `toml:"aliases"`
}

// ValidateName reports why name is not a legal skill name. A legal name is
// kebab-case — lowercase letters, digits and single hyphens — at most
// MaxNameLen characters, with no leading, trailing or doubled hyphen.
func ValidateName(name string) error {
	if name == "" {
		return errors.New("skill name is empty")
	}
	if len(name) > MaxNameLen {
		return fmt.Errorf("skill name %q is %d characters, limit is %d", name, len(name), MaxNameLen)
	}
	if !namePattern.MatchString(name) {
		return fmt.Errorf("skill name %q is not kebab-case", name)
	}
	return nil
}

// LoadRetirements reads the retired-name table at path. A missing file is an
// error, so deleting the table cannot silently disable the check.
func LoadRetirements(path string) (Retirements, error) {
	var file retirementsFile
	if _, err := toml.DecodeFile(path, &file); err != nil {
		return nil, fmt.Errorf("load retirements %s: %w", path, err)
	}
	if file.Aliases == nil {
		return Retirements{}, nil
	}
	return file.Aliases, nil
}

// CheckRetirements reports when name is a retired name, or when a retired name
// appears in it as a whole hyphen-separated component: with hexagonal-arch
// retired, hexagonal-arch-v2 is rejected too. Prose that mentions a retired
// name is out of scope; only names a skill can take are checked.
func CheckRetirements(name string, retired Retirements) error {
	if next, ok := retired[name]; ok {
		return fmt.Errorf("skill name %q is retired; it was renamed to %q", name, next)
	}
	parts := strings.Split(name, "-")
	for i := range parts {
		for j := i + 1; j <= len(parts); j++ {
			component := strings.Join(parts[i:j], "-")
			if next, ok := retired[component]; ok {
				return fmt.Errorf(
					"skill name %q contains retired name %q, renamed to %q",
					name, component, next,
				)
			}
		}
	}
	return nil
}

// ValidateRetirements rejects a table the registry cannot honor: an illegal
// key or target, an empty target, a self-reference, or a chain.
func ValidateRetirements(retired Retirements) error {
	for old, next := range retired {
		if err := ValidateName(old); err != nil {
			return fmt.Errorf("retired name: %w", err)
		}
		if next == "" {
			return fmt.Errorf("retired name %q has no target", old)
		}
		if err := ValidateName(next); err != nil {
			return fmt.Errorf("target of %q: %w", old, err)
		}
		if next == old {
			return fmt.Errorf("retired name %q targets itself", old)
		}
		if after, ok := retired[next]; ok {
			return fmt.Errorf("retired name %q -> %q is a chain: %q is also retired (-> %q)", old, next, next, after)
		}
	}
	return nil
}

// CatalogueIssues reports every naming problem under root: an illegal name, a
// retired name, a directory whose frontmatter name disagrees with it, or a name
// used twice. Sorted, so the output is stable.
func CatalogueIssues(root string, retired Retirements) ([]string, error) {
	var issues []string
	seen := make(map[string]string)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || !d.IsDir() {
			return walkErr
		}
		if !isSkillRoot(path) {
			return nil
		}

		name := filepath.Base(path)
		if err := ValidateName(name); err != nil {
			issues = append(issues, path+": "+err.Error())
		}
		if err := CheckRetirements(name, retired); err != nil {
			issues = append(issues, path+": "+err.Error())
		}
		issues = append(issues, frontmatterNameIssues(path, name)...)
		if first, ok := seen[name]; ok {
			issues = append(issues, fmt.Sprintf("%s: skill name %q already used by %s", path, name, first))
		} else {
			seen[name] = path
		}
		return fs.SkipDir
	})
	if err != nil {
		return nil, fmt.Errorf("walk %s: %w", root, err)
	}
	slices.Sort(issues)
	return issues, nil
}

// frontmatterNameIssues checks that name is declared in the frontmatter and
// matches the directory that holds it. A missing name is reported once.
func frontmatterNameIssues(path, name string) []string {
	fm, err := ReadFrontmatter(path)
	if err != nil {
		return []string{path + ": " + err.Error()}
	}
	if fm.Name == "" {
		return []string{path + ": frontmatter has no name"}
	}
	if fm.Name != name {
		return []string{fmt.Sprintf("%s: frontmatter name %q does not match the directory", path, fm.Name)}
	}
	return nil
}

// isSkillRoot reports whether dir holds a SKILL.md, which makes it a skill root.
func isSkillRoot(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, SkillFileName))
	return err == nil
}
