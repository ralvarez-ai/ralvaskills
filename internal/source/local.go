package source

import (
	"context"
	"fmt"
	"path/filepath"

	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

// Local resolves skills from the local ralvaskills repo clone.
type Local struct {
	spec     Spec
	repoPath string
	fs       fsSource
}

// newLocal returns a Local source rooted at spec.URL.
func newLocal(spec Spec) *Local {
	repoPath := spec.URL
	root := filepath.Join(repoPath, skill.SkillsFolderName)
	return &Local{
		spec:     spec,
		repoPath: repoPath,
		fs: fsSource{
			root: root,
			src:  skill.SourceLocal,
			notFound: func(name string) error {
				return fmt.Errorf("%w: local skill %q not found in repo at %s", ErrNotFound, name, repoPath)
			},
		},
	}
}

// Spec returns the descriptor the source was built from.
func (l *Local) Spec() Spec {
	return l.spec
}

// SkillsRoot returns the absolute path to the skills/ directory.
func (l *Local) SkillsRoot() string {
	return l.fs.root
}

// All walks the local skills directory and returns every discovered skill.
func (l *Local) All(ctx context.Context) ([]skill.Skill, error) {
	return l.fs.all(ctx)
}

// Find returns the local skill with the given name.
func (l *Local) Find(ctx context.Context, name string) (skill.Skill, error) {
	return l.fs.find(ctx, name)
}
