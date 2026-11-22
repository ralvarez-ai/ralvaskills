package source

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal/git"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

type (
	// Path resolves skills from a directory the user declared with `rsk source
	// add <dir>`. Unlike the local clone it is walked as given, so it has no
	// skills/ folder; it honors roots and excludes, and applies the rules a git
	// skill does: a valid name, no symlink followed or kept, a name used once.
	// There is no network, so there is nothing to trust.
	Path struct {
		spec  Spec
		globs []*regexp.Regexp
	}

	// Skipped is a directory that looked like a skill but was left out, and why.
	Skipped struct {
		Path   string `json:"path"`
		Reason string `json:"reason"`
	}
)

// newPath validates spec and builds the resolver.
func newPath(spec Spec) (*Path, error) {
	if !filepath.IsAbs(spec.URL) {
		return nil, fmt.Errorf("source %s: path %q must be absolute", spec.ID, spec.URL)
	}
	for _, root := range spec.Roots {
		if err := git.ValidateRepoPath(root); err != nil {
			return nil, fmt.Errorf("source %s: %w", spec.ID, err)
		}
	}
	globs := make([]*regexp.Regexp, 0, len(spec.Excludes))
	for _, pattern := range spec.Excludes {
		re, err := compileGlob(pattern)
		if err != nil {
			return nil, fmt.Errorf("source %s: exclude %q: %w", spec.ID, pattern, err)
		}
		globs = append(globs, re)
	}
	return &Path{spec: spec, globs: globs}, nil
}

// Spec returns the descriptor the source was built from.
func (p *Path) Spec() Spec {
	return p.spec
}

// All lists every skill in the directory.
func (p *Path) All(_ context.Context) ([]skill.Skill, error) {
	st, err := p.scan()
	if err != nil {
		return nil, err
	}
	if err = firstClash(clashesOf(st.paths)); err != nil {
		return nil, fmt.Errorf("source %s: %w", p.spec.ID, err)
	}
	return st.skills, nil
}

// Scan is a dry-run discovery: what resolves, what is left out, and which names
// clash, which All would refuse.
func (p *Path) Scan(_ context.Context) (ScanResult, error) {
	st, err := p.scan()
	if err != nil {
		return ScanResult{}, err
	}
	clashes := clashesOf(st.paths)
	res := ScanResult{Skills: []ScannedSkill{}, Skipped: st.skipped, Clashes: clashes}
	for i := range st.skills {
		s := &st.skills[i]
		if len(st.paths[s.Name]) > 1 {
			continue
		}
		fm := frontmatterAt(s.Path)
		res.Skills = append(res.Skills, ScannedSkill{
			Name: s.Name, Path: st.paths[s.Name][0], Version: s.Version, Description: s.Description,
			Issues: issuesOf(s.Name, &fm),
		})
	}
	return res, nil
}

// Skipped lists the directories that were left out, with the reason.
func (p *Path) Skipped(_ context.Context) ([]Skipped, error) {
	st, err := p.scan()
	return st.skipped, err
}

// Find returns the skill with the given name.
func (p *Path) Find(ctx context.Context, name string) (skill.Skill, error) {
	skills, err := p.All(ctx)
	if err != nil {
		return skill.Skill{}, err
	}
	for i := range skills {
		if skills[i].Name == name {
			return skills[i], nil
		}
	}
	return skill.Skill{}, fmt.Errorf("%w: source %s has no skill %q", ErrNotFound, p.spec.ID, name)
}

type scanState struct {
	skills  []skill.Skill
	skipped []Skipped
	paths   map[string][]string // name to the directories that provide it
}

func (p *Path) scan() (scanState, error) {
	roots := p.spec.Roots
	if len(roots) == 0 {
		roots = []string{"."}
	}
	st := scanState{paths: map[string][]string{}}
	for _, root := range roots {
		if err := p.confine(root); err != nil {
			return scanState{}, err
		}
		walkRoot := filepath.Join(p.spec.URL, filepath.FromSlash(root))
		err := filepath.WalkDir(walkRoot, func(dir string, d fs.DirEntry, err error) error {
			return p.visit(dir, d, err, &st)
		})
		if err != nil {
			return scanState{}, fmt.Errorf("source %s: %w", p.spec.ID, err)
		}
	}
	sort.Slice(st.skills, func(i, j int) bool { return st.skills[i].Name < st.skills[j].Name })
	sort.Slice(st.skipped, func(i, j int) bool { return st.skipped[i].Path < st.skipped[j].Path })
	return st, nil
}

// visit handles one directory of the walk: it prunes an excluded one, records a
// skill or the reason one is left out, and does not look inside a skill.
func (p *Path) visit(dir string, d fs.DirEntry, err error, st *scanState) error {
	if err != nil || !d.IsDir() {
		return err
	}
	rel, err := filepath.Rel(p.spec.URL, dir)
	if err != nil {
		return err
	}
	rel = filepath.ToSlash(rel)
	if rel != "." && p.excluded(rel) {
		return fs.SkipDir
	}
	found, skip, handled := p.readSkill(dir, rel)
	if skip != nil {
		st.skipped = append(st.skipped, *skip)
	}
	switch {
	case !handled:
		return nil
	case skip != nil:
	default:
		st.paths[found.Name] = append(st.paths[found.Name], rel)
		if len(st.paths[found.Name]) == 1 {
			st.skills = append(st.skills, found)
		}
	}
	return fs.SkipDir
}

// confine refuses a root that passes through a symlink, so a root the user
// declared cannot be pointed outside the directory afterwards.
func (p *Path) confine(root string) error {
	if root == "." {
		return nil
	}
	cur := p.spec.URL
	for seg := range strings.SplitSeq(root, "/") {
		cur = filepath.Join(cur, seg)
		info, err := os.Lstat(cur)
		if err != nil {
			return fmt.Errorf("source %s: root %s: %w", p.spec.ID, root, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("source %s: root %s passes through a symlink", p.spec.ID, root)
		}
	}
	return nil
}

func (p *Path) excluded(rel string) bool {
	for _, re := range p.globs {
		if matchGlob(re, rel) {
			return true
		}
	}
	return false
}

// readSkill reports whether dir is a skill directory (handled), and if so
// either the skill or the reason it is left out.
func (p *Path) readSkill(dir, rel string) (skill.Skill, *Skipped, bool) {
	info, err := os.Lstat(filepath.Join(dir, skillFile))
	if err != nil {
		return skill.Skill{}, nil, false
	}
	left := func(reason string) (skill.Skill, *Skipped, bool) {
		return skill.Skill{}, &Skipped{Path: rel, Reason: reason}, true
	}
	if rel == "." {
		// The skill's name is its directory's, and the root has none to give.
		skip := &Skipped{Path: rel, Reason: skillFile + " at the root has no directory name"}
		return skill.Skill{}, skip, false
	}
	if !info.Mode().IsRegular() {
		return left(skillFile + " is not a regular file")
	}
	if reason := unsafeContent(dir); reason != "" {
		return left(reason)
	}
	name := path.Base(rel)
	if nameErr := skill.ValidateName(name); nameErr != nil {
		return left(nameErr.Error())
	}
	if info.Size() > maxSkillFile {
		return left(skillFile + " is too large")
	}
	f, err := os.Open(filepath.Join(dir, skillFile))
	if err != nil {
		return left(skillFile + " is unreadable: " + err.Error())
	}
	defer func() { _ = f.Close() }()
	fm := skill.ParseFrontmatter(io.LimitReader(f, maxSkillFile))
	return skill.Skill{
		Name: name, Version: fm.Version, Description: fm.Description, Path: dir, Source: skill.SourcePath,
		AllowedTools: skill.NormalizeAllowedTools(fm.AllowedTools),
	}, nil, true
}

// unsafeContent names the first thing inside a skill directory that a skill may
// not hold: a symlink or any entry that is neither a file nor a directory.
func unsafeContent(dir string) string {
	reason := ""
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		switch {
		case err != nil:
			reason = "unreadable: " + err.Error()
		case d.Type()&os.ModeSymlink != 0:
			reason = "contains a symlink (" + relTo(dir, p) + ")"
		case !d.IsDir() && !d.Type().IsRegular():
			reason = "contains a special file (" + relTo(dir, p) + ")"
		default:
			return nil
		}
		return fs.SkipAll
	})
	if err != nil && reason == "" {
		reason = "unreadable: " + err.Error()
	}
	return reason
}

func relTo(base, p string) string {
	rel, err := filepath.Rel(base, p)
	if err != nil {
		return p
	}
	return filepath.ToSlash(rel)
}

// ScanSkillDir reads the one skill in dir, for `rsk install --from`. It applies
// the rules a path source does: dir must be a real directory (not a link) with a
// regular SKILL.md and no symlink inside, and its name must be a valid skill
// name. A directory holding several skills is a source, not a skill.
func ScanSkillDir(dir string) (skill.Skill, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return skill.Skill{}, fmt.Errorf("resolve %s: %w", dir, err)
	}
	info, err := os.Lstat(abs)
	switch {
	case err != nil:
		return skill.Skill{}, fmt.Errorf("--from %s: %w", dir, err)
	case info.Mode()&os.ModeSymlink != 0:
		return skill.Skill{}, fmt.Errorf("--from %s is a symlink: pass the real directory", dir)
	case !info.IsDir():
		return skill.Skill{}, fmt.Errorf("--from %s is not a directory", dir)
	}
	found, skip, handled := (&Path{}).readSkill(abs, filepath.Base(abs))
	switch {
	case !handled:
		return skill.Skill{}, fmt.Errorf(
			"%s has no %s: --from takes one skill directory; for a directory of skills use 'rsk source add %s'",
			dir, skillFile, dir,
		)
	case skip != nil:
		return skill.Skill{}, fmt.Errorf("--from %s: %s", dir, skip.Reason)
	}
	return found, nil
}
