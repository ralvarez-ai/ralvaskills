package source

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
	"go.ralva.dev/ralvaskills/v3/internal/fsx"
	"go.ralva.dev/ralvaskills/v3/internal/git"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

const (
	// DefaultIndexTTL is how long a git source's name index counts as fresh.
	DefaultIndexTTL = 24 * time.Hour

	maxSkillFile = 256 << 10
	indexVersion = 1
	clockSkew    = 5 * time.Minute
	privateFile  = 0o600
	skillFile    = "SKILL.md"
	regularFile  = "100644"
	execFile     = "100755"
	indexTemp    = ".rsk-index-*.tmp"
)

var versionPattern = regexp.MustCompile(`^[0-9A-Za-z][0-9A-Za-z.+-]*$`)

// repoOps is what Git needs from a bare clone; git.Repo is the production
// implementation, and tests supply an in-memory one.
//
//nolint:interfacebloat // the operations a git source performs on its clone
type repoOps interface {
	Fetch(ctx context.Context, remote git.RemoteURL, ref string) (string, error)
	ListTree(ctx context.Context, commit string, roots []string) ([]git.TreeEntry, error)
	TreeID(ctx context.Context, commit, dir string) (string, error)
	ReadBlobs(ctx context.Context, commit string, paths []string, maxBytes int64) (map[string]git.Blob, error)
	Archive(ctx context.Context, commit, dir, dest string, lim git.Limits) error
}

type (
	// Git resolves skills from a declared git repository. Every public method
	// asks the trust gate first, so an untrusted source is never read or
	// contacted. Names come from an index of the repository at one commit, kept
	// on disk; a skill is extracted into the cache only when it is resolved.
	Git struct {
		spec    Spec
		remote  git.RemoteURL
		repo    repoOps
		root    string
		owner   string
		name    string
		gate    func(id, url string) error
		offline bool
		ttl     time.Duration
		now     func() time.Time
		globs   []*regexp.Regexp
	}

	// LockedSkill is what rsk.lock records about a git skill: enough to fetch
	// the same content again without consulting the branch.
	LockedSkill struct {
		Name   string
		Rev    string
		Commit string
		Tree   string
		Path   string
	}

	gitIndex struct {
		Version   int           `json:"version"`
		ID        string        `json:"id"`
		URL       string        `json:"url"`
		Rev       string        `json:"rev,omitempty"`
		Roots     []string      `json:"roots,omitempty"`
		Excludes  []string      `json:"excludes,omitempty"`
		FetchedAt time.Time     `json:"fetched_at"`
		Commit    string        `json:"commit"`
		Skills    []gitSkill    `json:"skills"`
		Skipped   []skippedPath `json:"skipped,omitempty"`
	}

	gitSkill struct {
		Name         string   `json:"name"`
		Path         string   `json:"path"`
		Tree         string   `json:"tree"`
		Version      string   `json:"version,omitempty"`
		Description  string   `json:"description,omitempty"`
		AllowedTools []string `json:"allowed_tools,omitempty"`
	}

	// skippedPath is a directory that looked like a skill but was left out, and
	// why: the reason is how a user learns what to fix in a repository.
	skippedPath struct {
		Path   string `json:"path"`
		Reason string `json:"reason"`
	}
)

// newGit validates spec and builds the resolver. A nil repo means the real bare
// clone under the cache root.
func newGit(spec Spec, opts Options, repo repoOps) (*Git, error) {
	if opts.Gate == nil {
		return nil, fmt.Errorf("source %s: a git source needs a trust gate", spec.ID)
	}
	if opts.CacheRoot == "" {
		return nil, fmt.Errorf("source %s: a git source needs a cache location", spec.ID)
	}
	id, err := skill.ParseSourceID(spec.ID)
	if err != nil || id != spec.ID {
		return nil, fmt.Errorf("source %s: the identity must be canonical @owner/repo", spec.ID)
	}
	remote, err := git.ParseRemoteURL(spec.URL)
	if err != nil {
		return nil, fmt.Errorf("source %s: %w", spec.ID, err)
	}
	if spec.Rev != "" {
		if err = git.ValidateRef(spec.Rev); err != nil {
			return nil, fmt.Errorf("source %s: %w", spec.ID, err)
		}
	}
	for _, root := range spec.Roots {
		if err = git.ValidateRepoPath(root); err != nil {
			return nil, fmt.Errorf("source %s: %w", spec.ID, err)
		}
	}
	globs := make([]*regexp.Regexp, 0, len(spec.Excludes))
	for _, pattern := range spec.Excludes {
		re, globErr := compileGlob(pattern)
		if globErr != nil {
			return nil, fmt.Errorf("source %s: exclude %q: %w", spec.ID, pattern, globErr)
		}
		globs = append(globs, re)
	}

	owner, name, _ := strings.Cut(strings.TrimPrefix(spec.ID, "@"), "/")
	g := &Git{
		spec: spec, remote: remote, repo: repo, root: opts.CacheRoot, owner: owner, name: name,
		gate: opts.Gate, offline: opts.Offline, ttl: opts.IndexTTL, now: opts.Now, globs: globs,
	}
	if g.ttl == 0 {
		g.ttl = DefaultIndexTTL
	}
	if g.now == nil {
		g.now = time.Now
	}
	if g.repo == nil {
		g.repo = git.Repo{Dir: filepath.Join(g.root, "git", owner, name+".git")}
	}
	return g, nil
}

// Spec returns the descriptor the source was built from.
func (g *Git) Spec() Spec {
	return g.spec
}

// All lists every skill in the source at its indexed commit.
func (g *Git) All(ctx context.Context) ([]skill.Skill, error) {
	idx, err := g.index(ctx)
	if err != nil {
		return nil, err
	}
	skills := make([]skill.Skill, len(idx.Skills))
	for i := range idx.Skills {
		skills[i] = g.toSkill(&idx.Skills[i], "", idx.Commit, g.spec.Rev)
	}
	return skills, nil
}

// Find returns the skill with the given name, extracting it into the cache at
// its indexed commit if it is not there yet.
func (g *Git) Find(ctx context.Context, name string) (skill.Skill, error) {
	idx, err := g.index(ctx)
	if err != nil {
		return skill.Skill{}, err
	}
	for i := range idx.Skills {
		if idx.Skills[i].Name != name {
			continue
		}
		dest, extractErr := g.materializeAt(ctx, name, idx.Skills[i].Path, idx.Commit)
		if extractErr != nil {
			return skill.Skill{}, fmt.Errorf("source %s: extract %q: %w", g.spec.ID, name, extractErr)
		}
		return g.toSkill(&idx.Skills[i], dest, idx.Commit, g.spec.Rev), nil
	}
	return skill.Skill{}, fmt.Errorf("%w: source %s has no skill %q", ErrNotFound, g.spec.ID, name)
}

func (g *Git) indexPath() string {
	return filepath.Join(g.root, "index", g.owner, g.name+".json")
}

func (g *Git) extractDir(name, commit string) string {
	return filepath.Join(g.root, "sources", g.owner, g.name, name, commit)
}

// index returns the source's name index. The gate comes first. A fresh cached
// index answers without touching the network; a stale or missing one is
// rebuilt, and if that fails the lookup fails, so a source that cannot be
// reached never silently answers from old data. Offline accepts any existing
// index and fails when there is none.
func (g *Git) index(ctx context.Context) (gitIndex, error) {
	if err := g.gate(g.spec.ID, g.spec.URL); err != nil {
		return gitIndex{}, err
	}
	cached, err := g.loadIndex()
	usable := err == nil && g.matches(&cached)

	switch {
	case g.offline && usable:
		return cached, nil
	case g.offline:
		return gitIndex{}, fmt.Errorf("%w: %s has no index on this machine; run without --offline to fetch it",
			ErrNeverFetched, g.spec.ID)
	case usable && g.fresh(&cached):
		return cached, nil
	}

	idx, err := g.refresh(ctx)
	if err != nil {
		return gitIndex{}, fmt.Errorf("cannot consult source %s: %w", g.spec.ID, err)
	}
	return idx, nil
}

func (g *Git) matches(idx *gitIndex) bool {
	return idx.Version == indexVersion && idx.ID == g.spec.ID && idx.URL == g.spec.URL && idx.Rev == g.spec.Rev &&
		slices.Equal(idx.Roots, g.spec.Roots) && slices.Equal(idx.Excludes, g.spec.Excludes)
}

func (g *Git) fresh(idx *gitIndex) bool {
	age := g.now().Sub(idx.FetchedAt)
	return age < g.ttl && age > -clockSkew
}

func (g *Git) loadIndex() (gitIndex, error) {
	data, err := os.ReadFile(g.indexPath())
	if err != nil {
		return gitIndex{}, err
	}
	var idx gitIndex
	if err = json.Unmarshal(data, &idx); err != nil {
		return gitIndex{}, fmt.Errorf("parse %s: %w", g.indexPath(), err)
	}
	return idx, nil
}

func (g *Git) saveIndex(idx *gitIndex) error {
	file := g.indexPath()
	if err := os.MkdirAll(filepath.Dir(file), fsperm.Dir); err != nil {
		return fmt.Errorf("create %s: %w", filepath.Dir(file), err)
	}
	return fsx.WriteAtomicPerm(file, indexTemp, privateFile, func(w io.Writer) error {
		return json.NewEncoder(w).Encode(idx)
	})
}

// refresh fetches the declared rev, indexes it, and stores the index.
func (g *Git) refresh(ctx context.Context) (gitIndex, error) {
	rev := g.spec.Rev
	if rev == "" {
		rev = "HEAD"
	}
	commit, err := g.repo.Fetch(ctx, g.remote, rev)
	if err != nil {
		return gitIndex{}, err
	}
	entries, err := g.repo.ListTree(ctx, commit, g.spec.Roots)
	if err != nil {
		return gitIndex{}, err
	}
	skills, skipped, err := g.buildSkills(ctx, commit, entries)
	if err != nil {
		return gitIndex{}, err
	}
	idx := gitIndex{
		Version: indexVersion, ID: g.spec.ID, URL: g.spec.URL, Rev: g.spec.Rev, Roots: g.spec.Roots,
		Excludes: g.spec.Excludes, FetchedAt: g.now().UTC(), Commit: commit, Skills: skills, Skipped: skipped,
	}
	if err = g.saveIndex(&idx); err != nil {
		return gitIndex{}, err
	}
	return idx, nil
}

// buildSkills finds the skills in a tree listing. A skill is a directory with a
// regular SKILL.md; like skill.Walk, a skill directory is not searched further.
// A directory with a bad name, a symlink or submodule inside it, or an unreadable
// SKILL.md is left out and recorded with the reason. Two skills with one name
// are an error.
func (g *Git) buildSkills(
	ctx context.Context,
	commit string,
	entries []git.TreeEntry,
) ([]gitSkill, []skippedPath, error) {
	found := g.findSkillDirs(entries)
	skipped := found.skipped
	unsafe := unsafeEntries(entries, found.roots)

	var candidates, files []string
	for _, dir := range found.sorted {
		if !found.roots[dir] {
			continue
		}
		if reason, bad := unsafe[dir]; bad {
			skipped = append(skipped, skippedPath{Path: dir, Reason: reason})
			continue
		}
		if err := skill.ValidateName(path.Base(dir)); err != nil {
			skipped = append(skipped, skippedPath{Path: dir, Reason: err.Error()})
			continue
		}
		candidates = append(candidates, dir)
		files = append(files, dir+"/"+skillFile)
	}

	blobs, err := g.repo.ReadBlobs(ctx, commit, files, maxSkillFile)
	if err != nil {
		return nil, nil, err
	}
	skills, more, err := readSkills(candidates, files, blobs, found.trees)
	if err != nil {
		return nil, nil, err
	}
	skipped = append(skipped, more...)
	sort.Slice(skills, func(i, j int) bool { return skills[i].Name < skills[j].Name })
	sort.Slice(skipped, func(i, j int) bool { return skipped[i].Path < skipped[j].Path })
	return skills, skipped, nil
}

// skillDirs is what a tree listing says about where skills are.
type skillDirs struct {
	trees   map[string]string // directory path to tree id
	sorted  []string          // every directory with a SKILL.md that is not excluded
	roots   map[string]bool   // the outermost of those: the skill directories
	skipped []skippedPath
}

func (g *Git) findSkillDirs(entries []git.TreeEntry) skillDirs {
	found := skillDirs{trees: map[string]string{}, roots: map[string]bool{}}
	for i := range entries {
		e := &entries[i]
		if e.Type == "tree" {
			found.trees[e.Path] = e.ObjectID
			continue
		}
		if path.Base(e.Path) != skillFile {
			continue
		}
		dir := path.Dir(e.Path)
		switch {
		case e.Mode != regularFile && e.Mode != execFile:
			found.skipped = append(found.skipped, skippedPath{Path: dir, Reason: skillFile + " is not a regular file"})
		case dir == ".":
			found.skipped = append(found.skipped,
				skippedPath{Path: e.Path, Reason: skillFile + " at the repository root has no directory name"})
		case !g.excluded(dir):
			found.sorted = append(found.sorted, dir)
		}
	}

	sort.Strings(found.sorted)
	for _, dir := range found.sorted {
		nested := false
		for p := path.Dir(dir); p != "."; p = path.Dir(p) {
			nested = nested || found.roots[p]
		}
		if !nested {
			found.roots[dir] = true
		}
	}
	return found
}

// readSkills turns the SKILL.md files of the candidate directories into index
// entries, leaving out a missing or oversized file and rejecting a name used
// twice.
func readSkills(
	candidates, files []string, blobs map[string]git.Blob, trees map[string]string,
) ([]gitSkill, []skippedPath, error) {
	var skills []gitSkill
	var skipped []skippedPath
	byName := map[string]string{}
	for i, dir := range candidates {
		b := blobs[files[i]]
		switch {
		case b.Missing:
			skipped = append(skipped, skippedPath{Path: dir, Reason: skillFile + " is missing"})
			continue
		case b.Data == nil:
			skipped = append(skipped, skippedPath{Path: dir, Reason: skillFile + " is too large"})
			continue
		}
		name := path.Base(dir)
		if first, dup := byName[name]; dup {
			return nil, nil, fmt.Errorf("%w: %q is in %s and %s", ErrDuplicate, name, first, dir)
		}
		byName[name] = dir
		fm := skill.ParseFrontmatter(bytes.NewReader(b.Data))
		skills = append(skills, gitSkill{
			Name: name, Path: dir, Tree: trees[dir], Version: fm.Version, Description: fm.Description,
			AllowedTools: skill.NormalizeAllowedTools(fm.AllowedTools),
		})
	}
	return skills, skipped, nil
}

// unsafeEntries maps each kept skill directory that contains a symlink or a
// submodule to the reason it is excluded.
func unsafeEntries(entries []git.TreeEntry, kept map[string]bool) map[string]string {
	unsafe := map[string]string{}
	for i := range entries {
		e := &entries[i]
		var what string
		switch e.Mode {
		case git.ModeSymlink:
			what = "symlink"
		case git.ModeSubmodule:
			what = "submodule"
		default:
			continue
		}
		for p := path.Dir(e.Path); p != "."; p = path.Dir(p) {
			if kept[p] {
				if _, seen := unsafe[p]; !seen {
					unsafe[p] = fmt.Sprintf("contains a %s (%s)", what, e.Path)
				}
				break
			}
		}
	}
	return unsafe
}

func (g *Git) excluded(dir string) bool {
	for _, re := range g.globs {
		if matchGlob(re, dir) {
			return true
		}
	}
	return false
}

// materializeAt extracts the skill at commit into the cache, once. If the clone
// no longer holds the commit's objects it fetches that commit by id and tries
// again, unless offline or the archive was refused as unsafe.
func (g *Git) materializeAt(ctx context.Context, name, dir, commit string) (string, error) {
	dest := g.extractDir(name, commit)
	if _, err := os.Stat(filepath.Join(dest, skillFile)); err == nil {
		return dest, nil
	}
	if err := os.RemoveAll(dest); err != nil {
		return "", fmt.Errorf("clear %s: %w", dest, err)
	}

	err := g.repo.Archive(ctx, commit, dir, dest, git.Limits{})
	if err == nil {
		return dest, nil
	}
	if g.offline || errors.Is(err, git.ErrUnsafeEntry) || errors.Is(err, git.ErrTooLarge) {
		return "", err
	}
	if _, fetchErr := g.repo.Fetch(ctx, g.remote, commit); fetchErr != nil {
		return "", fmt.Errorf("fetch %s: %w", commit, fetchErr)
	}
	if err = g.repo.Archive(ctx, commit, dir, dest, git.Limits{}); err != nil {
		return "", err
	}
	return dest, nil
}

func (g *Git) toSkill(sk *gitSkill, dir, commit, rev string) skill.Skill {
	return skill.Skill{
		Name: sk.Name, Version: sk.Version, Description: sk.Description, Path: dir,
		Source: skill.DeclaredSource(g.spec.ID), AllowedTools: sk.AllowedTools,
		Git: &skill.GitRef{URL: g.spec.URL, Rev: rev, Commit: commit, Tree: sk.Tree, Path: sk.Path},
	}
}
