package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
)

var (
	// ErrInvalidRef is returned for a git ref the runner will not pass on.
	ErrInvalidRef = errors.New("invalid git ref")

	// ErrInvalidCommit is returned for a value that is not a full object id.
	ErrInvalidCommit = errors.New("invalid commit id")

	// ErrInvalidPath is returned for a repository path that is not relative and
	// confined.
	ErrInvalidPath = errors.New("invalid repository path")
)

const (
	// ModeSymlink is the tree mode of a symbolic link.
	ModeSymlink = "120000"

	// ModeSubmodule is the tree mode of a submodule (a gitlink).
	ModeSubmodule = "160000"

	pinPrefix = "refs/rsk/"
)

var (
	refPattern    = regexp.MustCompile(`^[A-Za-z0-9._/@+-]{1,255}$`)
	commitPattern = regexp.MustCompile(`^([0-9a-f]{40}|[0-9a-f]{64})$`)
)

// Repo is a bare clone that rsk owns. Every operation goes through the hardened
// runner; the zero value of its unexported field is the production setting.
type Repo struct {
	// Dir is the bare repository's directory.
	Dir string

	allowFile bool
}

// TreeEntry is one entry of a recursive tree listing. Directories are not
// listed: `ls-tree -r` reports blobs and submodules.
type TreeEntry struct {
	Mode     string
	Type     string
	ObjectID string
	Path     string
}

func validateRef(ref string) error {
	if !refPattern.MatchString(ref) || strings.HasPrefix(ref, "-") || ref == "@" ||
		strings.Contains(ref, "..") || strings.HasSuffix(ref, ".lock") || strings.Contains(ref, "//") {
		return fmt.Errorf("%w: %q", ErrInvalidRef, ref)
	}
	return nil
}

func validateCommit(commit string) error {
	if !commitPattern.MatchString(commit) {
		return fmt.Errorf("%w: %q", ErrInvalidCommit, commit)
	}
	return nil
}

func validateRepoPath(p string) error {
	if p == "" || strings.HasPrefix(p, "-") || strings.HasPrefix(p, "/") || strings.ContainsAny(p, "\\\x00\n") {
		return fmt.Errorf("%w: %q", ErrInvalidPath, p)
	}
	for seg := range strings.SplitSeq(p, "/") {
		if seg == "" || seg == "." || seg == ".." {
			return fmt.Errorf("%w: %q", ErrInvalidPath, p)
		}
	}
	return nil
}

// Fetch brings ref (a branch, a tag or a full object id) from remote into the
// bare clone, at depth 1, and returns the commit it resolved to. The commit is
// pinned under refs/rsk/ so garbage collection cannot drop it.
func (r Repo) Fetch(ctx context.Context, remote RemoteURL, ref string) (string, error) {
	if err := validateRef(ref); err != nil {
		return "", err
	}
	if err := r.ensureInit(ctx); err != nil {
		return "", err
	}
	run := r.runner()
	if _, err := run.run(ctx, r.Dir, "fetch",
		"fetch", "--quiet", "--no-tags", "--no-recurse-submodules", "--depth=1", "--", remote.String(), ref,
	); err != nil {
		return "", err
	}
	out, err := run.run(ctx, r.Dir, "rev-parse", "rev-parse", "--verify", "FETCH_HEAD^{commit}")
	if err != nil {
		return "", err
	}
	commit := strings.TrimSpace(string(out))
	if err = validateCommit(commit); err != nil {
		return "", err
	}
	if _, err = run.run(ctx, r.Dir, "update-ref", "update-ref", pinPrefix+commit, commit); err != nil {
		return "", err
	}
	return commit, nil
}

// TreeID returns the object id of the directory dir at commit: the content
// identity of a skill, which an unrelated commit elsewhere in the repository
// does not change. An empty dir means the repository root.
func (r Repo) TreeID(ctx context.Context, commit, dir string) (string, error) {
	if err := validateCommit(commit); err != nil {
		return "", err
	}
	spec := commit + "^{tree}"
	if dir != "" {
		if err := validateRepoPath(dir); err != nil {
			return "", err
		}
		spec = commit + ":" + dir
	}
	run := r.runner()
	out, err := run.run(ctx, r.Dir, "rev-parse", "rev-parse", "--verify", spec)
	if err != nil {
		return "", err
	}
	id := strings.TrimSpace(string(out))
	kind, err := run.run(ctx, r.Dir, "cat-file", "cat-file", "-t", id)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(string(kind)) != "tree" {
		return "", fmt.Errorf("%w: %q is not a directory", ErrInvalidPath, dir)
	}
	return id, nil
}

// ListTree lists every file and submodule under roots at commit, without
// checking anything out. With no roots it lists the whole tree. Callers decide
// what to do with the modes: ModeSymlink and ModeSubmodule entries are never
// something to extract.
func (r Repo) ListTree(ctx context.Context, commit string, roots []string) ([]TreeEntry, error) {
	if err := validateCommit(commit); err != nil {
		return nil, err
	}
	args := []string{"ls-tree", "-r", "-z", "--full-tree", commit}
	if len(roots) > 0 {
		args = append(args, "--")
		for _, root := range roots {
			if err := validateRepoPath(root); err != nil {
				return nil, err
			}
			args = append(args, root)
		}
	}
	out, err := r.runner().run(ctx, r.Dir, "ls-tree", args...)
	if err != nil {
		return nil, err
	}
	return parseTree(out)
}

func parseTree(out []byte) ([]TreeEntry, error) {
	var entries []TreeEntry
	for rec := range bytes.SplitSeq(out, []byte{0}) {
		if len(rec) == 0 {
			continue
		}
		meta, name, ok := bytes.Cut(rec, []byte{'\t'})
		fields := strings.Fields(string(meta))
		if !ok || len(fields) != 3 {
			return nil, fmt.Errorf("unexpected ls-tree record %q", rec)
		}
		entries = append(entries, TreeEntry{Mode: fields[0], Type: fields[1], ObjectID: fields[2], Path: string(name)})
	}
	return entries, nil
}

// Archive extracts the directory dir of commit into dest, with the directory
// prefix stripped. It refuses a symlink, a submodule, any path that leaves dest,
// and anything over lim; on refusal nothing is left at dest. dest must not
// already exist.
func (r Repo) Archive(ctx context.Context, commit, dir, dest string, lim Limits) error {
	if err := validateCommit(commit); err != nil {
		return err
	}
	if err := validateRepoPath(dir); err != nil {
		return err
	}
	return r.runner().stream(ctx, r.Dir, "archive",
		[]string{"archive", "--format=tar", commit + ":" + dir},
		func(src io.Reader) error { return extractTar(src, dest, lim) },
	)
}

func (r Repo) runner() runner {
	return runner{allowFile: r.allowFile}
}

func (r Repo) ensureInit(ctx context.Context) error {
	if _, err := os.Stat(filepath.Join(r.Dir, "HEAD")); err == nil {
		return nil
	}
	abs, err := filepath.Abs(r.Dir)
	if err != nil {
		return fmt.Errorf("resolve %s: %w", r.Dir, err)
	}
	if err = os.MkdirAll(filepath.Dir(abs), fsperm.Dir); err != nil {
		return fmt.Errorf("create %s: %w", filepath.Dir(abs), err)
	}
	_, err = r.runner().run(ctx, "", "init", "init", "--bare", "--quiet", "--", abs)
	return err
}
