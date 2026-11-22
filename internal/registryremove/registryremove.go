// Package registryremove edits registry/index.json to drop a skill, a single
// version, or an inclusive version range, then deletes the matching
// "<skill>@vX.Y.Z" GitHub tags and releases so the registry and GitHub stay
// in sync. Without --yes it only prints the plan. The registry-remove command
// is its flag shim.
package registryremove

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.ralva.dev/ralvaskills/v3/internal/fsx"
)

const (
	tempFilePattern   = ".rsk-rm-*.tmp"
	gitHubRepoDefault = "ralvarezdev/ralvaskills"
)

var semverRe = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

type (
	// index is the top-level shape of registry/index.json.
	index struct {
		Version     int                    `json:"version"`
		GeneratedAt string                 `json:"generated_at"`
		Aliases     map[string]string      `json:"aliases,omitempty"`
		Skills      map[string]*skillEntry `json:"skills"`
	}

	// skillEntry records a single skill's metadata and version history.
	skillEntry struct {
		Name        string                   `json:"name"`
		Description string                   `json:"description"`
		Personal    bool                     `json:"personal,omitempty"`
		Deprecated  bool                     `json:"deprecated,omitempty"`
		Latest      string                   `json:"latest"`
		Versions    map[string]*versionEntry `json:"versions"`
	}

	// versionEntry records a single published version of a skill.
	versionEntry struct {
		Version     string `json:"version"`
		PublishedAt string `json:"published_at"`
		ArchiveURL  string `json:"archive_url"`
	}
)

// GitHubRepoDefault is the repository whose tags and releases are deleted
// when Options.Repo is empty.
const GitHubRepoDefault = gitHubRepoDefault

// Options selects what Run removes.
type Options struct {
	// IndexPath is the registry index to edit.
	IndexPath string
	// Skill is the skill to remove from; required.
	Skill string
	// Version removes only this exact version; exclusive with Min and Max.
	Version string
	// Min and Max bound an inclusive version range.
	Min, Max string
	// Repo is the GitHub repo (owner/name) whose tags and releases are deleted.
	Repo string
	// SkipGitHub edits the index only and leaves GitHub untouched.
	SkipGitHub bool
	// DryRun prints the plan and changes nothing.
	DryRun bool
}

// Run validates the options, prints the removal plan and, unless DryRun is set,
// deletes the GitHub tags and rewrites the index.
func Run(o Options) error {
	return run(o.IndexPath, o.Skill, o.Version, o.Min, o.Max, o.Repo, o.SkipGitHub, o.DryRun)
}

func run(indexPath, skillName, version, minVersion, maxVersion, repo string, skipGitHub, dryRun bool) error {
	if skillName == "" {
		return errors.New("--skill is required")
	}
	if version != "" && (minVersion != "" || maxVersion != "") {
		return errors.New("--version cannot be combined with --min/--max")
	}
	for _, v := range []string{version, minVersion, maxVersion} {
		if v != "" && !semverRe.MatchString(v) {
			return fmt.Errorf("invalid version %q — expected X.Y.Z", v)
		}
	}

	idx, err := loadIndex(indexPath)
	if err != nil {
		return err
	}

	entry, ok := idx.Skills[skillName]
	if !ok {
		return fmt.Errorf("skill %q not found in %s", skillName, indexPath)
	}

	toRemove, err := selectVersions(entry, version, minVersion, maxVersion)
	if err != nil {
		return fmt.Errorf("%s: %w", skillName, err)
	}

	tags := make([]string, len(toRemove))
	for i, v := range toRemove {
		tags[i] = fmt.Sprintf("%s@v%s", skillName, v)
	}
	sort.Strings(tags)

	wholeSkill := len(toRemove) == len(entry.Versions)
	printPlan(skillName, tags, wholeSkill, dryRun)

	if dryRun {
		return nil
	}

	if !skipGitHub {
		if err = deleteGitHubTags(tags, repo); err != nil {
			return err
		}
	}

	applyRemoval(idx, skillName, entry, toRemove, wholeSkill)
	if err = writeIndex(indexPath, idx); err != nil {
		return err
	}

	for _, tag := range tags {
		fmt.Println(tag)
	}
	return nil
}

// selectVersions resolves which of entry's versions match the requested
// exact version or [min, max] range. An empty version and empty range select
// every version, i.e. the whole skill.
func selectVersions(entry *skillEntry, version, minVersion, maxVersion string) ([]string, error) {
	if version != "" {
		if _, ok := entry.Versions[version]; !ok {
			return nil, fmt.Errorf("version %s not found", version)
		}
		return []string{version}, nil
	}

	if minVersion == "" && maxVersion == "" {
		all := make([]string, 0, len(entry.Versions))
		for v := range entry.Versions {
			all = append(all, v)
		}
		return all, nil
	}

	var matched []string
	for v := range entry.Versions {
		if minVersion != "" && cmpSemver(v, minVersion) < 0 {
			continue
		}
		if maxVersion != "" && cmpSemver(v, maxVersion) > 0 {
			continue
		}
		matched = append(matched, v)
	}
	if len(matched) == 0 {
		return nil, fmt.Errorf("no versions in range [%s, %s]", displayBound(minVersion), displayBound(maxVersion))
	}
	return matched, nil
}

func displayBound(v string) string {
	if v == "" {
		return "*"
	}
	return v
}

// applyRemoval deletes the selected versions from entry, drops the skill
// entirely if none remain, and otherwise recomputes Latest.
func applyRemoval(idx *index, skillName string, entry *skillEntry, toRemove []string, wholeSkill bool) {
	if wholeSkill {
		delete(idx.Skills, skillName)
		return
	}

	for _, v := range toRemove {
		delete(entry.Versions, v)
	}

	latest := ""
	for v := range entry.Versions {
		if latest == "" || cmpSemver(v, latest) > 0 {
			latest = v
		}
	}
	entry.Latest = latest
}

func printPlan(skillName string, tags []string, wholeSkill, dryRun bool) {
	prefix := "will remove"
	if dryRun {
		prefix = "[dry-run] would remove"
	}
	for _, tag := range tags {
		fmt.Fprintf(os.Stderr, "%s %s\n", prefix, tag)
	}
	if wholeSkill {
		fmt.Fprintf(os.Stderr, "%s: skill entry for %q dropped entirely\n", prefix, skillName)
	}
}

// deleteGitHubTags deletes the GitHub release and tag for each entry in tags
// via the gh CLI. It aborts on the first failure, before any index.json
// changes are written, so a failed run can be safely retried.
func deleteGitHubTags(tags []string, repo string) error {
	if len(tags) == 0 {
		return nil
	}
	if _, err := exec.LookPath("gh"); err != nil {
		return errors.New("gh CLI not found on PATH — install it or pass --skip-github")
	}
	for _, tag := range tags {
		fmt.Fprintf(os.Stderr, "→ deleting release/tag %s\n", tag)
		cmd := exec.CommandContext(
			context.Background(),
			"gh",
			"release",
			"delete",
			tag,
			"--repo",
			repo,
			"--yes",
			"--cleanup-tag",
		)
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("delete release %s: %w", tag, err)
		}
	}
	return nil
}

func loadIndex(path string) (*index, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var idx index
	if err = json.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if idx.Skills == nil {
		idx.Skills = make(map[string]*skillEntry)
	}
	return &idx, nil
}

func writeIndex(path string, idx *index) error {
	idx.GeneratedAt = time.Now().UTC().Format(time.RFC3339)
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	data = append(data, '\n')
	return fsx.WriteAtomic(path, tempFilePattern, func(w io.Writer) error {
		_, writeErr := w.Write(data)
		return writeErr
	})
}

// cmpSemver compares two X.Y.Z version strings, returning -1, 0, or 1.
// Malformed components compare as 0, since every version reaching this point
// was already validated against semverRe upstream by generate-registry.
func cmpSemver(a, b string) int {
	pa, pb := strings.Split(a, "."), strings.Split(b, ".")
	for i := range 3 {
		na, nb := semverPart(pa, i), semverPart(pb, i)
		if na != nb {
			if na < nb {
				return -1
			}
			return 1
		}
	}
	return 0
}

func semverPart(parts []string, i int) int {
	if i >= len(parts) {
		return 0
	}
	n, err := strconv.Atoi(parts[i])
	if err != nil {
		return 0
	}
	return n
}
