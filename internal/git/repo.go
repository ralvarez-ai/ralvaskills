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
	"strconv"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
)

var (
	// ErrInvalidRef is returned for a git ref the runner will not pass on.
	ErrInvalidRef = errors.New("invalid git ref")

	// ErrInvalidCommit is returned for a value that is not a full object id.
	ErrInvalidCommit = errors.New("invalid commit id")

	// ErrRefNotFound is returned when a remote has no such branch or tag.
	ErrRefNotFound = errors.New("ref not found on the remote")

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

// TreeEntry is one entry of a recursive tree listing: a blob, a submodule, or a
// directory (Type "tree", with the directory's own object id).
type TreeEntry struct {
	Mode     string
	Type     string
	ObjectID string
	Path     string
}

// ValidateRef reports why ref (a branch, tag or object id) is not safe to pass
// to git, or nil.
func ValidateRef(ref string) error {
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

// ValidateRepoPath reports why p is not a relative, confined path inside a
// repository, or nil.
func ValidateRepoPath(p string) error {
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
	if err := ValidateRef(ref); err != nil {
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
		if err := ValidateRepoPath(dir); err != nil {
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

// ListTree lists every directory, file and submodule under roots at commit,
// without checking anything out. With no roots it lists the whole tree. Callers decide
// what to do with the modes: ModeSymlink and ModeSubmodule entries are never
// something to extract.
func (r Repo) ListTree(ctx context.Context, commit string, roots []string) ([]TreeEntry, error) {
	if err := validateCommit(commit); err != nil {
		return nil, err
	}
	args := []string{"ls-tree", "-r", "-t", "-z", "--full-tree", commit}
	if len(roots) > 0 {
		args = append(args, "--")
		for _, root := range roots {
			if err := ValidateRepoPath(root); err != nil {
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

// Blob is the content of one file read from a commit.
type Blob struct {
	// Data is the file content. It is nil when the file is missing or larger
	// than the limit given to ReadBlobs.
	Data []byte

	// Size is the file's size in bytes, when it exists.
	Size int64

	// Missing is true when the commit has no such file.
	Missing bool
}

// ReadBlobs reads the files at paths from commit in one git process. A file
// larger than maxBytes is reported with its Size and no Data, so a hostile
// repository cannot make rsk hold an arbitrarily large file in memory.
func (r Repo) ReadBlobs(ctx context.Context, commit string, paths []string, maxBytes int64) (map[string]Blob, error) {
	if err := validateCommit(commit); err != nil {
		return nil, err
	}
	var input bytes.Buffer
	for _, p := range paths {
		if err := ValidateRepoPath(p); err != nil {
			return nil, err
		}
		input.WriteString(commit + ":" + p + "\n")
	}
	if input.Len() == 0 {
		return map[string]Blob{}, nil
	}
	out, err := r.runner().runInput(ctx, r.Dir, "cat-file", input.Bytes(), "cat-file", "--batch")
	if err != nil {
		return nil, err
	}
	return parseBlobs(out, paths, maxBytes)
}

func parseBlobs(out []byte, paths []string, maxBytes int64) (map[string]Blob, error) {
	blobs := make(map[string]Blob, len(paths))
	rest := out
	for _, p := range paths {
		header, after, ok := bytes.Cut(rest, []byte{'\n'})
		if !ok {
			return nil, fmt.Errorf("cat-file output ended before %s", p)
		}
		fields := strings.Fields(string(header))
		if len(fields) == 2 && fields[1] == "missing" {
			blobs[p] = Blob{Missing: true}
			rest = after
			continue
		}
		if len(fields) != 3 || fields[1] != "blob" {
			return nil, fmt.Errorf("unexpected cat-file header %q for %s", header, p)
		}
		size, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil || size < 0 || int64(len(after)) < size+1 {
			return nil, fmt.Errorf("bad cat-file size in %q for %s", header, p)
		}
		blob := Blob{Size: size}
		if size <= maxBytes {
			blob.Data = append([]byte(nil), after[:size]...)
		}
		blobs[p] = blob
		rest = after[size+1:]
	}
	return blobs, nil
}

// Archive extracts the directory dir of commit into dest, with the directory
// prefix stripped. It refuses a symlink, a submodule, any path that leaves dest,
// and anything over lim; on refusal nothing is left at dest. dest must not
// already exist.
func (r Repo) Archive(ctx context.Context, commit, dir, dest string, lim Limits) error {
	if err := validateCommit(commit); err != nil {
		return err
	}
	if err := ValidateRepoPath(dir); err != nil {
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

// LsRemote reports the commit ref (a branch or tag; empty means HEAD) points to
// on remote, without fetching anything. It is how a source is tested for
// reachability. It contacts the host, so a caller must have checked that the
// source is trusted.
func LsRemote(ctx context.Context, remote RemoteURL, ref string) (string, error) {
	return runner{}.lsRemote(ctx, remote, ref)
}

func (r runner) lsRemote(ctx context.Context, remote RemoteURL, ref string) (string, error) {
	if ref == "" {
		ref = "HEAD"
	} else if err := ValidateRef(ref); err != nil {
		return "", err
	}
	out, err := r.run(ctx, "", "ls-remote", "ls-remote", "--exit-code", "--", remote.String(), ref)
	if err != nil {
		var gitErr *Error
		if errors.As(err, &gitErr) && strings.Contains(gitErr.Err.Error(), "exit status 2") {
			return "", fmt.Errorf("%w: %s", ErrRefNotFound, ref)
		}
		return "", err
	}
	line, _, _ := strings.Cut(string(out), "\n")
	commit, _, _ := strings.Cut(line, "\t")
	if err = validateCommit(commit); err != nil {
		return "", err
	}
	return commit, nil
}
