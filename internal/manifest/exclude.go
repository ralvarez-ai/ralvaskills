package manifest

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
)

// localFiles are the files that carry private skills, relative to the project
// root: the local layers, the generated digest and Claude Code's local import.
var localFiles = []string{
	".rsk/" + LocalModFileName, ".rsk/" + LocalLockFileName, ".rsk/SKILLS.local.md", "CLAUDE.local.md",
}

// ExcludeLocalFiles appends the local files to the repository's
// .git/info/exclude, not to the shared .gitignore: whether a developer keeps
// private skills is not the project's business. It is idempotent, and does
// nothing outside a git repository, where nothing can be committed. links are
// extra paths, relative to the project root, to exclude with them: the installed
// links of private skills, which are named after the skill.
func ExcludeLocalFiles(projectRoot string, links ...string) error {
	root, excludePath, ok := findGitExclude(projectRoot)
	if !ok {
		return nil
	}
	prefix := ""
	if rel, err := filepath.Rel(root, projectRoot); err == nil && rel != "." {
		prefix = filepath.ToSlash(rel) + "/"
	}

	have := map[string]bool{}
	if data, readErr := os.ReadFile(excludePath); readErr == nil {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			have[strings.TrimSpace(scanner.Text())] = true
		}
	}
	var add []string
	for _, file := range append(slices.Clone(localFiles), links...) {
		if line := "/" + prefix + file; !have[line] && !have[file] {
			add = append(add, line)
		}
	}
	if len(add) == 0 {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(excludePath), fsperm.Dir); err != nil {
		return err
	}
	f, err := os.OpenFile(excludePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fsperm.File)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	_, err = f.WriteString("\n# rsk: private skills stay on this machine\n" + strings.Join(add, "\n") + "\n")
	return err
}

// findGitExclude walks up from dir to the repository that holds it and returns
// its root and the path of its info/exclude. A worktree keeps that file in the
// common git directory, which its .git file points at. It reads files only, so
// it needs no git binary.
func findGitExclude(dir string) (root, exclude string, ok bool) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", "", false
	}
	for cur := abs; ; cur = filepath.Dir(cur) {
		dotGit := filepath.Join(cur, ".git")
		if info, statErr := os.Stat(dotGit); statErr == nil {
			gitDir := dotGit
			if !info.IsDir() {
				if gitDir, ok = gitDirFromFile(cur, dotGit); !ok {
					return "", "", false
				}
			}
			if common, readErr := os.ReadFile(filepath.Join(gitDir, "commondir")); readErr == nil {
				gitDir = filepath.Join(gitDir, strings.TrimSpace(string(common)))
			}
			return cur, filepath.Join(gitDir, "info", "exclude"), true
		}
		if filepath.Dir(cur) == cur {
			return "", "", false
		}
	}
}

// gitDirFromFile reads the "gitdir: <path>" line of a worktree's .git file.
func gitDirFromFile(base, file string) (string, bool) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", false
	}
	target, found := strings.CutPrefix(strings.TrimSpace(string(data)), "gitdir:")
	if !found {
		return "", false
	}
	target = strings.TrimSpace(target)
	if !filepath.IsAbs(target) {
		target = filepath.Join(base, target)
	}
	return target, true
}
