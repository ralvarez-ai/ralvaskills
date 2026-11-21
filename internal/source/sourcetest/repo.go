// Package sourcetest provides an in-memory git clone for tests of code that sits
// above internal/source, so they can exercise a git source without a network.
package sourcetest

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
	"go.ralva.dev/ralvaskills/v3/internal/git"
)

// Repo is a repository holding skills at skills/<name>. A name maps to the
// version its SKILL.md declares; an empty version leaves the line out, as the
// official skills do.
type Repo struct {
	// Commit is what every fetch resolves to.
	Commit string

	// Skills maps a skill name to its declared version.
	Skills map[string]string

	// FetchErr, when set, makes every fetch fail.
	FetchErr error

	// Fetches records each ref asked for.
	Fetches []string
}

func skillMD(name, version string) string {
	body := "---\nname: " + name + "\n"
	if version != "" {
		body += "version: " + version + "\n"
	}
	return body + "description: test skill " + name + "\n---\n"
}

// Fetch implements source.RepoOps.
func (r *Repo) Fetch(_ context.Context, _ git.RemoteURL, ref string) (string, error) {
	r.Fetches = append(r.Fetches, ref)
	if r.FetchErr != nil {
		return "", r.FetchErr
	}
	return r.Commit, nil
}

// ListTree implements source.RepoOps. It ignores roots and lists every skill.
func (r *Repo) ListTree(context.Context, string, []string) ([]git.TreeEntry, error) {
	entries := make([]git.TreeEntry, 0, 1+2*len(r.Skills))
	entries = append(entries, git.TreeEntry{Mode: "040000", Type: "tree", ObjectID: "t-skills", Path: "skills"})
	for name := range r.Skills {
		dir := "skills/" + name
		entries = append(entries,
			git.TreeEntry{Mode: "040000", Type: "tree", ObjectID: "t-" + dir, Path: dir},
			git.TreeEntry{Mode: "100644", Type: "blob", ObjectID: "b-" + dir, Path: dir + "/SKILL.md"},
		)
	}
	return entries, nil
}

// TreeID implements source.RepoOps.
func (r *Repo) TreeID(_ context.Context, _, dir string) (string, error) {
	return "t-" + dir, nil
}

// ReadBlobs implements source.RepoOps.
func (r *Repo) ReadBlobs(_ context.Context, _ string, paths []string, _ int64) (map[string]git.Blob, error) {
	out := make(map[string]git.Blob, len(paths))
	for _, p := range paths {
		name := strings.TrimSuffix(strings.TrimPrefix(p, "skills/"), "/SKILL.md")
		version, ok := r.Skills[name]
		if !ok {
			out[p] = git.Blob{Missing: true}
			continue
		}
		body := skillMD(name, version)
		out[p] = git.Blob{Data: []byte(body), Size: int64(len(body))}
	}
	return out, nil
}

// Archive implements source.RepoOps: it writes the skill's SKILL.md to dest.
func (r *Repo) Archive(_ context.Context, _, dir, dest string, _ git.Limits) error {
	name := strings.TrimPrefix(dir, "skills/")
	if err := os.MkdirAll(dest, fsperm.Dir); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dest, "SKILL.md"), []byte(skillMD(name, r.Skills[name])), fsperm.File)
}
