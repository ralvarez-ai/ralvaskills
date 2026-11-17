package skill

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Frontmatter is the subset of a SKILL.md YAML frontmatter that rsk reads.
type Frontmatter struct {
	Name         string
	Version      string
	Description  string
	AllowedTools []string
}

const (
	// FrontmatterDelimiter separates YAML metadata from the skill body.
	FrontmatterDelimiter = "---"

	// NamePrefix is the YAML frontmatter key that declares a skill's name. It
	// must equal the directory the skill lives in.
	NamePrefix = "name:"

	// VersionPrefix is the YAML frontmatter key used to store a skill's version.
	VersionPrefix = "version:"

	// DescriptionPrefix is the YAML frontmatter key used to store a skill's
	// description.
	DescriptionPrefix = "description:"

	// AllowedToolsPrefix is the YAML frontmatter key that lists the tool
	// permissions a skill needs to run without asking. It accepts the client's
	// own rule syntax, Tool(specifier), plus a comma or space separated string.
	AllowedToolsPrefix = "allowed-tools:"

	// MaxAllowedTools caps how many permissions one skill may declare. A skill
	// that needs more is asking for too much: split it or drop the entry.
	MaxAllowedTools = 5

	// MaxAllowedToolLength caps a single entry, so a pattern cannot smuggle a
	// whole command line into a rule.
	MaxAllowedToolLength = 200
)

// Walk discovers all skills under root. A directory that contains SKILL.md is
// a skill root; Walk does not descend into it further. The returned skills all
// have Source set to source.
func Walk(root string, source Source) ([]Skill, error) {
	var skills []Skill
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return err
		}
		if _, statErr := os.Stat(filepath.Join(path, SkillFileName)); statErr == nil {
			s := parseSkill(root, path, source)
			skills = append(skills, s)
			return fs.SkipDir
		}
		return nil
	})
	return skills, err
}

// ReadFrontmatter parses the version, description and allowed-tools from the
// SKILL.md in dir. A missing description is not an error; a missing version is.
func ReadFrontmatter(dir string) (Frontmatter, error) {
	return readFrontmatter(filepath.Join(dir, SkillFileName))
}

// ReadVersion reads the version field from SKILL.md in dir.
func ReadVersion(dir string) (string, error) {
	fm, err := ReadFrontmatter(dir)
	if err != nil {
		return "", err
	}
	return fm.Version, nil
}

func parseSkill(root, dir string, source Source) Skill {
	fm, err := readFrontmatter(filepath.Join(dir, SkillFileName))
	if err != nil {
		fm = Frontmatter{}
	}

	rel, relErr := filepath.Rel(root, dir)
	if relErr != nil {
		rel = dir
	}

	return Skill{
		Name:         filepath.Base(dir),
		Version:      fm.Version,
		Description:  fm.Description,
		AllowedTools: NormalizeAllowedTools(fm.AllowedTools),
		Path:         dir,
		Source:       source,
		IsPersonal:   IsPersonalPath(rel),
	}
}

// readFrontmatter reads the frontmatter keys rsk uses from path. allowed-tools
// accepts both a YAML block list and an inline string ("Read Grep" or
// "Read, Bash(git status *)"); an unindented key ends the block.
func readFrontmatter(path string) (Frontmatter, error) {
	f, err := os.Open(path)
	if err != nil {
		return Frontmatter{}, err
	}
	defer func() { _ = f.Close() }()

	var fm Frontmatter
	scanner := bufio.NewScanner(f)
	inFrontmatter := false
	inAllowedToolsList := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == FrontmatterDelimiter {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			break
		}
		if !inFrontmatter {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if inAllowedToolsList {
			if item, isItem := strings.CutPrefix(trimmed, "-"); isItem {
				fm.AllowedTools = append(fm.AllowedTools, unquote(item))
				continue
			}
			inAllowedToolsList = false
		}
		if v, ok := strings.CutPrefix(trimmed, AllowedToolsPrefix); ok {
			inAllowedToolsList = fm.readAllowedTools(v)
			continue
		}
		fm.readScalar(line)
	}
	if fm.Version == "" {
		return Frontmatter{}, fmt.Errorf("version field not found in %s", path)
	}
	return fm, nil
}

// readAllowedTools records an inline allowed-tools value and reports whether a
// block list follows it instead.
func (fm *Frontmatter) readAllowedTools(value string) bool {
	inline := strings.TrimSpace(value)
	if inline == "" {
		return true
	}
	fm.AllowedTools = append(fm.AllowedTools, splitRules(inline)...)
	return false
}

// readScalar records the simple string keys rsk understands.
func (fm *Frontmatter) readScalar(line string) {
	if v, ok := strings.CutPrefix(line, NamePrefix); ok {
		fm.Name = unquote(v)
	}
	if v, ok := strings.CutPrefix(line, VersionPrefix); ok {
		fm.Version = unquote(v)
	}
	if v, ok := strings.CutPrefix(line, DescriptionPrefix); ok {
		fm.Description = unquote(v)
	}
}

// splitRules splits an inline allowed-tools value on commas and on whitespace
// outside parentheses, so "Read Grep" is two entries and
// "Bash(git status *)" stays one.
func splitRules(value string) []string {
	var (
		entries []string
		current strings.Builder
		depth   int
	)
	flush := func() {
		if entry := strings.TrimSpace(current.String()); entry != "" {
			entries = append(entries, unquote(entry))
		}
		current.Reset()
	}
	for _, r := range value {
		switch {
		case r == '(':
			depth++
			current.WriteRune(r)
		case r == ')':
			if depth > 0 {
				depth--
			}
			current.WriteRune(r)
		case depth == 0 && (r == ',' || r == ' ' || r == '\t'):
			flush()
		default:
			current.WriteRune(r)
		}
	}
	flush()
	return entries
}

// unquote trims surrounding whitespace and the YAML quotes around value.
func unquote(value string) string {
	return strings.Trim(strings.TrimSpace(value), `"'`)
}

// NormalizeAllowedTools trims each entry, drops empty and duplicate entries,
// and returns them sorted.
func NormalizeAllowedTools(entries []string) []string {
	if len(entries) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(entries))
	kept := make([]string, 0, len(entries))
	for _, entry := range entries {
		trimmed := strings.TrimSpace(entry)
		if trimmed == "" {
			continue
		}
		if _, dup := seen[trimmed]; dup {
			continue
		}
		seen[trimmed] = struct{}{}
		kept = append(kept, trimmed)
	}
	slices.Sort(kept)
	return kept
}

// ValidateAllowedTools rejects a declaration rsk will not write: too many
// entries, an entry past the length cap, a malformed Tool(specifier), or a
// pattern broad enough to allow anything.
func ValidateAllowedTools(entries []string) error {
	if len(entries) > MaxAllowedTools {
		return fmt.Errorf("%d entries, at most %d allowed", len(entries), MaxAllowedTools)
	}
	for _, entry := range entries {
		if len(entry) > MaxAllowedToolLength {
			return fmt.Errorf("%q is longer than %d characters", entry, MaxAllowedToolLength)
		}
		toolName, spec, ok := SplitRule(entry)
		if !ok {
			return fmt.Errorf("%q is not Tool(specifier), for example Bash(git status *)", entry)
		}
		if toolName == "" {
			return fmt.Errorf("%q names no tool", entry)
		}
		if spec == "*" || spec == "**" {
			return fmt.Errorf("%q matches every command; name the command or a prefix", entry)
		}
		if strings.HasPrefix(spec, "*") {
			return fmt.Errorf("%q starts with a wildcard; name the command or a prefix", entry)
		}
	}
	return nil
}

// SplitRule splits a Tool(specifier) entry into its tool name and specifier.
// The specifier is everything between the first "(" and the closing ")".
func SplitRule(entry string) (toolName, spec string, ok bool) {
	entry = strings.TrimSpace(entry)
	open := strings.Index(entry, "(")
	if open <= 0 || !strings.HasSuffix(entry, ")") {
		return "", "", false
	}
	toolName = strings.TrimSpace(entry[:open])
	spec = strings.TrimSpace(entry[open+1 : len(entry)-1])
	if toolName == "" || spec == "" {
		return "", "", false
	}
	return toolName, spec, true
}
