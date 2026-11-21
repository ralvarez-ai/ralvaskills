package git

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

// MinVersion is the oldest git the hardened operations are written against:
// fetching by object id, `archive <tree-ish>:<dir>` and the protocol.*.allow
// settings.
var MinVersion = [3]int{2, 30, 0}

var versionPattern = regexp.MustCompile(`git version (\d+)\.(\d+)(?:\.(\d+))?`)

// Version reports the installed git's version.
func Version(ctx context.Context) ([3]int, error) {
	out, err := exec.CommandContext(ctx, "git", "--version").Output()
	if err != nil {
		return [3]int{}, fmt.Errorf("run git --version: %w", err)
	}
	return parseVersion(string(out))
}

func parseVersion(s string) ([3]int, error) {
	m := versionPattern.FindStringSubmatch(s)
	if m == nil {
		return [3]int{}, fmt.Errorf("cannot read a git version from %q", s)
	}
	var v [3]int
	for i := range v {
		if m[i+1] == "" {
			continue
		}
		n, err := strconv.Atoi(m[i+1])
		if err != nil {
			return [3]int{}, fmt.Errorf("git version %q: %w", s, err)
		}
		v[i] = n
	}
	return v, nil
}

// MeetsMinimum reports whether v is at least MinVersion.
func MeetsMinimum(v [3]int) bool {
	return versionAtLeast(v, MinVersion)
}

func versionAtLeast(have, want [3]int) bool {
	for i := range have {
		if have[i] != want[i] {
			return have[i] > want[i]
		}
	}
	return true
}
