package skill

import (
	"path/filepath"
	"strings"
)

// gitLinkParts is the number of path segments under <cache>/sources that name a
// git extraction: owner, repo, skill name and commit.
const gitLinkParts = 4

// DetectSource determines the origin of a linked skill by matching its path
// against the known cache roots. It falls back to SourceRegistry in registry
// mode (repoPath == "") and SourceLocal in local-clone mode.
func DetectSource(skillPath, repoPath, officialCachePath, registryCachePath string) Source {
	sp := filepath.ToSlash(filepath.Clean(skillPath))
	if repoPath != "" && strings.HasPrefix(sp, filepath.ToSlash(filepath.Clean(repoPath))+"/") {
		return SourceLocal
	}
	if registryCachePath != "" && strings.HasPrefix(sp, filepath.ToSlash(filepath.Clean(registryCachePath))+"/") {
		return SourceRegistry
	}
	if officialCachePath != "" && strings.HasPrefix(sp, filepath.ToSlash(filepath.Clean(officialCachePath))+"/") {
		return SourceOfficial
	}
	if repoPath == "" {
		return SourceRegistry
	}
	return SourceLocal
}

// ParseGitLink reads the identity, name and commit out of the target of a link
// into the git cache: <root>/sources/<owner>/<repo>/<name>/<commit>.
func ParseGitLink(dest, cacheRoot string) (id, name, commit string, ok bool) {
	prefix := filepath.Join(cacheRoot, "sources") + string(filepath.Separator)
	rest, found := strings.CutPrefix(filepath.Clean(dest), prefix)
	if !found {
		return "", "", "", false
	}
	parts := strings.Split(rest, string(filepath.Separator))
	if len(parts) != gitLinkParts {
		return "", "", "", false
	}
	return "@" + parts[0] + "/" + parts[1], parts[2], parts[3], true
}
