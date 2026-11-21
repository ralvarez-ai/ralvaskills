package git

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
)

var (
	// ErrUnsafeEntry is returned when an archive holds an entry that is not a
	// plain file or directory, or whose path leaves the destination.
	ErrUnsafeEntry = errors.New("unsafe archive entry")

	// ErrTooLarge is returned when an archive exceeds a Limits bound.
	ErrTooLarge = errors.New("archive exceeds its size limits")
)

const (
	defaultMaxTotalBytes = 32 << 20
	defaultMaxFileBytes  = 8 << 20
	defaultMaxFiles      = 2000
	execFileMode         = 0o755
)

// Limits bound what an extraction may write. A zero field takes its default.
type Limits struct {
	MaxTotalBytes int64
	MaxFileBytes  int64
	MaxFiles      int
}

func (l Limits) withDefaults() Limits {
	if l.MaxTotalBytes == 0 {
		l.MaxTotalBytes = defaultMaxTotalBytes
	}
	if l.MaxFileBytes == 0 {
		l.MaxFileBytes = defaultMaxFileBytes
	}
	if l.MaxFiles == 0 {
		l.MaxFiles = defaultMaxFiles
	}
	return l
}

// extractor writes the entries of one archive into root and counts what it
// has written against the limits.
type extractor struct {
	root  string
	lim   Limits
	files int
	total int64
}

// extractTar writes the entries of a tar stream into dest. Only regular files
// and directories are accepted: a symlink, hard link, device or any path that
// leaves dest refuses the whole archive. The result is built in a sibling temp
// directory and renamed into place, so a refused or oversized archive leaves
// nothing behind.
func extractTar(r io.Reader, dest string, lim Limits) (err error) {
	if _, statErr := os.Lstat(dest); statErr == nil {
		return fmt.Errorf("extract to %s: destination exists", dest)
	}
	parent := filepath.Dir(dest)
	if err = os.MkdirAll(parent, fsperm.Dir); err != nil {
		return fmt.Errorf("create %s: %w", parent, err)
	}
	tmp, err := os.MkdirTemp(parent, ".extract-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer func() {
		if err != nil {
			_ = os.RemoveAll(tmp)
		}
	}()

	x := &extractor{root: tmp, lim: lim.withDefaults()}
	tr := tar.NewReader(r)
	for {
		hdr, nextErr := tr.Next()
		if errors.Is(nextErr, io.EOF) {
			break
		}
		if nextErr != nil {
			return fmt.Errorf("read archive: %w", nextErr)
		}
		if err = x.add(tr, hdr); err != nil {
			return err
		}
	}

	if err = os.Rename(tmp, dest); err != nil {
		return fmt.Errorf("move extraction into place: %w", err)
	}
	return nil
}

// add validates one entry and writes it.
func (x *extractor) add(r io.Reader, hdr *tar.Header) error {
	if hdr.Typeflag == tar.TypeXGlobalHeader {
		return nil
	}
	target, rel, err := entryTarget(x.root, hdr)
	if err != nil || rel == "" {
		return err
	}
	if hdr.Typeflag == tar.TypeDir {
		if err = os.MkdirAll(target, fsperm.Dir); err != nil {
			return fmt.Errorf("create %s: %w", rel, err)
		}
		return nil
	}

	x.files++
	x.total += hdr.Size
	if x.files > x.lim.MaxFiles || hdr.Size > x.lim.MaxFileBytes || x.total > x.lim.MaxTotalBytes {
		return fmt.Errorf("%w: %s", ErrTooLarge, rel)
	}
	if err = writeEntry(target, r, hdr); err != nil {
		return fmt.Errorf("write %s: %w", rel, err)
	}
	return nil
}

// entryTarget validates an entry and returns the path it would be written to,
// plus its slash-separated path relative to the root ("" for the root itself).
func entryTarget(root string, hdr *tar.Header) (target, rel string, err error) {
	if hdr.Typeflag != tar.TypeReg && hdr.Typeflag != tar.TypeDir {
		return "", "", fmt.Errorf("%w: %q has type %q", ErrUnsafeEntry, hdr.Name, string(hdr.Typeflag))
	}
	name := hdr.Name
	if name == "" || strings.ContainsAny(name, "\\\x00") || path.IsAbs(name) {
		return "", "", fmt.Errorf("%w: %q", ErrUnsafeEntry, name)
	}
	for seg := range strings.SplitSeq(name, "/") {
		if seg == ".." {
			return "", "", fmt.Errorf("%w: %q", ErrUnsafeEntry, name)
		}
	}
	rel = path.Clean(name)
	if rel == "." {
		return "", "", nil
	}
	target = filepath.Join(root, filepath.FromSlash(rel))
	if !strings.HasPrefix(target, root+string(filepath.Separator)) {
		return "", "", fmt.Errorf("%w: %q", ErrUnsafeEntry, name)
	}
	return target, rel, nil
}

func writeEntry(target string, r io.Reader, hdr *tar.Header) error {
	if err := os.MkdirAll(filepath.Dir(target), fsperm.Dir); err != nil {
		return err
	}
	mode := os.FileMode(fsperm.File)
	if hdr.Mode&0o111 != 0 {
		mode = execFileMode
	}
	f, err := os.OpenFile(target, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	if _, err = io.CopyN(f, r, hdr.Size); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}
