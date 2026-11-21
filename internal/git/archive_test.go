package git

import (
	"archive/tar"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type tarEntry struct {
	hdr  tar.Header
	body string
}

func buildTar(t *testing.T, entries ...tarEntry) *bytes.Reader {
	t.Helper()
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for i := range entries {
		e := &entries[i]
		h := e.hdr
		if h.Typeflag == tar.TypeReg || h.Typeflag == 0 {
			h.Typeflag = tar.TypeReg
			h.Size = int64(len(e.body))
		}
		if err := tw.WriteHeader(&h); err != nil {
			t.Fatal(err)
		}
		if h.Typeflag == tar.TypeReg {
			if _, err := tw.Write([]byte(e.body)); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(buf.Bytes())
}

func reg(name, body string, mode int64) tarEntry {
	return tarEntry{hdr: tar.Header{Name: name, Mode: mode, Typeflag: tar.TypeReg}, body: body}
}

func TestExtractTarWritesFilesAndDirectories(t *testing.T) {
	t.Parallel()
	dest := filepath.Join(t.TempDir(), "out")
	src := buildTar(t,
		tarEntry{hdr: tar.Header{Name: "pax_global_header", Typeflag: tar.TypeXGlobalHeader,
			PAXRecords: map[string]string{"comment": "0123456789abcdef"}}},
		tarEntry{hdr: tar.Header{Name: "refs/", Typeflag: tar.TypeDir, Mode: 0o755}},
		reg("SKILL.md", "hello", 0o644),
		reg("refs/a.md", "a", 0o666),
		reg("scripts/run.sh", "#!/bin/sh\n", 0o755),
	)
	if err := extractTar(src, dest, Limits{}); err != nil {
		t.Fatal(err)
	}

	if got, _ := os.ReadFile(filepath.Join(dest, "SKILL.md")); string(got) != "hello" {
		t.Errorf("SKILL.md = %q", got)
	}
	if info, err := os.Stat(filepath.Join(dest, "scripts", "run.sh")); err != nil || info.Mode().Perm()&0o100 == 0 {
		t.Errorf("run.sh lost its exec bit: %v %v", info, err)
	}
	if info, _ := os.Stat(filepath.Join(dest, "refs", "a.md")); info.Mode().Perm()&0o022 != 0 {
		t.Errorf("a.md is group/world writable: %v", info.Mode())
	}
	if _, err := os.Stat(filepath.Join(dest, "pax_global_header")); err == nil {
		t.Error("the pax global header was extracted as a file")
	}
}

func TestExtractTarRefusesUnsafeEntries(t *testing.T) {
	t.Parallel()
	tests := map[string]tarEntry{
		"parent traversal":   reg("../evil", "x", 0o644),
		"nested traversal":   reg("a/../../evil", "x", 0o644),
		"absolute path":      reg("/etc/evil", "x", 0o644),
		"symlink":            {hdr: tar.Header{Name: "link", Typeflag: tar.TypeSymlink, Linkname: "/etc/passwd"}},
		"hard link":          {hdr: tar.Header{Name: "hard", Typeflag: tar.TypeLink, Linkname: "SKILL.md"}},
		"character device":   {hdr: tar.Header{Name: "dev", Typeflag: tar.TypeChar}},
		"fifo":               {hdr: tar.Header{Name: "fifo", Typeflag: tar.TypeFifo}},
		"empty name":         reg("", "x", 0o644),
		"backslash traverse": reg("..\\evil", "x", 0o644),
	}
	for name, e := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			dest := filepath.Join(root, "out")
			err := extractTar(buildTar(t, reg("SKILL.md", "ok", 0o644), e), dest, Limits{})
			if !errors.Is(err, ErrUnsafeEntry) {
				t.Fatalf("err = %v, want ErrUnsafeEntry", err)
			}
			if _, statErr := os.Stat(dest); statErr == nil {
				t.Error("a refused archive left a destination behind")
			}
			if _, statErr := os.Stat(filepath.Join(root, "evil")); statErr == nil {
				t.Error("an entry escaped the destination")
			}
			entries, _ := os.ReadDir(root)
			if len(entries) != 0 {
				t.Errorf("temp files left in %s: %v", root, entries)
			}
		})
	}
}

func TestExtractTarEnforcesLimits(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		lim     Limits
		entries []tarEntry
	}{
		"one file too big": {Limits{MaxFileBytes: 4}, []tarEntry{reg("a", "12345", 0o644)}},
		"total too big":    {Limits{MaxTotalBytes: 6}, []tarEntry{reg("a", "1234", 0o644), reg("b", "1234", 0o644)}},
		"too many files":   {Limits{MaxFiles: 1}, []tarEntry{reg("a", "1", 0o644), reg("b", "1", 0o644)}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			dest := filepath.Join(t.TempDir(), "out")
			err := extractTar(buildTar(t, tc.entries...), dest, tc.lim)
			if !errors.Is(err, ErrTooLarge) {
				t.Fatalf("err = %v, want ErrTooLarge", err)
			}
			if _, statErr := os.Stat(dest); statErr == nil {
				t.Error("an oversized archive left a destination behind")
			}
		})
	}
}

func TestExtractTarRefusesAnExistingDestination(t *testing.T) {
	t.Parallel()
	dest := t.TempDir()
	if err := extractTar(buildTar(t, reg("a", "1", 0o644)), dest, Limits{}); err == nil ||
		!strings.Contains(err.Error(), "exists") {
		t.Errorf("err = %v, want a refusal because the destination exists", err)
	}
}
