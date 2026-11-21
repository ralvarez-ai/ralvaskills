package main

import (
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/source"
	"go.ralva.dev/ralvaskills/v3/internal/source/sourcetest"
)

// withOfficialSkills makes the official source read an in-memory clone holding
// skills (name to declared version) for the length of the test, so nothing
// reaches GitHub. The seam is process-wide, so the test using it is sequential.
func withOfficialSkills(t *testing.T, skills map[string]string) *sourcetest.Repo {
	t.Helper()
	repo := &sourcetest.Repo{Commit: strings.Repeat("c", 40), Skills: skills}
	sourceRepoFor = func(source.Spec) source.RepoOps { return repo }
	t.Cleanup(func() { sourceRepoFor = nil })
	return repo
}
